package locales

// === SUCCESS MESSAGES ===
var SuccessID = map[string]string{
	// BASIC LISTS
	"SUCCESS-BASIC-0001": "Mantepp Bro",

	// AUTH LISTS
	"SUCCESS-AUTH-0001": "Pendaftaran berhasil",
	"SUCCESS-AUTH-0002": "Login berhasil",

	//	DATABASE LISTS
	"SUCCESS-DB-0001": "Sukses migrasi database",

	//	BLOG LISTS
	"SUCCESS-BLOG-0001": "Blog berhasil dibuat",
}

// === ERROR MESSAGES ===
var ErrorID = map[string]string{
	// ERROR 400 LISTS
	"ERROR-400001": "Pengguna telah diaktifkan",
	"ERROR-400002": "Anda telah menyelesaikan profil Anda",
	"ERROR-400003": "Harap isi form dengan benar",

	// ERROR 401 LISTS
	"ERROR-401001": "Email sudah terdaftar",
	"ERROR-401002": "Anda belum terdaftar",
	"ERROR-401003": "Kredensial tidak cocok",

	//	ERROR 500 LISTS
	"ERROR-50001": "Gagal migrasi database",
	"ERROR-50002": "Migrate key tidak valid",
	"ERROR-50003": "Upps, terjadi kesalahan",
}
