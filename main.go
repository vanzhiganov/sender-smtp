package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
)

// Version Версия приложения
var Version = "0.1.0"

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

func main() {
	switch arg := os.Args[1:][0]; arg {
	case "runserver":
		log.Println("Sender SMTP Rest API v" + Version)
		config.getConf()

		if config.Sentry.Enabled {
			sentry.Init(sentry.ClientOptions{Dsn: config.Sentry.DSN})
		}

		handleRequests(config)
	case "version":
		fmt.Println(Version)
	default:
		fmt.Printf("Unexpected argument: %s.\n", arg)
	}
}
