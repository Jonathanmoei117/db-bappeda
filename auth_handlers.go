package main

import (
	"log"
	"net/http"
	"os"
	
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// jwtKey diambil dari environment variable untuk keamanan.
var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// LoginRequest merepresentasikan body JSON yang diharapkan saat login.
type LoginRequest struct {
	NIP      string `json:"nip" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Claims merepresentasikan data (payload) yang akan disimpan di dalam JWT.
type Claims struct {
	ID      uint   `json:"id"`
	IDOPD   uint   `json:"id_opd"` // IDOPD disertakan untuk user dari OPD.
	NIP     string `json:"nip"`
	Nama    string `json:"nama"`
	Jabatan string `json:"jabatan"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

// LoginHandler memproses permintaan login untuk UserOPD dan UserPemda.
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NIP dan Password tidak boleh kosong"})
		return
	}

	// Coba cari pengguna di tabel user_opd terlebih dahulu.
	var userOPD UserOPD
	// [DIUBAH] Gunakan Preload("OPD") untuk mengambil data OPD terkait secara otomatis
	if err := DB.Preload("OPD").Where("nip = ?", req.NIP).First(&userOPD).Error; err == nil {
		if bcrypt.CompareHashAndPassword([]byte(userOPD.Password), []byte(req.Password)) == nil {
			log.Println("[LOGIN SUCCESS] Role: OPD, Nama:", userOPD.Nama, ", OPD:", userOPD.OPD.NamaOPD)
			// Kirim nama OPD ke fungsi generateToken
			generateTokenAndRespond(c, userOPD.ID, userOPD.IDOPD, userOPD.NIP, userOPD.Nama, userOPD.Jabatan, "opd", userOPD.OPD.NamaOPD)
			return
		}
	}

	// Jika tidak ditemukan atau password salah, coba cari di tabel user_pemda.
	var userPemda UserPemda
	if err := DB.Where("nip = ?", req.NIP).First(&userPemda).Error; err == nil {
		if bcrypt.CompareHashAndPassword([]byte(userPemda.Password), []byte(req.Password)) == nil {
			log.Println("[LOGIN SUCCESS] Role: Pemda, Nama:", userPemda.Nama)
			// Untuk Pemda, nama OPD bisa string kosong
			generateTokenAndRespond(c, userPemda.ID, 0, userPemda.NIP, userPemda.Nama, userPemda.Jabatan, "pemda", "Pemerintah Daerah")
			return
		}
	}

	// Jika setelah semua pengecekan pengguna tidak ditemukan atau password salah.
	log.Println("[LOGIN FAILED] NIP:", req.NIP)
	c.JSON(http.StatusUnauthorized, gin.H{"error": "NIP atau Password salah"})
}

func LogoutHandler(c *gin.Context) {
	// Atur cookie opd_token dengan MaxAge negatif untuk menghapusnya
	c.SetCookie(
		"opd_token",
		"",
		-1, // MaxAge = -1 akan langsung menghapus cookie
		"/",
		"localhost",
		false, // false untuk HTTP di development
		true,  // httpOnly
	)

	// Atur juga cookie pemda_token untuk memastikan semua bersih
	c.SetCookie(
		"pemda_token",
		"",
		-1,
		"/",
		"localhost",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Logout berhasil"})
}

// generateTokenAndRespond membuat JWT dan mengirimkannya sebagai respons JSON.
func generateTokenAndRespond(c *gin.Context, id uint, idOPD uint, nip, nama, jabatan, role string,opdName string) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		ID:      id,
		IDOPD:   idOPD,
		NIP:     nip,
		Nama:    nama,
		Jabatan: jabatan,
		Role:    role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat token"})
		return
	}

	// **[BARU]** Tentukan nama cookie dan path redirect berdasarkan role.
	var cookieName string
	var redirectPath string
		domain := os.Getenv("APP_DOMAIN")
		if domain == "" {
			domain = "localhost"
		}
		isSecure := os.Getenv("GIN_MODE") == "release"
	if role == "pemda" {
		cookieName = "pemda_token"
		redirectPath = "/pemda/dashboard"
	} else { // Asumsi jika bukan pemda, pasti opd
		cookieName = "opd_token"
		redirectPath = "/opd/dashboard"
	}

	// **[BARU]** Set token sebagai HttpOnly cookie di browser client.
	// Ini adalah langkah kunci untuk Next.js middleware.
	c.SetCookie(
		cookieName,
		tokenString,
		3600*24,      // maxAge
		"/",          // path
		domain,   // <-- Menjadi dinamis
        isSecure,    // secure (false untuk development HTTP di localhost)
		true,         // httpOnly
	)

	// **[BARU]** Kirim respons ke client berisi URL redirect dan data user.
	// Kita tidak perlu lagi mengirim token di body JSON karena sudah ada di cookie.
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"redirect": redirectPath, // Frontend akan menggunakan ini untuk navigasi
		"user": gin.H{
			"id":      id,
			"id_opd":  idOPD,
			"nip":     nip,
			"nama":    nama,
			"jabatan": jabatan,
			"role":    role,
			"opd":     opdName, // <-- TAMBAHKAN NAMA OPD DI SINI
		},
	})
}

// AuthMiddleware adalah middleware untuk memverifikasi JWT dan hak akses (role).
func AuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		  tokenString, err := c.Cookie("opd_token")
        if err != nil {
            tokenString, err = c.Cookie("pemda_token")
            if err != nil {
                c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token otentikasi tidak ditemukan di cookie"})
                return
            }
        }
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau kedaluwarsa"})
			return
		}

		// Cek apakah role pengguna diizinkan mengakses endpoint ini.
		isAllowed := false
		for _, role := range allowedRoles {
			if claims.Role == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki hak akses untuk sumber daya ini"})
			return
		}

		// Simpan data user dari token ke dalam context untuk digunakan di handler selanjutnya.
		c.Set("user", claims)
		c.Next()
	}
}