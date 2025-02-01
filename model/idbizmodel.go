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

type Transaksi struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID          string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	NamaPemesan     string             `json:"nama_pemesan,omitempty" bson:"nama_pemesan,omitempty"`
	DesainID        string             `json:"desain_id,omitempty" bson:"desain_id,omitempty"`
	NamaDesain      string             `json:"nama_desain,omitempty" bson:"nama_desain,omitempty"`
	Harga           string             `json:"harga,omitempty" bson:"harga,omitempty"`
	StatusPesanan   string             `json:"status_pesanan,omitempty" bson:"status_pesanan,omitempty"`
	CatatanPesanan  string             `json:"catatan_pesanan,omitempty" bson:"catatan_pesanan,omitempty"`
	BuktiPembayaran string             `json:"bukti_pembayaran,omitempty" bson:"bukti_pembayaran,omitempty"`
	TanggalPesanan  time.Time          `json:"tanggal_pesanan,omitempty" bson:"tanggal_pesanan,omitempty"`
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
