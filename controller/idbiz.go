package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/utils"

	// "github.com/joho/godotenv"
	"github.com/kimseokgis/backend-ai/helper"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gocroot/helper/watoken"
	"github.com/gocroot/model"
)

// Portofolio
func CreatePortofolio(respw http.ResponseWriter, req *http.Request) {
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		_, err = watoken.Decode(config.PUBLICKEY, at.GetLoginFromHeader(req))
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

	var portofolio model.Portofolio
	if err := json.NewDecoder(req.Body).Decode(&portofolio); err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	_, err = atdb.InsertOneDoc(config.Mongoconn, "portofolio", portofolio)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	at.WriteJSON(respw, http.StatusOK, model.Response{
		Response: "Portofolio berhasil disimpan",
		Info:     "",
		Status:   "Success",
		Location: "",
	})
}

func GetAllPortofolio(respw http.ResponseWriter, req *http.Request) {
	// Decode token untuk otentikasi
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		_, err = watoken.Decode(config.PUBLICKEY, at.GetLoginFromHeader(req))
		if err != nil {
			at.WriteJSON(respw, http.StatusForbidden, model.Response{
				Status:   "Error: Token Tidak Valid",
				Response: err.Error(),
			})
			return
		}
	}

	var resp itmodel.Response

	// Mengambil semua dokumen dari koleksi "portofolio"
	portofolios, err := atdb.GetAllDoc[[]model.Portofolio](config.Mongoconn, "portofolio", bson.M{})
	if err != nil {
		resp.Response = err.Error()
		helper.WriteJSON(respw, http.StatusInternalServerError, resp)
		return
	}
	helper.WriteJSON(respw, http.StatusOK, portofolios)
}

// UpdatePortofolio
func UpdatePortofolio(respw http.ResponseWriter, req *http.Request) {
	// Decode token untuk otentikasi
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		_, err = watoken.Decode(config.PUBLICKEY, at.GetLoginFromHeader(req))
		if err != nil {
			at.WriteJSON(respw, http.StatusForbidden, model.Response{
				Status:   "Error: Token Tidak Valid",
				Response: err.Error(),
			})
			return
		}
	}

	// Ambil ID dari parameter URL
	id := req.URL.Query().Get("id")

	// Decode request body ke dalam struct Portofolio
	var portofolio model.Portofolio
	if err := json.NewDecoder(req.Body).Decode(&portofolio); err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	// Update dokumen portofolio berdasarkan ID
	filter := bson.M{"_id": id}
	update := bson.M{"$set": portofolio}

	_, err = atdb.UpdateOneDoc(config.Mongoconn, "portofolio", filter, update)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	at.WriteJSON(respw, http.StatusOK, model.Response{
		Response: "Portofolio berhasil diperbarui",
		Info:     "",
		Status:   "Success",
		Location: "",
	})
}

// DeletePortofolio
func DeletePortofolio(respw http.ResponseWriter, req *http.Request) {
	// Decode token untuk otentikasi
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		_, err = watoken.Decode(config.PUBLICKEY, at.GetLoginFromHeader(req))
		if err != nil {
			at.WriteJSON(respw, http.StatusForbidden, model.Response{
				Status:   "Error: Token Tidak Valid",
				Response: err.Error(),
			})
			return
		}
	}

	// Ambil ID dari parameter URL
	id := req.URL.Query().Get("id")

	// Hapus dokumen portofolio berdasarkan ID
	filter := bson.M{"_id": id}

	_, err = atdb.DeleteOneDoc(config.Mongoconn, "portofolio", filter)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	at.WriteJSON(respw, http.StatusOK, model.Response{
		Response: "Portofolio berhasil dihapus",
		Info:     "",
		Status:   "Success",
		Location: "",
	})
}

func GetPortofolioByID(respw http.ResponseWriter, req *http.Request) {
	// Decode token untuk otentikasi
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		_, err = watoken.Decode(config.PUBLICKEY, at.GetLoginFromHeader(req))
		if err != nil {
			at.WriteJSON(respw, http.StatusForbidden, model.Response{
				Status:   "Error: Token Tidak Valid",
				Response: err.Error(),
			})
			return
		}
	}

	// Ambil ID dari parameter URL
	id := req.URL.Query().Get("id")

	// Ambil dokumen portofolio berdasarkan ID
	filter := bson.M{"_id": id}

	portofolio, err := atdb.GetOneDoc[model.Portofolio](config.Mongoconn, "portofolio", filter)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	at.WriteJSON(respw, http.StatusOK, portofolio)
}

func GetPortofolioByKategori(respw http.ResponseWriter, req *http.Request) {
	// Decode token untuk otentikasi
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		_, err = watoken.Decode(config.PUBLICKEY, at.GetLoginFromHeader(req))
		if err != nil {
			at.WriteJSON(respw, http.StatusForbidden, model.Response{
				Status:   "Error: Token Tidak Valid",
				Response: err.Error(),
			})
			return
		}
	}

	// Ambil kategori dari parameter URL
	kategori := req.URL.Query().Get("kategori")

	// Ambil dokumen portofolio berdasarkan kategori
	filter := bson.M{"kategori": kategori}

	portofolios, err := atdb.GetAllDoc[[]model.Portofolio](config.Mongoconn, "portofolio", filter)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	at.WriteJSON(respw, http.StatusOK, portofolios)
}

