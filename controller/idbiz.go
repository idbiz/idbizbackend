package controller

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atdb"

	// "github.com/joho/godotenv"
	"github.com/kimseokgis/backend-ai/helper"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"

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

// Pesanan
func CreatePesanan(respw http.ResponseWriter, req *http.Request) {
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

	var pesanan model.Pesanan
	if err := json.NewDecoder(req.Body).Decode(&pesanan); err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	// Jika tanggal pesanan ada dalam body JSON, kita akan memformatnya.
	if pesanan.TanggalPesanan.IsZero() {
		// Jika tidak ada tanggal, set tanggal saat ini
		pesanan.TanggalPesanan = time.Now()
	} else {
		// Parse tanggal dari JSON dengan format yang sesuai: "DD-MM-YYYY HH:MM"
		parsedDate, err := time.Parse("02-01-2006 15:04", pesanan.TanggalPesanan.Format("02-01-2006 15:04"))
		if err != nil {
			at.WriteJSON(respw, http.StatusInternalServerError, model.Response{Status: "Error", Response: err.Error()})
			return
		}
		pesanan.TanggalPesanan = parsedDate
	}

	_, err = atdb.InsertOneDoc(config.Mongoconn, "pesanan", pesanan)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	at.WriteJSON(respw, http.StatusOK, model.Response{
		Response: "Pesanan berhasil disimpan",
		Info:     "",
		Status:   "Success",
		Location: "",
	})
}

func GetAllPesanan(respw http.ResponseWriter, req *http.Request) {
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

	// Mengambil semua dokumen dari koleksi "pesanan"
	pesanans, err := atdb.GetAllDoc[[]model.Pesanan](config.Mongoconn, "pesanan", bson.M{})
	if err != nil {
		resp.Response = err.Error()
		helper.WriteJSON(respw, http.StatusInternalServerError, resp)
		return
	}
	helper.WriteJSON(respw, http.StatusOK, pesanans)
}

// UpdatePesanan
func UpdatePesanan(respw http.ResponseWriter, req *http.Request) {
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

	// Decode request body ke dalam struct Pesanan
	var pesanan model.Pesanan
	if err := json.NewDecoder(req.Body).Decode(&pesanan); err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	// Jika tanggal pesanan ada dalam body JSON, kita akan memformatnya.
	if pesanan.TanggalPesanan.IsZero() {
		// Jika tidak ada tanggal, set tanggal saat ini
		pesanan.TanggalPesanan = time.Now()
	} else {
		// Parse tanggal dari JSON dengan format yang sesuai: "DD-MM-YYYY HH:MM"
		parsedDate, err := time.Parse("02-01-2006 15:04", pesanan.TanggalPesanan.Format("02-01-2006 15:04"))
		if err != nil {
			at.WriteJSON(respw, http.StatusInternalServerError, model.Response{Status: "Error", Response: err.Error()})
			return
		}
		pesanan.TanggalPesanan = parsedDate
	}

	// Update pesanan berdasarkan ID
	filter := bson.M{"_id": id}
	update := bson.M{"$set": pesanan}

	_, err = atdb.UpdateOneDoc(config.Mongoconn, "pesanan", filter, update)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	at.WriteJSON(respw, http.StatusOK, model.Response{
		Response: "Pesanan berhasil diperbarui",
		Info:     "",
		Status:   "Success",
		Location: "",
	})
}

// DeletePesanan
func DeletePesanan(respw http.ResponseWriter, req *http.Request) {
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

	// Hapus pesanan berdasarkan ID
	filter := bson.M{"_id": id}

	_, err = atdb.DeleteOneDoc(config.Mongoconn, "pesanan", filter)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	at.WriteJSON(respw, http.StatusOK, model.Response{
		Response: "Pesanan berhasil dihapus",
		Info:     "",
		Status:   "Success",
		Location: "",
	})
}

func GetPesananByID(respw http.ResponseWriter, req *http.Request) {
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

	// Ambil pesanan berdasarkan ID
	filter := bson.M{"_id": id}

	pesanan, err := atdb.GetOneDoc[model.Pesanan](config.Mongoconn, "pesanan", filter)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	at.WriteJSON(respw, http.StatusOK, pesanan)
}

// Pembayaran
func CreatePembayaran(respw http.ResponseWriter, req *http.Request) {
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

	var pembayaran model.Pembayaran
	if err := json.NewDecoder(req.Body).Decode(&pembayaran); err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	_, err = atdb.InsertOneDoc(config.Mongoconn, "pembayaran", pembayaran)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	at.WriteJSON(respw, http.StatusOK, model.Response{
		Response: "Pembayaran berhasil disimpan",
		Info:     "",
		Status:   "Success",
		Location: "",
	})
}

