package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// ========= HELPER FUNCTIONS =========

func responseJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func responseError(w http.ResponseWriter, code int, message string) {
	responseJSON(w, code, map[string]string{"error": message})
}

// ========= CRUD HANDLERS: OPD (KHUSUS SUPER ADMIN) =========

// CreateOPD digunakan oleh Super Admin untuk mendaftarkan OPD baru ke sistem.
func CreateOPD(w http.ResponseWriter, r *http.Request) {
	var opd OPD
	if err := json.NewDecoder(r.Body).Decode(&opd); err != nil {
		responseError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}

	if err := DB.Create(&opd).Error; err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseJSON(w, http.StatusCreated, opd)
}
// GetAllOPD digunakan untuk mengambil daftar OPD, misalnya untuk dropdown form.
func GetAllOPD(w http.ResponseWriter, r *http.Request) {
	var opds []OPD
	if err := DB.Find(&opds).Error; err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseJSON(w, http.StatusOK, opds)
}

// ========= HANDLERS REGISTRASI PENGGUNA (OPD & PEMDA) =========

// CreateUserOPD menangani registrasi user OPD. NIP digunakan sebagai pengganti username.
func CreateUserOPD(w http.ResponseWriter, r *http.Request) {
	var user UserOPD
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		responseError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}

	// Hashing password harus dilakukan di aplikasi nyata
	// user.Password = hashPassword(user.Password)

	if err := DB.Create(&user).Error; err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Kosongkan password sebelum mengirim response
	user.Password = ""
	responseJSON(w, http.StatusCreated, user)
}

// CreateUserPemda menangani registrasi user Pemda. NIP digunakan sebagai pengganti username.
func CreateUserPemda(w http.ResponseWriter, r *http.Request) {
	var user UserPemda
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		responseError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}

	// Hashing password harus dilakukan di aplikasi nyata
	// user.Password = hashPassword(user.Password)

	if err := DB.Create(&user).Error; err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Kosongkan password sebelum mengirim response
	user.Password = ""
	responseJSON(w, http.StatusCreated, user)
}

// ========= CRUD: LAYANAN PEMBANGUNAN =========

func CreateLayananPembangunan(w http.ResponseWriter, r *http.Request) {
	var layanan LayananPembangunan
	if err := json.NewDecoder(r.Body).Decode(&layanan); err != nil {
		responseError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}
	// Simulasi: ID user didapat dari token otentikasi
	userOPDIDFromToken := uint(1)
	layanan.UserOPDID = userOPDIDFromToken

	if err := DB.Create(&layanan).Error; err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	DB.Preload("UserOPD.OPD").First(&layanan, layanan.ID)
	responseJSON(w, http.StatusCreated, layanan)
}

func GetAllLayananPembangunan(w http.ResponseWriter, r *http.Request) {
	var layanans []LayananPembangunan
	query := DB.Preload("UserOPD.OPD").Preload("ValidatorPemda").Order("created_at desc")

	if opdID := r.URL.Query().Get("opd_id"); opdID != "" {
		query = query.Joins("JOIN user_opds ON user_opds.id = layanan_pembangunans.user_opd_id").
			Where("user_opds.opd_id = ?", opdID)
	}

	if err := query.Find(&layanans).Error; err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseJSON(w, http.StatusOK, layanans)
}

func GetLayananPembangunanByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var layanan LayananPembangunan
	err := DB.Preload("UserOPD.OPD").Preload("ValidatorPemda").First(&layanan, id).Error
	if err != nil {
		responseError(w, http.StatusNotFound, "Layanan Pembangunan tidak ditemukan")
		return
	}
	responseJSON(w, http.StatusOK, layanan)
}

func UpdateLayananPembangunan(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var layanan LayananPembangunan
	if err := DB.First(&layanan, id).Error; err != nil {
		responseError(w, http.StatusNotFound, "Layanan Pembangunan tidak ditemukan")
		return
	}

	// Cek otorisasi & status
	userOPDIDFromToken := uint(1) // Simulasi dari token
	if layanan.UserOPDID != userOPDIDFromToken {
		responseError(w, http.StatusForbidden, "Akses ditolak")
		return
	}
	if layanan.StatusValidasi == "Disetujui" {
		responseError(w, http.StatusBadRequest, "Data yang sudah disetujui tidak dapat diubah")
		return
	}

	var input LayananPembangunan
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responseError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}
	DB.Model(&layanan).Updates(input)
	responseJSON(w, http.StatusOK, layanan)
}

