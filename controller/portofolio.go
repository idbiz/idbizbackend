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

// Insert Portfolio
func InsertPortofolio(respw http.ResponseWriter, req *http.Request) {

	Category := req.FormValue("category")
	DesignTitle := req.FormValue("design_title")
	DesignDesc := req.FormValue("design_desc")
	DesignImage := req.FormValue("design_image")

	PortofolioInput := model.Portofolio{
		Category:    model.DesignCategory{Category: Category},
		DesignTitle: DesignTitle,
		DesignDesc:  DesignDesc,
		DesignImage: DesignImage,
	}

	dataPortofolio, err := atdb.InsertOneDoc(config.Mongoconn, "portofolio", PortofolioInput)
	if err != nil {
		var respn model.Response
		respn.Status = "Gagal Insert Database"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotModified, respn)
		return
	}

	response := map[string]interface{}{
		"message": "Portofolio berhasil dibuat",
		"status":  "success",
		"data":    dataPortofolio,
	}

	at.WriteJSON(respw, http.StatusOK, response)
}

// Get Portofolio By Id
func GetPortofolioById(respw http.ResponseWriter, req *http.Request) {
	portofolioID := req.URL.Query().Get("id")
	if portofolioID == "" {
		var respn model.Response
		respn.Status = "Error: ID portfolio tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(portofolioID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID portfolio tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	filter := bson.M{"_id": objectID}
	dataPortofolio, err := atdb.GetOneDoc[model.Portofolio](config.Mongoconn, "portfolio", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Portfolio tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	data := model.Portofolio{
		ID: dataPortofolio.ID,
		// Category:    dataPortofolio.Category,
		Category:    model.DesignCategory{Category: dataPortofolio.Category.Category},
		DesignTitle: dataPortofolio.DesignTitle,
		DesignDesc:  dataPortofolio.DesignDesc,
		DesignImage: dataPortofolio.DesignImage,
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Portopolio ditemukan",
		"data":    data,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

// Get All Portofolio
func GetAllPortofolio(respw http.ResponseWriter, req *http.Request) {
	data, err := atdb.GetAllDoc[[]model.Portofolio](config.Mongoconn, "portofolio", primitive.M{})
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Data portofolio tidak ditemukan"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	var portofolios []map[string]interface{}
	for _, portofolio := range data {

		portofolios = append(portofolios, map[string]interface{}{
			"fullname":     model.DesignCategory{Category: portofolio.Category.Category},
			"design_title": portofolio.DesignTitle,
			"design_desc":  portofolio.DesignDesc,
			"design_image": portofolio.DesignImage,
		})
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Data Portofolio berhasil diambil",
		"data":    portofolios,
	}

	at.WriteJSON(respw, http.StatusOK, response)
}
