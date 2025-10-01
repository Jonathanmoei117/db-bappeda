package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ========= HANDLERS REGISTRASI PENGGUNA (OPD & PEMDA) =========

func CreateUserOPD(c *gin.Context) {
	var user UserOPD
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body tidak valid"})
		return
	}

	// Di aplikasi nyata, HASHING PASSWORD WAJIB DILAKUKAN DI SINI!
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	// if err != nil { ... }
	// user.Password = string(hashedPassword)

	if err := DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Jangan pernah kirim password kembali ke client
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

func CreateUserPemda(c *gin.Context) {
	var user UserPemda
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body tidak valid"})
		return
	}

	// HASHING PASSWORD WAJIB
	if err := DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

// ========= CRUD HANDLERS: OPD (KHUSUS SUPER ADMIN) =========
func CreateOPD(c *gin.Context) {
	var opd OPD
	if err := c.ShouldBindJSON(&opd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body tidak valid"})
		return
	}

	if err := DB.Create(&opd).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, opd)
}

func GetAllOPD(c *gin.Context) {
	var opds []OPD
	if err := DB.Find(&opds).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, opds)
}

// ========= HELPER UNTUK AMBIL FILE & FORM VALUE =========

func bindLayananFromMultipartForm(c *gin.Context, layanan interface{}) error {
	// Parsing form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		return errors.New("Gagal parsing form: " + err.Error())
	}

	// Handle file upload terlebih dahulu (jika ada)
	var uploadedFilename *string
	_, handler, err := c.Request.FormFile("berkas_pengajuan")
	if err == nil { // Jika ada file yang diupload
		os.MkdirAll("./uploads", os.ModePerm)
		filePath := "./uploads/" + handler.Filename
		if err := c.SaveUploadedFile(handler, filePath); err != nil {
			return errors.New("Gagal menyimpan file: " + err.Error())
		}
		uploadedFilename = &handler.Filename
	}

	// Binding data umum berdasarkan interface
	switch v := layanan.(type) {
	case *LayananPembangunan:
		// HAPUS SEMUA LOGIKA user_opd_id DARI SINI
		v.JudulKegiatan = c.PostForm("judul_kegiatan")
		v.Deskripsi = c.PostForm("deskripsi")
		v.NamaPemohon = c.PostForm("nama_pemohon")
		v.InstansiPemohon = c.PostForm("instansi_pemohon")
		if nip := c.PostForm("nip_pemohon"); nip != "" {
			v.NIPPemohon = &nip
		}
		if mulai := c.PostForm("periode_mulai"); mulai != "" {
			v.PeriodeMulai = &mulai
		}
		if selesai := c.PostForm("periode_selesai"); selesai != "" {
			v.PeriodeSelesai = &selesai
		}
		v.BerkasPengajuanPath = uploadedFilename // Set path file

	case *LayananAdministrasi:
		// HAPUS SEMUA LOGIKA user_opd_id DARI SINI
		v.JudulKegiatan = c.PostForm("judul_kegiatan")
		v.Deskripsi = c.PostForm("deskripsi")
		v.NamaPemohon = c.PostForm("nama_pemohon")
		v.InstansiPemohon = c.PostForm("instansi_pemohon")
		if nip := c.PostForm("nip_pemohon"); nip != "" {
			v.NIPPemohon = &nip
		}
		if mulai := c.PostForm("periode_mulai"); mulai != "" {
			v.PeriodeMulai = &mulai
		}
		if selesai := c.PostForm("periode_selesai"); selesai != "" {
			v.PeriodeSelesai = &selesai
		}
		v.BerkasPengajuanPath = uploadedFilename // Set path file

	case *LayananInformasiPengaduan:
		// HAPUS SEMUA LOGIKA user_opd_id DARI SINI
		v.JudulKegiatan = c.PostForm("judul_kegiatan")
		v.Deskripsi = c.PostForm("deskripsi")
		v.NamaPemohon = c.PostForm("nama_pemohon")
		v.InstansiPemohon = c.PostForm("instansi_pemohon")
		if nip := c.PostForm("nip_pemohon"); nip != "" {
			v.NIPPemohon = &nip
		}
		v.BerkasPengajuanPath = uploadedFilename // Set path file
	}

	return nil
}

// ========= CRUD: LAYANAN PEMBANGUNAN =========

