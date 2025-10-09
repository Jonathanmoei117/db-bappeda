package main

import "time"

// ValidasiRequest adalah struct untuk menampung body request
// saat Pemda melakukan validasi pengajuan.
type ValidasiRequest struct {
	StatusValidasi   string `json:"status_validasi" binding:"required"`
	KeteranganValidasi string `json:"keterangan_validasi"`
}

//================================================================================
// TABEL MASTER DATA
//================================================================================

// OPD merepresentasikan data master Organisasi Perangkat Daerah.
// Tabel: opd
type OPD struct {
	ID        uint   `gorm:"column:id_opd;primaryKey" json:"id_opd"`
	NamaOPD   string `gorm:"column:nama_opd;unique;not null;type:varchar(255)" json:"nama_opd"`
	AlamatOPD string `gorm:"column:alamat_opd;type:text" json:"alamat_opd"`

	// Relasi (sebagai Parent)
	UserOPDs        []UserOPD        `gorm:"foreignKey:IDOPD" json:"-"`
	JenisPelayanans []JenisPelayanan `gorm:"foreignKey:IDOPD" json:"-"`
}

// JenisPelayanan merepresentasikan Standar Pelayanan yang dimiliki oleh setiap OPD.
// Tabel: jenis_pelayanan
type JenisPelayanan struct {
	ID                 uint      `gorm:"column:id_jenis_pelayanan;primaryKey" json:"id_jenis_pelayanan"`
	IDOPD              uint      `gorm:"column:id_opd;not null" json:"id_opd"`
	NamaStandar        string    `gorm:"column:nama_standar;not null;type:varchar(255)" json:"nama_standar"`
	WaktuPelayanan     string    `gorm:"column:waktu_pelayanan;type:varchar(255)" json:"waktu_pelayanan"`
	BiayaTarif         string    `gorm:"column:biaya_tarif;type:varchar(255)" json:"biaya_tarif"`
	ProdukPelayanan    string    `gorm:"column:produk_pelayanan;type:varchar(255)" json:"produk_pelayanan"`
	Fasilitas          string    `gorm:"column:fasilitas;type:text" json:"fasilitas"` // DIUBAH
	JaminanPelayanan   string    `gorm:"column:jaminan_pelayanan;type:text" json:"jaminan_pelayanan"`
	JaminanKeselamatan string    `gorm:"column:jaminan_keselamatan;type:text" json:"jaminan_keselamatan"`
	EvaluasiKinerja    string    `gorm:"column:evaluasi_kinerja;type:text" json:"evaluasi_kinerja"`
	JumlahPelaksana    int       `gorm:"column:jumlah_pelaksana" json:"jumlah_pelaksana"`
	JangkaWaktu        string    `gorm:"column:jangka_waktu;type:varchar(255)" json:"jangka_waktu"` // DITAMBAHKAN
	CreatedAt          time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relasi (sebagai Child dan Parent)
	OPD             OPD             `gorm:"foreignKey:IDOPD" json:"opd"`
	FormPengajuans  []FormPengajuan `gorm:"foreignKey:IDJenisPelayanan" json:"-"`
}

//================================================================================
// TABEL PENGGUNA
//================================================================================

