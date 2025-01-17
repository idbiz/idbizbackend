package controller

import (
	"encoding/json"
	// "fmt"
	// "io"
	"net/http"
	"time"
	// "strings"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/kimseokgis/backend-ai/helper"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"

	// "github.com/gocroot/helper/ghupload"
	"github.com/gocroot/helper/watoken"
	"github.com/gocroot/model"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

// Portofolio
func CreatePortofolio(respw http.ResponseWriter, req *http.Request) {
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		_, err = watoken.Decode(config.PUBLICKEY, at.GetLoginFromHeader(req))
		if err != nil {
			var respn model.Response
			respn.Status = "Error: Token Tidak Valid"
			respn.Info = at.GetSecretFromHeader(req)
			respn.Location = "Decode Token Error"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusForbidden, respn)
			return
		}
	}

	var portofolio model.Portofolio
	if err := json.NewDecoder(req.Body).Decode(&portofolio); err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	_, err = atdb.InsertOneDoc(config.Mongoconn, "portofolio", portofolio)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	at.WriteJSON(respw, http.StatusOK, model.Response{
		Response: "Portofolio berhasil disimpan",
		Info:     "",
		Status:   "Success",
		Location: "",
	})
}

func GetAllPortofolio(respw http.ResponseWriter, req *http.Request) {
	// Decode token untuk otentikasi
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		_, err = watoken.Decode(config.PUBLICKEY, at.GetLoginFromHeader(req))
		if err != nil {
			at.WriteJSON(respw, http.StatusForbidden, model.Response{
				Status:   "Error: Token Tidak Valid",
				Response: err.Error(),
			})
			return
		}
	}

	var resp itmodel.Response

	// Mengambil semua dokumen dari koleksi "portofolio"
	portofolios, err := atdb.GetAllDoc[[]model.Portofolio](config.Mongoconn, "portofolio", bson.M{})
	if err != nil {
		resp.Response = err.Error()
		helper.WriteJSON(respw, http.StatusInternalServerError, resp)
		return
	}
	helper.WriteJSON(respw, http.StatusOK, portofolios)
}

// UpdatePortofolio
func UpdatePortofolio(respw http.ResponseWriter, req *http.Request) {
	// Decode token untuk otentikasi
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		_, err = watoken.Decode(config.PUBLICKEY, at.GetLoginFromHeader(req))
		if err != nil {
			at.WriteJSON(respw, http.StatusForbidden, model.Response{
				Status:   "Error: Token Tidak Valid",
				Response: err.Error(),
			})
			return
		}
	}

	// Ambil ID dari parameter URL
	id := req.URL.Query().Get("id")

	// Decode request body ke dalam struct Portofolio
	var portofolio model.Portofolio
	if err := json.NewDecoder(req.Body).Decode(&portofolio); err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	// Update dokumen portofolio berdasarkan ID
	filter := bson.M{"_id": id}
	update := bson.M{"$set": portofolio}

	_, err = atdb.UpdateOneDoc(config.Mongoconn, "portofolio", filter, update)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	at.WriteJSON(respw, http.StatusOK, model.Response{
		Response: "Portofolio berhasil diperbarui",
		Info:     "",
		Status:   "Success",
		Location: "",
	})
}

// DeletePortofolio
func DeletePortofolio(respw http.ResponseWriter, req *http.Request) {
	// Decode token untuk otentikasi
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		_, err = watoken.Decode(config.PUBLICKEY, at.GetLoginFromHeader(req))
		if err != nil {
			at.WriteJSON(respw, http.StatusForbidden, model.Response{
				Status:   "Error: Token Tidak Valid",
				Response: err.Error(),
			})
			return
		}
	}

	// Ambil ID dari parameter URL
	id := req.URL.Query().Get("id")

	// Hapus dokumen portofolio berdasarkan ID
	filter := bson.M{"_id": id}

	_, err = atdb.DeleteOneDoc(config.Mongoconn, "portofolio", filter)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	at.WriteJSON(respw, http.StatusOK, model.Response{
		Response: "Portofolio berhasil dihapus",
		Info:     "",
		Status:   "Success",
		Location: "",
	})
}

// Pesanan
func CreatePesanan(respw http.ResponseWriter, req *http.Request) {
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		_, err = watoken.Decode(config.PUBLICKEY, at.GetLoginFromHeader(req))
		if err != nil {
			var respn model.Response
			respn.Status = "Error: Token Tidak Valid"
			respn.Info = at.GetSecretFromHeader(req)
			respn.Location = "Decode Token Error"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusForbidden, respn)
			return
		}
	}

	var pesanan model.Pesanan
	if err := json.NewDecoder(req.Body).Decode(&pesanan); err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	// Jika tanggal pesanan ada dalam body JSON, kita akan memformatnya.
	if pesanan.TanggalPesanan.IsZero() {
		// Jika tidak ada tanggal, set tanggal saat ini
		pesanan.TanggalPesanan = time.Now()
	} else {
		// Parse tanggal dari JSON dengan format yang sesuai: "DD-MM-YYYY HH:MM"
		parsedDate, err := time.Parse("02-01-2006 15:04", pesanan.TanggalPesanan.Format("02-01-2006 15:04"))
		if err != nil {
			at.WriteJSON(respw, http.StatusInternalServerError, model.Response{Status: "Error", Response: err.Error()})
			return
		}
		pesanan.TanggalPesanan = parsedDate
	}

	_, err = atdb.InsertOneDoc(config.Mongoconn, "pesanan", pesanan)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	at.WriteJSON(respw, http.StatusOK, model.Response{
		Response: "Pesanan berhasil disimpan",
		Info:     "",
		Status:   "Success",
		Location: "",
	})
}

