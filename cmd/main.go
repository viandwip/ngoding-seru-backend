package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/oktaviandwip/musalabel-backend/internal/routers"
	"github.com/oktaviandwip/musalabel-backend/pkg"
)

func main() {
	db, err := pkg.Posql()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server running on port:" + os.Getenv("PORT"))

	router := routers.New(db)
	server := pkg.Server(router)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
