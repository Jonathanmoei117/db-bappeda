package main

import (
	"encoding/json"
	"fmt"
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


// ========= CRUD HANDLERS: LAYANAN PEMBANGUNAN =========
// (Semua fungsi CRUD untuk LayananPembangunan, termasuk ValidateLayananPembangunan, diletakkan di sini)

func GetAllLayananPembangunan(w http.ResponseWriter, r *http.Request) {
	var layanans []LayananPembangunan
	if err := DB.Preload("UserOPD").Preload("ValidatorPemda").Order("created_at desc").Find(&layanans).Error; err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseJSON(w, http.StatusOK, layanans)
}

func GetLayananPembangunanByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var layanan LayananPembangunan
	if err := DB.Preload("UserOPD").Preload("ValidatorPemda").First(&layanan, id).Error; err != nil {
		responseError(w, http.StatusNotFound, "Layanan Pembangunan tidak ditemukan")
		return
	}
	responseJSON(w, http.StatusOK, layanan)
}

func CreateLayananPembangunan(w http.ResponseWriter, r *http.Request) {
	var layanan LayananPembangunan
	if err := json.NewDecoder(r.Body).Decode(&layanan); err != nil {
		responseError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}
	layanan.UserOPDID = 1 

	if err := DB.Create(&layanan).Error; err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseJSON(w, http.StatusCreated, layanan)
}

func UpdateLayananPembangunan(w http.ResponseWriter, r *http.Request) {
    // Simulasi: Dapatkan info user yang sedang login dari token
    userIDFromToken := uint(1) // Anggap yang login adalah UserOPD dengan ID 1
    userRoleFromToken := "user_opd"

    id, _ := strconv.Atoi(mux.Vars(r)["id"])
    var layanan LayananPembangunan
    if err := DB.First(&layanan, id).Error; err != nil {
        responseError(w, http.StatusNotFound, "Layanan Pembangunan tidak ditemukan")
        return
    }

    if userRoleFromToken == "user_opd" && layanan.UserOPDID != userIDFromToken {
        responseError(w, http.StatusForbidden, "Akses ditolak: Anda bukan pemilik data ini.")
        return
    }
    if layanan.StatusValidasi == "Disetujui" {
        responseError(w, http.StatusForbidden, "Akses ditolak: Data yang sudah disetujui tidak dapat diubah.")
        return
    }

    var input LayananPembangunan
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        responseError(w, http.StatusBadRequest, "Request body tidak valid")
        return
    }
    if err := DB.Model(&layanan).Updates(input).Error; err != nil {
        responseError(w, http.StatusInternalServerError, err.Error())
        return
    }
    responseJSON(w, http.StatusOK, layanan)
}

func DeleteLayananPembangunan(w http.ResponseWriter, r *http.Request) {
    // Simulasi: Dapatkan info user yang sedang login dari token
    userRoleFromToken := "user_pemda" // Hanya pemda yang bisa hapus

    if userRoleFromToken != "user_pemda" {
        responseError(w, http.StatusForbidden, "Akses ditolak: Anda tidak memiliki izin untuk menghapus data.")
        return
    }

    id, _ := strconv.Atoi(mux.Vars(r)["id"])
    if res := DB.Delete(&LayananPembangunan{}, id); res.RowsAffected == 0 {
        responseError(w, http.StatusNotFound, "Layanan Pembangunan tidak ditemukan")
        return
    }
    responseJSON(w, http.StatusNoContent, nil)
}

