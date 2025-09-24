package main

import (
	"encoding/base64" // Import ini sekarang akan digunakan
	"fmt"
	"time"

	"golang.org/x/crypto/argon2"
)

func Seed() {
	fmt.Println("===== SEEDING DATA =====")

	// SALT
	saltOPD := []byte("salt-argon2-untuk-opd-123")
	saltPemda := []byte("salt-argon2-untuk-pemda-456")

	// HASH PASSWORD
	encodedHashOPD := base64.StdEncoding.EncodeToString(
		argon2.IDKey([]byte("password123"), saltOPD, 1, 64*1024, 4, 32),
	)
	encodedHashPemda := base64.StdEncoding.EncodeToString(
		argon2.IDKey([]byte("password456"), saltPemda, 1, 64*1024, 4, 32),
	)

	// ================== USER ==================
	userOPD := UserOPD{
		Nama:     "Budi Staf",
		Username: "budi.opd",
		Password: encodedHashOPD,
		NIP:      "199001012020121001",
		Jabatan:  "Staf Perencanaan",
	}
	DB.FirstOrCreate(&userOPD, UserOPD{Username: "budi.opd"})

	userPemda := UserPemda{
		Nama:     "Anisa Kepala",
		Username: "anisa.pemda",
		Password: encodedHashPemda,
		NIP:      "198505052015012002",
		Jabatan:  "Kepala Bidang Perencanaan",
	}
	DB.FirstOrCreate(&userPemda, UserPemda{Username: "anisa.pemda"})

	userOpdID := userOPD.ID
	userPemdaID := userPemda.ID
	now := time.Now()
	validatedAt := now.Add(-24 * time.Hour)
	
    // ... (sisa kode seeder untuk layanan) ...
	// 2. Buat Data Dummy untuk Layanan Pembangunan (4 data)
	// =======================================================
	fmt.Println("Seeding Layanan Pembangunan...")
	layanansPembangunan := []LayananPembangunan{
		{ // Data Disetujui
			UserOPDID:          userOpdID,
			JenisLayanan:       "Perencanaan Pembangunan Daerah",
			JudulKegiatan:      "Asistensi Penyusunan Renstra Diskominfo 2025-2029",
			Deskripsi:          "Memberikan pendampingan dalam penyusunan dokumen strategis untuk Diskominfo.",
			Status:             "Selesai",
			TanggalKegiatan:    now.AddDate(0, 0, -10),
			Lokasi:             "Kantor Diskominfo",
			StatusValidasi:     "Disetujui",
			IDValidatorPemda:   &userPemdaID,
			KeteranganValidasi: "Data sudah lengkap dan sesuai.",
			TanggalValidasi:    &validatedAt,
		},
		{ // Data Ditolak
			UserOPDID:          userOpdID,
			JenisLayanan:       "Evaluasi Pembangunan Daerah",
			JudulKegiatan:      "Laporan Evaluasi Kinerja Pembangunan Triwulan 3",
			Deskripsi:          "Analisis capaian indikator makro dan program prioritas pada Triwulan 3.",
			Status:             "Revisi",
			TanggalKegiatan:    now.AddDate(0, 0, -5),
			Lokasi:             "Internal Bappelitbangda",
			StatusValidasi:     "Ditolak",
			IDValidatorPemda:   &userPemdaID,
			KeteranganValidasi: "Mohon revisi bagian capaian indikator. Data belum sinkron dengan laporan keuangan.",
			TanggalValidasi:    &validatedAt,
		},
		{ // Data Menunggu Validasi
			UserOPDID:       userOpdID,
			JenisLayanan:    "Penelitian dan Pengembangan",
			JudulKegiatan:   "Kajian Potensi Ekonomi Kreatif di Kota Madiun",
			Deskripsi:       "Proposal untuk melakukan penelitian mendalam mengenai subsektor ekonomi kreatif yang potensial.",
			Status:          "Diajukan",
			TanggalKegiatan: now,
			Lokasi:          "Wilayah Kota Madiun",
			StatusValidasi:  "Menunggu Validasi", // Default
		},
		{ // Data Disetujui Lainnya
			UserOPDID:          userOpdID,
			JenisLayanan:       "Data dan Informasi Pembangunan",
			JudulKegiatan:      "Permintaan Data Laju Pertumbuhan Ekonomi 5 Tahun Terakhir",
			Deskripsi:          "Data LPE diminta oleh BPS untuk keperluan survei nasional.",
			Status:             "Selesai",
			TanggalKegiatan:    now.AddDate(0, 0, -20),
			Lokasi:             "Internal Bappelitbangda",
			StatusValidasi:     "Disetujui",
			IDValidatorPemda:   &userPemdaID,
			KeteranganValidasi: "Data disetujui untuk diberikan.",
			TanggalValidasi:    &validatedAt,
		},
	}
	for _, layanan := range layanansPembangunan {
		DB.FirstOrCreate(&layanan, "judul_kegiatan = ?", layanan.JudulKegiatan)
	}

	// 3. Buat Data Dummy untuk Layanan Administrasi (3 data)
	// ======================================================
	fmt.Println("Seeding Layanan Administrasi...")
	layanansAdministrasi := []LayananAdministrasi{
		{ // Data Disetujui
			UserOPDID:          userOpdID,
			JenisFasilitasi:    "Magang",
			NamaPemohon:        "Rendi Pratama",
			NamaInstansi:       "UNIPMA Madiun",
			JudulKegiatan:      "Praktik Kerja Lapangan di Bidang Perencanaan",
			StatusPermohonan:   "Diterima",
			StatusValidasi:     "Disetujui",
			IDValidatorPemda:   &userPemdaID,
			KeteranganValidasi: "Ditempatkan di Bidang Perencanaan, mulai tanggal 1 Oktober 2025.",
			TanggalValidasi:    &validatedAt,
		},
		{ // Data Ditolak
			UserOPDID:          userOpdID,
			JenisFasilitasi:    "Izin Penelitian",
			NamaPemohon:        "Sarah Amelia",
			NamaInstansi:       "UGM Yogyakarta",
			JudulKegiatan:      "Analisis Dampak Sosial Pembangunan Mall",
			StatusPermohonan:   "Ditolak",
			StatusValidasi:     "Ditolak",
			IDValidatorPemda:   &userPemdaID,
			KeteranganValidasi: "Topik penelitian tidak relevan dengan prioritas pembangunan daerah saat ini.",
			TanggalValidasi:    &validatedAt,
		},
		{ // Data Menunggu Validasi
			UserOPDID:        userOpdID,
			JenisFasilitasi:  "Permohonan Narasumber",
			NamaPemohon:      "HIMA PWK ITS",
			NamaInstansi:     "ITS Surabaya",
			JudulKegiatan:    "Webinar Nasional 'Smart City Planning'",
			StatusPermohonan: "Diajukan",
			StatusValidasi:   "Menunggu Validasi",
		},
	}
	for _, layanan := range layanansAdministrasi {
		DB.FirstOrCreate(&layanan, "judul_kegiatan = ? AND nama_pemohon = ?", layanan.JudulKegiatan, layanan.NamaPemohon)
	}

	// 4. Buat Data Dummy untuk Layanan Informasi & Pengaduan (3 data)
	// =================================================================
	fmt.Println("Seeding Layanan Informasi & Pengaduan...")
	layanansInfo := []LayananInformasiPengaduan{
		{ // Data Disetujui
			UserOPDID:          userOpdID,
			JenisPermintaan:    "Informasi Publik",
			KodeRegistrasi:     "INFO-2025-001",
			NamaPemohon:        "Warga Madiun",
			KontakPemohon:      "warga@email.com",
			DetailIsi:          "Mohon data rincian anggaran Bappelitbangda TA 2025.",
			Status:             "Telah Dibalas",
			StatusValidasi:     "Disetujui",
			IDValidatorPemda:   &userPemdaID,
			KeteranganValidasi: "Jawaban sudah sesuai dengan ketentuan UU KIP.",
			TanggalValidasi:    &validatedAt,
		},
		{ // Data Menunggu Validasi
			UserOPDID:        userOpdID,
			JenisPermintaan:  "Pengaduan",
			KodeRegistrasi:   "PENG-2025-001",
			NamaPemohon:      "Sarah Amelia",
			KontakPemohon:    "sarah.a@email.com",
			DetailIsi:        "Pengajuan izin penelitian saya belum ada kabar setelah 2 minggu, padahal sudah ditolak.",
			Status:           "Diterima",
			StatusValidasi:   "Menunggu Validasi",
		},
		{ // Data Ditolak
			UserOPDID:          userOpdID,
			JenisPermintaan:    "Informasi Publik",
			KodeRegistrasi:     "INFO-2025-002",
			NamaPemohon:        "Perusahaan Riset",
			KontakPemohon:      "riset@corp.com",
			DetailIsi:          "Meminta data mentah hasil survei kepuasan masyarakat.",
			Status:             "Ditolak",
			StatusValidasi:     "Ditolak",
			IDValidatorPemda:   &userPemdaID,
			KeteranganValidasi: "Data mentah termasuk informasi pribadi yang dikecualikan menurut UU KIP.",
			TanggalValidasi:    &validatedAt,
		},
	}
	for _, layanan := range layanansInfo {
		DB.FirstOrCreate(&layanan, "kode_registrasi = ?", layanan.KodeRegistrasi)
	}

}