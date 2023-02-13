package main

import (
	"net/http"

	"github.com/url-shortner/server"
)

func main() {
	server := server.CreateServer()
	http.ListenAndServe(":8081", server)
}
