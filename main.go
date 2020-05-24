package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
)

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
	log.Println("Sender SMTP Rest API v0.1")
	config.getConf()

	if config.Sentry.Enabled {
		sentry.Init(sentry.ClientOptions{Dsn: config.Sentry.DSN})
	}

	handleRequests(config)
}