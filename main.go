package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Zignd/bewitched-marmot/sites/mangahere"
	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/mangahere/search", func(res http.ResponseWriter, req *http.Request) {
		fmt.Println(time.Now(), "/mangahere/search ...")
		query := req.URL.Query().Get("query")
		mangas, err := mangahere.Search(query)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			fmt.Println(time.Now(), "/mangahere/search", http.StatusInternalServerError)
			return
		}

		resBody, err := json.Marshal(mangas)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			fmt.Println(time.Now(), "/mangahere/search", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(resBody)
		fmt.Println(time.Now(), "/mangahere/search", http.StatusOK)
	})
	mux.HandleFunc("/mangahere/getmanga", func(res http.ResponseWriter, req *http.Request) {
		fmt.Println(time.Now(), "/mangahere/getmanga ...")
		mangaURL := req.URL.Query().Get("mangaURL")
		manga, err := mangahere.GetManga(mangaURL)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			fmt.Println(time.Now(), "/mangahere/getmanga", http.StatusInternalServerError)
			return
		}

		resBody, err := json.Marshal(manga)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			fmt.Println(time.Now(), "/mangahere/getmanga", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(resBody)
		fmt.Println(time.Now(), "/mangahere/getmanga", http.StatusOK)
	})
	mux.HandleFunc("/mangahere/getchapter", func(res http.ResponseWriter, req *http.Request) {
		fmt.Println(time.Now(), "/mangahere/getchapter ...")
		chapterURL := req.URL.Query().Get("chapterURL")
		chapter, err := mangahere.GetChapter(chapterURL)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			fmt.Println(time.Now(), "/mangahere/getchapter", http.StatusInternalServerError)
			return
		}

		resBody, err := json.Marshal(chapter)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			fmt.Println(time.Now(), "/mangahere/getchapter", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(resBody)
		fmt.Println(time.Now(), "/mangahere/getchapter", http.StatusOK)
	})
	handler := cors.Default().Handler(mux)
	fmt.Println("Application is running at: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
