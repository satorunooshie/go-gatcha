package token

import (
	"github.com/google/uuid"
	"log"
)

func IssueToken() string {
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
	}
	return uuid.String()
}
