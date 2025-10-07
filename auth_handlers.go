package main

import (
	"log"
	"net/http"
	"os"
	"strings"
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
	if err := DB.Where("nip = ?", req.NIP).First(&userOPD).Error; err == nil {
		// User OPD ditemukan, verifikasi password.
		if bcrypt.CompareHashAndPassword([]byte(userOPD.Password), []byte(req.Password)) == nil {
			log.Println("[LOGIN SUCCESS] Role: OPD, Nama:", userOPD.Nama)
			generateTokenAndRespond(c, userOPD.ID, userOPD.IDOPD, userOPD.NIP, userOPD.Nama, userOPD.Jabatan, "opd")
			return
		}
	}

	// Jika tidak ditemukan atau password salah, coba cari di tabel user_pemda.
	var userPemda UserPemda
	if err := DB.Where("nip = ?", req.NIP).First(&userPemda).Error; err == nil {
		// User Pemda ditemukan, verifikasi password.
		if bcrypt.CompareHashAndPassword([]byte(userPemda.Password), []byte(req.Password)) == nil {
			log.Println("[LOGIN SUCCESS] Role: Pemda, Nama:", userPemda.Nama)
			// IDOPD di-set 0 karena User Pemda tidak terikat pada satu OPD.
			generateTokenAndRespond(c, userPemda.ID, 0, userPemda.NIP, userPemda.Nama, userPemda.Jabatan, "pemda")
			return
		}
	}

	// Jika setelah semua pengecekan pengguna tidak ditemukan atau password salah.
	log.Println("[LOGIN FAILED] NIP:", req.NIP)
	c.JSON(http.StatusUnauthorized, gin.H{"error": "NIP atau Password salah"})
}

// generateTokenAndRespond membuat JWT dan mengirimkannya sebagai respons JSON.
func generateTokenAndRespond(c *gin.Context, id uint, idOPD uint, nip, nama, jabatan, role string) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	// Kirim respons ke client berisi token dan data user.
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":      id,
			"id_opd":  idOPD,
			"nip":     nip,
			"nama":    nama,
			"jabatan": jabatan,
			"role":    role,
		},
	})
}

// AuthMiddleware adalah middleware untuk memverifikasi JWT dan hak akses (role).
func AuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header tidak ditemukan"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Format token tidak valid, harus 'Bearer <token>'"})
			return
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