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

// ========= CRUD HANDLERS: JENIS PELAYANAN (STANDAR PELAYANAN) =========

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

	// Set default status validasi
	standar.StatusValidasi = "Menunggu Validasi"

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
	// Preload OPD dan ValidatorPemda agar informasi ikut terambil
	if err := DB.Preload("OPD").Preload("ValidatorPemda").Find(&standar).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, standar)
}

// GetStandarPelayananByOPD: Mendapatkan semua standar pelayanan berdasarkan ID OPD
func GetStandarPelayananByOPD(c *gin.Context) {
	idOpdStr := c.Param("id_opd")
	if idOpdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID OPD tidak boleh kosong"})
		return
	}

	var standarPelayanan []JenisPelayanan

	err := DB.Where("id_opd = ?", idOpdStr).
		Preload("OPD").
		Preload("ValidatorPemda").
		Find(&standarPelayanan).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data dari database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": standarPelayanan})
}

// ValidateJenisPelayanan: Memvalidasi standar pelayanan (Hanya oleh User Pemda)
func ValidateJenisPelayanan(c *gin.Context) {
	id := c.Param("id")
	var standar JenisPelayanan

	// 1. Cari standar
	if err := DB.First(&standar, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Standar Pelayanan tidak ditemukan"})
		return
	}

	// 2. Cek status
	if standar.StatusValidasi != "Menunggu Validasi" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Standar ini sudah divalidasi sebelumnya"})
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
	standar.StatusValidasi = req.StatusValidasi

	if req.KeteranganValidasi != "" {
		standar.KeteranganValidasi = &req.KeteranganValidasi
	} else {
		standar.KeteranganValidasi = nil
	}

	standar.IDValidatorPemda = &validatorIDFromToken
	standar.TanggalValidasi = &now

	// 6. Simpan
	DB.Save(&standar)

	// 7. Response
	DB.Preload("ValidatorPemda").Preload("OPD").First(&standar, standar.ID)
	c.JSON(http.StatusOK, standar)
}




// ========= CRUD HANDLERS: FORM PEMOHON (MASTER DATA) =========


// CreateFormPemohon: Membuat data master pemohon (oleh User OPD)
func CreateFormPemohon(c *gin.Context) {
	var form FormPemohon
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body tidak valid"})
		return
	}

	// Ambil ID User OPD dan IDOPD dari token
	userClaims, _ := c.Get("user")
	claims := userClaims.(*Claims)
	form.IDUserOPDInput = claims.ID // Set petugas yang menginput
	form.IDOPD = claims.IDOPD 	 // <-- PERUBAHAN: Set OPD tempat mendaftar

	if err := DB.Create(&form).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Preload relasi UserOPDInput dan OPD baru
	DB.Preload("UserOPDInput").Preload("OPD").First(&form, form.ID)
	c.JSON(http.StatusCreated, form)
}

// GetAllFormPemohon: Mendapatkan semua data master pemohon
func GetAllFormPemohon(c *gin.Context) {
	var forms []FormPemohon

	// Tambahkan Preload("OPD")
	if err := DB.Preload("UserOPDInput").Preload("OPD").Find(&forms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, forms)
}

// GetFormPemohonByID: Mendapatkan detail satu pemohon
func GetFormPemohonByID(c *gin.Context) {
	id := c.Param("id")
	var form FormPemohon

	// Tambahkan Preload("OPD")
	if err := DB.Preload("UserOPDInput").Preload("OPD").First(&form, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data pemohon tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, form)
}

// UpdateFormPemohon: Memperbarui data master pemohon
func UpdateFormPemohon(c *gin.Context) {
	id := c.Param("id")
	var form FormPemohon

	if err := DB.First(&form, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data pemohon tidak ditemukan"})
		return
	}

	// Otorisasi: (Opsional) Cek apakah user yang mengedit adalah yang menginput
	// ... (logika otorisasi Anda tetap sama)

	var input FormPemohon
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body tidak valid"})
		return
	}

	// Update data
	// IDUserOPDInput dan IDOPD tidak diubah, tetap pencatat awal
	form.NamaLengkap = input.NamaLengkap
	form.NIK = input.NIK
	form.Alamat = input.Alamat
	form.NomorHP = input.NomorHP
	form.Email = input.Email
	// form.IDOPD = input.IDOPD // Sebaiknya jangan diupdate, biarkan tetap

	if err := DB.Save(&form).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data"})
		return
	}

	// Tambahkan Preload("OPD")
	DB.Preload("UserOPDInput").Preload("OPD").First(&form, form.ID)
	c.JSON(http.StatusOK, form)
}

