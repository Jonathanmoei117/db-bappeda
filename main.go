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

	// Router Gin. `gin.Default()` sudah termasuk logger dan recovery middleware.
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


    // =======================================================
    // --- ROUTE KHUSUS ADMIN / SUPERUSER (jika ada) ---
    // Di sini kita asumsikan 'pemda' juga bertindak sebagai admin
    // =======================================================
    adminRoutes := api.Group("/")
    adminRoutes.Use(AuthMiddleware("pemda")) // Hanya pemda/admin yang boleh
    {
        // Route untuk mengelola data master OPD
        adminRoutes.POST("/opd", CreateOPD)
        adminRoutes.GET("/opd", GetAllOPD)

        // Route untuk mendaftarkan user baru
        register := adminRoutes.Group("/register")
        {
            register.POST("/opd", CreateUserOPD)
            register.POST("/pemda", CreateUserPemda)
        }
    }


    // =======================================================
    // --- ROUTE UNTUK PENGAJUAN OLEH OPD ---
    // =======================================================
    opdRoutes := api.Group("/")
    opdRoutes.Use(AuthMiddleware("opd")) // Hanya OPD yang boleh
    {
        layananOPD := opdRoutes.Group("/layanan")
        {
            // OPD hanya bisa membuat dan meng-update pengajuan
            layananOPD.POST("/pembangunan", CreateLayananPembangunan)
            layananOPD.PUT("/pembangunan/:id", UpdateLayananPembangunan)
            
            layananOPD.POST("/administrasi", CreateLayananAdministrasi)
            layananOPD.PUT("/administrasi/:id", UpdateLayananAdministrasi)

            layananOPD.POST("/informasi", CreateLayananInformasiPengaduan)
            layananOPD.PUT("/informasi/:id", UpdateLayananInformasiPengaduan)
        }
        
        // OPD hanya bisa melihat laporannya sendiri
        opdRoutes.GET("/user-opd/:id/layanan/pembangunan", GetLayananPembangunanByUserOPD)
        opdRoutes.GET("/user-opd/:id/layanan/administrasi", GetLayananAdministrasiByUserOPD)
        opdRoutes.GET("/user-opd/:id/layanan/informasi", GetLayananInformasiPengaduanByUserOPD)
    }


    // =======================================================
    // --- ROUTE UNTUK PEMDA (Melihat & Validasi) ---
    // =======================================================
    pemdaRoutes := api.Group("/")
    pemdaRoutes.Use(AuthMiddleware("pemda")) // Hanya Pemda yang boleh
    {
        layananPemda := pemdaRoutes.Group("/layanan")
        {
            // Pemda bisa melihat semua data (GET all)
            layananPemda.GET("/pembangunan", GetAllLayananPembangunan)
            layananPemda.GET("/administrasi", GetAllLayananAdministrasi)
            layananPemda.GET("/informasi", GetAllLayananInformasiPengaduan)

            // Pemda bisa melakukan validasi
            layananPemda.POST("/pembangunan/:id/validate", ValidateLayananPembangunan)
            layananPemda.POST("/administrasi/:id/validate", ValidateLayananAdministrasi)
            layananPemda.POST("/informasi/:id/validate", ValidateLayananInformasiPengaduan)
        }
    }


    // =======================================================
    // --- ROUTE YANG BISA DIAKSES KEDUA ROLE (Contoh: GetByID) ---
    // =======================================================
    sharedRoutes := api.Group("/")
    sharedRoutes.Use(AuthMiddleware("opd", "pemda")) // OPD & Pemda boleh
    {
        layananShared := sharedRoutes.Group("/layanan")
        {
            // Keduanya bisa lihat detail, keamanan kepemilikan sudah ada di dalam handler
            layananShared.GET("/pembangunan/:id", GetLayananPembangunanByID)
            layananShared.DELETE("/pembangunan/:id", DeleteLayananPembangunan)
            
            layananShared.GET("/administrasi/:id", GetLayananAdministrasiByID)
            layananShared.DELETE("/administrasi/:id", DeleteLayananAdministrasi)
            
            layananShared.GET("/informasi/:id", GetLayananInformasiPengaduanByID)
            layananShared.DELETE("/informasi/:id", DeleteLayananInformasiPengaduan)
        }
    }

	// Middleware CORS
	
	// Start server
	port := ":8080"
	log.Println("🚀 Server running on http://localhost" + port)
	r.Run(port) // Cara menjalankan server di Gin
}
