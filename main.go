package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors" // Impor CORS untuk Gin
	"github.com/gin-gonic/gin"    // Impor Gin
	"github.com/joho/godotenv"
)

func main() {
	// log.Fatal("SERVER SENGAJA DIMATIKAN UNTUK TES")
	// Load .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("⚠ No .env file found, using system env")
	}

	// Init DB + AutoMigrate
	InitDB()

	// Jalankan seeder jika ada argumen "seed"
	if len(os.Args) > 1 && os.Args[1] == "seed" {
		Seed()
		fmt.Println("🌱 Seeding selesai!")
		return
	}

	// Router Gin. gin.Default() sudah termasuk logger dan recovery middleware.
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Static file server untuk folder "uploads"
	r.Static("/uploads", "./uploads")

	// Grup utama untuk semua endpoint di bawah /api
	api := r.Group("/api")

	// =======================================================
	// --- ROUTE PUBLIK (Tidak Perlu Login) ---
	// =======================================================
	api.POST("/login", LoginHandler)
	api.POST("/logout", LogoutHandler)

	// Rute GET /standar-pelayanan tetap publik agar semua user bisa melihat standar yang tersedia
	api.GET("/standar-pelayanan", GetAllJenisPelayanan) // Tampilan publik / daftar master

	// =======================================================
	// --- ROUTE KHUSUS ADMIN / SUPERUSER (PEMDA) ---
	// =======================================================
	adminRoutes := api.Group("/")
	adminRoutes.Use(AuthMiddleware("pemda")) // Hanya pemda/admin yang boleh
	{
		// 1. Route untuk mengelola data master OPD
		adminRoutes.POST("/opd", CreateOPD)
		adminRoutes.GET("/opd", GetAllOPD)

		// 2. Route untuk mendaftarkan user baru
		register := adminRoutes.Group("/register")
		{
			register.POST("/opd", CreateUserOPD)
			register.POST("/pemda", CreateUserPemda)
		}

		// 3. Route Validasi Standar Pelayanan (BARU)
		adminRoutes.POST("/standar-pelayanan/:id/validate", ValidateJenisPelayanan)

		// 4. Route Pemda untuk melihat SEMUA Form Pengajuan
		adminRoutes.GET("/pengajuan", GetAllFormPengajuan)
	}

	// =======================================================
	// --- ROUTE KHUSUS OPD ---
	// =======================================================
	opdRoutes := api.Group("/")
	opdRoutes.Use(AuthMiddleware("opd"))
	{
		// 1. ROUTE MASTER: OPD membuat standar pelayanan mereka sendiri
		opdRoutes.POST("/standar-pelayanan", CreateJenisPelayanan)
		opdRoutes.GET("/standar-pelayanan/opd/:id_opd", GetStandarPelayananByOPD)
		opdRoutes.GET("/user/:id/pengajuan", GetFormPengajuanByUserOPD)

		// 2. ROUTE TRANSAKSI: Create/Update Form Pengajuan
		opdRoutes.POST("/pengajuan", CreateFormPengajuan)
		opdRoutes.PUT("/pengajuan/:id", UpdateFormPengajuan)
		opdRoutes.DELETE("/pengajuan/:id", DeleteFormPengajuan)

		// 3. ROUTE MASTER PEMOHON: OPD mengelola data master pemohon (BARU)
		pemohonRoutes := opdRoutes.Group("/form-pemohon")
		{
			pemohonRoutes.POST("/", CreateFormPemohon)
			pemohonRoutes.GET("/", GetAllFormPemohon)
			pemohonRoutes.GET("/:id", GetFormPemohonByID)
			pemohonRoutes.PUT("/:id", UpdateFormPemohon)
			pemohonRoutes.DELETE("/:id", DeleteFormPemohon)
		}
	}

	// =======================================================
	// --- ROUTE YANG BISA DIAKSES KEDUA ROLE (Get Detail) ---
	// =======================================================
	sharedRoutes := api.Group("/") // <-- Grup ini tidak lagi "/pengajuan"
	sharedRoutes.Use(AuthMiddleware("opd", "pemda")) // OPD & Pemda boleh
	{
		// Keduanya bisa lihat detail pengajuan
		sharedRoutes.GET("/pengajuan/:id", GetFormPengajuanByID)
	}

	// Start server
	port := ":8080"
	log.Println("🚀 Server running on http://localhost" + port)
	r.Run(port) // Cara menjalankan server di Gin
}