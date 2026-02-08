package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var tmpl *template.Template
var urls = make(map[string]string)

func main() {
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

func handleShorten(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	longURL := r.FormValue("url")
	log.Println("Получен URL:", longURL)

	key := generateKey()
	urls[key] = longURL

	scheme := "http"
	if r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	shorturl := scheme + "://" + r.Host + "/wow/" + key

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

	longURL, ok := urls[key]
	if !ok {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}

func generateKey() string {

	rand.Seed(time.Now().UnixNano())

	// Генерируем число в диапазоне [0, 999]
	num := rand.Intn(1000)

	return fmt.Sprintf("%03d", num)
}
