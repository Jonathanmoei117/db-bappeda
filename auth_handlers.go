package main

import (
	"log"
	"net/http"
	"os"
	"strings" // <-- TAMBAHKAN IMPORT INI
	"time"

	"github.com/gin-gonic/gin" // Impor Gin
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Kunci rahasia untuk menandatangani token JWT.
// SANGAT PENTING: Simpan ini di file .env Anda, jangan di hardcode!
var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// LoginRequest adalah struct untuk menampung data dari body request login.
type LoginRequest struct {
	NIP      string `json:"nip" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Claims adalah struct untuk data yang akan kita simpan di dalam token JWT.
type Claims struct {
	ID      uint   `json:"id"`
	NIP     string `json:"nip"`
	Nama    string `json:"nama"`
	Jabatan string `json:"jabatan"`
	Role    string `json:"role"` // 'opd' atau 'pemda'
	jwt.RegisteredClaims
}

// LoginHandler menangani proses login untuk kedua role.
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NIP dan Password tidak boleh kosong"})
		return
	}

	log.Println("--- [LOGIN ATTEMPT] Mencoba mencari NIP", req.NIP, "di tabel 'user_opd' ---")
	var userOPD UserOPD
	err := DB.Where("nip = ?", req.NIP).First(&userOPD).Error
	if err == nil {
		log.Println("--- [LOGIN INFO] User DITEMUKAN di 'user_opd'. Mencocokkan password... ---")
		if err := bcrypt.CompareHashAndPassword([]byte(userOPD.Password), []byte(req.Password)); err == nil {
			log.Println("--- [LOGIN SUCCESS] Password cocok untuk user OPD. Membuat token... ---")
			generateTokenAndRespond(c, userOPD.ID, userOPD.NIP, userOPD.Nama, userOPD.Jabatan, "opd")
			return
		}
		log.Println("--- [LOGIN FAILED] Password user OPD tidak cocok. ---")
	}

	log.Println("--- [LOGIN INFO] User TIDAK DITEMUKAN di 'user_opd'. Mencoba mencari di 'user_pemda'... ---")
	var userPemda UserPemda
	err = DB.Where("nip = ?", req.NIP).First(&userPemda).Error
	if err == nil {
		log.Println("--- [LOGIN INFO] User DITEMUKAN di 'user_pemda'. Mencocokkan password... ---")
		if err := bcrypt.CompareHashAndPassword([]byte(userPemda.Password), []byte(req.Password)); err == nil {
			log.Println("--- [LOGIN SUCCESS] Password cocok untuk user Pemda. Membuat token... ---")
			generateTokenAndRespond(c, userPemda.ID, userPemda.NIP, userPemda.Nama, userPemda.Jabatan, "pemda")
			return
		}
		log.Println("--- [LOGIN FAILED] Password user Pemda tidak cocok. ---")
	}

	log.Println("--- [LOGIN FAILED] User TIDAK DITEMUKAN di kedua tabel, atau password salah. ---")
	c.JSON(http.StatusUnauthorized, gin.H{"error": "NIP atau Password salah"})
}

// generateTokenAndRespond adalah helper untuk membuat JWT dan mengirimkannya sebagai respons.
func generateTokenAndRespond(c *gin.Context, id uint, nip, nama, jabatan, role string) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		ID:      id,
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":      id,
			"nip":     nip,
			"nama":    nama,
			"jabatan": jabatan,
			"role":    role,
		},
	})
}

// =======================================================
// TAMBAHKAN FUNGSI MIDDLEWARE DI SINI
// =======================================================
func AuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil token dari header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header tidak ditemukan"})
			return
		}

		// Header biasanya formatnya "Bearer {token}", kita buang "Bearer "-nya
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Format token tidak valid"})
			return
		}

		// 2. Validasi token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau kedaluwarsa"})
			return
		}

		// 3. Cek apakah role user diizinkan
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

		// 4. Jika semua berhasil, simpan info user di context dan lanjutkan ke handler berikutnya
		c.Set("user", claims)
		c.Next()
	}
}