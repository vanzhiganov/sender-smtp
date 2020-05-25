package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

// Version Версия приложения
var Version = "0.1.0"

// DBC ...
var DBC *sql.DB

func readAPIBody(w http.ResponseWriter, r *http.Request) (RequestStruct, error) {
	var request RequestStruct

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &request); err != nil {
		log.Println(err)
		// unprocessable entity
		w.WriteHeader(422)
	}
	return request, err
}

// Existing code from above
func handleRequests(config Configuration) {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// конечная точка для отправки сообения
	myRouter.HandleFunc("/api/v1", APIV1Sendmail).Methods("POST")
	// TODO: добавить коненую точку для мониторинга
	// myRouter.HandleFunc("/api/v1/status", APIV1Status).Methods("GET")

	// TODO: останавливить приложение, если не в настройках нет параметра application.listen
	log.Fatal(http.ListenAndServe(config.Application.Listen, myRouter))
}

func createTable(db *sql.DB) {
	createSMTPTableSQL := `CREATE TABLE smtp (
		"project_id" varchar(36) NOT NULL PRIMARY KEY,
		"server" varchar(64),
		"port" int,
		"sender_login" varchar(128),
		"sender_password" varchar(128)
	);` // SQL Statement for Create Table

	createTemplatesTableSQL := `CREATE TABLE templates (
		"project_id" varchar(36) NOT NULL PRIMARY KEY,
		"template" TEXT
	);`

	log.Println("Create student table...")
	smtp, err := db.Prepare(createSMTPTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	smtp.Exec() // Execute SQL Statements
	log.Println("student table created")

	log.Println("Create student table...")
	statement, err := db.Prepare(createTemplatesTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("student table created")

}

func main() {
	switch arg := os.Args[1:][0]; arg {
	case "runserver":
		log.Println("Sender SMTP Rest API v" + Version)
		config.getConf()

		if config.Sentry.Enabled {
			sentry.Init(sentry.ClientOptions{Dsn: config.Sentry.DSN})
		}

		DBC, _ = sql.Open("sqlite3", config.Application.DB)

		handleRequests(config)
	case "initdb":
		config.getConf()

		os.Remove(config.Application.DB) // I delete the file to avoid duplicated records. SQLite is a file based database.

		log.Println("Creating sqlite-database.db...")
		file, err := os.Create(config.Application.DB) // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Println("sqlite-database.db created")

		sqliteDatabase, _ := sql.Open("sqlite3", config.Application.DB) // Open the created SQLite File
		defer sqliteDatabase.Close()                                    // Defer Closing the database
		createTable(sqliteDatabase)                                     // Create Database Tables
	case "version":
		fmt.Println(Version)
	default:
		fmt.Printf("Unexpected argument: %s.\n", arg)
	}
}
