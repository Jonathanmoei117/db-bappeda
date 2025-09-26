package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("⚠️ No .env file found, using system env")
	}

	// Init DB + AutoMigrate
	InitDB()

	// Jalankan seeder jika ada argumen "seed"
	if len(os.Args) > 1 && os.Args[1] == "seed" {
		Seed()
		fmt.Println("🌱 Seeding selesai!")
		return
	}

	// Router
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/opd", CreateOPD).Methods("POST")
    api.HandleFunc("/opd", GetAllOPD).Methods("GET")

	// ---- Routes untuk Registrasi User ----
	api.HandleFunc("/register/opd", CreateUserOPD).Methods("POST")
	api.HandleFunc("/register/pemda", CreateUserPemda).Methods("POST")

	// ---- Routes untuk Layanan Pembangunan ----
	api.HandleFunc("/layanan/pembangunan", CreateLayananPembangunan).Methods("POST")
	api.HandleFunc("/layanan/pembangunan", GetAllLayananPembangunan).Methods("GET")
	api.HandleFunc("/layanan/pembangunan/{id}", GetLayananPembangunanByID).Methods("GET")
	api.HandleFunc("/layanan/pembangunan/{id}", UpdateLayananPembangunan).Methods("PUT")
	api.HandleFunc("/layanan/pembangunan/{id}/validate", ValidateLayananPembangunan).Methods("POST")

	// ---- Routes untuk Layanan Administrasi ----
	api.HandleFunc("/layanan/administrasi", CreateLayananAdministrasi).Methods("POST")
	api.HandleFunc("/layanan/administrasi", GetAllLayananAdministrasi).Methods("GET")
	api.HandleFunc("/layanan/administrasi/{id}", GetLayananAdministrasiByID).Methods("GET")
	api.HandleFunc("/layanan/administrasi/{id}", UpdateLayananAdministrasi).Methods("PUT")
	api.HandleFunc("/layanan/administrasi/{id}/validate", ValidateLayananAdministrasi).Methods("POST")

	// ---- Routes untuk Layanan Informasi & Pengaduan ----
	api.HandleFunc("/layanan/informasi", CreateLayananInformasi).Methods("POST")
	api.HandleFunc("/layanan/informasi", GetAllLayananInformasi).Methods("GET")
	api.HandleFunc("/layanan/informasi/{id}", GetLayananInformasiByID).Methods("GET")
	api.HandleFunc("/layanan/informasi/{id}", UpdateLayananInformasi).Methods("PUT")
	api.HandleFunc("/layanan/informasi/{id}/validate", ValidateLayananInformasi).Methods("POST")


	// Middleware CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Origin", "Content-Type", "Authorization"}),
	)

	// Start server
	port := ":8080"
	log.Println("🚀 Server running on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, corsHandler(r)))
}
