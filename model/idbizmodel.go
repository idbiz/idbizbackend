package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Portofolio struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	NamaDesain string             `json:"nama_desain,omitempty" bson:"nama_desain,omitempty"`
	Deskripsi  string             `json:"deskripsi,omitempty" bson:"deskripsi,omitempty"`
	Gambar     string             `json:"gambar,omitempty" bson:"gambar,omitempty"`
	Kategori   string             `json:"kategori,omitempty" bson:"kategori,omitempty"`
	Harga      string             `json:"harga,omitempty" bson:"harga,omitempty"`
}

type Pesanan struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	NamaPemesan    string             `json:"nama_pemesan,omitempty" bson:"nama_pemesan,omitempty"`
	DaftarDesain   []Portofolio       `json:"daftar_desain,omitempty" bson:"daftar_desain,omitempty"`
	TanggalPesanan time.Time          `json:"tanggal_pesanan,omitempty" bson:"tanggal_pesanan,omitempty"`
	StatusPesanan  string             `json:"status_pesanan,omitempty" bson:"status_pesanan,omitempty"`
	Pembayaran     string             `json:"pembayaran,omitempty" bson:"pembayaran,omitempty"`
	CatatanPesanan string             `json:"catatan_pesanan,omitempty" bson:"catatan_pesanan,omitempty"`
	TotalHarga     string             `json:"total_harga,omitempty" bson:"total_harga,omitempty"`
}

type ItemPesanan struct {
	PortofolioID string `json:"portofolio_id,omitempty" bson:"portofolio_id,omitempty"`
	NamaDesain   string `json:"nama_desain,omitempty" bson:"nama_desain,omitempty"`
	Jumlah       int    `json:"jumlah,omitempty" bson:"jumlah,omitempty"`
	HargaSatuan  string `json:"harga_satuan,omitempty" bson:"harga_satuan,omitempty"`
	SubTotal     string `json:"sub_total,omitempty" bson:"sub_total,omitempty"`
}

type Pembayaran struct {
	ID                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	PortfolioID       string             `json:"portfolio_id,omitempty" bson:"portfolio_id,omitempty"`
	NamaDesain        string             `json:"nama_desain,omitempty" bson:"nama_desain,omitempty"`
	Harga             string             `json:"harga,omitempty" bson:"harga,omitempty"`
	UserID            string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	NamaPemesan       string             `json:"nama_pemesan,omitempty" bson:"nama_pemesan,omitempty"`
	MetodePembayaran  string             `json:"metode_pembayaran,omitempty" bson:"metode_pembayaran,omitempty"` // Contoh: Transfer Bank, E-Wallet, dll.
	StatusPembayaran  string             `json:"status_pembayaran,omitempty" bson:"status_pembayaran,omitempty"` // Contoh: Pending, Berhasil, Gagal.
	TotalBayar        string             `json:"total_bayar,omitempty" bson:"total_bayar,omitempty"`
	TanggalPembayaran time.Time          `json:"tanggal_pembayaran,omitempty" bson:"tanggal_pembayaran,omitempty"`
	BuktiPembayaran   string             `json:"bukti_pembayaran,omitempty" bson:"bukti_pembayaran,omitempty"` // URL atau path ke file bukti pembayaran.
}

type UploadRequest struct {
	UserID          string `json:"user_id,omitempty" bson:"user_id,omitempty"`
	PembayarnID     string `json:"pembayaran_id,omitempty" bson:"pembayaran_id,omitempty"`
	BuktiPembayaran string `json:"bukti_pembayaran,omitempty" bson:"bukti_pembayaran,omitempty"`
	FileName        string `json:"file_name,omitempty" bson:"file_name,omitempty"`
}

type GithubUploadRequest struct {
	Message string `json:"message,omitempty" bson:"message,omitempty"`
	Content string `json:"content,omitempty" bson:"content,omitempty"`
}
