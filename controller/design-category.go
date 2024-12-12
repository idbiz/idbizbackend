package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	// "github.com/gocroot/helper/watoken"
	// "github.com/gocroot/helper/ghupload"
	"github.com/gocroot/model"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

// Insert Design Category
func InsertDesignCategory(respw http.ResponseWriter, req *http.Request) {

	DesignCategory := req.FormValue("category")

	CategoryInput := model.DesignCategory{
		Category: DesignCategory,
	}

	dataCategory, err := atdb.InsertOneDoc(config.Mongoconn, "design-category", CategoryInput)
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

// Get DesignCategory By Id
func GetDesignCategoryById(respw http.ResponseWriter, req *http.Request) {
	categoryID := req.URL.Query().Get("id")
	if categoryID == "" {
		var respn model.Response
		respn.Status = "Error: ID category tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID category tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	filter := bson.M{"_id": objectID}
	dataCategory, err := atdb.GetOneDoc[model.DesignCategory](config.Mongoconn, "design-category", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Category tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	data := model.DesignCategory{
		ID:       dataCategory.ID,
		Category: dataCategory.Category,
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Category ditemukan",
		"data":    data,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

// Get All DesignCategory
func GetAllDesignCategory(respw http.ResponseWriter, req *http.Request) {
	data, err := atdb.GetAllDoc[[]model.DesignCategory](config.Mongoconn, "design-category", primitive.M{})
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Data category tidak ditemukan"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	var categories []map[string]interface{}
	for _, category := range data {

		categories = append(categories, map[string]interface{}{
			"id":       category.ID,
			"category": category.Category,
		})
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Data design category berhasil diambil",
		"data":    categories,
	}

	at.WriteJSON(respw, http.StatusOK, response)
}

// Update Design Category
func UpdateDesignCategory(respw http.ResponseWriter, req *http.Request) {
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))

	if err != nil {
		payload, err = watoken.Decode(config.PUBLICKEY, at.GetLoginFromHeader(req))

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

	categoryID := req.URL.Query().Get("id")
	if categoryID == "" {
		var respn model.Response
		respn.Status = "Error: ID Category tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Category tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	filter := bson.M{"_id": objectID}
	_, err = atdb.GetOneDoc[model.DesignCategory](config.Mongoconn, "design-category", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Category tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	var requestBody struct {
		// Icon         string `json:"icon"`
		NameCategory string `json:"category"`
	}
	err = json.NewDecoder(req.Body).Decode(&requestBody)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal membaca data JSON"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	updateData := bson.M{}
	if requestBody.NameCategory != "" {
		updateData["category"] = requestBody.NameCategory
	}
	// if requestBody.Icon != "" {
	// 	updateData["icon"] = requestBody.Icon
	// }

	update := bson.M{"$set": updateData}
	_, err = atdb.UpdateOneDoc(config.Mongoconn, "design-category", filter, update)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal mengupdate category"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotModified, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "category berhasil diupdate",
		"data":    updateData,
		"name":    payload.Alias,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

func DeleteDataDesignCategory(respw http.ResponseWriter, req *http.Request) {
	// payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	// if err != nil {
	// 	payload, err = watoken.Decode(config.PUBLICKEY, at.GetLoginFromHeader(req))
	// 	if err != nil {
	// 		var respn model.Response
	// 		respn.Status = "Error: Token Tidak Valid"
	// 		respn.Response = err.Error()
	// 		at.WriteJSON(respw, http.StatusForbidden, respn)
	// 		return
	// 	}
	// }

	categoryId := req.URL.Query().Get("categoryId")
	if categoryId == "" {
		var respn model.Response
		respn.Status = "Error: ID Category tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(categoryId)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Pemesanan tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	filter := bson.M{"_id": objectID}
	deleteData, err := atdb.DeleteOneDoc(config.Mongoconn, "design-category", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal menghapus data category"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Data category berhasil dihapus",
		"data": map[string]interface{}{
			// "user":    payload.Id,
			"category_id": objectID.Hex(),
			"deleted":     deleteData.DeletedCount,
		},
	}
	at.WriteJSON(respw, http.StatusOK, response)
}
