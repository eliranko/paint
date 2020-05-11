package main

import (
	"encoding/json"
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

func handleGetCanvases(w http.ResponseWriter, r *http.Request) {
	res, err := getCanvases()
	if err != nil {
		log.Println("could not get canvases ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.NewEncoder(w).Encode(res); err != nil { // TODO: gzip
		log.Println("could not serialize canvases ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func handleGetCanvas(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	canvas, err := getCanvas(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	if err := json.NewEncoder(w).Encode(canvas); err != nil { // TODO: gzip
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func handleSaveCanvasRequest(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("could not read request body ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	canvas := &Canvas{}
	if err := json.Unmarshal(buf, canvas); err != nil {
		log.Println("could not unmarshal request body ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	canvas.UUID = uuid.New().String()
	if err = saveCanvas(canvas); err != nil {
		log.Println("could not save canvas ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func startServer() {
	r := mux.NewRouter()
	r.HandleFunc("/api/canvas", handleGetCanvases).Methods("GET")
	r.HandleFunc("/api/canvas/{id}", handleGetCanvas).Methods("GET")
	r.HandleFunc("/api/canvas", handleSaveCanvasRequest).Methods("POST")
	r.Use(loggingMiddleware)

	log.Println("listening on :", viper.GetString("port"))
	log.Fatal(http.ListenAndServe("localhost:"+viper.GetString("port"), r))
}