func ValidateLayananPembangunan(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var layanan LayananPembangunan
	if err := DB.First(&layanan, id).Error; err != nil {
		responseError(w, http.StatusNotFound, "Layanan Pembangunan tidak ditemukan")
		return
	}

	if layanan.StatusValidasi != "Menunggu Validasi" {
		responseError(w, http.StatusBadRequest, "Layanan ini sudah divalidasi")
		return
	}

	var req ValidasiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}

	validatorIDFromToken := uint(1) // Simulasi dari token
	now := time.Now()

	layanan.StatusValidasi = req.StatusValidasi
	layanan.KeteranganValidasi = req.KeteranganValidasi
	layanan.IDValidatorPemda = &validatorIDFromToken
	layanan.TanggalValidasi = &now

	DB.Save(&layanan)
	DB.Preload("ValidatorPemda").First(&layanan, layanan.ID)
	responseJSON(w, http.StatusOK, layanan)
}

// ========= CRUD: LAYANAN ADMINISTRASI =========

func CreateLayananAdministrasi(w http.ResponseWriter, r *http.Request) {
	var layanan LayananAdministrasi
	if err := json.NewDecoder(r.Body).Decode(&layanan); err != nil {
		responseError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}
	userOPDIDFromToken := uint(1)
	layanan.UserOPDID = userOPDIDFromToken

	if err := DB.Create(&layanan).Error; err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	DB.Preload("UserOPD.OPD").First(&layanan, layanan.ID)
	responseJSON(w, http.StatusCreated, layanan)
}

func GetAllLayananAdministrasi(w http.ResponseWriter, r *http.Request) {
	var layanans []LayananAdministrasi
	query := DB.Preload("UserOPD.OPD").Preload("ValidatorPemda").Order("created_at desc")

	if opdID := r.URL.Query().Get("opd_id"); opdID != "" {
		query = query.Joins("JOIN user_opds ON user_opds.id = layanan_administrasis.user_opd_id").
			Where("user_opds.opd_id = ?", opdID)
	}

	if err := query.Find(&layanans).Error; err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseJSON(w, http.StatusOK, layanans)
}

func GetLayananAdministrasiByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var layanan LayananAdministrasi
	err := DB.Preload("UserOPD.OPD").Preload("ValidatorPemda").First(&layanan, id).Error
	if err != nil {
		responseError(w, http.StatusNotFound, "Layanan Administrasi tidak ditemukan")
		return
	}
	responseJSON(w, http.StatusOK, layanan)
}

func UpdateLayananAdministrasi(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var layanan LayananAdministrasi
	if err := DB.First(&layanan, id).Error; err != nil {
		responseError(w, http.StatusNotFound, "Layanan Administrasi tidak ditemukan")
		return
	}

	userOPDIDFromToken := uint(1)
	if layanan.UserOPDID != userOPDIDFromToken {
		responseError(w, http.StatusForbidden, "Akses ditolak")
		return
	}
	if layanan.StatusValidasi == "Disetujui" {
		responseError(w, http.StatusBadRequest, "Data yang sudah disetujui tidak dapat diubah")
		return
	}

	var input LayananAdministrasi
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responseError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}
	DB.Model(&layanan).Updates(input)
	responseJSON(w, http.StatusOK, layanan)
}

