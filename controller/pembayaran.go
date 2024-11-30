package controller

import (
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

func CreatePembayaran(respw http.ResponseWriter, req *http.Request) {
	// Parse form data
	err := req.ParseForm()
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal memproses form data"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Ambil data dari form
	DesignSelected := req.FormValue("design_selected")
	OrderDescription := req.FormValue("order_description")
	CardFullname := req.FormValue("card_fullname")
	CardNumber := req.FormValue("card_number")
	CardExpiration := req.FormValue("card_expiration")
	CVV := req.FormValue("cvv")
	Price := req.FormValue("price")

	// Validasi data input
	if DesignSelected == "" || OrderDescription == "" || CardFullname == "" || CardNumber == "" || CardExpiration == "" || CVV == "" || Price == "" {
		var respn model.Response
		respn.Status = "Error: Semua field wajib diisi"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	// Buat objek untuk disimpan
	PembayaranInput := model.Pembayaran{
		DesignSelected:   DesignSelected,
		OrderDescription: OrderDescription,
		CardFullname:     CardFullname,
		CardNumber:       CardNumber,
		CardExpiration:   CardExpiration,
		CVV:              CVV,
		Price:            Price,
	}

	// Masukkan ke database
	dataPembayaran, err := atdb.InsertOneDoc(config.Mongoconn, "pembayaran", PembayaranInput)
	if err != nil {
		var respn model.Response
		respn.Status = "Gagal Insert Database"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	// Respons sukses
	response := map[string]interface{}{
		"status":  "success",
		"message": "Transaksi pembayaran berhasil",
		"data":    dataPembayaran,
	}

	at.WriteJSON(respw, http.StatusOK, response)
}

// Get Pemesanan By Id
func GetPembayaranById(respw http.ResponseWriter, req *http.Request) {
	pembayaranID := req.URL.Query().Get("id")
	if pembayaranID == "" {
		var respn model.Response
		respn.Status = "Error: ID transaksi tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(pembayaranID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID transaksi tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	filter := bson.M{"_id": objectID}
	dataPembayaran, err := atdb.GetOneDoc[model.Pembayaran](config.Mongoconn, "pembayaran", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: transaksi tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	data := model.Pembayaran{
		ID:               dataPembayaran.ID,
		DesignSelected:   dataPembayaran.DesignSelected,
		OrderDescription: dataPembayaran.OrderDescription,
		CardFullname:     dataPembayaran.CardFullname,
		CardNumber:       dataPembayaran.CardNumber,
		CardExpiration:   dataPembayaran.CardExpiration,
		CVV:              dataPembayaran.CVV,
		Price:            dataPembayaran.Price,
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Transaksi ditemukan",
		"data":    data,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}
