package controller

import (
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

// Insert Feedback
func InsertFeedback(respw http.ResponseWriter, req *http.Request) {

	Category := req.FormValue("category")
	Comments := req.FormValue("comments")
	Image := req.FormValue("image")

	FeedbackInput := model.Feedback{
		Category: model.FeedbackCategory{Category: Category},
		Comments: Comments,
		Image:    Image,
	}

	dataFeedback, err := atdb.InsertOneDoc(config.Mongoconn, "feedback", FeedbackInput)
	if err != nil {
		var respn model.Response
		respn.Status = "Gagal Insert Database"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotModified, respn)
		return
	}

	response := map[string]interface{}{
		"message": "Feedback berhasil dibuat",
		"status":  "success",
		"data":    dataFeedback,
	}

	at.WriteJSON(respw, http.StatusOK, response)
}

// Get All Feedback
func GetAllFeedback(respw http.ResponseWriter, req *http.Request) {
	data, err := atdb.GetAllDoc[[]model.Feedback](config.Mongoconn, "feedback", primitive.M{})
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Data feedback tidak ditemukan"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	var feedbacks []map[string]interface{}
	for _, feedback := range data {

		feedbacks = append(feedbacks, map[string]interface{}{
			"category":     model.FeedbackCategory{Category: feedback.Category.Category},
			"comments": feedback.Comments,
			"image":  feedback.Image,
		})
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Data Feedback berhasil diambil",
		"data":    feedbacks,
	}

	at.WriteJSON(respw, http.StatusOK, response)
}