func ValidateLayananAdministrasi(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var layanan LayananAdministrasi
	if err := DB.First(&layanan, id).Error; err != nil {
		responseError(w, http.StatusNotFound, "Layanan Administrasi tidak ditemukan")
		return
	}

	if layanan.StatusValidasi != "Menunggu Validasi" {
		responseError(w, http.StatusBadRequest, "Layanan ini sudah divalidasi")
		return
	}

	var req ValidasiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}

	validatorIDFromToken := uint(1)
	now := time.Now()

	layanan.StatusValidasi = req.StatusValidasi
	layanan.KeteranganValidasi = req.KeteranganValidasi
	layanan.IDValidatorPemda = &validatorIDFromToken
	layanan.TanggalValidasi = &now

	DB.Save(&layanan)
	DB.Preload("ValidatorPemda").First(&layanan, layanan.ID)
	responseJSON(w, http.StatusOK, layanan)
}

// ========= CRUD: LAYANAN INFORMASI & PENGADUAN =========

func CreateLayananInformasi(w http.ResponseWriter, r *http.Request) {
	var layanan LayananInformasiPengaduan
	if err := json.NewDecoder(r.Body).Decode(&layanan); err != nil {
		responseError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}
	userOPDIDFromToken := uint(1)
	layanan.UserOPDID = userOPDIDFromToken

	if err := DB.Create(&layanan).Error; err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	DB.Preload("UserOPD.OPD").First(&layanan, layanan.ID)
	responseJSON(w, http.StatusCreated, layanan)
}

func GetAllLayananInformasi(w http.ResponseWriter, r *http.Request) {
	var layanans []LayananInformasiPengaduan
	query := DB.Preload("UserOPD.OPD").Preload("ValidatorPemda").Order("created_at desc")

	if opdID := r.URL.Query().Get("opd_id"); opdID != "" {
		query = query.Joins("JOIN user_opds ON user_opds.id = layanan_informasi_pengaduans.user_opd_id").
			Where("user_opds.opd_id = ?", opdID)
	}

	if err := query.Find(&layanans).Error; err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseJSON(w, http.StatusOK, layanans)
}

func GetLayananInformasiByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var layanan LayananInformasiPengaduan
	err := DB.Preload("UserOPD.OPD").Preload("ValidatorPemda").First(&layanan, id).Error
	if err != nil {
		responseError(w, http.StatusNotFound, "Layanan Informasi & Pengaduan tidak ditemukan")
		return
	}
	responseJSON(w, http.StatusOK, layanan)
}

func UpdateLayananInformasi(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var layanan LayananInformasiPengaduan
	if err := DB.First(&layanan, id).Error; err != nil {
		responseError(w, http.StatusNotFound, "Layanan Informasi & Pengaduan tidak ditemukan")
		return
	}

	userOPDIDFromToken := uint(1)
	if layanan.UserOPDID != userOPDIDFromToken {
		responseError(w, http.StatusForbidden, "Akses ditolak")
		return
	}
	if layanan.StatusValidasi == "Disetujui" {
		responseError(w, http.StatusBadRequest, "Data yang sudah disetujui tidak dapat diubah")
		return
	}

	var input LayananInformasiPengaduan
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responseError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}
	DB.Model(&layanan).Updates(input)
	responseJSON(w, http.StatusOK, layanan)
}

func ValidateLayananInformasi(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var layanan LayananInformasiPengaduan
	if err := DB.First(&layanan, id).Error; err != nil {
		responseError(w, http.StatusNotFound, "Layanan Informasi & Pengaduan tidak ditemukan")
		return
	}

	if layanan.StatusValidasi != "Menunggu Validasi" {
		responseError(w, http.StatusBadRequest, "Layanan ini sudah divalidasi")
		return
	}

	var req ValidasiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}

	validatorIDFromToken := uint(1)
	now := time.Now()

	layanan.StatusValidasi = req.StatusValidasi
	layanan.KeteranganValidasi = req.KeteranganValidasi
	layanan.IDValidatorPemda = &validatorIDFromToken
	layanan.TanggalValidasi = &now

	DB.Save(&layanan)
	DB.Preload("ValidatorPemda").First(&layanan, layanan.ID)
	responseJSON(w, http.StatusOK, layanan)
}

//
// CATATAN: Fungsi Delete sengaja tidak diimplementasikan untuk menjaga integritas data.
// Biasanya, data layanan tidak dihapus secara fisik (hard delete), melainkan ditandai
// sebagai "dibatalkan" atau "tidak aktif" (soft delete).
//

