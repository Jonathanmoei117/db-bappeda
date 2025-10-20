package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Catatan: Asumsi semua struct model (OPD, UserOPD, dll.) ada di package 'main'
// dan variabel koneksi GORM 'DB' sudah dideklarasikan.
// Implementasi FirstOrCreate di sini mengasumsikan GORM v2.

// Seed akan mengisi database dengan data awal untuk keperluan development.
func Seed() {
	fmt.Println("===== MEMULAI PROSES SEEDING DATA (V8) =====")

	// ==================================================================
	// LANGKAH 1: Persiapan Password Hash
	// ==================================================================
	passwordOPD := "opd123"
	hashedPasswordOPD, err := bcrypt.GenerateFromPassword([]byte(passwordOPD), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Gagal hash password OPD: %v", err)
	}

	passwordPemda := "pemda123"
	hashedPasswordPemda, err := bcrypt.GenerateFromPassword([]byte(passwordPemda), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Gagal hash password Pemda: %v", err)
	}

	// ==================================================================
	// LANGKAH 2: Buat Master Data (OPD, User Pemda, User OPD)
	// ==================================================================

	// --- OPD (3) ---
	opdBappeda := OPD{NamaOPD: "BAPPEDA", AlamatOPD: "Jl. Pahlawan No. 1"}
	DB.FirstOrCreate(&opdBappeda, OPD{NamaOPD: "BAPPEDA"})

	opdDinkes := OPD{NamaOPD: "Dinas Kesehatan", AlamatOPD: "Jl. Dr. Soetomo No. 5"}
	DB.FirstOrCreate(&opdDinkes, OPD{NamaOPD: "Dinas Kesehatan"})

	opdPerkim := OPD{NamaOPD: "Dinas PERKIM", AlamatOPD: "Jl. Mastrip No. 10"}
	DB.FirstOrCreate(&opdPerkim, OPD{NamaOPD: "Dinas PERKIM"})

	// --- User Pemda (1) ---
	userPemda := UserPemda{
		Nama:     "Dr. Anisa Wijayanti",
		Password: string(hashedPasswordPemda),
		NIP:      "198505052015012002",
		Jabatan:  "Kepala Bidang Verifikasi",
	}
	DB.FirstOrCreate(&userPemda, UserPemda{NIP: "198505052015012002"})

	// --- User OPD (3) ---
	userBappeda := UserOPD{
		IDOPD:    opdBappeda.ID,
		Nama:     "Budi Santoso",
		Password: string(hashedPasswordOPD),
		NIP:      "199001012020121001",
		Jabatan:  "Staf Perencanaan",
	}
	DB.FirstOrCreate(&userBappeda, UserOPD{NIP: "199001012020121001"})

	userDinkes := UserOPD{
		IDOPD:    opdDinkes.ID,
		Nama:     "Citra Lestari",
		Password: string(hashedPasswordOPD),
		NIP:      "199203152021012003",
		Jabatan:  "Staf Administrasi Kesehatan",
	}
	DB.FirstOrCreate(&userDinkes, UserOPD{NIP: "199203152021012003"})

	userPerkim := UserOPD{
		IDOPD:    opdPerkim.ID,
		Nama:     "Ahmad Sahroni",
		Password: string(hashedPasswordOPD),
		NIP:      "199407102022021005",
		Jabatan:  "Staf Pendataan PSU",
	}
	DB.FirstOrCreate(&userPerkim, UserOPD{NIP: "199407102022021005"})

	log.Println("🌱 Seeding Master Data (3 OPD, 1 Pemda, 3 User OPD) selesai!")

	// ==================================================================
	// LANGKAH 3: Buat Jenis Pelayanan (STANDAR) - 1 per OPD, Belum Divalidasi
	// ==================================================================

	// --- Standar milik OPD BAPPEDA (1) ---
	standarBappeda1 := JenisPelayanan{
		IDOPD:           opdBappeda.ID,
		NamaStandar:     "Rekomendasi Izin Prinsip Pembangunan",
		DasarHukum:      "Perda No. 5 Tahun 2020 tentang RTRW",
		Persyaratan:     "1. Fotokopi KTP\n2. Fotokopi Sertifikat Tanah\n3. Proposal Rencana Pembangunan",
		WaktuPelayanan:  "14 Hari Kerja",
		BiayaTarif:      "Rp 0 (Sesuai Perda)",
		ProdukPelayanan: "Surat Rekomendasi Izin Prinsip",
		StatusValidasi:  "Menunggu Validasi", // Sesuai permintaan
	}
	DB.FirstOrCreate(&standarBappeda1, JenisPelayanan{NamaStandar: standarBappeda1.NamaStandar})

	// --- Standar milik OPD Dinas Kesehatan (1) ---
	standarDinkes1 := JenisPelayanan{
		IDOPD:           opdDinkes.ID,
		NamaStandar:     "Penerbitan Surat Izin Praktik (SIP) Dokter",
		DasarHukum:      "UU No. 29 Tahun 2004 tentang Praktik Kedokteran",
		Persyaratan:     "1. Fotokopi KTP\n2. Pas Foto 4x6\n3. Surat Tanda Registrasi (STR)",
		WaktuPelayanan:  "7 Hari Kerja",
		BiayaTarif:      "Gratis",
		ProdukPelayanan: "Surat Izin Praktik (SIP) Dokter",
		StatusValidasi:  "Menunggu Validasi", // Sesuai permintaan
	}
	DB.FirstOrCreate(&standarDinkes1, JenisPelayanan{NamaStandar: standarDinkes1.NamaStandar})

	// --- Standar milik OPD Dinas PERKIM (1) ---
	standarPerkim1 := JenisPelayanan{
		IDOPD:           opdPerkim.ID,
		NamaStandar:     "Permohonan Bantuan Prasarana, Sarana, dan Utilitas (PSU) Perumahan",
		DasarHukum:      "Permen PUPR No. 03/PRT/M/2018",
		Persyaratan:     "1. Proposal dari Pengembang\n2. Site Plan yang Disetujui\n3. Data Calon Penerima Manfaat",
		WaktuPelayanan:  "30 Hari Kerja (Verifikasi Lapangan)",
		BiayaTarif:      "Gratis",
		ProdukPelayanan: "SK Penetapan Penerima Bantuan PSU",
		StatusValidasi:  "Menunggu Validasi", // Sesuai permintaan
	}
	DB.FirstOrCreate(&standarPerkim1, JenisPelayanan{NamaStandar: standarPerkim1.NamaStandar})

	log.Println("📃 Seeding Jenis Pelayanan (Standar) selesai! (Total 3 Standar, status 'Menunggu Validasi')")

	// ==================================================================
	// LANGKAH 4: Form Pemohon (Dikosongkan)
	// ==================================================================
	log.Println("⏩ Seeding Form Pemohon (Master) dilompati sesuai permintaan.")

	// ==================================================================
	// LANGKAH 5: Form Pengajuan (Dikosongkan)
	// ==================================================================
	log.Println("⏩ Seeding Form Pengajuan (Transaksi) dilompati sesuai permintaan.")

	fmt.Println("===== PROSES SEEDING SELESAI =====")
}