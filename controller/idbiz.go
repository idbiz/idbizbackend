package controller

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
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
	log.Println("[INFO] Memulai proses UpdatePortofolio")

	// Decode token untuk otentikasi
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		_, err = watoken.Decode(config.PUBLICKEY, at.GetLoginFromHeader(req))
		if err != nil {
			log.Println("[ERROR] Token tidak valid")
			at.WriteJSON(respw, http.StatusForbidden, model.Response{
				Status:   "Error: Token Tidak Valid",
				Response: err.Error(),
			})
			return
		}
	}

	// Ambil ID dari URL Path
	id := req.URL.Query().Get("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("[ERROR] ID tidak valid sebagai ObjectID: %s", err.Error())
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status:   "Error",
			Response: "ID tidak valid",
		})
		return
	}

	// Decode request body ke dalam struct Portofolio
	var portofolio struct {
		NamaDesain string `json:"nama_desain"`
		Deskripsi  string `json:"deskripsi"`
		Gambar     string `json:"gambar"`
		Kategori   string `json:"kategori"`
		Harga      string `json:"harga"`
	}
	if err := json.NewDecoder(req.Body).Decode(&portofolio); err != nil {
		log.Printf("[ERROR] Gagal mendekode request body: %s", err.Error())
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status:   "Error",
			Response: err.Error(),
		})
		return
	}

	// Hapus field `_id` dari update
	updateData := bson.M{
		"nama_desain": portofolio.NamaDesain,
		"deskripsi":   portofolio.Deskripsi,
		"gambar":      portofolio.Gambar,
		"kategori":    portofolio.Kategori,
		"harga":       portofolio.Harga,
	}

	// update := bson.M{"$set": updateData}

	// Update dokumen portofolio berdasarkan ID
	collection := config.Mongoconn.Collection("portofolio")
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"$set": updateData})
	if err != nil {
		log.Printf("[ERROR] Gagal memperbarui portofolio: %s", err.Error())
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{
			Status:   "Error",
			Response: err.Error(),
		})
		return
	}

	at.WriteJSON(respw, http.StatusOK, model.Response{
		Response: "Portofolio berhasil diperbarui",
		Status:   "Success",
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
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("[ERROR] ID tidak valid sebagai ObjectID: %s", err.Error())
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status:   "Error",
			Response: "ID tidak valid",
		})
		return
	}

	// Hapus dokumen portofolio berdasarkan ID
	filter := bson.M{"_id": objID}

	_, err = atdb.DeleteOneDoc(config.Mongoconn, "portofolio", filter)
	if err != nil {
		log.Printf("[ERROR] Gagal menghapus portofolio: %s", err.Error())
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{
			Status:   "Error",
			Response: err.Error(),
		})
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
	log.Println("[INFO] Memulai proses GetPortofolioByID")

	// Decode token untuk otentikasi
	log.Println("[INFO] Melakukan decoding token untuk otentikasi")
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		log.Println("[WARNING] Token pertama tidak valid, mencoba token kedua")
		_, err = watoken.Decode(config.PUBLICKEY, at.GetLoginFromHeader(req))
		if err != nil {
			log.Println("[ERROR] Token tidak valid")
			at.WriteJSON(respw, http.StatusForbidden, model.Response{
				Status:   "Error: Token Tidak Valid",
				Response: err.Error(),
			})
			return
		}
	}

	// Ambil ID dari parameter URL
	id := req.URL.Query().Get("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("[ERROR] ID tidak valid sebagai ObjectID: %s", err.Error())
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status:   "Error",
			Response: "ID tidak valid",
		})
		return
	}

	filter := bson.M{"_id": objID}

	portofolio, err := atdb.GetOneDoc[model.Portofolio](config.Mongoconn, "portofolio", filter)
	if err != nil {
		log.Printf("[ERROR] Gagal mengambil portofolio: %s", err.Error())
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{
			Status:   "Error",
			Response: err.Error(),
		})
		return
	}

	log.Println("[INFO] Berhasil mengambil portofolio, mengirim respons")
	at.WriteJSON(respw, http.StatusOK, portofolio)
}

