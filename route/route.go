package route

import (
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/controller"
	"github.com/gocroot/helper/at"
)

func URL(w http.ResponseWriter, r *http.Request) {
	if config.SetAccessControlHeaders(w, r) {
		return // If it's a preflight request, return early.
	}
	config.SetEnv()

	var method, path string = r.Method, r.URL.Path
	switch {
	case method == "GET" && path == "/":
		controller.GetHome(w, r)
	//chat bot inbox
	case method == "POST" && at.URLParam(path, "/webhook/nomor/:nomorwa"):
		controller.PostInboxNomor(w, r)
	//masking list nmor official
	case method == "GET" && path == "/data/phone/all":
		controller.GetBotList(w, r)
	//akses data helpdesk layanan user
	case method == "GET" && path == "/data/user/helpdesk/all":
		controller.GetHelpdeskAll(w, r)
	case method == "GET" && path == "/data/user/helpdesk/masuk":
		controller.GetLatestHelpdeskMasuk(w, r)
	case method == "GET" && path == "/data/user/helpdesk/selesai":
		controller.GetLatestHelpdeskSelesai(w, r)
	//pamong desa data from api
	case method == "GET" && path == "/data/lms/user":
		controller.GetDataUserFromApi(w, r)
	//simpan testimoni dari pamong desa lms api
	case method == "POST" && path == "/data/lms/testi":
		controller.PostTestimoni(w, r)
		//get random 4 testi
	case method == "GET" && path == "/data/lms/random/testi":
		controller.GetRandomTesti4(w, r)
	//mendapatkan data sent item
	case method == "GET" && at.URLParam(path, "/data/peserta/sent/:id"):
		controller.GetSentItem(w, r)
	//simpan feedback unsubs user
	case method == "POST" && path == "/data/peserta/unsubscribe":
		controller.PostUnsubscribe(w, r)
	//generate token linked device
	case method == "PUT" && path == "/data/user":
		controller.PutTokenDataUser(w, r)
	//Menambhahkan data nomor sender untuk broadcast
	case method == "PUT" && path == "/data/sender":
		controller.PutNomorBlast(w, r)
	//mendapatkan data list nomor sender untuk broadcast
	case method == "GET" && path == "/data/sender":
		controller.GetDataSenders(w, r)
	//mendapatkan data list nomor sender yang kena blokir dari broadcast
	case method == "GET" && path == "/data/blokir":
		controller.GetDataSendersTerblokir(w, r)
	//mendapatkan data rekap pengiriman wa blast
	case method == "GET" && path == "/data/rekap":
		controller.GetRekapBlast(w, r)
	//mendapatkan data faq
	case method == "GET" && at.URLParam(path, "/data/faq/:id"):
		controller.GetFAQ(w, r)
	//legacy
	case method == "PUT" && path == "/data/user/task/doing":
		controller.PutTaskUser(w, r)
	case method == "GET" && path == "/data/user/task/done":
		controller.GetTaskDone(w, r)
	case method == "POST" && path == "/data/user/task/done":
		controller.PostTaskUser(w, r)
	case method == "GET" && path == "/data/pushrepo/kemarin":
		controller.GetYesterdayDistincWAGroup(w, r)

	//helpdesk
	//mendapatkan data tiket
	case method == "GET" && at.URLParam(path, "/data/tiket/closed/:id"):
		controller.GetClosedTicket(w, r)
	//simpan feedback tiket user
	case method == "POST" && path == "/data/tiket/rate":
		controller.PostMasukanTiket(w, r)
		// order
	case method == "POST" && at.URLParam(path, "/data/order/:namalapak"):
		controller.HandleOrder(w, r)

	//user data
	case method == "GET" && path == "/data/user":
		controller.GetDataUser(w, r)
	//user pendaftaran
	case method == "POST" && path == "/auth/register/users": //mendapatkan email gmail
		controller.RegisterGmailAuth(w, r)
	case method == "POST" && path == "/data/user":
		controller.PostDataUser(w, r)
	case method == "POST" && path == "/upload/profpic": //upload gambar profile
		controller.UploadProfilePictureHandler(w, r)
	case method == "POST" && path == "/data/user/bio":
		controller.PostDataBioUser(w, r)
		/* 	case method == "POST" && at.URLParam(path, "/data/user/wa/:nomorwa"):
		controller.PostDataUserFromWA(w, r) */
	//data proyek
	case method == "GET" && path == "/data/proyek":
		controller.GetDataProject(w, r)
	case method == "GET" && path == "/data/proyek/approved": //akses untuk manager
		controller.GetEditorApprovedProject(w, r)
	case method == "POST" && path == "/data/proyek":
		controller.PostDataProject(w, r)
	case method == "PUT" && path == "/data/metadatabuku":
		controller.PutMetaDataProject(w, r)
	case method == "PUT" && path == "/data/proyek/publishbuku": //publish buku isbn by manager
		controller.PutPublishProject(w, r)
	case method == "PUT" && path == "/data/proyek":
		controller.PutDataProject(w, r)
	case method == "DELETE" && path == "/data/proyek":
		controller.DeleteDataProject(w, r)
	case method == "GET" && path == "/data/proyek/anggota":
		controller.GetDataMemberProject(w, r)
	case method == "GET" && path == "/data/proyek/editor":
		controller.GetDataEditorProject(w, r)
	case method == "DELETE" && path == "/data/proyek/anggota":
		controller.DeleteDataMemberProject(w, r)
	case method == "POST" && path == "/data/proyek/anggota":
		controller.PostDataMemberProject(w, r)
	case method == "POST" && path == "/data/proyek/editor": //set editor oleh owner
		controller.PostDataEditorProject(w, r)
	case method == "PUT" && path == "/data/proyek/editor": //set approved oleh editor
		controller.PUtApprovedEditorProject(w, r)
	//upload cover,draft,pdf,sampul buku project
	case method == "POST" && at.URLParam(path, "/upload/coverbuku/:projectid"):
		controller.UploadCoverBukuWithParamFileHandler(w, r)
	case method == "POST" && at.URLParam(path, "/upload/draftbuku/:projectid"):
		controller.UploadDraftBukuWithParamFileHandler(w, r)
	case method == "POST" && at.URLParam(path, "/upload/draftpdfbuku/:projectid"):
		controller.UploadDraftBukuPDFWithParamFileHandler(w, r)
	case method == "POST" && at.URLParam(path, "/upload/sampulpdfbuku/:projectid"):
		controller.UploadSampulBukuPDFWithParamFileHandler(w, r)
	case method == "POST" && at.URLParam(path, "/upload/spk/:projectid"):
		controller.UploadSPKPDFWithParamFileHandler(w, r)
	case method == "POST" && at.URLParam(path, "/upload/spi/:projectid"):
		controller.UploadSPIPDFWithParamFileHandler(w, r)
	case method == "GET" && at.URLParam(path, "/download/draft/:path"): //downoad file draft
		controller.AksesFileRepoDraft(w, r)
	case method == "POST" && path == "/data/proyek/katalog": //post blog katalog
		controller.PostKatalogBuku(w, r)
	case method == "GET" && at.URLParam(path, "/download/dokped/spk/:namaproject"): //base64 namaproject
		controller.GetFileDraftSPK(w, r)
	case method == "GET" && at.URLParam(path, "/download/dokped/spkt/:namaproject"): //base64 namaproject
		controller.GetFileDraftSPKT(w, r)
	case method == "GET" && at.URLParam(path, "/download/dokped/spi/:path"): //base64 path sampul
		controller.GetFileDraftSPI(w, r)

	case method == "POST" && path == "/data/proyek/menu":
		controller.PostDataMenuProject(w, r)
	case method == "POST" && path == "/approvebimbingan":
		controller.ApproveBimbinganbyPoin(w, r)
	case method == "DELETE" && path == "/data/proyek/menu":
		controller.DeleteDataMenuProject(w, r)
	case method == "POST" && path == "/notif/ux/postlaporan":
		controller.PostLaporan(w, r)
	case method == "POST" && path == "/notif/ux/postfeedback":
		controller.PostFeedback(w, r)

	case method == "POST" && path == "/notif/ux/postmeeting":
		controller.PostMeeting(w, r)
	case method == "POST" && at.URLParam(path, "/notif/ux/postpresensi/:id"):
		controller.PostPresensi(w, r)
	case method == "POST" && at.URLParam(path, "/notif/ux/posttasklists/:id"):
		controller.PostTaskList(w, r)
	case method == "POST" && at.URLParam(path, "/webhook/nomor/:nomorwa"):
		controller.PostInboxNomor(w, r)

	// LMS
	case method == "GET" && path == "/lms/refresh/cookie":
		controller.RefreshLMSCookie(w, r)
	case method == "GET" && path == "/lms/count/user":
		controller.GetCountDocUser(w, r)

	// Google Auth
	case method == "POST" && path == "/auth/users":
		controller.Auth(w, r)
	case method == "POST" && path == "/auth/login":
		controller.GeneratePasswordHandler(w, r)
	case method == "POST" && path == "/auth/verify":
		controller.VerifyPasswordHandler(w, r)
	case method == "POST" && path == "/auth/resend":
		controller.ResendPasswordHandler(w, r)

	// AUTH
	// Register
	case method == "POST" && path == "/auth/register":
		controller.RegisterAkun(w, r)
	// Login
	case method == "POST" && path == "/auth/login/form":
		controller.LoginAkun(w, r)
	// Get User By ID
	case method == "GET" && path == "/auth/users/cust/id":
		controller.GetAkunCustomerByID(w, r)
	// Get All User
	case method == "GET" && path == "/auth/users/all":
		controller.GetAllAkun(w, r)
	// User Whatsauth
	case method == "GET" && path == "/auth/users":
		controller.GetUser(w, r)
	// Delete User
	case method == "DELETE" && path == "/delete/user":
		controller.DeleteUser(w, r)
	// Update User
	case method == "PUT" && path == "/update/user":
		controller.UpdateUser(w, r)
		
	// Logout
	case method == "POST" && path == "/auth/logout":
		controller.Logout(w, r)

	// PORTFOLIO
	// Insert Portfolio
	case method == "POST" && path == "/insert/portofolio":
		controller.CreatePortofolio(w, r)
	// Get All Portfolio
	case method == "GET" && path == "/portofolio":
		controller.GetAllPortofolio(w, r)
	// Get Portofolio By Id
	case method == "PUT" && path == "/update/portofolio":
		controller.UpdatePortofolio(w, r)
	// Delete Portofolio
	case method == "DELETE" && path == "/delete/portofolio":
		controller.DeletePortofolio(w, r)
	// Get Portofolio By Id
	case method == "GET" && path == "/portofolio/":
		controller.GetPortofolioByID(w, r)
	// Get Portfolio By Kategori
	case method == "GET" && path == "/portofolio/kategori":
		controller.GetPortofolioByKategori(w, r)

	// PEMBAYARAN
	// Insert Pembayaran
	case method == "POST" && path == "/insert/pembayaran":
		controller.TransaksiHandler(w, r)
	// Get Pembayaran By Id
	case method == "GET" && path == "/transaksi/id":
		controller.GetTransaksiByID(w, r)
	// Get All Pembayaran
	case method == "GET" && path == "/data/transaksi":
		controller.GetAllTransaksi(w, r)
	// Get Pembayaran By Status
	case method == "GET" && path == "/transaksi/status":
		controller.GetTransaksiByStatus(w, r)
	// Update Pembayaran
	case method == "PUT" && path == "/update/transaksi":
		controller.UpdateTransaksi(w, r)
	// Delete Pembayaran
	case method == "DELETE" && path == "/delete/transaksi":
		controller.DeleteTransaksi(w, r)
	// Get All Pembayaran
	// case method == "GET" && path == "/pembayaran":
	// 	controller.GetAllPembayaran(w, r)
	
	// ITEM PESANAN
	// Insert Item Pesanan
	// case method == "POST" && path == "/insert/itempesanan":
	// 	controller.CreateItemPesanan(w, r)
	// // Get All Item Pesanan
	// case method == "GET" && path == "/itempesanan":
	// 	controller.GetAllItemPesanan(w, r)
	// // Update Pesan Status
	// case method == "PUT" && path == "/update/itempesanan":
	// 	controller.UpdatePesanStatus(w, r)
	

	//GEO
	// Roads
	case method == "POST" && path == "/geo/roads":
		controller.GetRoads(w, r)
	// Region
	case method == "POST" && path == "/geo/region":
		controller.GetRegion(w, r)
	
	// login admin
	case method == "POST" && path == "/auth/login/admin":
		controller.LoginAkunAdmin(w, r)
	// register admin
	case method == "POST" && path == "/auth/register/admin":
		controller.RegisterAkunAdmin(w, r)
	// get all admin
	case method == "GET" && path == "/auth/admin":
		controller.GetAllAkunAdmin(w, r)
	
	// Upload
	case method == "POST" && path == "/upload":
		controller.UploadHandler(w, r)
	default:
		controller.NotFound(w, r)
	}
}