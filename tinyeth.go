package main

import (
	"database/sql"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"unicode"

	_ "github.com/go-sql-driver/mysql"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

// TinyEth - represents TinyEth logic
type TinyEth struct {
	database *sql.DB
}

func (t *TinyEth) serveURL(w http.ResponseWriter, r *http.Request) {
	urlCode := r.URL.Path[1:]
	if unicode.IsUpper(rune(urlCode[0])) {
		rows, _ := t.database.Query("SELECT address FROM RegistrationAliases WHERE url=?", urlCode)
		if rows.Next() {
			var address string
			rows.Scan(&address)
			io.WriteString(w, address)
		} else {
			w.WriteHeader(404)
			io.WriteString(w, "Could not find url.")
		}
	} else {
		rows, _ := t.database.Query("SELECT address FROM Registrations WHERE url=?", convertMnemonicToID(urlCode))
		if rows.Next() {
			var address string
			rows.Scan(&address)
			io.WriteString(w, address)
		} else {
			w.WriteHeader(404)
			io.WriteString(w, "Could not find url.")
		}
	}
}

func convertIDToMnemonic(i int) string {
	var digits []int
	for i > 0 {
		remainder := i % 26
		digits = append([]int{remainder}, digits...)
		i = i / 26
	}
	var text string
	for n := range digits {
		b := letterBytes[digits[n]]
		text = text + string(b)
	}
	return text
}

func convertMnemonicToID(s string) int {
	i := 0
	exponent := 0
	for p := len(s) - 1; p >= 0; p-- {
		q := strings.IndexByte(letterBytes, s[p])
		i = i + int(math.Pow(float64(26), float64(exponent)))*q
		exponent = exponent + 1
	}
	return i
}

func (t *TinyEth) registerURL(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
		return
	}

	address := r.Form.Get("address")
	alias := r.Form.Get("alias")
	if address == "" {
		io.WriteString(w, "Missing required argument: address")
		return
	}
	if alias == "" {
		// random registrations will start with a lowercase letter
		rows, err := t.database.Query("SELECT url FROM Registrations WHERE address=?", address)
		var row int
		if rows.Next() {
			rows.Scan(&row)
			io.WriteString(w, convertIDToMnemonic(row))
			return
		}

		ex, err := t.database.Exec("INSERT INTO Registrations (address) VALUES (?)", address)
		if err != nil {
			log.Fatal(err)
		}
		id, err := ex.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(w, convertIDToMnemonic(int(id)))
	} else {
		// alias registrations will start with an uppercase letter
		rows, err := t.database.Query("SELECT url FROM RegistrationAliases WHERE url=?", strings.ToLower(alias))
		if rows.Next() {
			io.WriteString(w, "This address has already been used.")
			return
		}

		_, err = t.database.Exec("INSERT INTO RegistrationAliases (address, url) VALUES (?, ?)", address, strings.ToLower(alias))
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(w, string(unicode.ToUpper(rune(alias[0])))+alias[1:])
	}
}

func (t *TinyEth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {

	} else if strings.HasPrefix(r.URL.Path, "/static") {
	} else if r.URL.Path == "/register" && r.Method == "POST" {
		t.registerURL(w, r)
	} else {
		t.serveURL(w, r)
	}
}

func main() {
	databaseURL := os.Getenv("DATABASE")
	if databaseURL == "" {
		databaseURL = "root@tcp(127.0.0.1:3306)/tinyeth"
	}
	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	tinyeth := &TinyEth{database: db}
	http.Handle("/", tinyeth)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
