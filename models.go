package main

import "time"

//================================================================================
// VALIDASI REQUEST STRUCT
//================================================================================

// ValidasiRequest adalah struct untuk menampung body request
// saat Pemda melakukan validasi standar pelayanan.
type ValidasiRequest struct {
	StatusValidasi  string `json:"status_validasi" binding:"required"`
	KeteranganValidasi string `json:"keterangan_validasi"`
}

//================================================================================
// TABEL OPD
//================================================================================

// OPD merepresentasikan data master Organisasi Perangkat Daerah.
// Tabel: opd (1)
type OPD struct {
	ID    uint  `gorm:"column:id_opd;primaryKey" json:"id_opd"`
	NamaOPD  string `gorm:"column:nama_opd;unique;not null;type:varchar(255)" json:"nama_opd"`
	AlamatOPD string `gorm:"column:alamat_opd;type:text" json:"alamat_opd"`

	// Relasi (sebagai Parent)
	UserOPDs  []UserOPD `gorm:"foreignKey:IDOPD" json:"-"`
	JenisPelayanans []JenisPelayanan `gorm:"foreignKey:IDOPD" json:"-"`
}

//================================================================================
// TABEL JENIS PELAYANAN
//================================================================================

// JenisPelayanan merepresentasikan Standar Pelayanan yang divalidasi Pemda.
// Tabel: jenis_pelayanan (2)
type JenisPelayanan struct {
	ID  uint `gorm:"column:id_jenis_pelayanan;primaryKey" json:"id_jenis_pelayanan"`
	IDOPD uint `gorm:"column:id_opd;not null" json:"id_opd"`
	IDValidatorPemda *uint `gorm:"column:id_validator_pemda" json:"id_validator_pemda"` // Relasi ke UserPemda

	// --- ATRIBUT STANDAR PELAYANAN BARU ---
	NamaStandar string `gorm:"column:nama_standar;unique;not null;type:varchar(255)" json:"nama_standar"` // <--- UNIQUE DITAMBAHKAN
	DasarHukum string `gorm:"column:dasar_hukum;type:text" json:"dasar_hukum"`
	Persyaratan string `gorm:"column:persyaratan;type:text" json:"persyaratan"`
	SistemMekanismeProsedurPath string `gorm:"column:sistem_mekanisme_prosedur_path;type:varchar(255)" json:"sistem_mekanisme_prosedur_path"`
	WaktuPelayanan string `gorm:"column:waktu_pelayanan;type:varchar(255)" json:"waktu_pelayanan"`
	BiayaTarif string`gorm:"column:biaya_tarif;type:varchar(255)" json:"biaya_tarif"`
	ProdukPelayanan string `gorm:"column:produk_pelayanan;type:varchar(255)" json:"produk_pelayanan"`
	Fasilitas string `gorm:"column:fasilitas;type:text" json:"fasilitas"`
	KompetensiPelaksana string `gorm:"column:kompetensi_pelaksana;type:text" json:"kompetensi_pelaksana"`
	PengawasanInternal string `gorm:"column:pengawasan_internal;type:text" json:"pengawasan_internal"`
	JumlahPelaksana int `gorm:"column:jumlah_pelaksana" json:"jumlah_pelaksana"`
	JaminanPelayanan  string `gorm:"column:jaminan_pelayanan;type:text" json:"jaminan_pelayanan"`
	SaranDanMasukan string `gorm:"column:saran_dan_masukan;type:text" json:"saran_dan_masukan"`
	JaminanKeamanan string `gorm:"column:jaminan_keamanan;type:text" json:"jaminan_keamanan"`
	EvaluasiKinerja string `gorm:"column:evaluasi_kinerja;type:text" json:"evaluasi_kinerja"`

	// --- KOLOM STATUS & WAKTU VALIDASI STANDAR ---
	StatusValidasi  string `gorm:"column:status_validasi;not null;default:'Menunggu Validasi';type:varchar(255)" json:"status_validasi"`
	KeteranganValidasi *string `gorm:"column:keterangan_validasi;type:text" json:"keterangan_validasi"`
	TanggalValidasi *time.Time `gorm:"column:tanggal_validasi" json:"tanggal_validasi"`

	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relasi
	OPD OPD `gorm:"foreignKey:IDOPD" json:"opd"`
	ValidatorPemda *UserPemda  `gorm:"foreignKey:IDValidatorPemda" json:"validator_pemda"` // Pointer karena bisa NULL
	FormPengajuans []FormPengajuan `gorm:"foreignKey:IDJenisPelayanan" json:"-"`
}

//================================================================================
// TABEL USER OPD & USER PEMDA
//================================================================================

