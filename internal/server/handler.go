package server

import (
	"io"
	"log"

	"google.com/Dorfieeee/bootdev-http-protocol/internal/request"
	"google.com/Dorfieeee/bootdev-http-protocol/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)

func (handlerError *HandlerError) Write(w io.Writer) {
	err := response.WriteStatusLine(w, response.StatusCode(handlerError.StatusCode))
	if err != nil {
		log.Printf("Error WriteStatusLine(): %v\n", err)
	}
	headers := response.GetDefaultHeaders(len(handlerError.Message))
	err = response.WriteHeaders(w, headers)
	if err != nil {
		log.Printf("Error WriteStatusLine(): %v\n", err)
	}
	w.Write([]byte(handlerError.Message))
}
