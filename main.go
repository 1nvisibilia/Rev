package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	envLoadErr := godotenv.Load()
	if envLoadErr != nil {
		log.Fatal("Cannot load .env file")
	}
	backend_server_port := os.Getenv("BACKEND_SERVER_PORT")
	proxy_starting_port := os.Getenv("PROXY_PORT")

	target, err := url.Parse("http://" + backend_server_port)

	if err != nil {
		log.Fatal("Cannot parse URL:", err)
	}

	proxyServer := httputil.NewSingleHostReverseProxy(target)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.Host = target.Host
		proxyServer.ServeHTTP(w, r)
	})

	log.Println("Starting server")
	log.Println("		Proxy Server		==>>		Actual Server")
	log.Println("		" + proxy_starting_port + "		==>>		" + backend_server_port)
	log.Fatal(http.ListenAndServe(proxy_starting_port, nil))
}
