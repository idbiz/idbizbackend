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

// Insert Pembayaran
func InsertPembayaran(respw http.ResponseWriter, req *http.Request) {

	OrderDescription := req.FormValue("order_description")
	CardFullname := req.FormValue("card_fullname")
	CardNumber := req.FormValue("card_number")
	CardExpiration := req.FormValue("card_expiration")
	CVV := req.FormValue("cvv")
	Price := req.FormValue("price")

	PembayaranInput := model.Pembayaran{
		OrderDescription: model.Pemesanan{OrderDescription: OrderDescription},
		CardFullname:     CardFullname,
		CardNumber:       CardNumber,
		CardExpiration:   CardExpiration,
		CVV:              CVV,
		Price:            Price,
	}

	dataPembayaran, err := atdb.InsertOneDoc(config.Mongoconn, "pembayaran", PembayaranInput)
	if err != nil {
		var respn model.Response
		respn.Status = "Gagal Insert Database"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotModified, respn)
		return
	}

	response := map[string]interface{}{
		"message": "Transaksi pembayaran berhasil dibuat",
		"status":  "success",
		"data":    dataPembayaran,
	}

	at.WriteJSON(respw, http.StatusOK, response)
}

// Get Pembayaran By Id
func GetPembayaranById(respw http.ResponseWriter, req *http.Request) {
	pembayaranID := req.URL.Query().Get("id")
	if pembayaranID == "" {
		var respn model.Response
		respn.Status = "Error: ID transaksi pembayaran tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(pembayaranID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID transaksi pembayaran tidak valid"
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
		OrderDescription: dataPembayaran.OrderDescription,
		CardFullname:     dataPembayaran.CardFullname,
		CardNumber:       dataPembayaran.CardNumber,
		CardExpiration:   dataPembayaran.CardExpiration,
		CVV:              dataPembayaran.CVV,
		Price:            dataPembayaran.Price,
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Transaksi pembayaran ditemukan",
		"data":    data,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

// Get All Pembayaran
func GetAllPembayaran(respw http.ResponseWriter, req *http.Request) {
	data, err := atdb.GetAllDoc[[]model.Pembayaran](config.Mongoconn, "pembayaran", primitive.M{})
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Data transaksi pembayaran tidak ditemukan"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	var pembayarans []map[string]interface{}
	for _, pembayaran := range data {

		pembayarans = append(pembayarans, map[string]interface{}{
			"order_description": model.Pemesanan{OrderDescription: pembayaran.OrderDescription.OrderDescription},
			"card_fullname":     pembayaran.CardFullname,
			"card_number":       pembayaran.CardNumber,
			"card_expiration":   pembayaran.CardExpiration,
			"cvv":               pembayaran.CVV,
			"price":             pembayaran.Price,
		})
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Data transaksi pembayaran berhasil diambil",
		"data":    pembayarans,
	}

	at.WriteJSON(respw, http.StatusOK, response)
}
