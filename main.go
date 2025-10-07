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
    }
    // =======================================================
    // --- ROUTE KHUSUS OPD ---
    // =======================================================
    opdRoutes := api.Group("/")
    opdRoutes.Use(AuthMiddleware("opd")) 
    {
        // 1. ROUTE MASTER: OPD membuat standar pelayanan mereka sendiri
        opdRoutes.POST("/standar-pelayanan", CreateJenisPelayanan) // <-- PERUBAHAN UTAMA

        // 2. ROUTE TRANSAKSI: Create/Update Form Pengajuan
        pengajuanOPD := opdRoutes.Group("/pengajuan")
        {
            pengajuanOPD.POST("/", CreateFormPengajuan)
            pengajuanOPD.PUT("/:id", UpdateFormPengajuan)
            pengajuanOPD.DELETE("/:id", DeleteFormPengajuan) 
            
            // OPD melihat laporannya sendiri
            opdRoutes.GET("/user/:id/pengajuan", GetFormPengajuanByUserOPD) 
        }
    }
    // =======================================================
    // --- ROUTE UNTUK PEMDA (Melihat & Validasi) ---
    // =======================================================
    pemdaRoutes := api.Group("/pengajuan")
    pemdaRoutes.Use(AuthMiddleware("pemda")) // Hanya Pemda yang boleh
    {
        // Pemda bisa melihat semua data (GET all)
        pemdaRoutes.GET("/", GetAllFormPengajuan) 

        // Pemda bisa melakukan validasi
        pemdaRoutes.POST("/:id/validate", ValidateFormPengajuan)
    }
    // =======================================================
    // --- ROUTE YANG BISA DIAKSES KEDUA ROLE (Get Detail) ---
    // =======================================================
    sharedRoutes := api.Group("/pengajuan")
    sharedRoutes.Use(AuthMiddleware("opd", "pemda")) // OPD & Pemda boleh
    {
        // Keduanya bisa lihat detail
        sharedRoutes.GET("/:id", GetFormPengajuanByID)
    }
	// Middleware CORS
	
	// Start server
	port := ":8080"
	log.Println("🚀 Server running on http://localhost" + port)
	r.Run(port) // Cara menjalankan server di Gin
}