func respondWithError(respw http.ResponseWriter, code int, message string) {
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(code)
	json.NewEncoder(respw).Encode(map[string]string{"error": message})
}

func isValidObjectID(id string) bool {
    if len(id) != 24 {
        return false
    }
    _, err := primitive.ObjectIDFromHex(id)
    return err == nil
}

func GetTransaksiByID(respw http.ResponseWriter, req *http.Request) {
    transaksiID := req.URL.Query().Get("id")
    if transaksiID == "" {
        respondWithError(respw, http.StatusBadRequest, "Pesanan ID harus disertakan")
        return
    }

    // Validasi ID apakah valid ObjectID
    if !isValidObjectID(transaksiID) {
        respondWithError(respw, http.StatusBadRequest, "Pesanan ID tidak valid")
        return
    }

    // Konversi ID menjadi ObjectID MongoDB
    objID, err := primitive.ObjectIDFromHex(transaksiID)
    if err != nil {
        respondWithError(respw, http.StatusBadRequest, "Pesanan ID tidak valid")
        return
    }

    // Filter berdasarkan ID
    filter := bson.M{"_id": objID}
    var pesanan []model.Transaksi
    pesanan, err = atdb.GetFilteredDocs[[]model.Transaksi](config.Mongoconn, "transaksi", filter, nil)
    if err != nil || len(pesanan) == 0 {
        if err == mongo.ErrNoDocuments || len(pesanan) == 0 {
            respondWithError(respw, http.StatusNotFound, "Pesanan tidak ditemukan")
        } else {
            respondWithError(respw, http.StatusInternalServerError, fmt.Sprintf("Terjadi kesalahan: %v", err))
        }
        return
    }

    // Response data pesanan
    respw.Header().Set("Content-Type", "application/json")
    respw.WriteHeader(http.StatusOK)
    json.NewEncoder(respw).Encode(map[string]interface{}{
        "status": "success",
        "data":   pesanan[0],
    })
}

// UploadHandler handles file uploads to GitHub
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse file from request
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file content
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file content", http.StatusInternalServerError)
		return
	}

	// Encode file to Base64 (required by GitHub API)
	encodedContent := base64.StdEncoding.EncodeToString(fileBytes)

	// Upload to GitHub and get public URL
	fileURL, err := utils.UploadToGithub(header.Filename, encodedContent)
	if err != nil {
		http.Error(w, fmt.Sprintf("Upload failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success and file URL
	response := map[string]string{
		"status":  "success",
		"message": "File uploaded successfully",
		"file":    header.Filename,
		"url":     fileURL, // New: Return the public GitHub URL
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// TransaksiHandler handles the creation of a new transaction
func TransaksiHandler(respw http.ResponseWriter, req *http.Request) {

	// Ensure request is POST
	if req.Method != http.MethodPost {
		at.WriteJSON(respw, http.StatusMethodNotAllowed, model.Response{
			Status:   "Error",
			Response: "Only POST method is allowed",
		})
		return
	}

	// Parse form data (limit file upload size to 10MB)
	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status:   "Error",
			Response: "Failed to parse form data: " + err.Error(),
		})
		return
	}

	// Extract form values
	userID := req.FormValue("user_id")
	namaPemesan := req.FormValue("nama_pemesan")
	desainID := req.FormValue("desain_id")
	namaDesain := req.FormValue("nama_desain")
	harga := req.FormValue("harga")
	statusPesanan := req.FormValue("status_pesanan")
	catatanPesanan := req.FormValue("catatan_pesanan")

	// Validate required fields
	if userID == "" || namaPemesan == "" || desainID == "" || namaDesain == "" || harga == "" || statusPesanan == "" {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status:   "Error",
			Response: "Missing required fields",
		})
		return
	}

	// Get file from request
	file, header, err := req.FormFile("bukti_pembayaran")
	if err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status:   "Error",
			Response: "Failed to read payment proof: " + err.Error(),
		})
		return
	}
	defer file.Close()

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{
			Status:   "Error",
			Response: "Failed to read file content: " + err.Error(),
		})
		return
	}

	// Encode file content to Base64 (needed for GitHub API)
	encodedContent := base64.StdEncoding.EncodeToString(fileBytes)

	// Upload payment proof to GitHub
	buktiPembayaranURL, err := utils.UploadToGithub(header.Filename, encodedContent)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{
			Status:   "Error",
			Response: fmt.Sprintf("Failed to upload payment proof: %v", err),
		})
		return
	}

	// Create transaction object
	transaksi := model.Transaksi{
		ID:              primitive.NewObjectID(),
		UserID:          userID,
		NamaPemesan:     namaPemesan,
		DesainID:        desainID,
		NamaDesain:      namaDesain,
		Harga:           harga,
		StatusPesanan:   statusPesanan,
		CatatanPesanan:  catatanPesanan,
		BuktiPembayaran: buktiPembayaranURL,
		TanggalPesanan:  time.Now(),
	}

	// Insert into MongoDB
	_, err = atdb.InsertOneDoc(config.Mongoconn, "transaksi", transaksi)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{
			Status:   "Error",
			Response: "Failed to save transaction: " + err.Error(),
		})
		return
	}

	// Respond with success
	at.WriteJSON(respw, http.StatusOK, model.Response{
		Status:   "Success",
		Response: "Transaction created successfully",
		Location: "",
		Info:     fmt.Sprintf("%+v", transaksi),
	})
}
