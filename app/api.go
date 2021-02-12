package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_"github.com/go-sql-driver/mysql"
	_"github.com/google/uuid"
	//jwt "github.com/dgrijalva/jwt-go"
	"io"
	"log"
	"net/http"
	"go-gatcha/app/config"
	"go-gatcha/app/error"
	"go-gatcha/app/token"
)

type User struct {
	Name string `json:"name"`
}
type Token struct {
	Token string `json:"token"`
}
/**
func getUuid() uuid.UUID {
	u, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
	}
	return u
}
 */
/**
	//tokenがjwtによるbase64urlEncoding(header) + '.' + base64urlEncoding(payload) + '.' + base64urlEncoding(signature)になっている
	//DBにtokenとして保存するならjwt認証の意味は薄い(ユーザーの情報を詰める必要はないし、さらに負荷がかかるので)
	//tokenをuuidにして、app/tokenに移植した
func issueToken(user User) (string, error) {
	var err error
	secret := "secret"
	// {Base64 encoded Header}.{Base64 encoded Payload}.{Signature}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims {
		"uuid": getUuid(),
		"name": user.Name,
		"iss": "__init__",
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatal(err)
	}
	return tokenString, nil
}
 */
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is handler\n")
}
func create(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		body := r.Body
		defer body.Close()

		buf := new(bytes.Buffer)
		io.Copy(buf, body)

		var user User
		err := json.Unmarshal(buf.Bytes(), &user)
		if err != nil {
			log.Fatal(err)
		}

		// tokenとしてuuidを作成
		token := token.IssueToken()
		/**
		token, err := issueToken(user)
		if err != nil {
			log.Fatal(err)
		}
		 */

		// Exec method is for without getting records
		_, err = db.Exec("INSERT INTO users(name, token) VALUES(?, ?)", user.Name, token)
		if err != nil {
			log.Fatal(err)
		}

		// for response
		t := Token{token}
		res, err := json.Marshal(t)
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(res)
		fmt.Fprintf(w, "POST method create() called: %v\n", user)
	default:
		//fmt.Fprint(w, "Method not allowed\n")
		w = error.StatusCode405(w)
	}
}
func get(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "GET":
		// get request header
		header := r.Header

		//get user name by token
		var name string
		err := db.QueryRow(
			"SELECT name FROM users WHERE token = ?",
			header.Get("x-token")).
			Scan(&name)
		if err != nil {
			log.Fatal(err)
		}

		// generate response json
		u := User{name}
		res, err := json.Marshal(u)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		/*
		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			panic(err.Error())
		}

		// get columns
		columns, err := rows.Columns()
		if err != nil {
			panic(err.Error())
		}
		fmt.Fprintf(w, "GET method get() called: %v\n", columns)
		*/
		fmt.Fprintf(w, "GET method get() called\n")
	default:
		//fmt.Fprint(w, "Method not allowed\n")
		w = error.StatusCode405(w)
	}
}
func update(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "PUT":
		// get request
		header := r.Header
		body := r.Body
		defer body.Close()

		// copy for conversion into json
		buf := new(bytes.Buffer)
		io.Copy(buf, body)

		// convert json body(byte) to object
		var user User
		json.Unmarshal(buf.Bytes(), &user)

		// update user name by token
		_, err := db.Exec("UPDATE users SET name = ? WHERE token = ?", user.Name, header.Get("x-token"))
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, "PUT method update() called\n")
	default:
		//fmt.Fprint(w, "Method not allowed\n")
		w = error.StatusCode405(w)
	}
}
func main() {
	/**
	// 多分packageと干渉している
	var err error
	*/
	log.Printf("Server listening on http://localhost:%s", config.Config.Port)
	db, err := sql.Open(config.Config.DriverName, config.Config.DataSourceName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	http.HandleFunc("/", handler)
	/*
	// var db *sql.DBを引数に持たすことでこれらは作用しなくなる
	http.HandleFunc("/user/create", create)
	http.HandleFunc("/user/get", get)
	http.HandleFunc("/user/update", update)
	 */
	http.HandleFunc("/user/create", func(w http.ResponseWriter, r *http.Request) {
		create(w, r, db)
	})
	http.HandleFunc("/user/get", func(w http.ResponseWriter, r *http.Request) {
		get(w, r, db)
	})
	http.HandleFunc("/user/update", func(w http.ResponseWriter, r *http.Request) {
		update(w, r, db)
	})
	log.Print(http.ListenAndServe(":" + config.Config.Port, nil))
}