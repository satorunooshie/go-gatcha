package error

import (
	"encoding/json"
	"log"
	"net/http"
)

type Message struct {
	Message string `json:"message"`
}

//func StatusCode405(w http.ResponseWriter) {
	/*
func StatusCode405(w http.ResponseWriter) http.ResponseWriter {
	// generate json message for response
	m := Message{"method not allowed"}
	res, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write(res)
	return w
	 */
	/*
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
}
	 */

func ErrorResponse(w http.ResponseWriter, statusCode int) {
	m := Message{http.StatusText(statusCode)}
	res, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(statusCode)
	w.Write(res)
}