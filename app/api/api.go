package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go-gatcha/app/api/handler/gatcha"
	"go-gatcha/app/api/handler/user"
	"go-gatcha/app/config"
	"go-gatcha/app/db"
	// jwt "github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is handler\n")
}
func main() {
	db := db.Connect()
	defer db.Close()

	http.HandleFunc("/", handler)
	http.HandleFunc("/user/create", func(w http.ResponseWriter, r *http.Request) {
		user.Create(w, r, db)
	})
	http.HandleFunc("/user/get", func(w http.ResponseWriter, r *http.Request) {
		user.Get(w, r, db)
	})
	http.HandleFunc("/user/update", func(w http.ResponseWriter, r *http.Request) {
		user.Update(w, r, db)
	})
	http.HandleFunc("/gatcha/draw", func(w http.ResponseWriter, r *http.Request) {
		gatcha.Draw(w, r, db)
	})
	log.Print(http.ListenAndServe(":" + config.Config.Port, nil))
}