package controller

import (
	// "bytes"
	// "context"
	// "encoding/json"
	// "fmt"
	"net/http"
	// "strings"
	// "time"

	"github.com/gocroot/config"
	"github.com/gocroot/model"

	// "github.com/whatsauth/itmodel"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gocroot/helper/at"
	// "github.com/gocroot/helper/atapi"
	"github.com/gocroot/helper/atdb"
	// "github.com/gocroot/helper/gcallapi"
	// "github.com/gocroot/helper/lms"
	// "github.com/gocroot/helper/report"
	// "github.com/gocroot/helper/watoken"
	// "github.com/gocroot/helper/whatsauth"
)

// Create a new portofolio
func CreatePortofolio(respw http.ResponseWriter, req *http.Request) {
	// var portofolio model.Portofolio
	// var respn model.Response
	// payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(r))
	// if err != nil {
	// 	respn.Status = "Error : Token Tidak Valid"
	// 	respn.Info = at.GetSecretFromHeader(r)
	// 	respn.Location = "Decode Token Error"
	// 	respn.Response = err.Error()
	// 	at.WriteJSON(w, http.StatusForbidden, respn)
	// 	return
	// }

	// //check eksistensi user
	// docuser, err := atdb.GetOneDoc[model.Userdomyikado](config.Mongoconn, "user", primitive.M{"phonenumber": payload.Id})
	// if err != nil {
	// 	docuser.PhoneNumber = payload.Id
	// 	docuser.Name = payload.Alias
	// 	at.WriteJSON(w, http.StatusNotFound, docuser)
	// 	return
	// }

	Category := req.FormValue("category")
	DesignTitle := req.FormValue("design_title")
	DesignDesc := req.FormValue("design_desc")
	DesignImage := req.FormValue("design_image")

	PortofolioInput := model.Portofolio{
		Category:    model.DesignCategory{},
		DesignTitle: DesignTitle,
		DesignDesc:  DesignDesc,
		DesignImage: DesignImage,
	}

	// _, err := atdb.InsertOneDoc(config.Mongoconn, "portofolio", newPortofolio)
	// if err != nil {
	// 	resp := model.Response{
	// 		Status:   "Error : Gagal insert ke portofolio",
	// 		Response: err.Error(),
	// 	}
	// 	at.WriteJSON(w, http.StatusNotFound, resp)
	// 	return
	// }

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
