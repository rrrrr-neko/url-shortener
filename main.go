package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

var tmpl *template.Template

// var urls = make(map[string]string)
var db *sql.DB

const createTableQuery = `
		CREATE TABLE IF NOT EXISTS urls(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			short_key TEXT NOT NULL UNIQUE,
			long_url TEXT NOT NULL,
			cookie_id TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`

func main() {
	rand.Seed(time.Now().UnixNano())

	var err error
	db, err = sql.Open("sqlite", "db/url_db.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("cannot create table:", err)
	}

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table';")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		rows.Scan(&name)
		log.Println("Table found:", name)
	}

	tmpl = template.Must(template.ParseFiles("templates/index.html"))

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/shorten", handleShorten)
	http.HandleFunc("/wow/", handleRedirect)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	data := struct {
		ShortURL string
		LongURL  string
	}{}
	tmpl.Execute(w, data)
}

func getOrStCookieID(w http.ResponseWriter, r *http.Request) string {
	c, err := r.Cookie("session_id")
	if err == nil && c.Value != "" {
		return c.Value
	}

	id := fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Int())
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    id,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 24 * 365, // год
	})

	return id
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	longURL := r.FormValue("url")
	log.Println("Получен URL:", longURL)

	key := generateKey()
	// urls[key] = longURL
	cookieID := getOrStCookieID(w, r)

	sql := "INSERT INTO urls (short_key, long_url, cookie_id) VALUES (?, ?, ?)"
	_, err := db.Exec(sql, key, longURL, cookieID)
	if err != nil {
		http.Error(w, "db insert failed", http.StatusInternalServerError)
		return
	}

	// scheme := "http"
	// if r.Header.Get("X-Forwarded-Proto") == "https" {
	// 	scheme = "https"
	// }
	// shorturl := scheme + "://" + r.Host + "/wow/" + key

	shorturl := "/wow/" + key
	data := struct {
		ShortURL string
		LongURL  string
	}{

		ShortURL: shorturl,
		LongURL:  longURL,
	}

	tmpl.Execute(w, data)

}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/wow/"):]

	// longURL, ok := urls[key]
	// if !ok {
	// 	http.NotFound(w, r)
	// 	return
	// }

	var longURL string
	err := db.QueryRow("SELECT long_url FROM urls WHERE short_key = ?", key).Scan(&longURL)

	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}

func generateKey() string {
	// Генерируем число в диапазоне [0, 999]
	num := rand.Intn(1000)
	return fmt.Sprintf("%03d", num)
}
