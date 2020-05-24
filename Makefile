.PHONY: build run clean

build:
	go build -o build/pkg-debian/usr/bin/sender-smtp-api

run:
	go run main.go struct_config.go struct_request.go struct_response.go config.go endpoint_sendmail_v1.go

clean:
	go clean ./...
