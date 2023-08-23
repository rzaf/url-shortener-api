package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rzaf/url-shortener-api/database"
	"github.com/rzaf/url-shortener-api/routes"
)

func main() {
	godotenv.Load(".env")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	addr := os.Getenv("ADDR")
	port := os.Getenv("PORT")
	if addr == "" {
		log.Fatalln("ADDR NOT FOUND")
	}
	if port == "" {
		log.Fatalln("PORT NOT FOUND")
	}
	database.Connect()
	fmt.Printf("listening at: %v:%v \n", addr, port)

	baseRouter := routes.GetRoutes()
	err := http.ListenAndServe(addr+":"+port, baseRouter)
	if err != nil {
		log.Fatalln(err)
	}
}
