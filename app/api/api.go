package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gacha/app/api/handler/gacha"
	"gacha/app/api/handler/user"
	"gacha/app/config"
	"gacha/app/db"
	"gacha/app/api/auth"

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
	http.HandleFunc("/gacha/draw", func(w http.ResponseWriter, r *http.Request) {
		gacha.Draw(w, r, db)
	})
	var private = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := &user.User{Name: "namae"}
		json.NewEncoder(w).Encode(user)
	})
	http.Handle("/private", auth.JwtMiddleware.Handler(private))
	http.Handle("/auth", auth.GetTokenHandler)
	log.Print(http.ListenAndServe(":" + config.Config.Port, nil))
}