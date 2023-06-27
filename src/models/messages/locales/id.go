package locales

// === BASIC MESSAGES ===
var BasicID = map[string]string{
	"BASIC-0001": "Upps, terjadi kesalahan",
}

// === SUCCESS MESSAGES ===
var SuccessID = map[string]string{
	// AUTH LISTS
	"SUCCESS-0001": "Pendaftaran berhasil",
	"SUCCESS-0002": "Login berhasil",
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
}
