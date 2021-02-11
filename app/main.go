package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID int
	Name string
	Token string
}
type Chara struct {
	ID int
	Name string
	RarityId int
}
type Rarity struct {
	ID int
	Name string
	Probability float64
}
// htmlに埋め込む構造体
type Embed struct {
	Title string
	Message string
	Users map[int]User
	Chara map[int]Chara
	Rarities map[int]Rarity
	Time time.Time
}
const (
	DriverName = "mysql"
	DataSourceName = "root:golang@tcp(localhost:3306)/golang_db"
)
var usr = make(map[int]User)
var rarity = make(map[int]Rarity)
var char = make(map[int]Chara)
var templates = make(map[string]*template.Template)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmp := Embed{"Hello Golang", "こんにちは", usr, char, rarity, time.Now()}
	if err := templates["index"].Execute(w, tmp); err != nil {
		log.Printf("failed to execute template: %v", err)
	}
}

func loadTemplate(name string) *template.Template {
	t, err := template.ParseFiles(
		"app/" + name + ".html",
		"app/template/header.html",
		"app/template/footer.html",
	)
	if err != nil {
		log.Fatalf("template error: %v", err)
	}
	return t
}

func main() {
	// database connection
	db, dbErr := sql.Open(DriverName, DataSourceName)
	if dbErr != nil {
		log.Print("error connencting to database :", dbErr)
	}
	defer db.Close()
	// charaテーブルの全てのレコードを取得するクエリを実行
	rows, queryErr := db.Query("SElECT * FROM chara")
	if queryErr != nil {
		log.Print("query error :", queryErr)
	}
	defer rows.Close()
	// ループを回してrowsからScanでデータを取得
	for rows.Next() {
		var c Chara
		if err := rows.Scan(&c.ID, &c.Name, &c.RarityId); err != nil {
			log.Print(err)
		}
		char[c.ID] = Chara{
			ID: c.ID,
			Name: c.Name,
			RarityId: c.RarityId,
		}
	}
	// Web server
	port := "7777"
	templates["index"] = loadTemplate("index")
	http.HandleFunc("/", handleIndex)
	log.Printf("Server listening on http://localhost:%s", port)
	log.Print(http.ListenAndServe(":" + port, nil))
}