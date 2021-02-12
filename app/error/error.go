package error

import (
	"encoding/json"
	"log"
	"net/http"
)

type Message struct {
	Message string `json:"message"`
}

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
}