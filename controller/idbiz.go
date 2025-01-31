package controller

import (
	"bytes"
	"context"
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
	"github.com/joho/godotenv"

	// "github.com/joho/godotenv"
	"github.com/kimseokgis/backend-ai/helper"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

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

func GetPembayaranByID(respw http.ResponseWriter, req *http.Request) {
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

	// Ambil ID dari parameter query
	id := req.URL.Query().Get("id")

	// Konversi ID ke ObjectID MongoDB
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, model.Response{
			Status:   "Error: ID tidak valid",
			Response: err.Error(),
		})
		return
	}

	var pembayaran model.Pembayaran
	filter := bson.M{"_id": objID}

	// Mengambil satu dokumen dari koleksi "pembayaran"
	pembayaran, err = atdb.GetOneDoc[model.Pembayaran](config.Mongoconn, "pembayaran", filter)
	if err != nil {
		helper.WriteJSON(respw, http.StatusNotFound, model.Response{
			Status:   "Error: Pembayaran tidak ditemukan",
			Response: err.Error(),
		})
		return
	}

	helper.WriteJSON(respw, http.StatusOK, pembayaran)
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


// Check if file exists on GitHub
func getFileSha(apiURL, token string) (string, error) {
    req, err := http.NewRequest("GET", apiURL, nil)
    if err != nil {
        return "", err
    }
    req.Header.Set("Authorization", "Bearer "+token)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusNotFound {
        // File doesn't exist
        return "", nil
    } else if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return "", fmt.Errorf("GitHub API error: %s\n%s", resp.Status, string(body))
    }

    var fileData map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&fileData)
    if err != nil {
        return "", err
    }

    sha, ok := fileData["sha"].(string)
    if !ok {
        return "", fmt.Errorf("failed to extract sha from response")
    }

    return sha, nil
}

