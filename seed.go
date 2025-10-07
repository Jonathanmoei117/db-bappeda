package main

import (
	"fmt"
	"time"
	"log"

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
		Nama: 	 "Dr. Anisa Wijayanti",
		Password: string(hashedPasswordPemda), 
		NIP: 	 "198505052015012002",
		Jabatan: "Kepala Bidang Verifikasi",
	}
	DB.FirstOrCreate(&userPemda, UserPemda{NIP: "198505052015012002"})

	// --- User OPD (2) ---
	userBappeda := UserOPD{
		IDOPD: 	 opdBappeda.ID,
		Nama: 	 "Budi Santoso",
		Password: string(hashedPasswordOPD), 
		NIP: 	 "199001012020121001",
		Jabatan: "Staf Perencanaan",
	}
	DB.FirstOrCreate(&userBappeda, UserOPD{NIP: "199001012020121001"})

	userDisdik := UserOPD{
		IDOPD: 	 opdDisdik.ID,
		Nama: 	 "Citra Lestari",
		Password: string(hashedPasswordOPD), 
		NIP: 	 "199203152021012003",
		Jabatan: "Staf Kurikulum",
	}
	DB.FirstOrCreate(&userDisdik, UserOPD{NIP: "199203152021012003"})

	log.Println("🌱 Seeding Master Data (OPD, User, Pemda) selesai!")

	// ==================================================================
	// LANGKAH 3: Buat Jenis Pelayanan (STANDAR)
	// ==================================================================
	
	// Standar milik OPD BAPPEDA
	standarBappeda := JenisPelayanan{
		IDOPD: 			opdBappeda.ID,
		NamaStandar: 	"Rekomendasi Izin Prinsip Pembangunan",
		WaktuPelayanan: "14 Hari Kerja",
		BiayaTarif: 	"Rp 0 (Sesuai Perda)",
		ProdukPelayanan: "Surat Rekomendasi",
		SaranaPrasarana: "Aplikasi Simpel",
	}
	DB.FirstOrCreate(&standarBappeda, JenisPelayanan{NamaStandar: "Rekomendasi Izin Prinsip Pembangunan"})

	// Standar milik OPD Dinas Pendidikan
	standarDisdik := JenisPelayanan{
		IDOPD: 			opdDisdik.ID,
		NamaStandar: 	"Legalisasi Ijazah & Verifikasi Keabsahan",
		WaktuPelayanan: "3 Hari Kerja",
		BiayaTarif: 	"Rp 10.000,- / Dokumen",
		ProdukPelayanan: "Ijazah Terlegalisir",
		SaranaPrasarana: "Loket Pelayanan Disdik",
	}
	DB.FirstOrCreate(&standarDisdik, JenisPelayanan{NamaStandar: "Legalisasi Ijazah & Verifikasi Keabsahan"})
	
	log.Println("📃 Seeding Jenis Pelayanan (Standar) selesai!")

	// ==================================================================
	// LANGKAH 4: Buat Form Pengajuan (TRANSAKSI)
	// ==================================================================
	
	validatedAt := time.Now().Add(-48 * time.Hour)
	keteranganSetuju := "Dokumen pendukung lengkap dan sah. Disetujui untuk diproses."
	keteranganTolak := "Izin Prinsip tidak relevan, diajukan kembali ke Dinas PUPR."
    
    // VARIABEL NIP DUMMY SEBAGAI STRING (Wajib menggunakan variabel dan & untuk pointer)
    nipPemohon2 := "3301010120000001" 
    nipPemohon3 := "3301010120000002"

	
	// --- Form 1: Menunggu Validasi (OPD BAPPEDA menggunakan Standar BAPPEDA) ---
	form1 := FormPengajuan{
		IDUserOPD: userBappeda.ID,
		IDJenisPelayanan: standarBappeda.ID,
		JenisLayanan: "Pembangunan",
		JudulKegiatan: "Pengajuan Izin Pembangunan Sekolah Dasar Swasta",
		Deskripsi: "Permohonan Izin Prinsip untuk Yayasan Cerdas Bangsa.",
		NamaPemohon: "Yayasan Cerdas Bangsa",
		InstansiPemohon: "Swasta",
		StatusProses: "Baru",
		StatusValidasi: "Menunggu Validasi",
	}
	DB.FirstOrCreate(&form1, FormPengajuan{JudulKegiatan: form1.JudulKegiatan})

	// --- Form 2: Sudah Divalidasi (OPD DISDIK menggunakan Standar DISDIK) ---
	form2 := FormPengajuan{
		IDUserOPD: userDisdik.ID,
		IDJenisPelayanan: standarDisdik.ID,
		IDValidatorPemda: &userPemda.ID, 
		TanggalValidasi: &validatedAt,
		StatusValidasi: "Disetujui",
		KeteranganValidasi: &keteranganSetuju,
		
		JenisLayanan: "Administrasi",
		JudulKegiatan: "Legalisasi Ijazah untuk Keperluan CPNS",
		Deskripsi: "Membutuhkan 5 lembar legalisir ijazah untuk pendaftaran.",
		NamaPemohon: "Andi Wijaya",
		NIPPemohon: &nipPemohon2, // PERBAIKAN: Gunakan alamat dari string
		InstansiPemohon: "Pribadi",
		StatusProses: "Selesai",
	}
	DB.FirstOrCreate(&form2, FormPengajuan{JudulKegiatan: form2.JudulKegiatan})
	
	// --- Form 3: Ditolak (OPD BAPPEDA menggunakan Standar BAPPEDA) ---
	form3 := FormPengajuan{
		IDUserOPD: userBappeda.ID,
		IDJenisPelayanan: standarBappeda.ID,
		IDValidatorPemda: &userPemda.ID,
		TanggalValidasi: &validatedAt,
		StatusValidasi: "Ditolak",
		KeteranganValidasi: &keteranganTolak,
		
		JenisLayanan: "Pembangunan",
		JudulKegiatan: "Kajian Pembangunan Pusat Komersial Baru",
		Deskripsi: "Mengkaji dampak ekonomi pembangunan mall baru di pusat kota.",
		NamaPemohon: "PT. Maju Bersama",
		NIPPemohon: &nipPemohon3, // PERBAIKAN: Gunakan alamat dari string
		InstansiPemohon: "Swasta",
		StatusProses: "Revisi",
	}
	DB.FirstOrCreate(&form3, FormPengajuan{JudulKegiatan: form3.JudulKegiatan})


	log.Println("📝 Seeding Form Pengajuan (Transaksi) selesai! (Total 3 Form)")
	fmt.Println("===== PROSES SEEDING SELESAI =====")
}