func UpdateTransaksi(respw http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] Memulai proses UpdateTransaksi")

	// Decode token untuk otentikasi
	_, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		_, err = watoken.Decode(config.PUBLICKEY, at.GetLoginFromHeader(req))
		if err != nil {
			log.Println("[ERROR] Token tidak valid")
			at.WriteJSON(respw, http.StatusForbidden, model.Response{
				Status:   "Error: Token Tidak Valid",
				Response: err.Error(),
			})
			return
		}
	}

	// Ambil ID dari URL Path
	id := req.URL.Query().Get("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("[ERROR] ID tidak valid sebagai ObjectID: %s", err.Error())
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status:   "Error",
			Response: "ID tidak valid",
		})
		return
	}

	// Decode request body ke dalam struct Transaksi
	var transaksi struct {
		UserID          string `json:"user_id"`
		NamaPemesan     string `json:"nama_pemesan"`
		DesainID        string `json:"desain_id"`
		NamaDesain      string `json:"nama_desain"`
		Harga           string `json:"harga"`
		StatusPesanan   string `json:"status_pesanan"`
		CatatanPesanan  string `json:"catatan_pesanan"`
		BuktiPembayaran string `json:"bukti_pembayaran"`
		TanggalPesanan  string `json:"tanggal_pesanan"`
	}

	if err := json.NewDecoder(req.Body).Decode(&transaksi); err != nil {
		log.Printf("[ERROR] Gagal mendekode request body: %s", err.Error())
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status:   "Error",
			Response: err.Error(),
		})
		return
	}
	
	parsedTime, err := time.Parse(time.RFC3339, transaksi.TanggalPesanan)
	if err != nil {
		log.Printf("[ERROR] Format tanggal_pesanan tidak valid: %s", err.Error())
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status:   "Error",
			Response: "Format tanggal_pesanan tidak valid, gunakan format RFC3339",
		})
		return
	}

	// Hapus field _id dari update
	updateData := bson.M{
		"user_id":          transaksi.UserID,
		"nama_pemesan":     transaksi.NamaPemesan,
		"desain_id":        transaksi.DesainID,
		"nama_desain":      transaksi.NamaDesain,
		"harga":            transaksi.Harga,
		"status_pesanan":   transaksi.StatusPesanan,
		"catatan_pesanan":  transaksi.CatatanPesanan,
		"bukti_pembayaran": transaksi.BuktiPembayaran,
		"tanggal_pesanan":  parsedTime,
	}

	if transaksi.TanggalPesanan != "" {
		updateData["tanggal_pesanan"] = parsedTime
	}

	// Update dokumen transaksi berdasarkan ID
	collection := config.Mongoconn.Collection("transaksi")
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"$set": updateData})
	if err != nil {
		log.Printf("[ERROR] Gagal memperbarui transaksi: %s", err.Error())
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{
			Status:   "Error",
			Response: err.Error(),
		})
		return
	}

	log.Println("[INFO] Transaksi berhasil diperbarui")
	at.WriteJSON(respw, http.StatusOK, model.Response{
		Response: "Transaksi berhasil diperbarui",
		Status:   "Success",
	})
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

func GetTransaksiByStatus(respw http.ResponseWriter, req *http.Request) {
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
	
	status := req.URL.Query().Get("status")
	if status == "" {
		respondWithError(respw, http.StatusBadRequest, "Status pesanan harus disertakan")
		return
	}

	// Filter transaksi berdasarkan status
	filter := bson.M{"status_pesanan": status}

	var transaksi []model.Transaksi
	transaksi, err = atdb.GetFilteredDocs[[]model.Transaksi](config.Mongoconn, "transaksi", filter, nil)
	if err != nil || len(transaksi) == 0 {
		if err == mongo.ErrNoDocuments || len(transaksi) == 0 {
			respondWithError(respw, http.StatusNotFound, "Tidak ada transaksi dengan status tersebut")
		} else {
			respondWithError(respw, http.StatusInternalServerError, fmt.Sprintf("Terjadi kesalahan: %v", err))
		}
		return
	}

	// Response data transaksi
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(http.StatusOK)
	json.NewEncoder(respw).Encode(map[string]interface{}{
		"status": "success",
		"data":   transaksi,
	})
}