func CreateLayananPembangunan(c *gin.Context) {
	// Ambil data user yang login dari context (sudah di-set oleh AuthMiddleware)
	userClaims, _ := c.Get("user")
	claims := userClaims.(*Claims)

	// UserOPDID sekarang diambil dari token, bukan form. JAUH LEBIH AMAN.
	var layanan = LayananPembangunan{
		UserOPDID: claims.ID, 
	}
	
	// Bind sisa data dari form ke struct 'layanan'
	if err := bindLayananFromMultipartForm(c, &layanan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set nilai default saat pembuatan
	layanan.StatusProses = "Baru"
	layanan.StatusValidasi = "Menunggu Validasi"
	layanan.CreatedAt = time.Now()

	if err := DB.Create(&layanan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, layanan)
}

func GetAllLayananPembangunan(c *gin.Context) {
	 log.Println("--- 1. MASUK KE HANDLER GetAllLayananPembangunan ---")

    var layanans []LayananPembangunan
    log.Println("--- 2. Variabel 'layanans' berhasil dibuat ---")

    // Kita coba query TANPA Preload dulu untuk menyederhanakan
    tx := DB.Find(&layanans)
    log.Println("--- 3. Perintah DB.Find sudah dieksekusi ---")

    if tx.Error != nil {
        log.Println("!!! ERROR SAAT QUERY DATABASE:", tx.Error.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
        return
    }

    log.Println("--- 4. Query database BERHASIL. Jumlah data:", tx.RowsAffected, "---")

    c.JSON(http.StatusOK, layanans)
    log.Println("--- 5. Response JSON berhasil dikirim ---")
}

func GetLayananPembangunanByID(c *gin.Context) {
	layananID := c.Param("id")
	userClaims, _ := c.Get("user")
	claims := userClaims.(*Claims)

	var layanan LayananPembangunan
	if err := DB.First(&layanan, layananID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data layanan tidak ditemukan"})
		return
	}

	// Jika user adalah 'opd', cek apakah dia pemilik data ini
	if claims.Role == "opd" && layanan.UserOPDID != claims.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki hak akses untuk melihat data ini"})
		return
	}

	// Jika rolenya 'pemda' atau dia adalah pemilik, lanjutkan
	DB.Preload("UserOPD.OPD").First(&layanan, layananID)
	c.JSON(http.StatusOK, layanan)
}

func GetLayananPembangunanByUserOPD(c *gin.Context) {
	userOpdID := c.Param("id")
	var layanans []LayananPembangunan
	
	err := DB.Where("user_opd_id = ?", userOpdID).Preload("UserOPD.OPD").Find(&layanans).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, layanans)
}

