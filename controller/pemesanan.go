package controller

import (
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

func CreatePemesanan(respw http.ResponseWriter, req *http.Request) {
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

	file, header, err := req.FormFile("uploadReferences")
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
	pathFile := "uploadReferences/" + hashedFileName
	replace := true

	content, _, err := ghupload.GithubUpload(GitHubAccessToken, GitHubAuthorName, GitHubAuthorEmail, fileContent, githubOrg, githubRepo, pathFile, replace)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal mengupload gambar ke GitHub"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusInternalServerError, respn)
		return
	}

	gambarURL := *content.Content.HTMLURL

	Fullname := req.FormValue("fullname")
	Email := req.FormValue("email")
	PhoneNumber := req.FormValue("phone_number")
	DesignType := req.FormValue("design_type")
	OrderDescription := req.FormValue("order_description")

	PemesananInput := model.Pemesanan{
		Fullname:         Fullname,
		Email:            Email,
		PhoneNumber:      PhoneNumber,
		UploadReferences: gambarURL,
		DesignType:       DesignType,
		OrderDescription: OrderDescription,
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
		ID:               dataPemesanan.ID,
		Fullname:         dataPemesanan.Fullname,
		Email:            dataPemesanan.Email,
		PhoneNumber:      dataPemesanan.PhoneNumber,
		DesignType:       dataPemesanan.DesignType,
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
		// imageUrl := strings.Replace(pemesanan.UploadReferences, "github.com", "raw.githubusercontent.com", 1)
		// imageUrls := strings.Replace(imageUrl, "/blob/", "/", 1)

		// finalPrice := menu.Price
		// diskonValue := 0.00
		// potonganHarga := 0.00

		// if menu.Diskon != nil && menu.Diskon.Status == "Active" {
		// 	if menu.Diskon.JenisDiskon == "Persentase" {
		// 		diskonAmount := float64(menu.Price) * (float64(menu.Diskon.NilaiDiskon) / 100)
		// 		finalPrice = menu.Price - int(diskonAmount)
		// 		diskonValue = float64(menu.Diskon.NilaiDiskon)
		// 		potonganHarga = diskonAmount
		// 	} else if menu.Diskon.JenisDiskon == "Nominal" {
		// 		finalPrice = menu.Price - menu.Diskon.NilaiDiskon
		// 		if finalPrice < 0 {
		// 			finalPrice = 0
		// 		}
		// 		diskonValue = float64(menu.Diskon.NilaiDiskon)
		// 		potonganHarga = float64(menu.Diskon.NilaiDiskon)
		// 	}
		// }

		pemesanans = append(pemesanans, map[string]interface{}{
			"fullname":          pemesanan.Fullname,
			"email":             pemesanan.Email,
			"phone_number":      pemesanan.PhoneNumber,
			"design_type":       pemesanan.DesignType,
			"order_description": pemesanan.OrderDescription,
			// "image":             imageUrls,
		})
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Data pemesanan berhasil diambil",
		"data":    pemesanans,
	}

	at.WriteJSON(respw, http.StatusOK, response)
}