func GetAllPembayaran(respw http.ResponseWriter, req *http.Request) {
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

	// Mengambil semua dokumen dari koleksi "pembayaran"
	pembayarans, err := atdb.GetAllDoc[[]model.Pembayaran](config.Mongoconn, "pembayaran", bson.M{})
	if err != nil {
		resp.Response = err.Error()
		helper.WriteJSON(respw, http.StatusInternalServerError, resp)
		return
	}
	helper.WriteJSON(respw, http.StatusOK, pembayarans)
}

// ItemPesanan
func CreateItemPesanan(respw http.ResponseWriter, req *http.Request) {
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

	var itemPesanan model.ItemPesanan
	if err := json.NewDecoder(req.Body).Decode(&itemPesanan); err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	_, err = atdb.InsertOneDoc(config.Mongoconn, "item_pesanan", itemPesanan)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	at.WriteJSON(respw, http.StatusOK, model.Response{
		Response: "Item Pesanan berhasil disimpan",
		Info:     "",
		Status:   "Success",
		Location: "",
	})
}

func GetAllItemPesanan(respw http.ResponseWriter, req *http.Request) {
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

	// Mengambil semua dokumen dari koleksi "item_pesanan"
	itemPesanans, err := atdb.GetAllDoc[[]model.ItemPesanan](config.Mongoconn, "item_pesanan", bson.M{})
	if err != nil {
		resp.Response = err.Error()
		helper.WriteJSON(respw, http.StatusInternalServerError, resp)
		return
	}
	helper.WriteJSON(respw, http.StatusOK, itemPesanans)
}

// Update Pesan Status
func UpdatePesanStatus(respw http.ResponseWriter, req *http.Request) {
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

	// Decode request body ke dalam struct Pesanan
	var pesanan model.Pesanan
	if err := json.NewDecoder(req.Body).Decode(&pesanan); err != nil {
		at.WriteJSON(respw, http.StatusBadRequest, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	// Update pesanan berdasarkan ID
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status_pesanan": pesanan.StatusPesanan}}

	_, err = atdb.UpdateOneDoc(config.Mongoconn, "pesanan", filter, update)
	if err != nil {
		at.WriteJSON(respw, http.StatusInternalServerError, model.Response{Status: "Error", Response: err.Error()})
		return
	}

	at.WriteJSON(respw, http.StatusOK, model.Response{
		Response: "Status Pesanan berhasil diperbarui",
		Info:     "",
		Status:   "Success",
		Location: "",
	})
}

// UploadtoGithub
func UploadtoGithub(respw http.ResponseWriter, req *http.Request) {
	// Baca file dari form-data
	req.ParseMultipartForm(10 << 20) // Maksimal ukuran file 10 MB
	file, fileHeader, err := req.FormFile("file")
	if err != nil {
		http.Error(respw, "Gagal membaca file dari form-data: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Baca konten file
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(respw, "Gagal membaca konten file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode konten file ke Base64
	encodedContent := base64.StdEncoding.EncodeToString(fileBytes)

	// Pesan commit dari form-data (opsional)
	commitMessage := req.FormValue("message")
	if commitMessage == "" {
		commitMessage = "Menambahkan file baru"
	}

	// Nama file di GitHub
	fileName := fileHeader.Filename

	// URL endpoint untuk GitHub API
	githubAPIURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s",
		"dibiz",     // Ganti dengan username GitHub Anda
		"upload",    // Ganti dengan nama repository GitHub Anda
		fileName,    // Nama file yang diupload
	)

	// Siapkan payload untuk request
	payload := map[string]string{
		"message": commitMessage,
		"content": encodedContent,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(respw, "Gagal mempersiapkan payload: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Baca token dari environment variable
	token := os.Getenv("GH_ACCESS_TOKEN")
	if token == "" {
		http.Error(respw, "Token GitHub tidak ditemukan di environment variable", http.StatusInternalServerError)
		return
	}

	// Buat request HTTP PUT ke GitHub
	client := &http.Client{}
	request, err := http.NewRequest("PUT", githubAPIURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(respw, "Gagal membuat request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Tambahkan header Authorization dengan token
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")

	// Kirim request
	response, err := client.Do(request)
	if err != nil {
		http.Error(respw, "Gagal mengunggah file ke GitHub: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// Baca respons dari GitHub
	body, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(respw, "Gagal membaca respons dari GitHub: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Periksa status code GitHub API
	if response.StatusCode != http.StatusCreated {
		http.Error(respw, fmt.Sprintf("Gagal mengunggah file: %s\n%s", response.Status, string(body)), response.StatusCode)
		return
	}

	// Berikan respons sukses
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(http.StatusOK)
	respw.Write([]byte(`{"status":"success", "message":"File berhasil diunggah ke GitHub"}`))
}