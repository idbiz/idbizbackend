package controller

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/ghupload"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Insert Portfolio
func InsertPortofolio(respw http.ResponseWriter, req *http.Request) {
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
	GitHubAuthorName := "GhaidaFasya"
	GitHubAuthorEmail := "ghaidafasya5@gmail.com"
	githubOrg := "idbiz-img"
	githubRepo := "img"
	pathFile := "pemesanan/" + hashedFileName
	replace := true
	content, _, err := ghupload.GithubUpload(GitHubAccessToken, GitHubAuthorName, GitHubAuthorEmail, fileContent, githubOrg, githubRepo, pathFile, replace)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal mengupload gambar ke GitHub"
		respn.Response = err.Error()
		fmt.Println(err.Error())
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}
	upload_references := *content.Content.HTMLURL

	Category := req.FormValue("category_id")
	DesignTitle := req.FormValue("design_title")
	DesignDesc := req.FormValue("design_desc")
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

	PortofolioInput := model.Portofolio{
		Category:    categoryDoc,
		DesignTitle: DesignTitle,
		DesignDesc:  DesignDesc,
		DesignImage: upload_references,
	}

	fmt.Println(PortofolioInput)
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
	dataPortofolio, err := atdb.GetOneDoc[model.Portofolio](config.Mongoconn, "portofolio", filter)
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
		"message": "Portofolio ditemukan",
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
			"category":     model.DesignCategory{Category: portofolio.Category.Category},
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

func DeleteDataPortofolio(respw http.ResponseWriter, req *http.Request) {

	portofolioId := req.URL.Query().Get("portofolioId")
	if portofolioId == "" {
		var respn model.Response
		respn.Status = "Error: ID Portofolio tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(portofolioId)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID Portofolio tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	filter := bson.M{"_id": objectID}
	deleteData, err := atdb.DeleteOneDoc(config.Mongoconn, "portofolio", filter)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal menghapus data portofolio"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Portofolio berhasil dihapus",
		"data": map[string]interface{}{
			// "user":    payload.Id,
			"portofolio_id": objectID.Hex(),
			"deleted":       deleteData.DeletedCount,
		},
	}
	at.WriteJSON(respw, http.StatusOK, response)
}
