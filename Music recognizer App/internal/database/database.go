package database

import (
	"awesomeProject/internal/models"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load .env file!")
	}
}

func StartDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", os.Getenv("DATABASE"))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	return db, err
}

func closeDB(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func Insert(request models.Request) {
	db, _ := StartDB()
	defer closeDB(db)
	stmt, err := db.Prepare("INSERT INTO Request(ID,Email,Status,SongID) values(?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	result, err := stmt.Exec(request.RequestID, request.Email, request.Status, request.SongID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result.LastInsertId())
}

func SelectFromDB(value string, columnName string) ([]models.Request, error) {
	db, _ := StartDB()
	defer closeDB(db)
	var requests []models.Request
	query := "SELECT * FROM Request WHERE " + columnName + " = ?"
	rows, err := db.Query(query, value)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	for rows.Next() {
		request := models.Request{}
		err := rows.Scan(&request.RequestID, &request.Email, &request.Status, &request.SongID)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	return requests, nil
}

func Update(reqID string, value string, columnName string) {
	db, err := StartDB()
	if err != nil {
		panic(err.Error())
	}
	defer closeDB(db)
	updateStmt, err := db.Prepare("UPDATE Request SET " + columnName + " = ? " + "WHERE ID = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = updateStmt.Exec(value, reqID)
	if err != nil {
		panic(err.Error())
	}
}

func Delete(value string, columnName string) {
	db, err := StartDB()
	if err != nil {
		panic(err.Error())
	}
	defer closeDB(db)
	query := "DELETE FROM Request WHERE " + columnName + " = ?"
	result, err := db.Exec(query, value)
	if err != nil {
		panic(err)
	}
	_ = result
}
