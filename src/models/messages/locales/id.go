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
	"SUCCESS-DB-0002": "Sukses mengisi data ke database",

	//	BLOG LISTS
	"SUCCESS-BLOG-0001": "Blog berhasil dibuat",
	"SUCCESS-BLOG-0002": "Blog berhasil diubah",
	"SUCCESS-BLOG-0003": "Blog berhasil dihapus",
	"SUCCESS-BLOG-0004": "Blog berhasil di publikasi",
	"SUCCESS-BLOG-0005": "Slug blog berhasil diubah",
	"SUCCESS-BLOG-0006": "Blog berhasil diubah ke draft",
	"SUCCESS-BLOG-0007": "Berhasil mengambil data blog",
	"SUCCESS-BLOG-0008": "Berhasil mencari blog",
}

// === ERROR MESSAGES ===
var ErrorID = map[string]string{
	// ERROR 400 LISTS
	"ERROR-400001": "Pengguna telah diaktifkan",
	"ERROR-400002": "Anda telah menyelesaikan profil Anda",
	"ERROR-400003": "Harap isi form dengan benar",
	"ERROR-400004": "Haraf tambahkan sampul blog",

	// ERROR 401 LISTS
	"ERROR-401001": "Email sudah terdaftar",
	"ERROR-401002": "Anda belum terdaftar",
	"ERROR-401003": "Kredensial tidak cocok",

	// ERROR 403 LIST
	"ERROR-403001": "Tidak dapat mengedit blog, status invalid",

	// ERROR 404 LIST
	"ERROR-404001": "Data tidak ditemukan",

	//	ERROR 500 LISTS
	"ERROR-50001":       "Gagal migrasi infrastructures",
	"ERROR-50002":       "Migrate key tidak valid",
	"ERROR-50003":       "Upps, terjadi kesalahan",
	"ERROR-IMAGE-50001": "Ukuran gambar melebihi %f megabyte",
	"ERROR-IMAGE-50002": "Format gambar tidak valid",
	"ERROR-IMAGE-50003": "Ukuran gambar harus antara %dx%d dan %dx%d",
}
