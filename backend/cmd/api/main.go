package main

import (
	"log"
	"net/http"
	"os"

	"github.com/asallaram/cbb-analytics/internal/api"
	"github.com/asallaram/cbb-analytics/internal/storage"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	mongo, err := storage.NewMongoDB(mongoURI, "cbb_analytics")
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer mongo.Close()

	router := mux.NewRouter()
	h := api.NewHandler(mongo)

	router.HandleFunc("/api/games", h.GetGames).Methods("GET")
	router.HandleFunc("/api/games/{id}", h.GetGame).Methods("GET")
	router.HandleFunc("/api/games/{id}/plays", h.GetPlays).Methods("GET")
	router.HandleFunc("/api/games/{id}/stats", h.GetStats).Methods("GET")
	router.HandleFunc("/api/games/{id}/zones", h.GetZones).Methods("GET")
	router.HandleFunc("/api/games/{id}/insights", h.GetInsights).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	corsHandler := c.Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 API Server running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
