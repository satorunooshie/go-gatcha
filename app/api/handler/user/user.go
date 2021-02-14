package user

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"gacha/app/error"
	"gacha/app/token"
	"io"
	"log"
	"net/http"
	"time"
)

type M_User struct {
	ID int `db:"id"`
	Name string `db:"name"`
	Token string `db:"token"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
type User struct {
	Name string `json:"name"`
}
type Token struct {
	Token string `json:"token"`
}
func Create(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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

		// transaction start
		trn, trnErr := db.Begin()
		if trnErr != nil {
			log.Fatal(trnErr)
		}

		// Exec method is for without getting records
		_, exeErr := db.Exec("INSERT INTO users(name, token) VALUES(?, ?)", user.Name, token)
		if exeErr != nil {
			_ = trn.Rollback()
			log.Fatal(exeErr)
		}
		trnErr = trn.Commit()
		if trnErr != nil {
			log.Fatal(trnErr)
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
		// fmt.Fprint(w, "Method not allowed\n")
		// replace original status code with default function
		// w = error.StatusCode405(w)
		//error.StatusCode405(w)
		error.ErrorResponse(w, http.StatusMethodNotAllowed)
	}
}
func Get(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
		// fmt.Fprint(w, "Method not allowed\n")
		// replace original status code with default function
		// w = error.StatusCode405(w)
		//error.StatusCode405(w)
		error.ErrorResponse(w, http.StatusMethodNotAllowed)
	}
}
func Update(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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

		// transaction start
		trn, trnErr := db.Begin()
		if trnErr != nil {
			log.Fatal(trnErr)
		}
		// update user name by token
		_, exeErr := db.Exec("UPDATE users SET name = ? WHERE token = ?", user.Name, header.Get("x-token"))
		if exeErr != nil {
			_ = trn.Rollback()
			log.Fatal(exeErr)
		}
		trnErr = trn.Commit()
		if trnErr != nil {
			log.Fatal(trnErr)
		}

		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, "PUT method update() called\n")
	default:
		// fmt.Fprint(w, "Method not allowed\n")
		// replace original status code with default function
		// w = error.StatusCode405(w)
		//error.StatusCode405(w)
		error.ErrorResponse(w, http.StatusMethodNotAllowed)
	}
}
