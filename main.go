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
	"strconv"
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

var articles []Article

func Articles(begin int, pageSize int) (articles []Article, err error) {
	fmt.Println(begin)
	fmt.Println(pageSize)
	// Open database connection
	db, err := sql.Open("mysql", "root:xing@/blog")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Execute the query
	stmt, err := db.Prepare("SELECT cid,title,slug,text,created,modified FROM articles ORDER BY created DESC LIMIT ?,?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	rows, err := stmt.Query(begin, pageSize)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Fetch rows
	for rows.Next() {
		article := Article{}
		rows.Scan(&article.ID, &article.Title, &article.Slug, &article.Content, &article.Created, &article.Modified)
		articles = append(articles, article)
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello, world\n")
}

func handlerArticles(w http.ResponseWriter, r *http.Request) {
	var page, pageSize int
	vars := mux.Vars(r)
	page, _ = strconv.Atoi(vars["page"])
	pageSize, _ = strconv.Atoi(vars["pageSize"])
	begin := (page - 1) * pageSize

	w.Header().Set("Content-Type", "application/json")

	articles, err := Articles(begin, pageSize)
	if err != nil {
		panic(err.Error())
	}

	json.NewEncoder(w).Encode(articles)
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

func handlerArticleUpdate(w http.ResponseWriter,
	r *http.Request) {
	fmt.Fprintf(w, "Update")
}

func main() {
	r := mux.NewRouter()
	s := r.PathPrefix(PREFIX).Subrouter()
	s.HandleFunc("/", HomeHandler).Methods("GET")
	s.Path("/articles").HandlerFunc(handlerArticles).Methods("GET").Queries("page", "{page:[1-9][0-9]*}", "pageSize", "{pageSize:[1-9][0-9]*}")
	s.HandleFunc("/articles", handlerArticleAdd).Methods("POST")
	s.HandleFunc("/articles/{id}", handlerArticleDetail).Methods("GET")
	s.HandleFunc("/articles/{id}", handlerArticleDelete).Methods("DELETE")
	s.HandleFunc("/articles/{id}", handlerArticleUpdate).Methods("PUT")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8888", r))
}
