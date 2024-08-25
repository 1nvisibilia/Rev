package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	RPB "rev/proxy_observer"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	envLoadErr := godotenv.Load()
	if envLoadErr != nil {
		log.Fatal("Cannot load .env file")
	}
	backendServerPort := os.Getenv("BACKEND_SERVER_PORT")
	proxyStartingPort := os.Getenv("PROXY_PORT")

	target, err := url.Parse("http://" + backendServerPort)

	if err != nil {
		log.Fatal("Cannot parse URL:", err)
	}

	reverseProxyBalancer := RPB.NewReverseProxyBalancer()

	proxyServer := httputil.NewSingleHostReverseProxy(target)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.Host = target.Host

		allowRequest := reverseProxyBalancer.ProcessRequest(r)

		if allowRequest {
			proxyServer.ServeHTTP(w, r)
		} else {
			log.Println("Request has been blocked due to rate limit. IP: " + strings.Split(r.RemoteAddr, ":")[0])
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprintf(w, "Request has been blocked due to rate limit")
		}
	})

	go reverseProxyBalancer.MonitorCoolDownList()

	log.Println("Starting server")
	log.Println("		Proxy Server		==>>		Actual Server")
	log.Println("		" + proxyStartingPort + "		==>>		" + backendServerPort)
	log.Fatal(http.ListenAndServe(proxyStartingPort, nil))
}
