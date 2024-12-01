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

func InsertDesignCategory(respw http.ResponseWriter, req *http.Request) {
	// _, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))

	// if err != nil {
	// 	_, err = watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))

	// 	if err != nil {
	// 		var respn model.Response
	// 		respn.Status = "Error: Token Tidak Valid"
	// 		respn.Info = at.GetSecretFromHeader(req)
	// 		respn.Location = "Decode Token Error"
	// 		respn.Response = err.Error()
	// 		at.WriteJSON(respw, http.StatusForbidden, respn)
	// 		return
	// 	}
	// }

	// err = req.ParseMultipartForm(10 << 20)
	// if err != nil {
	// 	var respn model.Response
	// 	respn.Status = "Error: Gagal memproses form data"
	// 	respn.Response = err.Error()
	// 	at.WriteJSON(respw, http.StatusBadRequest, respn)
	// 	return
	// }

	// file, header, err := req.FormFile("uploadReferences")
	// if err != nil {
	// 	var respn model.Response
	// 	respn.Status = "Error: Gambar tidak ditemukan"
	// 	at.WriteJSON(respw, http.StatusBadRequest, respn)
	// 	return
	// }
	// defer file.Close()

	// fileContent, err := io.ReadAll(file)
	// if err != nil {
	// 	var respn model.Response
	// 	respn.Status = "Error: Gagal membaca file"
	// 	at.WriteJSON(respw, http.StatusInternalServerError, respn)
	// 	return
	// }

	// hashedFileName := ghupload.CalculateHash(fileContent) + header.Filename[strings.LastIndex(header.Filename, "."):]
	// GitHubAccessToken := config.GHAccessToken
	// GitHubAuthorName := "Rolly Maulana Awangga"
	// GitHubAuthorEmail := "awangga@gmail.com"
	// githubOrg := "idbiz"
	// githubRepo := "img"
	// // pathFile := "uploadReferences/" + hashedFileName
	// replace := true

	// content, _, err := ghupload.GithubUpload(GitHubAccessToken, GitHubAuthorName, GitHubAuthorEmail, fileContent, githubOrg, githubRepo, pathFile, replace)
	// if err != nil {
	// 	var respn model.Response
	// 	respn.Status = "Error: Gagal mengupload gambar ke GitHub"
	// 	respn.Response = err.Error()
	// 	at.WriteJSON(respw, http.StatusInternalServerError, respn)
	// 	return
	// }

	// gambarURL := *content.Content.HTMLURL

	DesignCategory := req.FormValue("category")

	CategoryInput := model.DesignCategory{
		Category:    DesignCategory,
		
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