func DeleteTransaksi(respw http.ResponseWriter, req *http.Request) {
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
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("[ERROR] ID tidak valid sebagai ObjectID: %s", err.Error())
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status:   "Error",
			Response: "ID tidak valid",
		})
		return
	}

	// Hapus dokumen portofolio berdasarkan ID
	filter := bson.M{"_id": objID}

	_, err = atdb.DeleteOneDoc(config.Mongoconn, "transaksi", filter)
	if err != nil {
		log.Printf("[ERROR] Gagal menghapus portofolio: %s", err.Error())
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{
			Status:   "Error",
			Response: err.Error(),
		})
		return
	}

	at.WriteJSON(respw, http.StatusOK, model.Response{
		Response: "Transaksi berhasil dihapus",
		Info:     "",
		Status:   "Success",
		Location: "",
	})
}


func GetAllTransaksi(respw http.ResponseWriter, req *http.Request) {
	var resp itmodel.Response
	orders, err := atdb.GetAllDoc[[]model.Transaksi](config.Mongoconn, "transaksi", bson.M{})
	if err != nil {
		resp.Response = err.Error()
		helper.WriteJSON(respw, http.StatusBadRequest, resp)
		return
	}
	helper.WriteJSON(respw, http.StatusOK, orders)
}

// UploadHandler handles file uploads to GitHub
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		at.WriteJSON(w, http.StatusMethodNotAllowed, model.Response{
			Status:   "Error",
			Response: "Only POST method is allowed",
		})
		return
	}

	// Parse file from request
	file, header, err := r.FormFile("file")
	if err != nil {
		at.WriteJSON(w, http.StatusBadRequest, model.Response{
			Status:   "Error",
			Response: "Failed to read file from request: " + err.Error(),
		})
		return
	}
	defer file.Close()

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, model.Response{
			Status:   "Error",
			Response: "Failed to read file content: " + err.Error(),
		})
		return
	}

	// Generate a unique ID
	uniqueID := primitive.NewObjectID().Hex()

	// Append unique ID to filename (e.g., "receipt.png" -> "receipt_<uniqueID>.png")
	originalFileName := header.Filename
	extIndex := len(originalFileName) - len(filepath.Ext(originalFileName))
	newFileName := originalFileName[:extIndex] + "_" + uniqueID + filepath.Ext(originalFileName)

	// Encode file content to Base64 (required by GitHub API)
	encodedContent := base64.StdEncoding.EncodeToString(fileBytes)

	// Get MongoDB connection from config
	db := config.Mongoconn

	// Upload to GitHub using credentials from MongoDB
	_, err = utils.UploadToGithub(newFileName, encodedContent, db)
	if err != nil {
		at.WriteJSON(w, http.StatusInternalServerError, model.Response{
			Status:   "Error",
			Response: fmt.Sprintf("Upload failed: %v", err),
		})
		return
	}

	// Respond with success and file URL
	at.WriteJSON(w, http.StatusOK, model.Response{
		Response: "File uploaded successfully",
		Info:     "",
		Status:   "Success",
		Location: "",
	})
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

	// Generate a unique ID for the filename
	uniqueID := primitive.NewObjectID().Hex()

	// Append unique ID to filename
	originalFileName := header.Filename
	extIndex := len(originalFileName) - len(filepath.Ext(originalFileName))
	newFileName := originalFileName[:extIndex] + "_" + uniqueID + filepath.Ext(originalFileName)

	// Encode file content to Base64 (needed for GitHub API)
	encodedContent := base64.StdEncoding.EncodeToString(fileBytes)

	// Get MongoDB connection from config
	db := config.Mongoconn

	// Upload payment proof to GitHub
	buktiPembayaranURL, err := utils.UploadToGithub(newFileName, encodedContent, db)
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
	_, err = atdb.InsertOneDoc(db, "transaksi", transaksi)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{
			Status:   "Error",
			Response: "Failed to save transaction: " + err.Error(),
		})
		return
	}

	// Respond with success
	response := map[string]interface{}{
		"status":        "Success",
		"message":       "Transaksi berhasil disimpan",
		"data":          transaksi,
	}

	at.WriteJSON(respw, http.StatusOK, response)
}
