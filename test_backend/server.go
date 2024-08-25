package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received :: " + r.Method + " /")
	fmt.Fprintf(w, "welcome to home directory")
}

func test1Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received :: " + r.Method + " /test1")
	fmt.Fprintf(w, "welcome to test1")
}

func test2Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received :: " + r.Method + " /test2")
	fmt.Fprintf(w, "welcome to test2")
}

func main() {
	envLoadErr := godotenv.Load()

	backend_server_port := os.Getenv("BACKEND_SERVER_PORT")

	if envLoadErr != nil {
		log.Fatal("Cannot load env")
	}

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/test1", test1Handler)
	http.HandleFunc("/test2", test2Handler)

	log.Println("Starting sample REST server on " + backend_server_port)
	log.Fatal(http.ListenAndServe(backend_server_port, nil))
}