// UploadtoGithub uploads a file to GitHub repository
func UploadtoGithub(respw http.ResponseWriter, req *http.Request) {
    // Load environment variables
    err := godotenv.Load("../.env")
    if err != nil {
        http.Error(respw, "Failed to load .env file: "+err.Error(), http.StatusInternalServerError)
        return
    }

    githubOwner := os.Getenv("GITHUB_OWNER")
    githubRepo := os.Getenv("GITHUB_REPO")
    githubToken := os.Getenv("GITHUB_TOKEN")

    // Validate environment variables
    if githubOwner == "" || githubRepo == "" || githubToken == "" {
        http.Error(respw, "Missing environment variables (GITHUB_OWNER, GITHUB_REPO, or GITHUB_TOKEN)", http.StatusInternalServerError)
        return
    }

    // Parse multipart form data
    req.ParseMultipartForm(10 << 20) // Max file size 10 MB
    file, fileHeader, err := req.FormFile("file")
    if err != nil {
        http.Error(respw, "Failed to read file from form-data: "+err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Read file content
    fileBytes, err := io.ReadAll(file)
    if err != nil {
        http.Error(respw, "Failed to read file content: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Encode file content to Base64
    encodedContent := base64.StdEncoding.EncodeToString(fileBytes)

    // Generate commit message
    commitMessage := req.FormValue("message")
    if commitMessage == "" {
        commitMessage = "Upload file: " + fileHeader.Filename
    }

    // Build API URL
    filePath := "upload/" + fileHeader.Filename
    apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", githubOwner, githubRepo, filePath)

    // Check if file exists
    sha, err := getFileSha(apiURL, githubToken)
    if err != nil {
        http.Error(respw, "Failed to check file existence: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Prepare payload
    payload := map[string]string{
        "message": commitMessage,
        "content": encodedContent,
    }
    if sha != "" {
        payload["sha"] = sha // Include sha if file exists
    }
    jsonData, err := json.Marshal(payload)
    if err != nil {
        http.Error(respw, "Failed to marshal payload: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Create HTTP PUT request
    reqGithub, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer(jsonData))
    if err != nil {
        http.Error(respw, "Failed to create HTTP request: "+err.Error(), http.StatusInternalServerError)
        return
    }
    reqGithub.Header.Set("Authorization", "Bearer "+githubToken)
    reqGithub.Header.Set("Content-Type", "application/json")

    // Execute the request
    client := &http.Client{}
    resp, err := client.Do(reqGithub)
    if err != nil {
        http.Error(respw, "Failed to send request to GitHub: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // Read response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        http.Error(respw, "Failed to read response body: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Check response status
    if resp.StatusCode != http.StatusCreated {
        http.Error(respw, fmt.Sprintf("GitHub API error: %s\n%s", resp.Status, string(body)), resp.StatusCode)
        return
    }

    // Respond with success
    respw.Header().Set("Content-Type", "application/json")
    respw.WriteHeader(http.StatusOK)
    respw.Write([]byte(`{"status":"success", "message":"File successfully uploaded to GitHub"}`))
}

func PostPesanan(respw http.ResponseWriter, req *http.Request) {
	// Load environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		http.Error(respw, "Failed to load .env file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	githubOwner := os.Getenv("GITHUB_OWNER")
	githubRepo := os.Getenv("GITHUB_REPO")
	githubToken := os.Getenv("GITHUB_TOKEN")

	// Validate environment variables
	if githubOwner == "" || githubRepo == "" || githubToken == "" {
		http.Error(respw, "Missing environment variables (GITHUB_OWNER, GITHUB_REPO, or GITHUB_TOKEN)", http.StatusInternalServerError)
		return
	}

	// Parse multipart form data
	req.ParseMultipartForm(10 << 20) // Max file size 10 MB
	file, fileHeader, err := req.FormFile("file")
	if err != nil {
		http.Error(respw, "Failed to read file from form-data: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(respw, "Failed to read file content: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode file content to Base64
	encodedContent := base64.StdEncoding.EncodeToString(fileBytes)

	// Generate file upload path
	filePath := "bukti_pembayaran/" + fileHeader.Filename

	// Build API URL
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", githubOwner, githubRepo, filePath)

	// Prepare payload for GitHub
	payload := map[string]string{
		"message": "Upload bukti pembayaran: " + fileHeader.Filename,
		"content": encodedContent,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		http.Error(respw, "Failed to marshal payload: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create HTTP PUT request for GitHub
	reqGithub, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(respw, "Failed to create HTTP request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	reqGithub.Header.Set("Authorization", "Bearer "+githubToken)
	reqGithub.Header.Set("Content-Type", "application/json")

	// Send request to GitHub
	client := &http.Client{}
	respGithub, err := client.Do(reqGithub)
	if err != nil {
		http.Error(respw, "Failed to upload file to GitHub: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer respGithub.Body.Close()

	// Check response from GitHub
	if respGithub.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(respGithub.Body)
		http.Error(respw, "GitHub upload failed: "+string(body), http.StatusInternalServerError)
		return
	}

	// Extract file URL from GitHub response
	var githubResp map[string]interface{}
	err = json.NewDecoder(respGithub.Body).Decode(&githubResp)
	if err != nil {
		http.Error(respw, "Failed to parse GitHub response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	buktiPembayaranURL, _ := githubResp["content"].(map[string]interface{})["html_url"].(string)

	// Parse additional form data into Pesanan struct
	var daftarDesain []model.Portofolio
	err = json.Unmarshal([]byte(req.FormValue("daftar_desain")), &daftarDesain)
	if err != nil {
		http.Error(respw, "Invalid daftar_desain format: "+err.Error(), http.StatusBadRequest)
		return
	}

	pesanan := model.Pesanan{
		ID:             primitive.NewObjectID(),
		NamaPemesan:    req.FormValue("nama_pemesan"),
		DaftarDesain:   daftarDesain,
		TanggalPesanan: time.Now(),
		StatusPesanan:  "Pending",
		Pembayaran:     buktiPembayaranURL,
		CatatanPesanan: req.FormValue("catatan_pesanan"),
		TotalHarga:     req.FormValue("total_harga"),
	}

	// Validate required fields
	if pesanan.NamaPemesan == "" || len(pesanan.DaftarDesain) == 0 || pesanan.TotalHarga == "" {
		http.Error(respw, "Missing required fields in Pesanan data", http.StatusBadRequest)
		return
	}

	// Insert pesanan into database
	// Define pesananCollection
	pesananCollection := config.Mongoconn.Collection("pesanan")
	
	result, err := pesananCollection.InsertOne(context.TODO(), pesanan)
	if err != nil {
		http.Error(respw, "Failed to save pesanan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with Pesanan data
	respw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(respw).Encode(result)
}