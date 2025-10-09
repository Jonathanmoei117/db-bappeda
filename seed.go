package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Catatan: Asumsi semua struct model (OPD, UserOPD, dll.) ada di package 'main'
// dan variabel koneksi GORM 'DB' sudah dideklarasikan.
// Implementasi FirstOrCreate di sini mengasumsikan GORM v2.

// Seed akan mengisi database dengan data awal untuk keperluan development.
func Seed() {
	fmt.Println("===== MEMULAI PROSES SEEDING DATA =====")

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
	// LANGKAH 2: Buat Master Data (OPD, User, Pemda)
	// ==================================================================

	// --- OPD (2) ---
	opdBappeda := OPD{NamaOPD: "BAPPEDA", AlamatOPD: "Jl. Pahlawan No. 1"}
	DB.FirstOrCreate(&opdBappeda, OPD{NamaOPD: "BAPPEDA"})

	opdDisdik := OPD{NamaOPD: "Dinas Pendidikan", AlamatOPD: "Jl. Pendidikan No. 2"}
	DB.FirstOrCreate(&opdDisdik, OPD{NamaOPD: "Dinas Pendidikan"})

	// --- User Pemda (1) ---
	userPemda := UserPemda{
		Nama:     "Dr. Anisa Wijayanti",
		Password: string(hashedPasswordPemda),
		NIP:      "198505052015012002",
		Jabatan:  "Kepala Bidang Verifikasi",
	}
	DB.FirstOrCreate(&userPemda, UserPemda{NIP: "198505052015012002"})

	// --- User OPD (2) ---
	userBappeda := UserOPD{
		IDOPD:    opdBappeda.ID,
		Nama:     "Budi Santoso",
		Password: string(hashedPasswordOPD),
		NIP:      "199001012020121001",
		Jabatan:  "Staf Perencanaan",
	}
	DB.FirstOrCreate(&userBappeda, UserOPD{NIP: "199001012020121001"})

	userDisdik := UserOPD{
		IDOPD:    opdDisdik.ID,
		Nama:     "Citra Lestari",
		Password: string(hashedPasswordOPD),
		NIP:      "199203152021012003",
		Jabatan:  "Staf Kurikulum",
	}
	DB.FirstOrCreate(&userDisdik, UserOPD{NIP: "199203152021012003"})

	log.Println("🌱 Seeding Master Data (OPD, User, Pemda) selesai!")

	// ==================================================================
	// LANGKAH 3: Buat Jenis Pelayanan (STANDAR) - 2 per OPD
	// ==================================================================

	// --- Standar milik OPD BAPPEDA (2) ---
	standarBappeda1 := JenisPelayanan{
		IDOPD:           opdBappeda.ID,
		NamaStandar:     "Rekomendasi Izin Prinsip Pembangunan",
		WaktuPelayanan:  "14 Hari Kerja",
		BiayaTarif:      "Rp 0 (Sesuai Perda)",
		ProdukPelayanan: "Surat Rekomendasi",
		Fasilitas:       "Aplikasi Simpel & Loket Pelayanan",
		JangkaWaktu:     "14 Hari Kerja",
	}
	DB.FirstOrCreate(&standarBappeda1, JenisPelayanan{IDOPD: opdBappeda.ID, NamaStandar: standarBappeda1.NamaStandar})

	standarBappeda2 := JenisPelayanan{
		IDOPD:           opdBappeda.ID,
		NamaStandar:     "Permohonan Data Spasial & Peta RTRW",
		WaktuPelayanan:  "5 Hari Kerja",
		BiayaTarif:      "Gratis",
		ProdukPelayanan: "File Peta Digital (SHP, PDF)",
		Fasilitas:       "Email & Google Drive",
		JangkaWaktu:     "5 Hari Kerja",
	}
	DB.FirstOrCreate(&standarBappeda2, JenisPelayanan{IDOPD: opdBappeda.ID, NamaStandar: standarBappeda2.NamaStandar})

	// --- Standar milik OPD Dinas Pendidikan (2) ---
	standarDisdik1 := JenisPelayanan{
		IDOPD:           opdDisdik.ID,
		NamaStandar:     "Legalisasi Ijazah & Verifikasi Keabsahan",
		WaktuPelayanan:  "3 Hari Kerja",
		BiayaTarif:      "Rp 10.000,- / Dokumen",
		ProdukPelayanan: "Ijazah Terlegalisir",
		Fasilitas:       "Loket Pelayanan Disdik",
		JangkaWaktu:     "3 Hari Kerja",
	}
	DB.FirstOrCreate(&standarDisdik1, JenisPelayanan{IDOPD: opdDisdik.ID, NamaStandar: standarDisdik1.NamaStandar})

	standarDisdik2 := JenisPelayanan{
		IDOPD:           opdDisdik.ID,
		NamaStandar:     "Pendaftaran Bantuan Siswa Miskin (BSM)",
		WaktuPelayanan:  "Sesuai Jadwal Periode",
		BiayaTarif:      "Gratis",
		ProdukPelayanan: "SK Penerima Bantuan",
		Fasilitas:       "Formulir Online & Sekolah",
		JangkaWaktu:     "1 Semester",
	}
	DB.FirstOrCreate(&standarDisdik2, JenisPelayanan{IDOPD: opdDisdik.ID, NamaStandar: standarDisdik2.NamaStandar})

	log.Println("📃 Seeding Jenis Pelayanan (Standar) selesai! (Total 4 Standar)")

	// ==================================================================
	// LANGKAH 4: Buat Form Pengajuan (TRANSAKSI) - 2 per OPD
	// ==================================================================

	validatedAt := time.Now().Add(-48 * time.Hour)
	keteranganSetuju := "Dokumen pendukung lengkap dan sah. Disetujui untuk diproses."
	keteranganTolak := "Data siswa tidak terdaftar di Dapodik. Pendaftaran ditolak."
	nipPemohon := "199508172020121005"

	// --- Form Pengajuan untuk BAPPEDA (2) ---
	formBappeda1 := FormPengajuan{
		IDUserOPD:        userBappeda.ID,
		IDJenisPelayanan: standarBappeda1.ID, // Menggunakan standar 1 Bappeda
		JudulKegiatan:    "Pengajuan Izin Pembangunan Sekolah Dasar Swasta",
		Deskripsi:        "Permohonan Izin Prinsip untuk Yayasan Cerdas Bangsa.",
		NamaPemohon:      "Yayasan Cerdas Bangsa",
		InstansiPemohon:  "Swasta",
		StatusProses:     "Baru",
		StatusValidasi:   "Menunggu Validasi",
	}
	DB.FirstOrCreate(&formBappeda1, FormPengajuan{JudulKegiatan: formBappeda1.JudulKegiatan})

	formBappeda2 := FormPengajuan{
		IDUserOPD:          userBappeda.ID,
		IDJenisPelayanan:   standarBappeda2.ID, // Menggunakan standar 2 Bappeda
		IDValidatorPemda:   &userPemda.ID,
		TanggalValidasi:    &validatedAt,
		StatusValidasi:     "Disetujui",
		KeteranganValidasi: &keteranganSetuju,
		JudulKegiatan:      "Request Peta RTRW Kecamatan Kartoharjo",
		Deskripsi:          "Data peta untuk keperluan analisis penelitian mahasiswa.",
		NamaPemohon:        "Ahmad Dahlan",
		InstansiPemohon:    "Universitas Merdeka Madiun",
		StatusProses:       "Selesai",
	}
	DB.FirstOrCreate(&formBappeda2, FormPengajuan{JudulKegiatan: formBappeda2.JudulKegiatan})

	// --- Form Pengajuan untuk Dinas Pendidikan (2) ---
	formDisdik1 := FormPengajuan{
		IDUserOPD:        userDisdik.ID,
		IDJenisPelayanan: standarDisdik1.ID, // Menggunakan standar 1 Disdik
		JudulKegiatan:    "Legalisasi Ijazah untuk Keperluan CPNS",
		Deskripsi:        "Membutuhkan 5 lembar legalisir ijazah untuk pendaftaran.",
		NamaPemohon:      "Rina Amelia",
		NIPPemohon:       &nipPemohon,
		InstansiPemohon:  "Pribadi",
		StatusProses:     "Diproses",
		StatusValidasi:   "Disetujui",
		IDValidatorPemda: &userPemda.ID,
		TanggalValidasi:  &validatedAt,
	}
	DB.FirstOrCreate(&formDisdik1, FormPengajuan{JudulKegiatan: formDisdik1.JudulKegiatan})

	formDisdik2 := FormPengajuan{
		IDUserOPD:          userDisdik.ID,
		IDJenisPelayanan:   standarDisdik2.ID, // Menggunakan standar 2 Disdik
		IDValidatorPemda:   &userPemda.ID,
		TanggalValidasi:    &validatedAt,
		StatusValidasi:     "Ditolak",
		KeteranganValidasi: &keteranganTolak,
		JudulKegiatan:      "Pendaftaran BSM untuk Siswa SDN 2 Kejuron",
		Deskripsi:          "Mengajukan bantuan untuk 15 siswa sesuai data terlampir.",
		NamaPemohon:        "Kepala Sekolah SDN 2 Kejuron",
		InstansiPemohon:    "SDN 2 Kejuron",
		StatusProses:       "Revisi",
	}
	DB.FirstOrCreate(&formDisdik2, FormPengajuan{JudulKegiatan: formDisdik2.JudulKegiatan})

	log.Println("📝 Seeding Form Pengajuan (Transaksi) selesai! (Total 4 Form)")
	fmt.Println("===== PROSES SEEDING SELESAI =====")
}