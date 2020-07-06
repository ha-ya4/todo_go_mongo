package main

import (
	"context"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	dbName  = "gotodo"
	colName = "todo"
)

var db database

type database struct {
	c *mongo.Client
}

func (d database) getCol(name string) *mongo.Collection {
	return d.c.Database(dbName).Collection(name)
}

func connect() (database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return database{}, err
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return database{}, err
	}
	return database{c: client}, err
}

type Todo struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
}

func getTodoAll(col *mongo.Collection) ([]Todo, error) {
	var todo []Todo
	cur, err := col.Find(context.Background(), bson.D{}, options.Find())
	if err != nil {
		return todo, err
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		t := Todo{}
		if err = cur.Decode(&t); err != nil {
			return todo, err
		}
		todo = append(todo, t)
	}

	if err = cur.Err(); err != nil {
		return todo, err
	}
	return todo, err
}

func (t Todo) insert(col *mongo.Collection) (*mongo.InsertOneResult, error) {
	d := bson.D{
		bson.E{Key: "id", Value: t.ID},
		bson.E{Key: "title", Value: t.Title},
		bson.E{Key: "comment", Value: t.Comment},
	}
	return col.InsertOne(context.Background(), d)
}

func (t Todo) update() error {
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
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	todo, err := getTodoAll(db.getCol(colName))
	if err != nil {
		return
	}
	st := struct {
		Todo []Todo
	}{
		Todo: todo,
	}

	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		return
	}

	if e := t.Execute(w, st); err != nil {
		log.Println(e)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	todo := &Todo{}
	if err = json.Unmarshal(b, todo); err != nil {
		return
	}
	if id, err := uuid.NewRandom(); err != nil {
		return
	} else {
		todo.ID = id.String()
	}

	if _, err = todo.insert(db.getCol(colName)); err != nil {
		return
	}

	todoSlice, err := getTodoAll(db.getCol(colName))
	if err != nil {
		return
	}
	st := struct {
		Todo []Todo
	}{
		Todo: todoSlice,
	}

	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		return
	}

	if e := t.Execute(w, st); err != nil {
		log.Println(e)
	}
}

func handlePut(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		return
	}

	if e := t.Execute(w, nil); err != nil {
		log.Println(e)
	}
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		return
	}

	if e := t.Execute(w, nil); err != nil {
		log.Println(e)
	}
}

func main() {
	log.SetFlags(log.Lshortfile)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources/"))))
	http.HandleFunc("/", handle)

	var err error
	db, err = connect()
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err = db.c.Disconnect(ctx); err != nil {
			log.Println(err)
		}
	}()

	log.Println("--start--")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println(err)
	}
}
