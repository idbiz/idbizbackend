package controller

import (
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"
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
	PembayaranInput := model.Pemesanan{
		// DesignSelected:   DesignSelected,
		// OrderDescription: OrderDescription,
		// CardFullname:     CardFullname,
		// CardNumber:       CardNumber,
		// CardExpiration: CardExpiration,
		// CVV:              CVV,
		// Price:            Price,
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
		"message": "Pemesanan berhasil ditambahkan",
		"data":    dataPembayaran,
	}

	at.WriteJSON(respw, http.StatusOK, response)
}