func ValidateLayananPembangunan(w http.ResponseWriter, r *http.Request) {
    // Simulasi: Dapatkan info user PEMDA yang sedang login dari token
    validatorIDFromToken := uint(1) // Anggap yang login adalah UserPemda dengan ID 1
    userRoleFromToken := "user_pemda"

    // 1. Cek Peran: Hanya user_pemda yang boleh validasi
    if userRoleFromToken != "user_pemda" {
        responseError(w, http.StatusForbidden, "Akses ditolak: Hanya peran Pemda yang bisa melakukan validasi.")
        return
    }

    // 2. Ambil data layanan dari database
    id, _ := strconv.Atoi(mux.Vars(r)["id"])
    var layanan LayananPembangunan
    if err := DB.First(&layanan, id).Error; err != nil {
        responseError(w, http.StatusNotFound, "Layanan Pembangunan tidak ditemukan")
        return
    }

    // 3. Cek apakah sudah pernah divalidasi
    if layanan.StatusValidasi != "Menunggu Validasi" {
        responseError(w, http.StatusBadRequest, "Data ini sudah divalidasi sebelumnya.")
        return
    }

    // 4. Baca input dari body request
    var input ValidasiRequest
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        responseError(w, http.StatusBadRequest, "Request body tidak valid")
        return
    }
    
    // 5. Pastikan input status valid
    if input.StatusValidasi != "Disetujui" && input.StatusValidasi != "Ditolak" {
        responseError(w, http.StatusBadRequest, "Status validasi harus 'Disetujui' atau 'Ditolak'")
        return
    }

    // 6. Update data di database
    now := time.Now()
    layanan.StatusValidasi = input.StatusValidasi
    layanan.KeteranganValidasi = input.KeteranganValidasi
    layanan.IDValidatorPemda = &validatorIDFromToken // Gunakan pointer
    layanan.TanggalValidasi = &now                  // Gunakan pointer

    if err := DB.Save(&layanan).Error; err != nil {
        responseError(w, http.StatusInternalServerError, err.Error())
        return
    }

    // 7. Jika disetujui, data "masuk ke laporan" (bisa trigger aksi lain)
    if layanan.StatusValidasi == "Disetujui" {
        fmt.Println("LOG: Data Pembangunan ID", layanan.ID, "disetujui dan siap untuk dilaporkan.")
    }
    
    responseJSON(w, http.StatusOK, layanan)
}

// ========= CRUD HANDLERS: LAYANAN ADMINISTRASI =========
// (Implementasi lengkap CRUD dan validasi untuk LayananAdministrasi, polanya sama persis dengan LayananPembangunan)


func GetAllLayananAdministrasi(w http.ResponseWriter, r *http.Request) {
	var layanans []LayananAdministrasi
	if err := DB.Preload("UserOPD").Preload("ValidatorPemda").Order("created_at desc").Find(&layanans).Error; err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseJSON(w, http.StatusOK, layanans)
}

func GetLayananAdministrasiByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var layanan LayananAdministrasi
	if err := DB.Preload("UserOPD").Preload("ValidatorPemda").First(&layanan, id).Error; err != nil {
		responseError(w, http.StatusNotFound, "Layanan Administrasi tidak ditemukan")
		return
	}
	responseJSON(w, http.StatusOK, layanan)
}

func CreateLayananAdministrasi(w http.ResponseWriter, r *http.Request) {
	var layanan LayananAdministrasi
	if err := json.NewDecoder(r.Body).Decode(&layanan); err != nil {
		responseError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}
	layanan.UserOPDID = 1

	if err := DB.Create(&layanan).Error; err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseJSON(w, http.StatusCreated, layanan)
}

func UpdateLayananAdministrasi(w http.ResponseWriter, r *http.Request) {
	// Simulasi: Dapatkan info user yang sedang login dari token
    userIDFromToken := uint(1) // Anggap yang login adalah UserOPD dengan ID 1
    userRoleFromToken := "user_opd"

    id, _ := strconv.Atoi(mux.Vars(r)["id"])
    var layanan LayananAdministrasi
    if err := DB.First(&layanan, id).Error; err != nil {
        responseError(w, http.StatusNotFound, "Layanan Administrasi tidak ditemukan")
        return
    }

    if userRoleFromToken == "user_opd" && layanan.UserOPDID != userIDFromToken {
        responseError(w, http.StatusForbidden, "Akses ditolak: Anda bukan pemilik data ini.")
        return
    }
    if layanan.StatusValidasi == "Disetujui" {
        responseError(w, http.StatusForbidden, "Akses ditolak: Data yang sudah disetujui tidak dapat diubah.")
        return
    }

    var input LayananAdministrasi
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        responseError(w, http.StatusBadRequest, "Request body tidak valid")
        return
    }
    if err := DB.Model(&layanan).Updates(input).Error; err != nil {
        responseError(w, http.StatusInternalServerError, err.Error())
        return
    }
    responseJSON(w, http.StatusOK, layanan)
}

