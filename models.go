package main

import "time"

// UserOPD merepresentasikan petugas OPD yang menginput data layanan.
type UserOPD struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Nama      string    `json:"nama"`
	Username  string    `gorm:"unique" json:"username"`
	Password  string    `json:"-"`
	NIP       string    `gorm:"unique" json:"nip"`
	Jabatan   string    `json:"jabatan"`
	CreatedAt time.Time `json:"created_at"`
}

// UserPemda merepresentasikan user dari Pemda yang melakukan validasi.
type UserPemda struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Nama      string    `json:"nama"`
	Username  string    `gorm:"unique" json:"username"`
	Password  string    `json:"-"`
	NIP       string    `gorm:"unique" json:"nip"`
	Jabatan   string    `json:"jabatan"`
	CreatedAt time.Time `json:"created_at"`
}

// LayananPembangunan merepresentasikan kegiatan/layanan terkait pembangunan daerah.
type LayananPembangunan struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	UserOPDID          uint       `json:"user_opd_id"`
	JenisLayanan       string     `json:"jenis_layanan"`
	JudulKegiatan      string     `json:"judul_kegiatan"`
	Deskripsi          string     `json:"deskripsi"`
	Status             string     `json:"status"`
	TanggalKegiatan    time.Time  `json:"tanggal_kegiatan"`
	Lokasi             string     `json:"lokasi"`
	Catatan            string     `json:"catatan"`
	CreatedAt          time.Time  `json:"created_at"`
	IDValidatorPemda   *uint      `json:"id_validator_pemda"`
	StatusValidasi     string     `gorm:"default:'Menunggu Validasi'" json:"status_validasi"`
	KeteranganValidasi string     `json:"keterangan_validasi"`
	TanggalValidasi    *time.Time `json:"tanggal_validasi"`
	UserOPD            UserOPD    `gorm:"foreignKey:UserOPDID" json:"user_opd"`
	ValidatorPemda     UserPemda  `gorm:"foreignKey:IDValidatorPemda" json:"validator_pemda"`
}

// LayananInformasiPengaduan merepresentasikan layanan informasi publik dan pengaduan.
type LayananInformasiPengaduan struct {
	ID                  uint       `gorm:"primaryKey" json:"id"`
	UserOPDID           uint       `json:"user_opd_id"`
	JenisPermintaan     string     `json:"jenis_permintaan"`
	KodeRegistrasi      string     `gorm:"unique" json:"kode_registrasi"`
	NamaPemohon         string     `json:"nama_pemohon"`
	KontakPemohon       string     `json:"kontak_pemohon"`
	DetailIsi           string     `json:"detail_isi"`
	Status              string     `json:"status"`
	CatatanTindakLanjut string     `json:"catatan_tindak_lanjut"`
	CreatedAt           time.Time  `json:"created_at"`
	IDValidatorPemda    *uint      `json:"id_validator_pemda"`
	StatusValidasi      string     `gorm:"default:'Menunggu Validasi'" json:"status_validasi"`
	KeteranganValidasi  string     `json:"keterangan_validasi"`
	TanggalValidasi     *time.Time `json:"tanggal_validasi"`
	UserOPD             UserOPD    `gorm:"foreignKey:UserOPDID" json:"user_opd"`
	ValidatorPemda      UserPemda  `gorm:"foreignKey:IDValidatorPemda" json:"validator_pemda"`
}

// LayananAdministrasi merepresentasikan layanan administrasi seperti izin magang/penelitian.
type LayananAdministrasi struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	UserOPDID          uint       `json:"user_opd_id"`
	JenisFasilitasi    string     `json:"jenis_fasilitasi"`
	NamaPemohon        string     `json:"nama_pemohon"`
	NamaInstansi       string     `json:"nama_instansi"`
	JudulKegiatan      string     `json:"judul_kegiatan"`
	SuratPengantarPath string     `json:"surat_pengantar_path"`
	PeriodeMulai       time.Time  `json:"periode_mulai"`
	PeriodeSelesai     time.Time  `json:"periode_selesai"`
	StatusPermohonan   string     `json:"status_permohonan"`
	CreatedAt          time.Time  `json:"created_at"`
	IDValidatorPemda   *uint      `json:"id_validator_pemda"`
	StatusValidasi     string     `gorm:"default:'Menunggu Validasi'" json:"status_validasi"`
	KeteranganValidasi string     `json:"keterangan_validasi"`
	TanggalValidasi    *time.Time `json:"tanggal_validasi"`
	UserOPD            UserOPD    `gorm:"foreignKey:UserOPDID" json:"user_opd"`
	ValidatorPemda     UserPemda  `gorm:"foreignKey:IDValidatorPemda" json:"validator_pemda"`
}

// ValidasiRequest adalah struct untuk request body saat validasi
type ValidasiRequest struct {
	StatusValidasi     string `json:"status_validasi"`
	KeteranganValidasi string `json:"keterangan_validasi"`
}