package main

import (
	"net/http"
	"fmt"
	"io"
	"github.com/gorilla/mux"
	"log"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const VERSION string = "v1"
const PREFIX string = "/api/" + VERSION

type Article struct {
	ID       int    `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	Slug     string `json:"slug,omitempty"`
	Content  string `json:"content,omitempty"`
	Created  int    `json:"created,omitempty"`
	Modified int    `json:"modified,omitempty"`
}

var article []Article

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello, world\n")
}

func handlerArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Open database connection
	db, err := sql.Open("mysql", "root:xing@/blog")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Execute the query
	rows, err := db.Query("SELECT id,title,slug,content,created,modified FROM articles")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Fetch rows
	for rows.Next() {
		var id, created, modified int
		var title, slug, content string
		rows.Scan(&id, &title, &slug, &content, &created, &modified)
		article = append(article, Article{id, title, slug, content, created, modified})
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	json.NewEncoder(w).Encode(article)
}

func handlerArticleDetail(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Detail")
}

func handlerArticleAdd(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Add")
}

func handlerArticleDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Delete")
}

func handlerArticleUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Update")
}

func main() {
	r := mux.NewRouter()
	s := r.PathPrefix(PREFIX).Subrouter()
	s.HandleFunc("/", HomeHandler).Methods("GET")
	s.HandleFunc("/articles", handlerArticles).Methods("GET")
	s.HandleFunc("/articles", handlerArticleAdd).Methods("POST")
	s.HandleFunc("/articles/{id}", handlerArticleDetail).Methods("GET")
	s.HandleFunc("/articles/{id}", handlerArticleDelete).Methods("DELETE")
	s.HandleFunc("/articles/{id}", handlerArticleUpdate).Methods("PUT")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8888", r))
}
