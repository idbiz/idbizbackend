package controller

import (
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
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

// Get All Pemesanan
func GetAllDesignCategory(respw http.ResponseWriter, req *http.Request) {
	data, err := atdb.GetAllDoc[[]model.DesignCategory](config.Mongoconn, "design-category", primitive.M{})
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Data category tidak ditemukan"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	var categorys []map[string]interface{}
	for _, category := range data {

		categorys = append(categorys, map[string]interface{}{
			"category": category.Category,
		})
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Data design category berhasil diambil",
		"data":    categorys,
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