func DeleteLayananAdministrasi(w http.ResponseWriter, r *http.Request) {
	 // Simulasi: Dapatkan info user yang sedang login dari token
    userRoleFromToken := "user_pemda" // Hanya pemda yang bisa hapus

    if userRoleFromToken != "user_pemda" {
        responseError(w, http.StatusForbidden, "Akses ditolak: Anda tidak memiliki izin untuk menghapus data.")
        return
    }

    id, _ := strconv.Atoi(mux.Vars(r)["id"])
    if res := DB.Delete(&LayananAdministrasi{}, id); res.RowsAffected == 0 {
        responseError(w, http.StatusNotFound, "Layanan Administrasi tidak ditemukan")
        return
    }
    responseJSON(w, http.StatusNoContent, nil)
}

func ValidateLayananAdministrasi(w http.ResponseWriter, r *http.Request) {
    // Logika sama persis, hanya beda tipe struct
    validatorIDFromToken := uint(1)
    userRoleFromToken := "user_pemda"

    if userRoleFromToken != "user_pemda" {
        responseError(w, http.StatusForbidden, "Akses ditolak: Hanya peran Pemda yang bisa melakukan validasi.")
        return
    }

    id, _ := strconv.Atoi(mux.Vars(r)["id"])
    var layanan LayananAdministrasi
    if err := DB.First(&layanan, id).Error; err != nil {
        responseError(w, http.StatusNotFound, "Layanan Administrasi tidak ditemukan")
        return
    }

    if layanan.StatusValidasi != "Menunggu Validasi" {
        responseError(w, http.StatusBadRequest, "Data ini sudah divalidasi sebelumnya.")
        return
    }

    var input ValidasiRequest
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        responseError(w, http.StatusBadRequest, "Request body tidak valid")
        return
    }
    
    if input.StatusValidasi != "Disetujui" && input.StatusValidasi != "Ditolak" {
        responseError(w, http.StatusBadRequest, "Status validasi harus 'Disetujui' atau 'Ditolak'")
        return
    }

    now := time.Now()
    layanan.StatusValidasi = input.StatusValidasi
    layanan.KeteranganValidasi = input.KeteranganValidasi
    layanan.IDValidatorPemda = &validatorIDFromToken
    layanan.TanggalValidasi = &now

    if err := DB.Save(&layanan).Error; err != nil {
        responseError(w, http.StatusInternalServerError, err.Error())
        return
    }
    
    responseJSON(w, http.StatusOK, layanan)
}


// ========= CRUD HANDLERS: LAYANAN INFORMASI & PENGADUAN =========
// (Implementasi lengkap CRUD dan validasi untuk LayananInformasiPengaduan, polanya sama persis)

func GetAllLayananInformasiPengaduan(w http.ResponseWriter, r *http.Request) {
	var layanans []LayananInformasiPengaduan
	if err := DB.Preload("UserOPD").Preload("ValidatorPemda").Order("created_at desc").Find(&layanans).Error; err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseJSON(w, http.StatusOK, layanans)
}

func GetLayananInformasiPengaduanByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var layanan LayananInformasiPengaduan
	if err := DB.Preload("UserOPD").Preload("ValidatorPemda").First(&layanan, id).Error; err != nil {
		responseError(w, http.StatusNotFound, "Data tidak ditemukan")
		return
	}
	responseJSON(w, http.StatusOK, layanan)
}

