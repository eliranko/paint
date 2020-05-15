package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"

	"github.com/gorilla/mux"
)

var httpTimeout = 5 * time.Second

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.RequestURI, r.Method)
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		next.ServeHTTP(w, r)
	})
}

func handleGetCanvases(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), httpTimeout)
	res, err := getCanvases(ctx)
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
	ctx, _ := context.WithTimeout(context.Background(), httpTimeout)
	id := mux.Vars(r)["id"]
	canvas, err := getCanvas(ctx, id)
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
	ctx, _ := context.WithTimeout(context.Background(), httpTimeout)
	if err = saveCanvas(ctx, canvas); err != nil {
		log.Println("could not save canvas ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

	canvas.Data = ""
	pushCanvases(canvas)
}

func startServer() {
	r := mux.NewRouter()
	r.Use(corsMiddleware)
	r.HandleFunc("/api/canvas", handleGetCanvases).Methods("GET")
	r.HandleFunc("/api/canvas/{id}", handleGetCanvas).Methods("GET")
	r.HandleFunc("/api/canvas", handleSaveCanvasRequest).Methods("POST")
	r.PathPrefix("/socket.io").Handler(getPushNotificationServer())
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./dist/")))
	r.Use(loggingMiddleware)

	log.Println("listening on :", viper.GetString("port"))
	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + viper.GetString("port"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
