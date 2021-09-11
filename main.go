package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Everyting is OK\n")
}

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "index\n")
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello\n")
}

func startListener() error {
	// create the router
	router := mux.NewRouter()
	router.Path("/").HandlerFunc(index).Methods("GET")
	router.Path("/health").HandlerFunc(health).Methods("GET")
	router.Path("/hello").HandlerFunc(hello).Methods("GET")

	// create server certificates
	certDir := "/etc/go-https-server"
	err := os.MkdirAll(certDir, os.ModePerm)
	if nil != err {
		return err
	}
	certfile := fmt.Sprintf("%v/cert.pem", certDir)
	keyfile := fmt.Sprintf("%v/key.pem", certDir)
	tlsCfg, err := NewServerCertificate(certfile, keyfile)
	if nil != err {
		return err
	}

	// configure the server
	srv := &http.Server{
		Addr: ":8090",
		Handler: handlers.CORS(
			handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}))(router),

		TLSConfig:    tlsCfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),

		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second}

	return srv.ListenAndServeTLS(certfile, keyfile)
}

func main() {
	log.Fatal(startListener())
}
