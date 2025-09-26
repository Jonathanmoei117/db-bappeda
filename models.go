package main

import "time"

// OPD merepresentasikan tabel master untuk Organisasi Perangkat Daerah.
type OPD struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	NamaOPD   string    `gorm:"unique;not null" json:"nama_opd"`
	AlamatOPD string    `json:"alamat_opd"`
	Users     []UserOPD `gorm:"foreignKey:OPDID" json:"-"`
}

// UserOPD merepresentasikan petugas OPD yang menginput data layanan.
type UserOPD struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OPDID     uint      `gorm:"not null" json:"opd_id"`
	Nama      string    `gorm:"not null" json:"nama"`
	NIP       string    `gorm:"unique;not null" json:"nip"` // NIP digunakan untuk login
	Password  string    `gorm:"not null" json:"-"`
	Jabatan   string    `json:"jabatan"`
	CreatedAt time.Time `json:"created_at"`

	// Relasi
	OPD                       OPD                       `gorm:"foreignKey:OPDID" json:"opd"`
	LayananPembangunans       []LayananPembangunan      `gorm:"foreignKey:UserOPDID" json:"-"`
	LayananAdministrasis      []LayananAdministrasi     `gorm:"foreignKey:UserOPDID" json:"-"`
	LayananInformasiPengaduans []LayananInformasiPengaduan `gorm:"foreignKey:UserOPDID" json:"-"`
}

// UserPemda merepresentasikan user dari Pemda yang melakukan validasi.
type UserPemda struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Nama      string    `gorm:"not null" json:"nama"`
	NIP       string    `gorm:"unique;not null" json:"nip"` // NIP digunakan untuk login
	Password  string    `gorm:"not null" json:"-"`
	Jabatan   string    `json:"jabatan"`
	CreatedAt time.Time `json:"created_at"`

	// Relasi
	ValidatedLayananPembangunans      []LayananPembangunan      `gorm:"foreignKey:IDValidatorPemda" json:"-"`
	ValidatedLayananAdministrasis      []LayananAdministrasi     `gorm:"foreignKey:IDValidatorPemda" json:"-"`
	ValidatedLayananInformasiPengaduans []LayananInformasiPengaduan `gorm:"foreignKey:IDValidatorPemda" json:"-"`
}

// ======== STRUCT LAYANAN (DIPISAH) ========

// LayananPembangunan merepresentasikan entitas di tabel Layanan_Pembangunan.
type LayananPembangunan struct {
	ID                  uint       `gorm:"primaryKey" json:"id"`
	UserOPDID           uint       `gorm:"not null" json:"user_opd_id"`
	JudulKegiatan       string     `gorm:"not null" json:"judul_kegiatan"`
	Deskripsi           string     `json:"deskripsi"`
	NamaPemohon         string     `gorm:"not null" json:"nama_pemohon"`
	NIPPemohon          *string    `json:"nip_pemohon"`
	InstansiPemohon     string     `gorm:"not null" json:"instansi_pemohon"`
	PeriodeMulai        *time.Time `json:"periode_mulai"`
	PeriodeSelesai      *time.Time `json:"periode_selesai"`
	BerkasPengajuanPath *string    `json:"berkas_pengajuan_path"`
	StatusProses        string     `gorm:"not null;default:'Baru'" json:"status_proses"`
	CreatedAt           time.Time  `json:"created_at"`

	// Kolom Validasi
	IDValidatorPemda   *uint      `json:"id_validator_pemda"`
	StatusValidasi     string     `gorm:"not null;default:'Menunggu Validasi'" json:"status_validasi"`
	KeteranganValidasi string     `json:"keterangan_validasi"`
	TanggalValidasi    *time.Time `json:"tanggal_validasi"`

	// Relasi
	UserOPD        UserOPD   `gorm:"foreignKey:UserOPDID" json:"user_opd"`
	ValidatorPemda UserPemda `gorm:"foreignKey:IDValidatorPemda" json:"validator_pemda"`
}

// LayananAdministrasi merepresentasikan entitas di tabel Layanan_Administrasi.
type LayananAdministrasi struct {
	ID                  uint       `gorm:"primaryKey" json:"id"`
	UserOPDID           uint       `gorm:"not null" json:"user_opd_id"`
	JudulKegiatan       string     `gorm:"not null" json:"judul_kegiatan"`
	Deskripsi           string     `json:"deskripsi"`
	NamaPemohon         string     `gorm:"not null" json:"nama_pemohon"`
	NIPPemohon          *string    `json:"nip_pemohon"`
	InstansiPemohon     string     `gorm:"not null" json:"instansi_pemohon"`
	PeriodeMulai        *time.Time `json:"periode_mulai"`
	PeriodeSelesai      *time.Time `json:"periode_selesai"`
	BerkasPengajuanPath *string    `json:"berkas_pengajuan_path"`
	StatusProses        string     `gorm:"not null;default:'Baru'" json:"status_proses"`
	CreatedAt           time.Time  `json:"created_at"`

	// Kolom Validasi
	IDValidatorPemda   *uint      `json:"id_validator_pemda"`
	StatusValidasi     string     `gorm:"not null;default:'Menunggu Validasi'" json:"status_validasi"`
	KeteranganValidasi string     `json:"keterangan_validasi"`
	TanggalValidasi    *time.Time `json:"tanggal_validasi"`

	// Relasi
	UserOPD        UserOPD   `gorm:"foreignKey:UserOPDID" json:"user_opd"`
	ValidatorPemda UserPemda `gorm:"foreignKey:IDValidatorPemda" json:"validator_pemda"`
}

// LayananInformasiPengaduan merepresentasikan entitas di tabel Layanan_Informasi_Pengaduan.
type LayananInformasiPengaduan struct {
	ID                  uint       `gorm:"primaryKey" json:"id"`
	UserOPDID           uint       `gorm:"not null" json:"user_opd_id"`
	JudulKegiatan       string     `gorm:"not null" json:"judul_kegiatan"`
	Deskripsi           string     `gorm:"not null" json:"deskripsi"`
	NamaPemohon         string     `gorm:"not null" json:"nama_pemohon"`
	NIPPemohon          *string    `json:"nip_pemohon"`
	InstansiPemohon     string     `gorm:"not null" json:"instansi_pemohon"`
	BerkasPengajuanPath *string    `json:"berkas_pengajuan_path"`
	StatusProses        string     `gorm:"not null;default:'Baru'" json:"status_proses"`
	CreatedAt           time.Time  `json:"created_at"`

	// Kolom Validasi
	IDValidatorPemda   *uint      `json:"id_validator_pemda"`
	StatusValidasi     string     `gorm:"not null;default:'Menunggu Validasi'" json:"status_validasi"`
	KeteranganValidasi string     `json:"keterangan_validasi"`
	TanggalValidasi    *time.Time `json:"tanggal_validasi"`

	// Relasi
	UserOPD        UserOPD   `gorm:"foreignKey:UserOPDID" json:"user_opd"`
	ValidatorPemda UserPemda `gorm:"foreignKey:IDValidatorPemda" json:"validator_pemda"`
}

// ValidasiRequest tetap digunakan untuk body request saat validasi.
type ValidasiRequest struct {
	StatusValidasi     string `json:"status_validasi" binding:"required"`
	KeteranganValidasi string `json:"keterangan_validasi"`
}