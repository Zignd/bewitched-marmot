package main

import (
	"log"
	"net/http"

	"encoding/json"

	"github.com/Zignd/bewitched-marmot/sites/mangahere"
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
	http.HandleFunc("/mangahere/getmanga", func(res http.ResponseWriter, req *http.Request) {
		mangaURL := req.URL.Query().Get("mangaURL")
		manga, err := mangahere.GetManga(mangaURL)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		resBody, err := json.Marshal(manga)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(resBody)
	})
	http.HandleFunc("/mangahere/getchapter", func(res http.ResponseWriter, req *http.Request) {
		chapterURL := req.URL.Query().Get("mangaURL")
		chapter, err := mangahere.GetChapter(chapterURL)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		resBody, err := json.Marshal(chapter)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(resBody)
	})
	log.Fatal(http.ListenAndServe("localhost:8081", nil))
}