// UserOPD merepresentasikan pengguna dari OPD yang bertugas menginput data.
// Tabel: user_opd
type UserOPD struct {
	ID        uint      `gorm:"column:id_user_opd;primaryKey" json:"id_user_opd"`
	IDOPD     uint      `gorm:"column:id_opd;not null" json:"id_opd"` // Foreign Key ke OPD
	Nama      string    `gorm:"column:nama;not null;type:varchar(255)" json:"nama"`
	NIP       string    `gorm:"column:nip;unique;not null;type:varchar(255)" json:"nip"` // Digunakan untuk login
	Password  string    `gorm:"column:password;not null;type:varchar(255)" json:"-"`     // Sembunyikan saat JSON output
	Jabatan   string    `gorm:"column:jabatan;type:varchar(255)" json:"jabatan"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relasi (sebagai Child dan Parent)
	OPD            OPD             `gorm:"foreignKey:IDOPD" json:"opd"`
	FormPengajuans []FormPengajuan `gorm:"foreignKey:IDUserOPD" json:"-"`
}

// UserPemda merepresentasikan pengguna dari Pemda yang bertugas sebagai validator.
// Tabel: user_pemda
type UserPemda struct {
	ID        uint      `gorm:"column:id_user_pemda;primaryKey" json:"id_user_pemda"`
	Nama      string    `gorm:"column:nama;not null;type:varchar(255)" json:"nama"`
	NIP       string    `gorm:"column:nip;unique;not null;type:varchar(255)" json:"nip"` // Digunakan untuk login
	Password  string    `gorm:"column:password;not null;type:varchar(255)" json:"-"`     // Sembunyikan saat JSON output
	Jabatan   string    `gorm:"column:jabatan;type:varchar(255)" json:"jabatan"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relasi (sebagai Parent)
	ValidatedFormPengajuans []FormPengajuan `gorm:"foreignKey:IDValidatorPemda" json:"-"`
}

//================================================================================
// TABEL TRANSAKSI
//================================================================================

// FormPengajuan merepresentasikan tabel transaksi utama untuk semua pengajuan layanan.
// Tabel: form_pengajuan
type FormPengajuan struct {
	ID uint `gorm:"column:id_form_pengajuan;primaryKey" json:"id_form_pengajuan"`

	// Foreign Keys
	IDUserOPD        uint  `gorm:"column:id_user_opd;not null" json:"id_user_opd"`
	IDJenisPelayanan uint  `gorm:"column:id_jenis_pelayanan;not null" json:"id_jenis_pelayanan"`
	IDValidatorPemda *uint `gorm:"column:id_validator_pemda" json:"id_validator_pemda"` // Pointer karena bisa NULL

	// Atribut Transaksi
	JudulKegiatan       string     `gorm:"column:judul_kegiatan;not null;type:varchar(255)" json:"judul_kegiatan"`
	Deskripsi           string     `gorm:"column:deskripsi;not null;type:text" json:"deskripsi"`
	NamaPemohon         string     `gorm:"column:nama_pemohon;not null;type:varchar(255)" json:"nama_pemohon"`
	NIPPemohon          *string    `gorm:"column:nip_pemohon;type:varchar(255)" json:"nip_pemohon"` // Pointer karena bisa NULL
	InstansiPemohon     string     `gorm:"column:instansi_pemohon;not null;type:varchar(255)" json:"instansi_pemohon"`
	PeriodeMulai        *time.Time `gorm:"column:periode_mulai;type:date" json:"periode_mulai"`            // Pointer karena bisa NULL
	PeriodeSelesai      *time.Time `gorm:"column:periode_selesai;type:date" json:"periode_selesai"`        // Pointer karena bisa NULL
	BerkasPengajuanPath *string    `gorm:"column:berkas_pengajuan_path;type:varchar(255)" json:"berkas_pengajuan_path"` // Pointer karena bisa NULL

	// Status & Waktu
	StatusProses       string     `gorm:"column:status_proses;not null;default:'Baru';type:varchar(255)" json:"status_proses"`
	StatusValidasi     string     `gorm:"column:status_validasi;not null;default:'Menunggu Validasi';type:varchar(255)" json:"status_validasi"`
	KeteranganValidasi *string    `gorm:"column:keterangan_validasi;type:text" json:"keterangan_validasi"` // Pointer karena bisa NULL
	CreatedAt          time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	TanggalValidasi    *time.Time `gorm:"column:tanggal_validasi" json:"tanggal_validasi"` // Pointer karena bisa NULL

	// Relasi
	UserOPD        UserOPD        `gorm:"foreignKey:IDUserOPD" json:"user_opd"`
	JenisPelayanan JenisPelayanan `gorm:"foreignKey:IDJenisPelayanan" json:"jenis_pelayanan"`
	ValidatorPemda *UserPemda     `gorm:"foreignKey:IDValidatorPemda" json:"validator_pemda"` // Pointer karena FK bisa NULL
}