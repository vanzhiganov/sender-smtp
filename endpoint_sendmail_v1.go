package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"mime/quotedprintable"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	uuid "github.com/satori/go.uuid"
)

// Sender ...
type Sender struct {
	Login    string
	Password string
}

var config Configuration

// NewSender ...
func NewSender(Username, Password string) Sender {
	return Sender{Username, Password}
}

// SendMail ...
func (sender Sender) SendMail(server string, port int, login string, password string, Dest []string, Subject, bodyMessage string) {

	msg := "From: " + sender.Login + "\n" +
		"To: " + strings.Join(Dest, ",") + "\n" +
		"Subject: " + Subject + "\n" + bodyMessage

	err := smtp.SendMail(server+":"+strconv.Itoa(port),
		smtp.PlainAuth("", login, password, server),
		sender.Login, Dest, []byte(msg))

	if err != nil {

		fmt.Printf("smtp error: %s", err)
		return
	}

	log.Println("Mail sent successfully!")
}

// WriteHTMLEmail ...
func (sender *Sender) WriteHTMLEmail(dest []string, subject, bodyMessage string) string {

	return sender.WriteEmail(dest, "text/html", subject, bodyMessage)
}

// WritePlainEmail ...
func (sender *Sender) WritePlainEmail(dest []string, subject, bodyMessage string) string {

	return sender.WriteEmail(dest, "text/plain", subject, bodyMessage)
}

// WriteEmail ...
func (sender Sender) WriteEmail(dest []string, contentType, subject, bodyMessage string) string {

	header := make(map[string]string)
	header["From"] = sender.Login

	receipient := ""

	for _, user := range dest {
		receipient = receipient + user
	}

	header["To"] = receipient
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = fmt.Sprintf("%s; charset=\"utf-8\"", contentType)
	header["Content-Transfer-Encoding"] = "quoted-printable"
	header["Content-Disposition"] = "inline"

	message := ""

	for key, value := range header {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	var encodedMessage bytes.Buffer

	finalMessage := quotedprintable.NewWriter(&encodedMessage)
	finalMessage.Write([]byte(bodyMessage))
	finalMessage.Close()

	message += "\r\n" + encodedMessage.String()

	return message
}

// APIV1Sendmail ...
func APIV1Sendmail(w http.ResponseWriter, r *http.Request) {
	var request RequestStruct
	var response ResponseStruct

	// Необходимо будет проверять ключи
	secretKey := r.Header.Get("X-Secret-Key")
	// ID проекта
	projectID := r.Header.Get("X-Project-ID")

	var smtpServer string
	var smtpPort int
	var smtpLogin string
	var smtpPassword string
	var smtpTemplate string

	//
	row, err0 := DBC.Query("SELECT server, port, sender_login, sender_password FROM smtp WHERE project_id='" + projectID + "' LIMIT 1;")
	if err0 != nil {
		log.Fatal(err0)
		w.WriteHeader(403)
		return
	}
	defer row.Close()
	for row.Next() {
		row.Scan(&smtpServer, &smtpPort, &smtpLogin, &smtpPassword)
	}

	//
	row, err1 := DBC.Query("SELECT template FROM templates WHERE project_id='" + projectID + "' LIMIT 1;")
	if err1 != nil {
		log.Fatal(err1)
		w.WriteHeader(403)
		return
	}
	defer row.Close()
	for row.Next() {
		row.Scan(&smtpTemplate)
	}

	// проверка секретного ключа в проекте
	if config.Application.SecretKey != secretKey {
		w.WriteHeader(403)
		return
	}

	request, err := readAPIBody(w, r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(422)
		return
	}

	//The receiver needs to be in slice as the receive supports multiple receiver
	Receiver := []string{request.To}
	subject := string(request.Subject)
	message := string(request.Message)

	response.ID = uuid.NewV1().String()

	//
	sender := NewSender(smtpLogin, smtpPassword)

	bodyMessage := ""

	if request.ContentType == "html" {
		doc := new(bytes.Buffer)

		t, _ := template.New("template").Parse(smtpTemplate)

		if err = t.Execute(doc, message); err != nil {
			// return err
			log.Panicln(err)
			w.WriteHeader(422)
		}

		bodyMessage = sender.WriteHTMLEmail(Receiver, subject, doc.String())
	} else {
		bodyMessage = sender.WritePlainEmail(Receiver, subject, message)
	}

	sender.SendMail(smtpServer, smtpPort, smtpLogin, smtpPassword, Receiver, subject, bodyMessage)

	data, _ := json.Marshal(response)

	fmt.Fprintf(w, string(data))
}
