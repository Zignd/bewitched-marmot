package main

import (
	"log"
	"net/http"

	"encoding/json"

	"github.com/zignd/bewitched-marmot/sites/mangahere"
)

func main() {
	http.HandleFunc("/mangahere/search", func(res http.ResponseWriter, req *http.Request) {
		query := req.URL.Query().Get("query")
		mangas, err := mangahere.Search(query)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		resBody, err := json.Marshal(mangas)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(resBody)
	})
	log.Fatal(http.ListenAndServe("localhost:8081", nil))
}
