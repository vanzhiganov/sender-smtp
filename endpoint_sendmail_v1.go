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
func (sender Sender) SendMail(conf Configuration, Dest []string, Subject, bodyMessage string) {

	msg := "From: " + sender.Login + "\n" +
		"To: " + strings.Join(Dest, ",") + "\n" +
		"Subject: " + Subject + "\n" + bodyMessage

	err := smtp.SendMail(conf.SMTP.Server+":"+strconv.Itoa(conf.SMTP.Port),
		smtp.PlainAuth("", sender.Login, sender.Password, conf.SMTP.Server),
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

	///
	var secretKey string

	// Необходимо будет проверять ключи
	secretKey = r.Header.Get("X-Secret-Key")

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
	sender := NewSender(config.SMTP.Sender.Login, config.SMTP.Sender.Password)

	bodyMessage := ""

	if request.ContentType == "html" {
		doc := new(bytes.Buffer)

		t, _ := template.ParseFiles(config.Application.TemplateFile)

		if err = t.Execute(doc, message); err != nil {
			// return err
			log.Panicln(err)
			w.WriteHeader(422)
		}

		bodyMessage = sender.WriteHTMLEmail(Receiver, subject, doc.String())
	} else {
		bodyMessage = sender.WritePlainEmail(Receiver, subject, message)
	}

	sender.SendMail(config, Receiver, subject, bodyMessage)

	data, _ := json.Marshal(response)

	fmt.Fprintf(w, string(data))
}
