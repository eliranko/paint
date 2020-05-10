package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.RequestURI, r.Method)
		next.ServeHTTP(w, r)
	})
}

func handleSaveRequest(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("could not read request body ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = saveCanvas(&Canvas{
		UUID: uuid.New().String(),
		Data: string(buf),
	}); err != nil {
		log.Println("could not save canvas ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func startServer() {
	r := mux.NewRouter()
	r.HandleFunc("/api/canvas", handleSaveRequest).Methods("POST")
	r.Use(loggingMiddleware)

	log.Println("listening on :", viper.GetString("port"))
	log.Fatal(http.ListenAndServe("localhost:"+viper.GetString("port"), r))
}
