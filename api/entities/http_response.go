package entities

import (
	"fmt"
	"net/http"
)

type HttpResponse struct {
	StatusCode int         `json:"statusCode,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	URL        string      `json:"url,omitempty"`
}

func newHttpResponse(statusCode int, data interface{}, url string) *HttpResponse {
	return &HttpResponse{StatusCode: statusCode, Data: data}
}

func formatError(message interface{}) map[string]string {
	mapError := make(map[string]string)
	mapError["error"] = fmt.Sprintf("%v", message)
	return mapError
}

func Ok(data interface{}) *HttpResponse {
	return newHttpResponse(http.StatusOK, data, "")
}

func OkWithUrl(data interface{}, url string) *HttpResponse {
	return newHttpResponse(http.StatusOK, data, url)
}

func Created(data interface{}) *HttpResponse {
	return newHttpResponse(http.StatusCreated, data, "")
}

func NoContent() *HttpResponse {
	return newHttpResponse(http.StatusNoContent, nil, "")
}

func BadRequest(data interface{}) *HttpResponse {
	return newHttpResponse(http.StatusBadRequest, formatError(data), "")
}

func Unauthorized(data interface{}) *HttpResponse {
	return newHttpResponse(http.StatusUnauthorized, formatError("Token inválido ou expirado"), "")
}

func NotFound(data interface{}) *HttpResponse {
	return newHttpResponse(http.StatusNotFound, formatError("Recurso não encontrado"), "")
}

func ServerError() *HttpResponse {
	return newHttpResponse(http.StatusInternalServerError, formatError("Ocorreu um erro inesperado"), "")
}

func InternalServerError(data interface{}) *HttpResponse {
	return newHttpResponse(http.StatusInternalServerError, formatError(data), "")
}
