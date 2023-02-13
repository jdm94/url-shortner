package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/url-shortner/DBHandler"
	"github.com/url-shortner/convertor"
	"github.com/url-shortner/generator"
)

type URL struct {
	Url string `json:"url"`
}

type Server struct {
	*mux.Router
}

func CreateServer() *Server {
	server := &Server{
		Router: mux.NewRouter(),
	}
	server.routes()
	return server
}

func (s *Server) routes() {
	s.HandleFunc("/generate-short-url", s.generateShortURL()).Methods("POST")
}

func (s *Server) generateShortURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var url URL
		if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("url decoded is %s", url.Url)
		var uuid int64 = generator.Generator()
		shortUrl := "https://my-short-link/" + convertor.Convertor(uuid)
		// fmt.Println(convertor.Convertor(a))
		db, _ := DBHandler.CreateDBConnection()
		DBHandler.CreateURLShortenerTable(db)
		DBHandler.InsertRow(db, uuid , url.Url, shortUrl)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(shortUrl); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
