package controller

import (
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	// "github.com/gocroot/helper/watoken"
	// "github.com/gocroot/helper/ghupload"
	"github.com/gocroot/model"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"

)

func InsertFeedbackCategory(respw http.ResponseWriter, req *http.Request) {
	
	FeedbackCategory := req.FormValue("category")

	CategoryInput := model.DesignCategory{
		Category:    FeedbackCategory,
		
	}

	dataCategory, err := atdb.InsertOneDoc(config.Mongoconn, "feedback-category", CategoryInput)
	if err != nil {
		var respn model.Response
		respn.Status = "Gagal Insert Database"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotModified, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Data Category berhasil ditambahkan",
		"data":    dataCategory,
	}

	at.WriteJSON(respw, http.StatusOK, response)
}