package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

var ctx = context.Background()

var rdb *redis.Client

func main() {
	initRedis()

	r := mux.NewRouter()
	r.HandleFunc("/{username}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Category: %v\n", username)

		val, err := rdb.Get(username).Result()
		if err != nil {

			err = rdb.Set(username, username, 0).Err()
			if err != nil {
				panic(err)
			}
		}

		fmt.Println(val)

	}).Methods(http.MethodGet)
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
