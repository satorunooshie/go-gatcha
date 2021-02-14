package gacha

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
	e "gacha/app/error"
)

type Rarity struct {
	Probability float64
	Characters []Chara
}

type Chara struct {
	ID string `json:"characterID"`
	Name string `json:"name"`
}

type Response struct {
	Results []Chara `json:"results"`
}

type Times struct {
	Times int `json:"times"`
}

func Draw(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		body := r.Body
		defer body.Close()

		// byte[]に変換するためのコピー
		buf := new(bytes.Buffer)
		io.Copy(buf, body)

		var times Times
		err := json.Unmarshal(buf.Bytes(), &times)
		if err != nil {
			log.Fatal(err)
		}

		rand.Seed(time.Now().UnixNano())
		var response Response

		for i := 0; i < times.Times; i++ {
			character, err := oneDraw(rand.Float64(), db)
			if err != nil {
				fmt.Fprint(w, err)
				return
			}
			response.Results = append(response.Results, *character)
		}

		res, err := json.Marshal(response)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	default:
		e.ErrorResponse(w, http.StatusMethodNotAllowed)
	}
}

func oneDraw(n float64, db *sql.DB) (*Chara, error) {
	config, err := getConfig(db)
	if err != nil {
		return nil, err
	}
	boundary := 0.0
	for _, rarity := range config {
		for _, chara := range rarity.Characters {
			boundary += rarity.Probability / float64(len(rarity.Characters))
			if n <= boundary {
				return &chara, nil
			}
		}
	}
	return nil, err
}

func getConfig(db *sql.DB) ([]Rarity, error) {
	var config []Rarity
	rows, err := db.Query("SELECT rarities.probability, GROUP_CONCAT(chara.id), GROUP_CONCAT(chara.name) FROM rarities JOIN chara ON rarities.id = chara.rarity_id GROUP BY rarities.id ORDER BY rarities.probability")
	if err != nil {
		return nil, err
	}
	/**
	+-------------+-------------------------+-----------------------------------------------------------------+
	| probability | GROUP_CONCAT(chara.id)  | GROUP_CONCAT(chara.name)                                        |
	+-------------+-------------------------+-----------------------------------------------------------------+
	|        0.01 | 15,16                   | super1,super2                                                   |
	|         0.1 | 17,18,19,20             | frequent1,frequent2,frequent3,frequent4                         |
	|        0.89 | 21,22,23,24,25,26,27,28 | common1,common2,common3,common4,common5,common6,common7,common8 |
	+-------------+-------------------------+-----------------------------------------------------------------+
	 */
	for rows.Next() {
		var r Rarity
		// 後でSplitするためにstring
		var charaIdsLinkedComma string
		var charaNamesLinkedComma string
		// ProbabilityはそのままRarity[]に詰める
		if err := rows.Scan(&r.Probability, &charaIdsLinkedComma, &charaNamesLinkedComma); err != nil {
			return nil, err
		}
		// データを整形してスライスに詰める
		charaIds := strings.Split(charaIdsLinkedComma, ",")
		charaNames := strings.Split(charaNamesLinkedComma, ",")
		// &r.Charactersに詰めるためのスライス
		var characters []Chara
		// IDを基軸に展開する
		for i, id := range charaIds {
			characters = append(characters, Chara{ID: id, Name: charaNames[i]})
		}
		// RarityごとにソートされたのでRarityに詰める
		r.Characters = characters
		config = append(config, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()
	return config, nil
}