func GetAllPesanan(respw http.ResponseWriter, req *http.Request) {
	// Decode token untuk otentikasi
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		_, err = watoken.Decode(config.PUBLICKEY, at.GetLoginFromHeader(req))
		if err != nil {
			at.WriteJSON(respw, http.StatusForbidden, model.Response{
				Status:   "Error: Token Tidak Valid",
				Response: err.Error(),
			})
			return
		}
	}

	var resp itmodel.Response

	// Mengambil semua dokumen dari koleksi "pesanan"
	pesanans, err := atdb.GetAllDoc[[]model.Pesanan](config.Mongoconn, "pesanan", bson.M{})
	if err != nil {
		resp.Response = err.Error()
		helper.WriteJSON(respw, http.StatusInternalServerError, resp)
		return
	}
	helper.WriteJSON(respw, http.StatusOK, pesanans)
}

// UpdatePesanan
func UpdatePesanan(respw http.ResponseWriter, req *http.Request) {
	// Decode token untuk otentikasi
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		_, err = watoken.Decode(config.PUBLICKEY, at.GetLoginFromHeader(req))
		if err != nil {
			at.WriteJSON(respw, http.StatusForbidden, model.Response{
				Status:   "Error: Token Tidak Valid",
				Response: err.Error(),
			})
			return
		}
	}

	// Ambil ID dari parameter URL
	id := req.URL.Query().Get("id")

	// Decode request body ke dalam struct Pesanan
	var pesanan model.Pesanan
	if err := json.NewDecoder(req.Body).Decode(&pesanan); err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	// Jika tanggal pesanan ada dalam body JSON, kita akan memformatnya.
	if pesanan.TanggalPesanan.IsZero() {
		// Jika tidak ada tanggal, set tanggal saat ini
		pesanan.TanggalPesanan = time.Now()
	} else {
		// Parse tanggal dari JSON dengan format yang sesuai: "DD-MM-YYYY HH:MM"
		parsedDate, err := time.Parse("02-01-2006 15:04", pesanan.TanggalPesanan.Format("02-01-2006 15:04"))
		if err != nil {
			at.WriteJSON(respw, http.StatusInternalServerError, model.Response{Status: "Error", Response: err.Error()})
			return
		}
		pesanan.TanggalPesanan = parsedDate
	}

	// Update pesanan berdasarkan ID
	filter := bson.M{"_id": id}
	update := bson.M{"$set": pesanan}

	_, err = atdb.UpdateOneDoc(config.Mongoconn, "pesanan", filter, update)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	at.WriteJSON(respw, http.StatusOK, model.Response{
		Response: "Pesanan berhasil diperbarui",
		Info:     "",
		Status:   "Success",
		Location: "",
	})
}

// DeletePesanan
func DeletePesanan(respw http.ResponseWriter, req *http.Request) {
	// Decode token untuk otentikasi
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		_, err = watoken.Decode(config.PUBLICKEY, at.GetLoginFromHeader(req))
		if err != nil {
			at.WriteJSON(respw, http.StatusForbidden, model.Response{
				Status:   "Error: Token Tidak Valid",
				Response: err.Error(),
			})
			return
		}
	}

	// Ambil ID dari parameter URL
	id := req.URL.Query().Get("id")

	// Hapus pesanan berdasarkan ID
	filter := bson.M{"_id": id}

	_, err = atdb.DeleteOneDoc(config.Mongoconn, "pesanan", filter)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	at.WriteJSON(respw, http.StatusOK, model.Response{
		Response: "Pesanan berhasil dihapus",
		Info:     "",
		Status:   "Success",
		Location: "",
	})
}

// Pembayaran
func CreatePembayaran(respw http.ResponseWriter, req *http.Request) {
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		_, err = watoken.Decode(config.PUBLICKEY, at.GetLoginFromHeader(req))
		if err != nil {
			var respn model.Response
			respn.Status = "Error: Token Tidak Valid"
			respn.Info = at.GetSecretFromHeader(req)
			respn.Location = "Decode Token Error"
			respn.Response = err.Error()
			at.WriteJSON(respw, http.StatusForbidden, respn)
			return
		}
	}

	var pembayaran model.Pembayaran
	if err := json.NewDecoder(req.Body).Decode(&pembayaran); err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	_, err = atdb.InsertOneDoc(config.Mongoconn, "pembayaran", pembayaran)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	at.WriteJSON(respw, http.StatusOK, model.Response{
		Response: "Pembayaran berhasil disimpan",
		Info:     "",
		Status:   "Success",
		Location: "",
	})
}

func GetAllPembayaran(respw http.ResponseWriter, req *http.Request) {
	// Decode token untuk otentikasi
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		_, err = watoken.Decode(config.PUBLICKEY, at.GetLoginFromHeader(req))
		if err != nil {
			at.WriteJSON(respw, http.StatusForbidden, model.Response{
				Status:   "Error: Token Tidak Valid",
				Response: err.Error(),
			})
			return
		}
	}

	var resp itmodel.Response

	// Mengambil semua dokumen dari koleksi "pembayaran"
	pembayarans, err := atdb.GetAllDoc[[]model.Pembayaran](config.Mongoconn, "pembayaran", bson.M{})
	if err != nil {
		resp.Response = err.Error()
		helper.WriteJSON(respw, http.StatusInternalServerError, resp)
		return
	}
	helper.WriteJSON(respw, http.StatusOK, pembayarans)
}