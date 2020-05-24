.PHONY: build run

build:
	go build

run:
	go run main.go struct_config.go struct_request.go struct_response.go config.go endpoint_sendmail_v1.go
