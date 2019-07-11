package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jadilet/ad_zh/handlers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jadilet/ad_zh/models"
)

func main() {
	db, err := models.InitDB()

	if err != nil {
		log.Panic(err)
	}

	env := &models.Env{DB: db}
	env.DB.SetConnMaxLifetime(time.Minute * 5)
	env.DB.SetMaxIdleConns(0)
	env.DB.SetMaxOpenConns(5)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler: handlers.New(env),
	}

	log.Printf("Starting server. Listening at %q", server.Addr)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("%v", err)
	} else {
		log.Printf("Sever closed!")
	}

	defer env.DB.Close()
}
