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

	// --- Endpoints untuk Layanan Pembangunan ---
	r.HandleFunc("/layanan-pembangunan", GetAllLayananPembangunan).Methods("GET")
	r.HandleFunc("/layanan-pembangunan/{id:[0-9]+}", GetLayananPembangunanByID).Methods("GET")
	r.HandleFunc("/layanan-pembangunan", CreateLayananPembangunan).Methods("POST")
	r.HandleFunc("/layanan-pembangunan/{id:[0-9]+}", UpdateLayananPembangunan).Methods("PUT")
	r.HandleFunc("/layanan-pembangunan/{id:[0-9]+}", DeleteLayananPembangunan).Methods("DELETE")

	// --- Endpoints untuk Layanan Administrasi ---
	r.HandleFunc("/layanan-administrasi", GetAllLayananAdministrasi).Methods("GET")
	r.HandleFunc("/layanan-administrasi/{id:[0-9]+}", GetLayananAdministrasiByID).Methods("GET")
	r.HandleFunc("/layanan-administrasi", CreateLayananAdministrasi).Methods("POST")
	r.HandleFunc("/layanan-administrasi/{id:[0-9]+}", UpdateLayananAdministrasi).Methods("PUT")
	r.HandleFunc("/layanan-administrasi/{id:[0-9]+}", DeleteLayananAdministrasi).Methods("DELETE")

	// --- Endpoints untuk Layanan Informasi & Pengaduan ---
	r.HandleFunc("/layanan-informasi-pengaduan", GetAllLayananInformasiPengaduan).Methods("GET")
	r.HandleFunc("/layanan-informasi-pengaduan/{id:[0-9]+}", GetLayananInformasiPengaduanByID).Methods("GET")
	r.HandleFunc("/layanan-informasi-pengaduan", CreateLayananInformasiPengaduan).Methods("POST")
	r.HandleFunc("/layanan-informasi-pengaduan/{id:[0-9]+}", UpdateLayananInformasiPengaduan).Methods("PUT")
	r.HandleFunc("/layanan-informasi-pengaduan/{id:[0-9]+}", DeleteLayananInformasiPengaduan).Methods("DELETE")


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