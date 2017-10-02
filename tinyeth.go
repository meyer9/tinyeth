package main

import (
	"database/sql"
	"html/template"
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
	static   http.Handler
}

func (t *TinyEth) resolveAlias(urlCode string) string {
	rows, _ := t.database.Query("SELECT address FROM RegistrationAliases WHERE url=?", urlCode)
	if rows.Next() {
		var address string
		rows.Scan(&address)
		return address
	}
	return ""
}

func (t *TinyEth) resolveRandom(urlCode string) string {
	rows, _ := t.database.Query("SELECT address FROM Registrations WHERE url=?", convertMnemonicToID(urlCode))
	if rows.Next() {
		var address string
		rows.Scan(&address)
		return address
	}
	return ""
}

func (t *TinyEth) getAddress(urlCode string) string {
	if unicode.IsUpper(rune(urlCode[0])) {
		return t.resolveAlias(urlCode)
	}
	address := t.resolveRandom(urlCode)
	if address == "" {
		address = t.resolveAlias(urlCode)
	}
	return address
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

	}
}

func (t TinyEth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.FileServer(http.Dir("./static")).ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/static") {
		t.static.ServeHTTP(w, r)
	} else if r.URL.Path == "/api/register" && r.Method == "POST" {
		t.registerURL(w, r)
	} else {
		address := t.getAddress(r.URL.Path[1:])
		template, _ := template.ParseFiles("template/url.html")

		if err := template.Execute(w, address); err != nil {
			log.Fatal(err)
		}
	}
}

func initTinyEth(database *sql.DB) *TinyEth {
	return &TinyEth{database: database, static: http.StripPrefix("/static", http.FileServer(http.Dir("static")))}
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
	tinyeth := initTinyEth(db)
	http.Handle("/", tinyeth)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