func UpdateLayananPembangunan(c *gin.Context) {
	id := c.Param("id")
	var layanan LayananPembangunan

	if err := DB.First(&layanan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	// Bind data dari form ke struct yang sudah ada
	if err := bindLayananFromMultipartForm(c, &layanan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := DB.Save(&layanan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data"})
		return
	}

	c.JSON(http.StatusOK, layanan)
}

func DeleteLayananPembangunan(c *gin.Context) {
	id := c.Param("id")
	if err := DB.Delete(&LayananPembangunan{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}

func ValidateLayananPembangunan(c *gin.Context) {
	// 1. Ambil ID layanan dari parameter URL
	id := c.Param("id")

	// 2. Deklarasikan variabel 'layanan'
	var layanan LayananPembangunan

	// 3. Cari layanan di database
	if err := DB.First(&layanan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Layanan Pembangunan tidak ditemukan"})
		return
	}

	// 4. Cek apakah sudah divalidasi sebelumnya
	if layanan.StatusValidasi != "Menunggu Validasi" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Layanan ini sudah divalidasi"})
		return
	}

	// 5. Bind request body JSON
	var req ValidasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body tidak valid"})
		return
	}

	// 6. Ambil data user 'pemda' yang sedang login dari context
	userClaims, _ := c.Get("user")
	claims := userClaims.(*Claims)

	// 7. Siapkan data untuk update
	validatorIDFromToken := claims.ID // Ambil ID dari token
	now := time.Now()

	// 8. Update field di struct 'layanan'
	layanan.StatusValidasi = req.StatusValidasi
	layanan.KeteranganValidasi = req.KeteranganValidasi
	layanan.IDValidatorPemda = &validatorIDFromToken
	layanan.TanggalValidasi = &now

	// 9. Simpan perubahan ke database
	DB.Save(&layanan)
	
	// 10. Ambil data terbaru dengan relasi untuk dikirim kembali
	DB.Preload("ValidatorPemda").First(&layanan, layanan.ID)
	c.JSON(http.StatusOK, layanan)
}

// ========= CRUD: LAYANAN ADMINISTRASI =========

func CreateLayananAdministrasi(c *gin.Context) {
	// Ambil data user yang login dari context (sudah di-set oleh AuthMiddleware)
	userClaims, _ := c.Get("user")
	claims := userClaims.(*Claims)

	// UserOPDID sekarang diambil dari token, bukan form. JAUH LEBIH AMAN.
	var layanan = LayananAdministrasi{
		UserOPDID: claims.ID, 
	}
	
	// Bind sisa data dari form ke struct 'layanan'
	if err := bindLayananFromMultipartForm(c, &layanan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set nilai default saat pembuatan
	layanan.StatusProses = "Baru"
	layanan.StatusValidasi = "Menunggu Validasi"
	layanan.CreatedAt = time.Now()

	if err := DB.Create(&layanan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, layanan)
}

func GetAllLayananAdministrasi(c *gin.Context) {
	 log.Println("--- 1. MASUK KE HANDLER GetAllLayananAdministrasi ---")

    var layanans []LayananAdministrasi
    log.Println("--- 2. Variabel 'layanans' berhasil dibuat ---")

    // Kita coba query TANPA Preload dulu untuk menyederhanakan
    tx := DB.Find(&layanans)
    log.Println("--- 3. Perintah DB.Find sudah dieksekusi ---")

    if tx.Error != nil {
        log.Println("!!! ERROR SAAT QUERY DATABASE:", tx.Error.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
        return
    }

    log.Println("--- 4. Query database BERHASIL. Jumlah data:", tx.RowsAffected, "---")

    c.JSON(http.StatusOK, layanans)
    log.Println("--- 5. Response JSON berhasil dikirim ---")
}

func GetLayananAdministrasiByID(c *gin.Context) {
	layananID := c.Param("id")
	userClaims, _ := c.Get("user")
	claims := userClaims.(*Claims)

	var layanan LayananAdministrasi
	if err := DB.First(&layanan, layananID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data layanan tidak ditemukan"})
		return
	}

	// Jika user adalah 'opd', cek apakah dia pemilik data ini
	if claims.Role == "opd" && layanan.UserOPDID != claims.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki hak akses untuk melihat data ini"})
		return
	}

	// Jika rolenya 'pemda' atau dia adalah pemilik, lanjutkan
	DB.Preload("UserOPD.OPD").First(&layanan, layananID)
	c.JSON(http.StatusOK, layanan)
}

func GetLayananAdministrasiByUserOPD(c *gin.Context) {
	userOpdID := c.Param("id")
	var layanans []LayananAdministrasi
	
	err := DB.Where("user_opd_id = ?", userOpdID).Preload("UserOPD.OPD").Find(&layanans).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, layanans)
}

func UpdateLayananAdministrasi(c *gin.Context) {
	id := c.Param("id")
	var layanan LayananAdministrasi

	if err := DB.First(&layanan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	// Bind data dari form ke struct yang sudah ada
	if err := bindLayananFromMultipartForm(c, &layanan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := DB.Save(&layanan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data"})
		return
	}

	c.JSON(http.StatusOK, layanan)
}

func DeleteLayananAdministrasi(c *gin.Context) {
	id := c.Param("id")
	if err := DB.Delete(&LayananAdministrasi{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}

func ValidateLayananAdministrasi(c *gin.Context) {
	// 1. Ambil ID layanan dari parameter URL
	id := c.Param("id")

	// 2. Deklarasikan variabel 'layanan'
	var layanan LayananAdministrasi

	// 3. Cari layanan di database
	if err := DB.First(&layanan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Layanan Pembangunan tidak ditemukan"})
		return
	}

	// 4. Cek apakah sudah divalidasi sebelumnya
	if layanan.StatusValidasi != "Menunggu Validasi" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Layanan ini sudah divalidasi"})
		return
	}

	// 5. Bind request body JSON
	var req ValidasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body tidak valid"})
		return
	}

	// 6. Ambil data user 'pemda' yang sedang login dari context
	userClaims, _ := c.Get("user")
	claims := userClaims.(*Claims)

	// 7. Siapkan data untuk update
	validatorIDFromToken := claims.ID // Ambil ID dari token
	now := time.Now()

	// 8. Update field di struct 'layanan'
	layanan.StatusValidasi = req.StatusValidasi
	layanan.KeteranganValidasi = req.KeteranganValidasi
	layanan.IDValidatorPemda = &validatorIDFromToken
	layanan.TanggalValidasi = &now

	// 9. Simpan perubahan ke database
	DB.Save(&layanan)
	
	// 10. Ambil data terbaru dengan relasi untuk dikirim kembali
	DB.Preload("ValidatorPemda").First(&layanan, layanan.ID)
	c.JSON(http.StatusOK, layanan)
}
// ========= CRUD: LAYANAN INFORMASI & PENGADUAN =========

func CreateLayananInformasiPengaduan(c *gin.Context) {
// Ambil data user yang login dari context (sudah di-set oleh AuthMiddleware)
	userClaims, _ := c.Get("user")
	claims := userClaims.(*Claims)

	// UserOPDID sekarang diambil dari token, bukan form. JAUH LEBIH AMAN.
	var layanan = LayananInformasiPengaduan{
		UserOPDID: claims.ID, 
	}
	
	// Bind sisa data dari form ke struct 'layanan'
	if err := bindLayananFromMultipartForm(c, &layanan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set nilai default saat pembuatan
	layanan.StatusProses = "Baru"
	layanan.StatusValidasi = "Menunggu Validasi"
	layanan.CreatedAt = time.Now()

	if err := DB.Create(&layanan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, layanan)
}

func GetAllLayananInformasiPengaduan(c *gin.Context) {
 log.Println("--- 1. MASUK KE HANDLER GetAllLayananInformasiPengaduan ---")

    var layanans []LayananInformasiPengaduan
    log.Println("--- 2. Variabel 'layanans' berhasil dibuat ---")

    // Kita coba query TANPA Preload dulu untuk menyederhanakan
    tx := DB.Find(&layanans)
    log.Println("--- 3. Perintah DB.Find sudah dieksekusi ---")

    if tx.Error != nil {
        log.Println("!!! ERROR SAAT QUERY DATABASE:", tx.Error.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
        return
    }

    log.Println("--- 4. Query database BERHASIL. Jumlah data:", tx.RowsAffected, "---")

    c.JSON(http.StatusOK, layanans)
    log.Println("--- 5. Response JSON berhasil dikirim ---")
}

func GetLayananInformasiPengaduanByID(c *gin.Context) {
	id := c.Param("id")
	var layanan LayananInformasiPengaduan
	
	if err := DB.Preload("UserOPD.OPD").First(&layanan, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, layanan)
}

func GetLayananInformasiPengaduanByUserOPD(c *gin.Context) {
	userOpdID := c.Param("id")
	var layanans []LayananInformasiPengaduan
	
	err := DB.Where("user_opd_id = ?", userOpdID).Preload("UserOPD.OPD").Find(&layanans).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, layanans)
}

func UpdateLayananInformasiPengaduan(c *gin.Context) {
	id := c.Param("id")
	var layanan LayananInformasiPengaduan

	if err := DB.First(&layanan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	// Bind data dari form ke struct yang sudah ada
	if err := bindLayananFromMultipartForm(c, &layanan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := DB.Save(&layanan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data"})
		return
	}

	c.JSON(http.StatusOK, layanan)
}

func DeleteLayananInformasiPengaduan(c *gin.Context) {
	id := c.Param("id")
	if err := DB.Delete(&LayananInformasiPengaduan{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}

func ValidateLayananInformasiPengaduan(c *gin.Context) {
	// 1. Ambil ID layanan dari parameter URL
	id := c.Param("id")

	// 2. Deklarasikan variabel 'layanan'
	var layanan LayananInformasiPengaduan

	// 3. Cari layanan di database
	if err := DB.First(&layanan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Layanan Pembangunan tidak ditemukan"})
		return
	}

	// 4. Cek apakah sudah divalidasi sebelumnya
	if layanan.StatusValidasi != "Menunggu Validasi" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Layanan ini sudah divalidasi"})
		return
	}

	// 5. Bind request body JSON
	var req ValidasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body tidak valid"})
		return
	}

	// 6. Ambil data user 'pemda' yang sedang login dari context
	userClaims, _ := c.Get("user")
	claims := userClaims.(*Claims)

	// 7. Siapkan data untuk update
	validatorIDFromToken := claims.ID // Ambil ID dari token
	now := time.Now()

	// 8. Update field di struct 'layanan'
	layanan.StatusValidasi = req.StatusValidasi
	layanan.KeteranganValidasi = req.KeteranganValidasi
	layanan.IDValidatorPemda = &validatorIDFromToken
	layanan.TanggalValidasi = &now

	// 9. Simpan perubahan ke database
	DB.Save(&layanan)
	
	// 10. Ambil data terbaru dengan relasi untuk dikirim kembali
	DB.Preload("ValidatorPemda").First(&layanan, layanan.ID)
	c.JSON(http.StatusOK, layanan)
}


//
// CATATAN: Fungsi Delete sengaja tidak diimplementasikan untuk menjaga integritas data.
// Biasanya, data layanan tidak dihapus secara fisik (hard delete), melainkan ditandai
// sebagai "dibatalkan" atau "tidak aktif" (soft delete).
//

