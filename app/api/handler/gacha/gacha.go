package gacha

import (
	"database/sql"
	"fmt"
	"net/http"
)

func Draw(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	n := 0.01
	var character string
	if n < 0.01 {
		character = "SSR"
	} else if n < 0.11 {
		character = "SR"
	} else {
		character = "R"
	}
	fmt.Fprint(w, character)
}
