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
			"category": model.FeedbackCategory{Category: feedback.Category.Category},
			"comments": feedback.Comments,
			"image":    feedback.Image,
		})
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Data Feedback berhasil diambil",
		"data":    feedbacks,
	}

	at.WriteJSON(respw, http.StatusOK, response)
}

// Get Feedback By Id
func GetFeedbackById(respw http.ResponseWriter, req *http.Request) {
	feedbackID := req.URL.Query().Get("id")
	if feedbackID == "" {
		var respn model.Response
		respn.Status = "Error: ID feedback tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(feedbackID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Feedback tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	filter := bson.M{"_id": objectID}
	dataFeedback, err := atdb.GetOneDoc[model.Feedback](config.Mongoconn, "feedback", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Feedback tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	data := model.Feedback{
		ID: dataFeedback.ID,
		// Category:    dataPortofolio.Category,
		Category: model.FeedbackCategory{Category: dataFeedback.Category.Category},
		Comments: dataFeedback.Comments,
		Image:    dataFeedback.Image,
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Feedback ditemukan",
		"data":    data,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}