// DeleteFormPemohon: Menghapus data master pemohon
func DeleteFormPemohon(c *gin.Context) {
	id := c.Param("id")
	if err := DB.Delete(&FormPemohon{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data pemohon berhasil dihapus"})
}


// ========= HELPER UNTUK AMBIL FILE & FORM VALUE (FORM PENGAJUAN) =========

// BindFormPengajuanFromMultipartForm diperbarui untuk ERD V8
func BindFormPengajuanFromMultipartForm(c *gin.Context, form *FormPengajuan) error {
	// Parsing form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		return errors.New("Gagal parsing form: " + err.Error())
	}

	// Handle file upload
	var uploadedFilename *string
	_, handler, err := c.Request.FormFile("dokumen_pengajuan") // <-- Nama field file baru
	if err == nil {
		os.MkdirAll("./uploads", os.ModePerm)
		// Tambahkan timestamp unik untuk menghindari nama file yang sama
		filePath := "./uploads/" + strconv.FormatInt(time.Now().Unix(), 10) + "_" + handler.Filename
		if err := c.SaveUploadedFile(handler, filePath); err != nil {
			return errors.New("Gagal menyimpan file: " + err.Error())
		}
		uploadedFilename = &filePath
	} else if form.DokumenPengajuanPath != nil {
		// Jika tidak ada file baru diupload, pertahankan path yang lama (untuk update)
		uploadedFilename = form.DokumenPengajuanPath
	}

	// --- BIND FIELD WAJIB ---

	// Konversi ID OPD (Baru)
	idOpdStr := c.PostForm("id_opd")
	idOpd, err := strconv.ParseUint(idOpdStr, 10, 64)
	if err != nil {
		return errors.New("ID OPD tidak valid atau kosong")
	}
	form.IDOPD = uint(idOpd)

	// Konversi ID Jenis Pelayanan
	idJenisStr := c.PostForm("id_jenis_pelayanan")
	idJenis, err := strconv.ParseUint(idJenisStr, 10, 64)
	if err != nil {
		return errors.New("ID Jenis Pelayanan tidak valid atau kosong")
	}
	form.IDJenisPelayanan = uint(idJenis)

	// Bind Data Pemohon (Baru)
	form.NamaPemohonLengkap = c.PostForm("nama_pemohon_lengkap")
	form.NIKPemohon = c.PostForm("nik_pemohon")
	form.JudulPengajuan = c.PostForm("judul_pengajuan")

	// Bind Checkbox
	isAgreedStr := c.PostForm("is_agreed")
	isAgreed, err := strconv.ParseBool(isAgreedStr)
	if err != nil {
		form.IsAgreed = false // Default jika parsing gagal
	}
	form.IsAgreed = isAgreed

	// --- BIND FIELD OPSIONAL (NULLABLE/TEXT/VARCHAR) ---
	form.AlamatPemohon = c.PostForm("alamat_pemohon")
	form.NomorHPPemohon = c.PostForm("nomor_hp_pemohon")
	form.EmailPemohon = c.PostForm("email_pemohon")
	form.DeskripsiSingkat = c.PostForm("deskripsi_singkat")

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

	form.DokumenPengajuanPath = uploadedFilename // <-- Nama field struct baru

	// Validasi Sederhana
	if form.NamaPemohonLengkap == "" || form.NIKPemohon == "" || form.JudulPengajuan == "" {
		return errors.New("nama Pemohon, NIK Pemohon, dan Judul Pengajuan tidak boleh kosong")
	}

	return nil
}

// =========== FORM PENGAJUAN (TRANSAKSI) KHUSUS OPD =================

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

	// StatusProses sudah memiliki default di models.go
	form.StatusProses = "Baru"
	// StatusValidasi DIHAPUS

	if err := DB.Create(&form).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data: " + err.Error()})
		return
	}

	// Ambil kembali data dengan relasi untuk response
	DB.Preload("UserOPD.OPD").Preload("JenisPelayanan.OPD").Preload("OPD").First(&form, form.ID)
	c.JSON(http.StatusCreated, form)
}

// GetAllFormPengajuan: Mendapatkan semua pengajuan (untuk Pemda/Admin)
func GetAllFormPengajuan(c *gin.Context) {
	var forms []FormPengajuan

	tx := DB.Preload("UserOPD.OPD").
		Preload("JenisPelayanan.OPD").
		Preload("OPD"). // Preload OPD yang dituju
		Find(&forms)

	if tx.Error != nil {
		log.Println("!!! ERROR SAAT QUERY DATABASE:", tx.Error.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": tx.Error.Error()})
		return
	}

	log.Println("--- Query GetAllFormPengajuan BERHASIL. Jumlah data:", tx.RowsAffected, "---")
	c.JSON(http.StatusOK, forms)
}

// GetFormPengajuanByUserOPD: Mendapatkan pengajuan berdasarkan user OPD ID
func GetFormPengajuanByUserOPD(c *gin.Context) {
	userOpdID := c.Param("id")
	var forms []FormPengajuan

	err := DB.Where("id_user_opd = ?", userOpdID).
		Preload("UserOPD.OPD").
		Preload("JenisPelayanan.OPD").
		Preload("OPD").
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

	// Preload semua relasi (ValidatorPemda DIHAPUS)
	if err := DB.Preload("UserOPD.OPD").Preload("JenisPelayanan.OPD").Preload("OPD").First(&form, formID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data pengajuan tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Otorisasi: (Opsional) Cek apakah user OPD yang login adalah pemilik data
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

	// 2. Cek status (Opsional: Cek StatusProses)
	if form.StatusProses == "Selesai" {
		 c.JSON(http.StatusForbidden, gin.H{"error": "Data ini sudah Selesai dan tidak dapat diubah"})
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
	DB.Preload("UserOPD.OPD").Preload("JenisPelayanan.OPD").Preload("OPD").First(&form, form.ID)
	c.JSON(http.StatusOK, form)
}

// DeleteFormPengajuan: Menghapus data pengajuan
func DeleteFormPengajuan(c *gin.Context) {
	id := c.Param("id")
	// Di dunia nyata, Anda mungkin ingin memeriksa status sebelum menghapus
	if err := DB.Delete(&FormPengajuan{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}

// ValidateFormPengajuan: DIHAPUS KARENA PINDAH KE JENIS PELAYANAN
// func ValidateFormPengajuan(c *gin.Context) { ... }