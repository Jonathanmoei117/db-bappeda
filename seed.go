package main

import (
	"fmt"
	"time"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Seed akan mengisi database dengan data awal untuk keperluan development.
func Seed() {
	fmt.Println("===== MEMULAI PROSES SEEDING DATA =====")

	// ==================================================================
	// LANGKAH 1: Buat Master Data OPD
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
	
	// --- OPD (2) ---
	opdBappeda := OPD{NamaOPD: "BAPPEDA", AlamatOPD: "Jl. Pahlawan No. 1"}
	DB.FirstOrCreate(&opdBappeda, OPD{NamaOPD: "BAPPEDA"})

	opdDisdik := OPD{NamaOPD: "Dinas Pendidikan", AlamatOPD: "Jl. Pendidikan No. 2"}
	DB.FirstOrCreate(&opdDisdik, OPD{NamaOPD: "Dinas Pendidikan"})

	// --- User Pemda (1) ---
	userPemda := UserPemda{
		Nama:     "Dr. Anisa Wijayanti",
		Password: string(hashedPasswordPemda), // DIUBAH: Simpan hash bcrypt
		NIP:      "198505052015012002",
		Jabatan:  "Kepala Bidang Verifikasi",
	}
	DB.FirstOrCreate(&userPemda, UserPemda{NIP: "198505052015012002"})

	// --- User OPD (2) ---
	userBappeda := UserOPD{
		OPDID:    opdBappeda.ID,
		Nama:     "Budi Santoso",
		Password: string(hashedPasswordOPD), // DIUBAH: Simpan hash bcrypt
		NIP:      "199001012020121001",
		Jabatan:  "Staf Perencanaan",
	}
	DB.FirstOrCreate(&userBappeda, UserOPD{NIP: "199001012020121001"})

	userDisdik := UserOPD{
		OPDID:    opdDisdik.ID,
		Nama:     "Citra Lestari",
		Password: string(hashedPasswordOPD), // DIUBAH: Simpan hash bcrypt
		NIP:      "199203152021012003",
		Jabatan:  "Staf Kurikulum",
	}
	DB.FirstOrCreate(&userDisdik, UserOPD{NIP: "199203152021012003"})

	log.Println("🌱 Seeding dengan password Bcrypt selesai!") // DIUBAH: Log message

	// ==================================================================
	// LANGKAH 3: Buat Data Layanan Dummy (10 per jenis layanan)
	// ==================================================================
	validatedAt := time.Now().Add(-48 * time.Hour)

	// --- 10 Data Layanan Pembangunan ---
	fmt.Println("--> Seeding 10 data Layanan Pembangunan...")
	pembangunans := []LayananPembangunan{
		{UserOPDID: userBappeda.ID, JudulKegiatan: "Kajian Awal Pembangunan Ruang Terbuka Hijau", Deskripsi: "Analisis kelayakan untuk taman kota baru di Kecamatan Taman.", InstansiPemohon: "Internal Bappeda", NamaPemohon: userBappeda.Nama, StatusValidasi: "Disetujui", KeteranganValidasi: "Data lengkap.", IDValidatorPemda: &userPemda.ID, TanggalValidasi: &validatedAt},
		{UserOPDID: userDisdik.ID, JudulKegiatan: "Rencana Rehabilitasi Gedung SDN 01 Kartoharjo", Deskripsi: "Pengajuan anggaran untuk perbaikan atap dan fasilitas toilet sekolah.", InstansiPemohon: "Internal Disdik", NamaPemohon: userDisdik.Nama, StatusValidasi: "Disetujui", KeteranganValidasi: "Sudah sesuai RAB.", IDValidatorPemda: &userPemda.ID, TanggalValidasi: &validatedAt},
		{UserOPDID: userBappeda.ID, JudulKegiatan: "Proposal Sistem Drainase Perkotaan", Deskripsi: "Pengajuan proposal untuk perbaikan sistem drainase di area rawan banjir.", InstansiPemohon: "Internal Bappeda", NamaPemohon: userBappeda.Nama, StatusValidasi: "Ditolak", KeteranganValidasi: "Lampiran teknis tidak ada.", IDValidatorPemda: &userPemda.ID, TanggalValidasi: &validatedAt},
		{UserOPDID: userDisdik.ID, JudulKegiatan: "Pengadaan Komputer untuk Laboratorium SMPN 2", Deskripsi: "Proposal pengadaan 40 unit komputer untuk meningkatkan fasilitas belajar.", InstansiPemohon: "Internal Disdik", NamaPemohon: userDisdik.Nama, StatusValidasi: "Menunggu Validasi"},
		{UserOPDID: userBappeda.ID, JudulKegiatan: "Studi Kelayakan Flyover Jalan Pahlawan", Deskripsi: "Mengkaji urgensi pembangunan flyover untuk mengurangi kemacetan.", InstansiPemohon: "Dinas PUPR", NamaPemohon: "Dr. Ir. Haryanto"},
		{UserOPDID: userDisdik.ID, JudulKegiatan: "Pembangunan Perpustakaan Digital Kota", Deskripsi: "Mengajukan konsep dan anggaran untuk perpustakaan modern berbasis teknologi.", InstansiPemohon: "Internal Disdik", NamaPemohon: userDisdik.Nama},
		{UserOPDID: userBappeda.ID, JudulKegiatan: "Program Bedah Rumah Warga Kurang Mampu 2025", Deskripsi: "Perencanaan dan pendataan calon penerima bantuan renovasi rumah.", InstansiPemohon: "Dinas Sosial", NamaPemohon: "Siti Aminah, S.Sos"},
		{UserOPDID: userDisdik.ID, JudulKegiatan: "Renovasi Lapangan Olahraga SMAN 3 Madiun", Deskripsi: "Perbaikan lapangan basket dan futsal beserta tribun penonton.", InstansiPemohon: "Internal Disdik", NamaPemohon: userDisdik.Nama},
		{UserOPDID: userBappeda.ID, JudulKegiatan: "Analisis Dampak Lingkungan (AMDAL) Pasar Besar", Deskripsi: "Studi AMDAL untuk rencana revitalisasi Pasar Besar Madiun.", InstansiPemohon: "DLH", NamaPemohon: "Ir. Endang P."},
		{UserOPDID: userDisdik.ID, JudulKegiatan: "Penambahan Ruang Kelas Baru (RKB) di SLB Kartini", Deskripsi: "Proposal penambahan 2 ruang kelas untuk siswa berkebutuhan khusus.", InstansiPemohon: "Internal Disdik", NamaPemohon: userDisdik.Nama},
	}
	for _, l := range pembangunans {
		DB.FirstOrCreate(&l, LayananPembangunan{JudulKegiatan: l.JudulKegiatan})
	}

	// --- 10 Data Layanan Administrasi ---
	fmt.Println("--> Seeding 10 data Layanan Administrasi...")
	administrasis := []LayananAdministrasi{
		{UserOPDID: userBappeda.ID, JudulKegiatan: "Izin Penelitian Mahasiswa S2 ITB", Deskripsi: "Penelitian mengenai dampak ekonomi digital terhadap UMKM di Madiun.", InstansiPemohon: "Institut Teknologi Bandung", NamaPemohon: "Ahmad Zulkifli", StatusValidasi: "Disetujui", KeteranganValidasi: "Surat dari kampus lengkap.", IDValidatorPemda: &userPemda.ID, TanggalValidasi: &validatedAt},
		{UserOPDID: userDisdik.ID, JudulKegiatan: "Permohonan Narasumber untuk Seminar Guru", Deskripsi: "Memohon Kepala Dinas sebagai pembicara utama dalam acara PGRI.", InstansiPemohon: "PGRI Kota Madiun", NamaPemohon: "Panitia Seminar", StatusValidasi: "Menunggu Validasi"},
		{UserOPDID: userBappeda.ID, JudulKegiatan: "Fasilitasi Magang dari Universitas Merdeka", Deskripsi: "Permohonan magang untuk 3 mahasiswa di bidang perencanaan kota.", InstansiPemohon: "Universitas Merdeka Madiun", NamaPemohon: "Fakultas Teknik"},
		{UserOPDID: userDisdik.ID, JudulKegiatan: "Izin Studi Banding dari Kabupaten Ngawi", Deskripsi: "Kunjungan dari Dinas Pendidikan Ngawi untuk mempelajari Kurikulum Merdeka.", InstansiPemohon: "Dinas Pendidikan Ngawi", NamaPemohon: "Kabid Kurikulum Ngawi"},
		{UserOPDID: userBappeda.ID, JudulKegiatan: "Permohonan Rekomendasi Beasiswa LPDP", Deskripsi: "Permohonan surat rekomendasi dari Kepala Bappeda untuk S3.", InstansiPemohon: "Pribadi", NamaPemohon: "Andi Pratama"},
		{UserOPDID: userDisdik.ID, JudulKegiatan: "Verifikasi Ijazah untuk Keperluan PNS", Deskripsi: "Legalisasi dan verifikasi ijazah atas nama Susi Susanti.", InstansiPemohon: "BKPSDM", NamaPemohon: "Susi Susanti"},
		{UserOPDID: userBappeda.ID, JudulKegiatan: "Surat Keterangan Telah Melaksanakan Riset", Deskripsi: "Penerbitan surat keterangan untuk mahasiswa UGM.", InstansiPemohon: "Universitas Gadjah Mada", NamaPemohon: "Rina Hartati"},
		{UserOPDID: userDisdik.ID, JudulKegiatan: "Izin Penggunaan Aula untuk Lomba Cerdas Cermat", Deskripsi: "Peminjaman Aula Dinas Pendidikan untuk final LCC tingkat SMP.", InstansiPemohon: "MGMP IPA SMP", NamaPemohon: "Ketua Panitia"},
		{UserOPDID: userBappeda.ID, JudulKegiatan: "Permohonan Data Spasial Wilayah Rawan Banjir", Deskripsi: "Data diminta oleh tim peneliti dari ITS Surabaya.", InstansiPemohon: "Institut Teknologi Sepuluh Nopember", NamaPemohon: "Dr. Bambang S."},
		{UserOPDID: userDisdik.ID, JudulKegiatan: "Pendaftaran Ulang Sekolah Swasta", Deskripsi: "Proses administrasi untuk pendaftaran ulang izin operasional Yayasan Pelita Harapan.", InstansiPemohon: "Yayasan Pelita Harapan", NamaPemohon: "Kepala Sekolah"},
	}
	for _, l := range administrasis {
		DB.FirstOrCreate(&l, LayananAdministrasi{JudulKegiatan: l.JudulKegiatan})
	}

	// --- 10 Data Layanan Informasi & Pengaduan ---
	fmt.Println("--> Seeding 10 data Layanan Informasi & Pengaduan...")
	informasis := []LayananInformasiPengaduan{
		{UserOPDID: userBappeda.ID, JudulKegiatan: "Permintaan Data PDRB 5 Tahun Terakhir", Deskripsi: "Data diminta oleh Bank Indonesia untuk laporan triwulanan.", InstansiPemohon: "Bank Indonesia", NamaPemohon: "BI Perwakilan Jatim", StatusValidasi: "Disetujui", KeteranganValidasi: "Data sudah dikirim via email.", IDValidatorPemda: &userPemda.ID, TanggalValidasi: &validatedAt},
		{UserOPDID: userDisdik.ID, JudulKegiatan: "Pengaduan Terkait Pungutan Liar di Sekolah X", Deskripsi: "Laporan dari wali murid mengenai adanya biaya tambahan yang tidak resmi.", InstansiPemohon: "Masyarakat", NamaPemohon: "Wali Murid (Anonim)", StatusValidasi: "Menunggu Validasi"},
		{UserOPDID: userBappeda.ID, JudulKegiatan: "Keluhan Mengenai Lampu Jalan Mati di Jl. Serayu", Deskripsi: "Laporan warga RT 02 RW 05, lampu sudah mati selama 2 minggu.", InstansiPemohon: "Masyarakat", NamaPemohon: "Bapak Sutrisno"},
		{UserOPDID: userDisdik.ID, JudulKegiatan: "Permintaan Informasi PPDB Jalur Zonasi", Deskripsi: "Orang tua murid meminta penjelasan detail mengenai aturan zonasi PPDB 2025.", InstansiPemohon: "Masyarakat", NamaPemohon: "Ibu Indah"},
		{UserOPDID: userBappeda.ID, JudulKegiatan: "Laporan Tumpukan Sampah di Dekat Jembatan Manguharjo", Deskripsi: "Sampah liar menumpuk dan menimbulkan bau tidak sedap.", InstansiPemohon: "Masyarakat", NamaPemohon: "Komunitas Peduli Sungai"},
		{UserOPDID: userDisdik.ID, JudulKegiatan: "Apresiasi atas Prestasi Siswa di Olimpiade Sains", Deskripsi: "Ucapan terima kasih dan apresiasi kepada Disdik atas dukungan kepada siswa berprestasi.", InstansiPemohon: "Masyarakat", NamaPemohon: "Forum Orang Tua Siswa"},
		{UserOPDID: userBappeda.ID, JudulKegiatan: "Informasi Progres Pembangunan Taman Kota", Deskripsi: "Warga menanyakan kapan taman kota baru akan selesai dibangun.", InstansiPemohon: "Masyarakat", NamaPemohon: "Karang Taruna Kel. Taman"},
		{UserOPDID: userDisdik.ID, JudulKegiatan: "Pengaduan Kualitas Makanan Kantin Sekolah Y", Deskripsi: "Beberapa siswa mengeluh sakit perut setelah jajan di kantin sekolah.", InstansiPemohon: "Masyarakat", NamaPemohon: "Orang Tua Murid"},
		{UserOPDID: userBappeda.ID, JudulKegiatan: "Permintaan Data Jumlah UMKM per Kecamatan", Deskripsi: "Data dibutuhkan untuk analisis potensi ekonomi oleh akademisi.", InstansiPemohon: "Universitas Widya Mandala", NamaPemohon: "Diana, S.E., M.M."},
		{UserOPDID: userDisdik.ID, JudulKegiatan: "Saran Penambahan Ekstrakurikuler Robotik", Deskripsi: "Usulan dari komunitas teknologi agar sekolah negeri memfasilitasi ekskul robotik.", InstansiPemohon: "Komunitas Madiun Coder", NamaPemohon: "Ketua Komunitas"},
	}
	for _, l := range informasis {
		DB.FirstOrCreate(&l, LayananInformasiPengaduan{JudulKegiatan: l.JudulKegiatan})
	}

	fmt.Println("===== PROSES SEEDING SELESAI =====")
}
