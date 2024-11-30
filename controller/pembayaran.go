package controller

import (
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"

	// "github.com/gocroot/helper/watoken"
	// "github.com/gocroot/helper/ghupload"
	"github.com/gocroot/model"
)

func CreatePembayaran(respw http.ResponseWriter, req *http.Request) {
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

	DesignSelected := req.FormValue("design_selected")
	OrderDescription := req.FormValue("order_description")
	CardFullname := req.FormValue("card_fullname")
	CardNumber := req.FormValue("card_number")
	CardExpiration := req.FormValue("card_expication")
	CVV := req.FormValue("cvv")
	Price := req.FormValue("price")

	PembayaranInput := model.Pemesanan{
		DesignSelected:   DesignSelected,
		OrderDescription: OrderDescription,
		CardFullname:     CardFullname,
		CardNumber:       CardNumber,
		CardExpiration:   CardExpiration,
		CVV:              CVV,
		Price:            Price,
	}

	dataPembayaran, err := atdb.InsertOneDoc(config.Mongoconn, "pembayaran", PembayaranInput)
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
		"data":    dataPembayaran,
	}

	at.WriteJSON(respw, http.StatusOK, response)
}
