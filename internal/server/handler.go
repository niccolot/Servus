package server

import (
	"io"
	"fmt"

	"Servus/internal/request"
	"Servus/internal/response"
)

type HandlerError struct {
	Code response.StatusCode
	Message string
}

func (h *HandlerError) WriteHandlerError(w io.Writer) {
	errString := "error: " + fmt.Sprint(h.Code) + " " + h.Message
	w.Write([]byte(errString))
}

type Handler func(w io.Writer, req *request.Request) *HandlerError 

