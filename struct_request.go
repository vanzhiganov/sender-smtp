package main

// RequestStruct структура запроса
type RequestStruct struct {
	ID          string `json:"id"`
	Subject     string `json:"subject"`
	Message     string `json:"message"`
	To          string `json:"to"`
	ContentType string `json:"content-type"` // both|plain|html
}
