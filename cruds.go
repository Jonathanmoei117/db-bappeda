package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
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

func CreateJenisPelayanan(c *gin.Context) {
	var standar JenisPelayanan
	if err := c.ShouldBindJSON(&standar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body tidak valid"})
		return
	}

	// Otorisasi sederhana: Pastikan IDOPD valid
	var opd OPD
	if err := DB.First(&opd, standar.IDOPD).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID OPD tidak valid"})
		return
	}
	
	if err := DB.Create(&standar).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Ambil data OPD untuk response
	DB.Preload("OPD").First(&standar, standar.ID)
	c.JSON(http.StatusCreated, standar)
}

func GetAllJenisPelayanan(c *gin.Context) {
	var standar []JenisPelayanan
	// Preload OPD agar informasi OPD ikut terambil
	if err := DB.Preload("OPD").Find(&standar).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, standar)
}


// ========= HELPER UNTUK AMBIL FILE & FORM VALUE =========

func BindFormPengajuanFromMultipartForm(c *gin.Context, form *FormPengajuan) error {
	// Parsing form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		return errors.New("Gagal parsing form: " + err.Error())
	}

	// Handle file upload
	var uploadedFilename *string
	_, handler, err := c.Request.FormFile("berkas_pengajuan")
	if err == nil { 
		os.MkdirAll("./uploads", os.ModePerm)
		filePath := "./uploads/" + handler.Filename
		if err := c.SaveUploadedFile(handler, filePath); err != nil {
			return errors.New("Gagal menyimpan file: " + err.Error())
		}
		// Set path file yang baru diupload
		uploadedFilename = &filePath 
	} else if form.BerkasPengajuanPath != nil {
		// Jika tidak ada file baru diupload, pertahankan path yang lama (untuk update)
		uploadedFilename = form.BerkasPengajuanPath 
	}

	// --- BIND FIELD WAJIB ---
	
	// Konversi ID Jenis Pelayanan (Wajib diisi)
	idJenisStr := c.PostForm("id_jenis_pelayanan")
	idJenis, err := strconv.ParseUint(idJenisStr, 10, 64)
	if err != nil {
		return errors.New("ID Jenis Pelayanan tidak valid atau kosong")
	}
	form.IDJenisPelayanan = uint(idJenis)
	
	form.JenisLayanan = c.PostForm("jenis_layanan")
	form.JudulKegiatan = c.PostForm("judul_kegiatan")
	form.Deskripsi = c.PostForm("deskripsi")
	form.NamaPemohon = c.PostForm("nama_pemohon")
	form.InstansiPemohon = c.PostForm("instansi_pemohon")

	// --- BIND FIELD OPSIONAL (NULLABLE) ---

	if nip := c.PostForm("nip_pemohon"); nip != "" {
		form.NIPPemohon = &nip
	} else {
		form.NIPPemohon = nil
	}

	// Bind Tanggal (Periode Mulai)
	if mulaiStr := c.PostForm("periode_mulai"); mulaiStr != "" {
		t, err := time.Parse("2006-01-02", mulaiStr)
		if err == nil {
			form.PeriodeMulai = &t
		}
	} else {
		form.PeriodeMulai = nil
	}

	// Bind Tanggal (Periode Selesai)
	if selesaiStr := c.PostForm("periode_selesai"); selesaiStr != "" {
		t, err := time.Parse("2006-01-02", selesaiStr)
		if err == nil {
			form.PeriodeSelesai = &t
		}
	} else {
		form.PeriodeSelesai = nil
	}
	
	form.BerkasPengajuanPath = uploadedFilename

	return nil
}

// =========== FORM PENGAJUAN KHUSUS OPD =================

func CreateFormPengajuan(c *gin.Context) {
	// Ambil ID User OPD dari token
	userClaims, _ := c.Get("user")
	claims := userClaims.(*Claims)

	var form FormPengajuan
	form.IDUserOPD = claims.ID // Ambil ID dari user yang login
	
	if err := BindFormPengajuanFromMultipartForm(c, &form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// StatusProses & StatusValidasi sudah memiliki default di models.go, 
	// tapi kita set lagi untuk kepastian
	form.StatusProses = "Baru"
	form.StatusValidasi = "Menunggu Validasi"

	if err := DB.Create(&form).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data: " + err.Error()})
		return
	}

	// Ambil kembali data dengan relasi untuk response
	DB.Preload("UserOPD.OPD").Preload("JenisPelayanan.OPD").First(&form, form.ID)
	c.JSON(http.StatusCreated, form)
}

