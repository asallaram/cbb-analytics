package main

import (
	"log"
	"net/http"

	"github.com/asallaram/cbb-analytics/internal/api"
	"github.com/asallaram/cbb-analytics/internal/storage"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	mongo, err := storage.NewMongoDB("mongodb+srv://aneeshsallaram_db_user:NYGiants1@cluster0.xumhdzd.mongodb.net/?appName=Cluster0", "cbb_analytics")
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

	log.Println("🚀 API Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
