package main

import (
	"awesomeProject/internal/database"
	"awesomeProject/internal/request"
	"awesomeProject/internal/utils"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"io"
	"log"
	"os"
)

func getFile(filename string) (io.ReadCloser, error) {
	client, bucketName, err := utils.Initialize()
	if err != nil {
		return nil, err
	}
	result, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("uploads/" + filename),
	})
	if err != nil {
		fmt.Println(err)
	}
	return result.Body, nil
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func consumer() {
	conn, err := amqp.Dial(os.Getenv("CLOUDAMQP_URL"))

	failOnError(err, "Failed to connect to RabbitMQ")
	defer utils.CloseConnection(conn)

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer utils.CloseChannel(ch)

	// Declare a queue
	q, err := ch.QueueDeclare(
		"music", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Consume messages from the queue
	messages, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	for msg := range messages {
		if err := handle(string(msg.Body)); err != nil {
			log.Print(err)
		}
	}
}

func isError(err error, reqID string) bool {
	if err != nil {
		log.Println(err)
		database.Update(reqID, "failure", "Status")
		return true
	}
	return false
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func handle(msg string) error {
	audio, _ := getFile(string(msg) + ".mp3")
	fmt.Println("New request come in!")
	track, err := request.SendRequestShazam(audio)
	if isError(err, msg) {
		return err
	}
	fmt.Println("track", track)
	songID, err := request.SendRequestSpotify(track)
	fmt.Println("songID", songID)
	if isError(err, msg) {
		return err
	}
	database.Update(msg, "ready", "Status")
	database.Update(msg, songID, "SongID")
	return nil
}