func CreateLayananInformasiPengaduan(w http.ResponseWriter, r *http.Request) {
	var layanan LayananInformasiPengaduan
	if err := json.NewDecoder(r.Body).Decode(&layanan); err != nil {
		responseError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}
	layanan.UserOPDID = 1

	if err := DB.Create(&layanan).Error; err != nil {
		responseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseJSON(w, http.StatusCreated, layanan)
}

func UpdateLayananInformasiPengaduan(w http.ResponseWriter, r *http.Request) {
	// Simulasi: Dapatkan info user yang sedang login dari token
    userIDFromToken := uint(1) // Anggap yang login adalah UserOPD dengan ID 1
    userRoleFromToken := "user_opd"

    id, _ := strconv.Atoi(mux.Vars(r)["id"])
    var layanan LayananInformasiPengaduan
    if err := DB.First(&layanan, id).Error; err != nil {
        responseError(w, http.StatusNotFound, "Layanan Administrasi tidak ditemukan")
        return
    }

    if userRoleFromToken == "user_opd" && layanan.UserOPDID != userIDFromToken {
        responseError(w, http.StatusForbidden, "Akses ditolak: Anda bukan pemilik data ini.")
        return
    }
    if layanan.StatusValidasi == "Disetujui" {
        responseError(w, http.StatusForbidden, "Akses ditolak: Data yang sudah disetujui tidak dapat diubah.")
        return
    }

    var input LayananInformasiPengaduan
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        responseError(w, http.StatusBadRequest, "Request body tidak valid")
        return
    }
    if err := DB.Model(&layanan).Updates(input).Error; err != nil {
        responseError(w, http.StatusInternalServerError, err.Error())
        return
    }
    responseJSON(w, http.StatusOK, layanan)
}

func DeleteLayananInformasiPengaduan(w http.ResponseWriter, r *http.Request) {
	 // Simulasi: Dapatkan info user yang sedang login dari token
    userRoleFromToken := "user_pemda" // Hanya pemda yang bisa hapus

    if userRoleFromToken != "user_pemda" {
        responseError(w, http.StatusForbidden, "Akses ditolak: Anda tidak memiliki izin untuk menghapus data.")
        return
    }

    id, _ := strconv.Atoi(mux.Vars(r)["id"])
    if res := DB.Delete(&LayananInformasiPengaduan{}, id); res.RowsAffected == 0 {
        responseError(w, http.StatusNotFound, "Layanan Administrasi tidak ditemukan")
        return
    }
    responseJSON(w, http.StatusNoContent, nil)
}

func ValidateLayananInformasiPengaduan(w http.ResponseWriter, r *http.Request) {
    // Logika sama persis, hanya beda tipe struct
    validatorIDFromToken := uint(1)
    userRoleFromToken := "user_pemda"

    if userRoleFromToken != "user_pemda" {
        responseError(w, http.StatusForbidden, "Akses ditolak: Hanya peran Pemda yang bisa melakukan validasi.")
        return
    }

    id, _ := strconv.Atoi(mux.Vars(r)["id"])
    var layanan LayananInformasiPengaduan
    if err := DB.First(&layanan, id).Error; err != nil {
        responseError(w, http.StatusNotFound, "Data tidak ditemukan")
        return
    }

    if layanan.StatusValidasi != "Menunggu Validasi" {
        responseError(w, http.StatusBadRequest, "Data ini sudah divalidasi sebelumnya.")
        return
    }

    var input ValidasiRequest
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        responseError(w, http.StatusBadRequest, "Request body tidak valid")
        return
    }
    
    if input.StatusValidasi != "Disetujui" && input.StatusValidasi != "Ditolak" {
        responseError(w, http.StatusBadRequest, "Status validasi harus 'Disetujui' atau 'Ditolak'")
        return
    }

    now := time.Now()
    layanan.StatusValidasi = input.StatusValidasi
    layanan.KeteranganValidasi = input.KeteranganValidasi
    layanan.IDValidatorPemda = &validatorIDFromToken
    layanan.TanggalValidasi = &now

    if err := DB.Save(&layanan).Error; err != nil {
        responseError(w, http.StatusInternalServerError, err.Error())
        return
    }
    
    responseJSON(w, http.StatusOK, layanan)
}


