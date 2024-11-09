package main

import (
	"net/http"
	"stt/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("cannot load env file")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	routes.Setup(r)

	println("app is running on http://localhost:3000")
	http.ListenAndServe(":3000", r)
}
