package controller

import (
	"io"
	"net/http"
	"strings"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/ghupload"
	"github.com/gocroot/helper/watoken"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

func InsertPemesanan(respw http.ResponseWriter, req *http.Request) {
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))

	if err != nil {
		var respn model.Response
		respn.Status = "Error: Token Tidak Valid"
		respn.Info = config.PublicKeyWhatsAuth
		respn.Location = "Decode Token Error"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusForbidden, respn)
		return

	}

	err = req.ParseMultipartForm(10 << 20)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal memproses form data"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	file, header, err := req.FormFile("upload_references")
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gambar tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal membaca file"
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	hashedFileName := ghupload.CalculateHash(fileContent) + header.Filename[strings.LastIndex(header.Filename, "."):]
	GitHubAccessToken := config.GHAccessToken
	GitHubAuthorName := "Rolly Maulana Awangga"
	GitHubAuthorEmail := "awangga@gmail.com"
	githubOrg := "idbiz"
	githubRepo := "img"
	pathFile := "pemesanan/" + hashedFileName
	replace := true

	content, _, err := ghupload.GithubUpload(GitHubAccessToken, GitHubAuthorName, GitHubAuthorEmail, fileContent, githubOrg, githubRepo, pathFile, replace)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal mengupload gambar ke GitHub"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	upload_references := *content.Content.HTMLURL

	Fullname := req.FormValue("fullname")
	Email := req.FormValue("email")
	PhoneNumber := req.FormValue("phone_number")
	Category := req.FormValue("category_id")
	OrderDescription := req.FormValue("order_description")

	objectCategoryID, err := primitive.ObjectIDFromHex(Category)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Kategori ID tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}
	categoryDoc, err := atdb.GetOneDoc[model.DesignCategory](config.Mongoconn, "design-category", primitive.M{"_id": objectCategoryID})
	if err != nil || categoryDoc.ID == primitive.NilObjectID {
		var respn model.Response
		respn.Status = "Error: Kategori tidak ditemukan"
		respn.Response = "ID yang dicari: " + Category
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	PemesananInput := model.Pemesanan{
		Fullname:         Fullname,
		Email:            Email,
		PhoneNumber:      PhoneNumber,
		Category:         categoryDoc,
		OrderDescription: OrderDescription,
		UploadReferences: upload_references,
	}

	dataPemesanan, err := atdb.InsertOneDoc(config.Mongoconn, "pemesanan", PemesananInput)
	if err != nil {
		var respn model.Response
		respn.Status = "Gagal Insert Database"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotModified, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Pemesanan berhasil ditambahkan",
		"data":    dataPemesanan,
	}

	at.WriteJSON(respw, http.StatusOK, response)
}

// Get Pemesanan By Id
func GetPemesananById(respw http.ResponseWriter, req *http.Request) {
	pemesananID := req.URL.Query().Get("id")
	if pemesananID == "" {
		var respn model.Response
		respn.Status = "Error: ID pemesanan tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(pemesananID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID pemesanan tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	filter := bson.M{"_id": objectID}
	dataPemesanan, err := atdb.GetOneDoc[model.Pemesanan](config.Mongoconn, "pemesanan", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Pemesanan tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	data := model.Pemesanan{
		ID:          dataPemesanan.ID,
		Fullname:    dataPemesanan.Fullname,
		Email:       dataPemesanan.Email,
		PhoneNumber: dataPemesanan.PhoneNumber,
		// Category:         model.DesignCategory{Category: dataPemesanan.Category.Category},
		Category:         dataPemesanan.Category,
		OrderDescription: dataPemesanan.OrderDescription,
		UploadReferences: dataPemesanan.UploadReferences,
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Pemesanan ditemukan",
		"data":    data,
	}
	at.WriteJSON(respw, http.StatusOK, response)
}

// Get All Pemesanan
func GetAllPemesanan(respw http.ResponseWriter, req *http.Request) {
	data, err := atdb.GetAllDoc[[]model.Pemesanan](config.Mongoconn, "pemesanan", primitive.M{})
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Data pemesanan tidak ditemukan"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	var pemesanans []map[string]interface{}
	for _, pemesanan := range data {

		pemesanans = append(pemesanans, map[string]interface{}{
			"fullname":          pemesanan.Fullname,
			"email":             pemesanan.Email,
			"phone_number":      pemesanan.PhoneNumber,
			"category":          pemesanan.Category,
			"order_description": pemesanan.OrderDescription,
			"upload_references": pemesanan.UploadReferences,
		})
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Data pemesanan berhasil diambil",
		"data":    pemesanans,
	}

	at.WriteJSON(respw, http.StatusOK, response)
}

func DeleteDataPemesanan(respw http.ResponseWriter, req *http.Request) {

	pemesananId := req.URL.Query().Get("pemesananId")
	if pemesananId == "" {
		var respn model.Response
		respn.Status = "Error: ID Pemesanan tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(pemesananId)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Pemesanan tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	filter := bson.M{"_id": objectID}
	deleteData, err := atdb.DeleteOneDoc(config.Mongoconn, "pemesanan", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal menghapus data pemesanan"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Pemesanan berhasil dihapus",
		"data": map[string]interface{}{
			// "user":    payload.Id,
			"pemesanan_id": objectID.Hex(),
			"deleted":      deleteData.DeletedCount,
		},
	}
	at.WriteJSON(respw, http.StatusOK, response)
}
