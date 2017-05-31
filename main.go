package main

import (
	"database/sql"
	"log"

	"github.com/RangelReale/osin"
	"github.com/zignd/bewitched-marmot/sites/mangahere"
	"github.com/zignd/bewitched-marmot/users"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/ory/osin-storage/storage/postgres"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/cors"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
)

type jsonError struct {
	Error string `json:"error"`
}

func main() {
	// Setting up the database to store OAuth 2 related data
	db, err := sql.Open("postgres", "postgres://marmot:1234@localhost:5432/marmot")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	storage := postgres.New(db)
	defer storage.Clone()
	if err = storage.CreateSchemas(); err != nil {
		log.Fatalf("failed to create database schemas: %v", err)
		return
	}

	// Settings up the OAuth 2.0 server
	sconfig := osin.NewServerConfig()
	sconfig.AllowedAccessTypes = osin.AllowedAccessType{
		osin.PASSWORD,
		osin.REFRESH_TOKEN,
	}
	server := osin.NewServer(sconfig, storage)

	// Setting up the endpoints
	app := iris.New()
	app.Adapt(iris.DevLogger(), httprouter.New(), cors.New(cors.Options{AllowedOrigins: []string{"*"}}))

	app.Post("/oauth/token", func(ctx *iris.Context) {
		ares := server.NewResponse()
		defer ares.Close()

		if areq := server.HandleAccessRequest(ares, ctx.Request); areq != nil {
			switch areq.Type {
			case osin.PASSWORD:
				user, err := users.GetUser(areq.Username, areq.Password)
				if err != nil {
					// TODO: Do something with the error value
					ares.SetError(osin.E_SERVER_ERROR, "")
					break
				}
				if user == nil {
					break
				}
				if !user.ConfirmedEmailAddress {
					ares.SetError("email_confirmation_pending", "The email confirmation is still pending. Check your inbox.")
					break
				}

				areq.Authorized = true
				server.FinishAccessRequest(ares, ctx.Request, areq)
			case osin.REFRESH_TOKEN:
				areq.Authorized = true
				server.FinishAccessRequest(ares, ctx.Request, areq)
			}
		}
		osin.OutputJSON(ares, ctx.ResponseWriter, ctx.Request)
	})

	app.Get("/oauth/info", func(ctx *iris.Context) {
		resp := server.NewResponse()
		defer resp.Close()

		if ir := server.HandleInfoRequest(resp, ctx.Request); ir != nil {
			server.FinishInfoRequest(resp, ctx.Request, ir)
		}
		osin.OutputJSON(resp, ctx.ResponseWriter, ctx.Request)
	})

	app.Get("/mangahere/search/:query", func(ctx *iris.Context) {
		query := ctx.Param("query")
		mangas, err := mangahere.Search(query)
		if err != nil {
			ctx.JSON(iris.StatusInternalServerError, &jsonError{
				Error: err.Error(),
			})
			return
		}

		err = ctx.JSON(iris.StatusOK, mangas)
		if err != nil {
			ctx.JSON(iris.StatusInternalServerError, &jsonError{
				Error: err.Error(),
			})
			return
		}
	})

	app.Get("/mangahere/getmanga/:mangaURL", func(ctx *iris.Context) {
		mangaURL := ctx.Param("mangaURL")
		manga, err := mangahere.GetManga(mangaURL)
		if err != nil {
			ctx.JSON(iris.StatusInternalServerError, &jsonError{
				Error: err.Error(),
			})
			return
		}

		err = ctx.JSON(iris.StatusOK, manga)
		if err != nil {
			ctx.JSON(iris.StatusInternalServerError, &jsonError{
				Error: err.Error(),
			})
			return
		}
	})

	app.Get("/mangahere/getchapter/:chapterURL", func(ctx *iris.Context) {
		chapterURL := ctx.Param("chapterURL")
		chapter, err := mangahere.GetChapter(chapterURL)
		if err != nil {
			ctx.JSON(iris.StatusInternalServerError, &jsonError{
				Error: err.Error(),
			})
			return
		}

		err = ctx.JSON(iris.StatusOK, chapter)
		if err != nil {
			ctx.JSON(iris.StatusInternalServerError, &jsonError{
				Error: err.Error(),
			})
			return
		}
	})

	app.Listen(":8080")
}