// GetAllFormPengajuan: Mendapatkan semua pengajuan (untuk Pemda/Admin)
func GetAllFormPengajuan(c *gin.Context) {
	log.Println("--- 1. MASUK KE HANDLER GetAllFormPengajuan ---")

    var forms []FormPengajuan
    log.Println("--- 2. Variabel 'forms' berhasil dibuat ---")

    // Kita coba query TANPA Preload dulu untuk menyederhanakan
    tx := DB.Find(&forms)
    log.Println("--- 3. Perintah DB.Find sudah dieksekusi ---")

    if tx.Error != nil {
        log.Println("!!! ERROR SAAT QUERY DATABASE:", tx.Error.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
        return
    }

    log.Println("--- 4. Query database BERHASIL. Jumlah data:", tx.RowsAffected, "---")

    c.JSON(http.StatusOK, forms)
    log.Println("--- 5. Response JSON berhasil dikirim ---")
}

// GetFormPengajuanByUserOPD: Mendapatkan pengajuan berdasarkan user OPD ID (untuk User Pemda atau User OPD itu sendiri)
func GetFormPengajuanByUserOPD(c *gin.Context) {
    userOpdID := c.Param("id")
    var forms []FormPengajuan
    
    // Cari Form Pengajuan berdasarkan IDUserOPD
    err := DB.Where("id_user_opd = ?", userOpdID).
		Preload("UserOPD.OPD").
		Preload("JenisPelayanan.OPD").
		Find(&forms).Error
		
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, forms)
}

// GetFormPengajuanByID: Mendapatkan detail pengajuan
func GetFormPengajuanByID(c *gin.Context) {
	formID := c.Param("id")
	var form FormPengajuan

	// Preload semua relasi
	if err := DB.Preload("UserOPD.OPD").Preload("JenisPelayanan.OPD").Preload("ValidatorPemda").First(&form, formID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data pengajuan tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Otorisasi: Cek apakah user OPD yang login adalah pemilik data
	// (Jika API ini dilindungi middleware Otorisasi)
	userClaims, exists := c.Get("user")
	if exists {
		claims := userClaims.(*Claims)
		// Jika role-nya OPD dan bukan pemilik data, tolak
		if claims.Role == "opd" && form.IDUserOPD != claims.ID {
			 c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki hak akses untuk melihat data ini"})
			 return
		}
	}

	c.JSON(http.StatusOK, form)
}

// UpdateFormPengajuan: Memperbarui data pengajuan (hanya oleh User OPD)
func UpdateFormPengajuan(c *gin.Context) {
	id := c.Param("id")
	var form FormPengajuan

	// 1. Cari data lama
	if err := DB.First(&form, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	// 2. Cek status validasi: Tidak boleh diubah jika sudah Disetujui/Ditolak
	if form.StatusValidasi != "Menunggu Validasi" && form.StatusValidasi != "Baru" {
		 c.JSON(http.StatusForbidden, gin.H{"error": "Data ini sudah divalidasi dan tidak dapat diubah"})
		 return
	}

	// 3. Bind data baru dari form ke struct lama
	if err := BindFormPengajuanFromMultipartForm(c, &form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 4. Simpan perubahan
	if err := DB.Save(&form).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data"})
		return
	}

	// 5. Response
	DB.Preload("UserOPD.OPD").Preload("JenisPelayanan.OPD").First(&form, form.ID)
	c.JSON(http.StatusOK, form)
}

// DeleteFormPengajuan: Menghapus data pengajuan
func DeleteFormPengajuan(c *gin.Context) {
	id := c.Param("id")
	// Biasanya ini adalah Soft Delete atau hanya diizinkan jika statusnya "Baru" atau "Ditolak"
	if err := DB.Delete(&FormPengajuan{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}

// ValidateFormPengajuan: Memvalidasi pengajuan (Hanya oleh User Pemda)
func ValidateFormPengajuan(c *gin.Context) {
	id := c.Param("id")
	var form FormPengajuan

	// 1. Cari form
	if err := DB.First(&form, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Form Pengajuan tidak ditemukan"})
		return
	}

	// 2. Cek status
	if form.StatusValidasi != "Menunggu Validasi" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Form ini sudah divalidasi sebelumnya"})
		return
	}

	// 3. Bind request body validasi
	var req ValidasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body tidak valid"})
		return
	}

	// 4. Ambil ID Validator dari token
	userClaims, _ := c.Get("user")
	claims := userClaims.(*Claims)
	validatorIDFromToken := claims.ID 
	now := time.Now()

	// 5. Update field
	form.StatusValidasi = req.StatusValidasi
	
	// KeteranganValidasi adalah *string (pointer)
	if req.KeteranganValidasi != "" {
		form.KeteranganValidasi = &req.KeteranganValidasi
	} else {
		form.KeteranganValidasi = nil
	}
	
	form.IDValidatorPemda = &validatorIDFromToken
	form.TanggalValidasi = &now

	// 6. Simpan
	DB.Save(&form)
	
	// 7. Response
	DB.Preload("ValidatorPemda").Preload("UserOPD.OPD").Preload("JenisPelayanan.OPD").First(&form, form.ID)
	c.JSON(http.StatusOK, form)
}