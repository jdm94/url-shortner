package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

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
	db *sql.DB
}
// This function create server and create DBconnection and initial table is also created
func CreateServer() *Server {
	server := &Server{
		Router: mux.NewRouter(),
	}
	server.db, _ = DBHandler.CreateDBConnection()
	DBHandler.CreateURLShortenerTable(server.db)
	server.routes()
	return server
}

func (s *Server) routes() {
	s.HandleFunc("/generate-short-url", s.generateShortURL()).Methods("POST")
}

// method that generatesShortURL 
func (s *Server) generateShortURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var url URL
		if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("url decoded is %s", url.Url)
		id, shortUrl, timeStamp, err := DBHandler.SearchLongUrl(s.db, url.Url)
		log.Printf("Id %d ShortUrl %s TimeStamp %d", id, shortUrl, timeStamp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		currentTimeStamp := time.Now().Unix()
		// var shortUrl string
		if(currentTimeStamp - timeStamp > 24*3600 || id == 0) {
			// delete using id
			if(id != 0) {
				DBHandler.DeleteRow(s.db, id)
			}
			var uuid int64 = generator.Generator()
			shortUrl = "https://my-short-link/" + convertor.Convertor(uuid)
			DBHandler.InsertRow(s.db, uuid , url.Url, shortUrl)
		} 
	

		// fmt.Println(convertor.Convertor(a))
		// db, _ := DBHandler.CreateDBConnection()
		
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(shortUrl); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
