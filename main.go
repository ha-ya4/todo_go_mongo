package main

import (
	"html/template"
	"log"
	"net/http"
)

type Todo struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
}

func (t Todo) get() error {
	return nil
}

func (t Todo) post() error {
	return nil
}

func (t Todo) put() error {
	return nil
}

func (t Todo) delete() error {
	return nil
}

func handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGet(w, r)
	case http.MethodPost:
		handlePost(w, r)
	case http.MethodPut:
		handlePut(w, r)
	case http.MethodDelete:
		handleDelete(w, r)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		log.Println(err)
	}

	todo := Todo{
		ID:      0,
		Title:   "hello",
		Comment: "world!",
	}

	if err = t.Execute(w, todo); err != nil {
		log.Println(err)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		log.Println(err)
	}

	if err = t.Execute(w, nil); err != nil {
		log.Println(err)
	}
}

func handlePut(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		log.Println(err)
	}

	if err = t.Execute(w, nil); err != nil {
		log.Println(err)
	}
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		log.Println(err)
	}

	if err = t.Execute(w, nil); err != nil {
		log.Println(err)
	}
}

func main() {
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
	http.HandleFunc("/", handle)
	log.Println("--start--")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println(err)
	}
}