// UserOPD merepresentasikan pengguna dari OPD yang bertugas menginput data dan memproses.
// Tabel: user_opd (3)
type UserOPD struct {
	ID uint`gorm:"column:id_user_opd;primaryKey" json:"id_user_opd"`
	IDOPD uint  `gorm:"column:id_opd;not null" json:"id_opd"` // Foreign Key ke OPD
	Nama string `gorm:"column:nama;not null;type:varchar(255)" json:"nama"`
	NIP  string `gorm:"column:nip;unique;not null;type:varchar(255)" json:"nip"`
	Password string `gorm:"column:password;not null;type:varchar(255)" json:"-"`
	Jabatan string `gorm:"column:jabatan;type:varchar(255)" json:"jabatan"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relasi (sebagai Child dan Parent)
	OPD OPD  `gorm:"foreignKey:IDOPD" json:"opd"`
	FormPengajuans []FormPengajuan `gorm:"foreignKey:IDUserOPD" json:"-"`
	FormPemohons []FormPemohon  `gorm:"foreignKey:IDUserOPDInput" json:"-"` // Relasi ke FormPemohon
}

// UserPemda merepresentasikan pengguna dari Pemda yang bertugas sebagai validator.
// Tabel: user_pemda (4)
type UserPemda struct {
	ID  uint `gorm:"column:id_user_pemda;primaryKey" json:"id_user_pemda"`
	Nama  string `gorm:"column:nama;not null;type:varchar(255)" json:"nama"`
	NIP   string `gorm:"column:nip;unique;not null;type:varchar(255)" json:"nip"`
	Password string `gorm:"column:password;not null;type:varchar(255)" json:"-"`
	Jabatan  string `gorm:"column:jabatan;type:varchar(255)" json:"jabatan"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relasi (sebagai Parent)
	ValidatedJenisPelayanans []JenisPelayanan `gorm:"foreignKey:IDValidatorPemda" json:"-"`
}

//================================================================================
// TABEL FORM PEMOHON
//================================================================================

// FormPemohon merepresentasikan Data Master Pemohon yang diinput/didaftarkan oleh User OPD.
// Tabel: form_pemohon (6)
type FormPemohon struct {
	ID  uint `gorm:"column:id_form_pemohon;primaryKey" json:"id_form_pemohon"`
	IDUserOPDInput uint `gorm:"column:id_user_opd_input;not null" json:"id_user_opd_input"` // Petugas OPD yang menginput
	
	// --- DATA PEMOHON ---
	NamaLengkap string `gorm:"column:nama_lengkap;not null;type:varchar(255)" json:"nama_lengkap"`
	NIK  string `gorm:"column:nik;unique;not null;type:varchar(255)" json:"nik"`
	Alamat string `gorm:"column:alamat;type:text" json:"alamat"`
	NomorHP  string `gorm:"column:nomor_hp;type:varchar(255)" json:"nomor_hp"`
	Email  string `gorm:"column:email;type:varchar(255)" json:"email"`
	
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relasi
	UserOPDInput UserOPD `gorm:"foreignKey:IDUserOPDInput" json:"user_opd_input"`
}


//================================================================================
// TABEL fORM PENGAJUAN LAYANAN OPD
//================================================================================

// FormPengajuan merepresentasikan tabel transaksi utama untuk semua pengajuan layanan.
// Tabel: form_pengajuan (5)
type FormPengajuan struct {
	ID uint `gorm:"column:id_form_pengajuan;primaryKey" json:"id_form_pengajuan"`

	// Foreign Keys
	IDOPD uint `gorm:"column:id_opd;not null" json:"id_opd"`
	IDJenisPelayanan  uint `gorm:"column:id_jenis_pelayanan;not null" json:"id_jenis_pelayanan"`
	IDUserOPD uint `gorm:"column:id_user_opd;not null" json:"id_user_opd"` // Petugas OPD yang memproses

	// --- DATA PEMOHON (EKSPLISIT DALAM TRANSAKSI) ---
	NamaPemohonLengkap string `gorm:"column:nama_pemohon_lengkap;not null;type:varchar(255)" json:"nama_pemohon_lengkap"`
	NIKPemohon  string `gorm:"column:nik_pemohon;not null;type:varchar(255)" json:"nik_pemohon"`
	AlamatPemohon string `gorm:"column:alamat_pemohon;type:text" json:"alamat_pemohon"`
	NomorHPPemohon  string `gorm:"column:nomor_hp_pemohon;type:varchar(255)" json:"nomor_hp_pemohon"`
	EmailPemohon  string `gorm:"column:email_pemohon;type:varchar(255)" json:"email_pemohon"`

	// --- DETAIL PENGAJUAN ---
	JudulPengajuan string `gorm:"column:judul_pengajuan;not null;type:varchar(255)" json:"judul_pengajuan"`
	DeskripsiSingkat string `gorm:"column:deskripsi_singkat;type:text" json:"deskripsi_singkat"`
	DokumenPengajuanPath *string `gorm:"column:dokumen_pengajuan_path;type:varchar(255)" json:"dokumen_pengajuan_path"`
	IsAgreed bool  `gorm:"column:is_agreed;not null;default:false" json:"is_agreed"`

	// --- ATRIBUT TAMBAHAN TRANSAKSI ---
	PeriodeMulai *time.Time `gorm:"column:periode_mulai;type:date" json:"periode_mulai"`
	PeriodeSelesai *time.Time `gorm:"column:periode_selesai;type:date" json:"periode_selesai"`

	// --- KOLOM STATUS & WAKTU ---
	StatusProses string `gorm:"column:status_proses;not null;default:'Baru';type:varchar(255)" json:"status_proses"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relasi
	OPD OPD  `gorm:"foreignKey:IDOPD" json:"opd"`
	JenisPelayanan JenisPelayanan `gorm:"foreignKey:IDJenisPelayanan" json:"jenis_pelayanan"`
	UserOPD UserOPD `gorm:"foreignKey:IDUserOPD" json:"user_opd"`
}