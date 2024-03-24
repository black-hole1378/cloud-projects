package utils

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"log"
	"mime/multipart"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func CloseConnection(conn *amqp.Connection) {
	err := conn.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func CloseChannel(conn *amqp.Channel) {
	err := conn.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func Publish(message string) error {
	conn, err := amqp.Dial(os.Getenv("CLOUDAMQP_URL"))
	if err != nil {
		return err
	}
	defer CloseConnection(conn)
	cha, err := conn.Channel()
	if err != nil {
		return err
	}
	defer CloseChannel(cha)
	q, err := cha.QueueDeclare(
		"music", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return err
	}
	err = cha.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	return err
}

func Initialize() (*s3.Client, string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		return nil, "", err
	}

	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_KEY")
	endpoint := os.Getenv("AWS_ENDPOINT")
	// Initialize S3 client with custom configuration
	cfg.Credentials = aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
		return aws.Credentials{
			AccessKeyID:     awsAccessKeyID,
			SecretAccessKey: awsSecretAccessKey,
		}, nil
	})

	cfg.BaseEndpoint = aws.String(endpoint)
	bucketName := os.Getenv("AWS_BUCKET_NAME")
	client := s3.NewFromConfig(cfg)
	return client, bucketName, nil
}

func UploadFile(file *multipart.FileHeader, fileName string) error {
	client, bucketName, err := Initialize()
	if err != nil {
		return err
	}
	destinationKey := "uploads/" + fileName

	fileContent, err := file.Open()
	defer func(file multipart.File) {
		err := fileContent.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(fileContent)

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(destinationKey),
		Body:   fileContent,
	})
	return err
}
