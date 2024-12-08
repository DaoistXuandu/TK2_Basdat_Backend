package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "tk_basdat"
	password = "123"
	dbname   = "tk_basdat"
)

var db *sql.DB

// -----------------------------------
// Structs untuk Login/Register/User
// -----------------------------------
type LoginRequestBody struct {
	NoHP string `json:"NoHP"`
	Pwd  string `json:"Pwd"`
}

type LoginResponseBody struct {
	Status  bool   `json:"status"`
	Role    int    `json:"role"`
	UserId  string `json:"userId"`
	Message string `json:"message"`
	Name    string `json:"name"`
}

type RegisterRequestBody struct {
	Role              int       `json:"role"`
	Nama              string    `json:"name"`
	JenisKelamin      string    `json:"sex"`
	NoHP              string    `json:"number"`
	Pwd               string    `json:"password"`
	TglLahir          time.Time `json:"date"`
	Alamat            string    `json:"address"`
	NamaBank          string    `json:"bank"`
	NomorRekening     string    `json:"noRek"`
	NPWP              string    `json:"npwp"`
	LinkFoto          string    `json:"link"`
	Rating            float64   `json:"rating"`
	JmlPsnananSelesai int       `json:"amount"`
}

type RegisterResponseBody struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

type GetUserRequestBody struct {
	User string `json:"user"`
	Role int    `json:"role"`
}

type GetUserResponseBody struct {
	Status            bool      `json:"status"`
	Message           string    `json:"message"`
	User              string    `json:"userid"`
	Role              int       `json:"role"`
	Nama              string    `json:"name"`
	JenisKelamin      string    `json:"sex"`
	NoHP              string    `json:"number"`
	Pwd               string    `json:"password"`
	TglLahir          time.Time `json:"date"`
	Alamat            string    `json:"address"`
	SaldoMyPay        float64   `json:"saldo"`
	Level             string    `json:"level"`
	NamaBank          string    `json:"bank"`
	NomorRekening     string    `json:"noRek"`
	NPWP              string    `json:"npwp"`
	LinkFoto          string    `json:"link"`
	Rating            float64   `json:"rating"`
	JmlPsnananSelesai int       `json:"amount"`
}

type UpdateUserRequestBody struct {
	User          string    `json:"user"`
	Role          int       `json:"role"`
	Nama          string    `json:"name"`
	JenisKelamin  string    `json:"sex"`
	NoHP          string    `json:"number"`
	TglLahir      time.Time `json:"date"`
	Alamat        string    `json:"address"`
	NamaBank      string    `json:"bank"`
	NomorRekening string    `json:"noRek"`
	NPWP          string    `json:"npwp"`
	LinkFoto      string    `json:"link"`
}

type UpdateUserResponseBody struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

// ------------------------------
// Structs untuk Testimoni
// ------------------------------
type Testimoni struct {
	IdTrPemesanan string `json:"idTrPemesanan"`
	Tgl           string `json:"tgl"`
	Teks          string `json:"teks"`
	Rating        int    `json:"rating"`
}

// ------------------------------
// Structs untuk Diskon & Voucher
// ------------------------------
type VoucherItem struct {
	Kode            string  `json:"kode"`
	Potongan        float64 `json:"potongan"`
	MinTrPemesanan  int     `json:"minTrPemesanan"`
	JmlHariBerlaku  int     `json:"jmlHariBerlaku"`
	KuotaPenggunaan int     `json:"kuotaPenggunaan"`
	Harga           float64 `json:"harga"`
}

type PromoItem struct {
	Kode            string    `json:"kode"`
	Potongan        float64   `json:"potongan"`
	MinTrPemesanan  int       `json:"minTrPemesanan"`
	TglAkhirBerlaku time.Time `json:"tglAkhirBerlaku"`
}

type GetDiskonResponse struct {
	Status  bool          `json:"status"`
	Message string        `json:"message"`
	Voucher []VoucherItem `json:"voucher"`
	Promo   []PromoItem   `json:"promo"`
}

type BuyVoucherRequest struct {
	UserID        string `json:"userId"`
	VoucherCode   string `json:"voucherCode"`
	MetodeBayarId string `json:"metodeBayarId"`
}

type BuyVoucherResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// ------------------------------
// Middleware CORS
// ------------------------------
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// -------------------------------------
// Fungsi Helper Testimoni
// -------------------------------------
func IsPesananSelesai(db *sql.DB, pemesananID string) (bool, error) {
	query := `
        SELECT COUNT(*) 
        FROM TR_PEMESANAN_STATUS tps
        JOIN STATUS_PESANAN sp ON tps.IdStatus = sp.Id
        WHERE tps.IdTrPemesanan = $1 AND sp.Status = 'Pesanan selesai'
    `
	var count int
	err := db.QueryRow(query, pemesananID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func IsPelangganPemesan(db *sql.DB, userID, pemesananID string) (bool, error) {
	query := `
        SELECT COUNT(*)
        FROM TR_PEMESANAN_JASA
        WHERE Id = $1 AND IdPelanggan = $2
    `
	var count int
	err := db.QueryRow(query, pemesananID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func CreateTestimoni(db *sql.DB, userID string, pemesananID string, teks string, rating int) error {
	isPemesan, err := IsPelangganPemesan(db, userID, pemesananID)
	if err != nil {
		return err
	}
	if !isPemesan {
		return fmt.Errorf("anda bukan pelanggan yang memesan jasa ini")
	}

	selesai, err := IsPesananSelesai(db, pemesananID)
	if err != nil {
		return err
	}
	if !selesai {
		return fmt.Errorf("pesanan belum selesai, tidak dapat memberikan testimoni")
	}

	tgl := time.Now().Format("2006-01-02")
	query := `
        INSERT INTO TESTIMONI (IdTrPemesanan, Tgl, Teks, Rating)
        VALUES ($1, $2, $3, $4)
    `
	_, err = db.Exec(query, pemesananID, tgl, teks, rating)
	if err != nil {
		return err
	}

	return nil
}

func GetTestimoniBySubkategori(db *sql.DB, subkategoriID string) ([]Testimoni, error) {
	query := `
    SELECT t.IdTrPemesanan, t.Tgl, t.Teks, t.Rating
    FROM TESTIMONI t
    JOIN TR_PEMESANAN_JASA pj ON t.IdTrPemesanan = pj.Id
    WHERE pj.IdKategoriJasa = $1
    `
	rows, err := db.Query(query, subkategoriID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Testimoni
	for rows.Next() {
		var t Testimoni
		err := rows.Scan(&t.IdTrPemesanan, &t.Tgl, &t.Teks, &t.Rating)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func DeleteTestimoni(db *sql.DB, userID, pemesananID, tgl string) error {
	isPemesan, err := IsPelangganPemesan(db, userID, pemesananID)
	if err != nil {
		return err
	}
	if !isPemesan {
		return fmt.Errorf("anda bukan pelanggan yang memesan jasa ini, tidak dapat menghapus testimoni")
	}

	query := `
        DELETE FROM TESTIMONI
        WHERE IdTrPemesanan = $1 AND Tgl = $2
    `
	_, err = db.Exec(query, pemesananID, tgl)
	if err != nil {
		return err
	}

	return nil
}

// Handler Testimoni
func createTestimoniHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	type createTestimoniReq struct {
		UserID       string `json:"userId"`
		PemesananID  string `json:"pemesananId"`
		Teks         string `json:"teks"`
		Rating       int    `json:"rating"`
	}

	var req createTestimoniReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = CreateTestimoni(db, req.UserID, req.PemesananID, req.Teks, req.Rating)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("Testimoni berhasil ditambahkan"))
}

func getTestimoniHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}
	subkategoriID := r.URL.Query().Get("subkategori_id")
	if subkategoriID == "" {
		http.Error(w, "subkategori_id is required", http.StatusBadRequest)
		return
	}

	testimonies, err := GetTestimoniBySubkategori(db, subkategoriID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(testimonies)
}

func deleteTestimoniHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	type deleteTestimoniReq struct {
		UserID      string `json:"userId"`
		PemesananID string `json:"pemesananId"`
		Tgl         string `json:"tgl"`
	}

	var req deleteTestimoniReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = DeleteTestimoni(db, req.UserID, req.PemesananID, req.Tgl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("Testimoni berhasil dihapus"))
}

// --------------------------------------
// Diskon & Voucher Handlers
// --------------------------------------
func getDiskonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	voucherQuery := `
    SELECT d.Kode, d.Potongan, d.MinTrPemesanan, v.JmlHariBerlaku, v.KuotaPenggunaan, v.Harga
    FROM VOUCHER v
    JOIN DISKON d ON v.Kode = d.Kode
    `

	rows, err := db.Query(voucherQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var voucherList []VoucherItem
	for rows.Next() {
		var v VoucherItem
		err := rows.Scan(&v.Kode, &v.Potongan, &v.MinTrPemesanan, &v.JmlHariBerlaku, &v.KuotaPenggunaan, &v.Harga)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		voucherList = append(voucherList, v)
	}

	promoQuery := `
    SELECT d.Kode, d.Potongan, d.MinTrPemesanan, p.TglAkhirBerlaku
    FROM PROMO p
    JOIN DISKON d ON p.Kode = d.Kode
    `
	promoRows, err := db.Query(promoQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer promoRows.Close()

	var promoList []PromoItem
	for promoRows.Next() {
		var p PromoItem
		err := promoRows.Scan(&p.Kode, &p.Potongan, &p.MinTrPemesanan, &p.TglAkhirBerlaku)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		promoList = append(promoList, p)
	}

	response := GetDiskonResponse{
		Status:  true,
		Message: "Berhasil mendapatkan daftar voucher dan promo",
		Voucher: voucherList,
		Promo:   promoList,
	}

	json.NewEncoder(w).Encode(response)
}

func buyVoucherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body BuyVoucherRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var potongan float64
	var minTr int
	var jmlHari int
	var kuota int
	var harga float64
	err = db.QueryRow(`
        SELECT d.Potongan, d.MinTrPemesanan, v.JmlHariBerlaku, v.KuotaPenggunaan, v.Harga
        FROM VOUCHER v
        JOIN DISKON d ON v.Kode = d.Kode
        WHERE v.Kode = $1
    `, body.VoucherCode).Scan(&potongan, &minTr, &jmlHari, &kuota, &harga)

	if err == sql.ErrNoRows {
		response := BuyVoucherResponse{
			Status:  false,
			Message: "voucher tidak ditemukan",
		}
		json.NewEncoder(w).Encode(response)
		return
	} else if err != nil {
		response := BuyVoucherResponse{
			Status:  false,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Id MyPay, sesuaikan dengan data di DB Anda
	myPayId := "e2ae7f92-eefb-47a7-aa1b-c7d157ab94d7"

	tglAwal := time.Now()
	tglAkhir := tglAwal.AddDate(0, 0, jmlHari)

	if body.MetodeBayarId != myPayId {
		_, err := db.Exec(`
            INSERT INTO TR_PEMBELIAN_VOUCHER (Id, TglAwal, TglAkhir, TelahDigunakan, IdPelanggan, IdVoucher, IdMetodeBayar)
            VALUES ($1, $2, $3, 0, $4, $5, $6)`,
			uuid.New(), tglAwal, tglAkhir, body.UserID, body.VoucherCode, body.MetodeBayarId)
		if err != nil {
			response := BuyVoucherResponse{
				Status:  false,
				Message: err.Error(),
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		response := BuyVoucherResponse{
			Status:  true,
			Message: "Voucher berhasil dibeli tanpa MyPay",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Jika metode bayar MyPay, cek saldo
	var saldo float64
	err = db.QueryRow(`SELECT SaldoMyPay FROM "user" WHERE Id = $1`, body.UserID).Scan(&saldo)
	if err == sql.ErrNoRows {
		response := BuyVoucherResponse{
			Status:  false,
			Message: "user tidak ditemukan",
		}
		json.NewEncoder(w).Encode(response)
		return
	} else if err != nil {
		response := BuyVoucherResponse{
			Status:  false,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	if saldo < harga {
		response := BuyVoucherResponse{
			Status:  false,
			Message: "saldo MyPay tidak cukup",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	newSaldo := saldo - harga
	_, err = db.Exec(`UPDATE "user" SET SaldoMyPay = $1 WHERE Id = $2`, newSaldo, body.UserID)
	if err != nil {
		response := BuyVoucherResponse{
			Status:  false,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err = db.Exec(`
        INSERT INTO TR_PEMBELIAN_VOUCHER (Id, TglAwal, TglAkhir, TelahDigunakan, IdPelanggan, IdVoucher, IdMetodeBayar)
        VALUES ($1, $2, $3, 0, $4, $5, $6)`,
		uuid.New(), tglAwal, tglAkhir, body.UserID, body.VoucherCode, body.MetodeBayarId)

	if err != nil {
		response := BuyVoucherResponse{
			Status:  false,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := BuyVoucherResponse{
		Status:  true,
		Message: "Voucher berhasil dibeli dengan MyPay",
	}
	json.NewEncoder(w).Encode(response)
}

// -----------------------------------------
// Endpoint Homepage, Subkategori, Pesan
// -----------------------------------------
func getHomepage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query(`
		SELECT k.id, k.namakategori, s.id, s.namasubkategori
		FROM kategori_jasa k
		LEFT JOIN subkategori_jasa s ON k.id = s.kategorijasaid`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var data []map[string]interface{}
	for rows.Next() {
		var kategoriID, subkategoriID string
		var kategoriNama, subkategoriNama string
		if err := rows.Scan(&kategoriID, &kategoriNama, &subkategoriID, &subkategoriNama); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data = append(data, map[string]interface{}{
			"kategori_id":      kategoriID,
			"kategori_nama":    kategoriNama,
			"subkategori_id":   subkategoriID,
			"subkategori_nama": subkategoriNama,
		})
	}
	json.NewEncoder(w).Encode(data)
}

func getSubkategori(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	rows, err := db.Query(`
		SELECT s.id, s.namasubkategori, s.deskripsi, sesi.sesi, sesi.harga
		FROM subkategori_jasa s
		LEFT JOIN sesi_layanan sesi ON s.id = sesi.subkategoriid
		WHERE s.id = $1`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var data []map[string]interface{}
	for rows.Next() {
		var subID, sesi int
		var subNama, subDeskripsi string
		var harga float64
		if err := rows.Scan(&subID, &subNama, &subDeskripsi, &sesi, &harga); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data = append(data, map[string]interface{}{
			"subkategori_id":        subID,
			"subkategori_nama":      subNama,
			"subkategori_deskripsi": subDeskripsi,
			"sesi":                  sesi,
			"harga":                 harga,
		})
	}
	json.NewEncoder(w).Encode(data)
}

func createPesanan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		UserID           string  `json:"user_id"`
		SesiID           int     `json:"sesi_id"`
		Tanggal          string  `json:"tanggal"`
		Diskon           float64 `json:"diskon"`
		MetodePembayaran string  `json:"metode_pembayaran"`
		Total            float64 `json:"total"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(`
		INSERT INTO pesanan (id_user, id_sesi, tanggal, diskon, metode_pembayaran, total, status)
		VALUES ($1, $2, $3, $4, $5, $6, 'Menunggu Pembayaran')`,
		body.UserID, body.SesiID, body.Tanggal, body.Diskon, body.MetodePembayaran, body.Total)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Order created"))
}

func main() {
	pgConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", pgConnStr)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	fmt.Println("Connected to the PostgreSQL database")

	// Endpoint user
	http.HandleFunc("/login", corsMiddleware(checkLogin))
	http.HandleFunc("/register", corsMiddleware(registerHandler))
	http.HandleFunc("/getUser", corsMiddleware(getUser))
	http.HandleFunc("/updateUser", corsMiddleware(updateUser))

	// Endpoint testimoni
	http.HandleFunc("/createTestimoni", corsMiddleware(createTestimoniHandler))
	http.HandleFunc("/getTestimoni", corsMiddleware(getTestimoniHandler))
	http.HandleFunc("/deleteTestimoni", corsMiddleware(deleteTestimoniHandler))

	// Endpoint diskon & voucher
	http.HandleFunc("/getDiskon", corsMiddleware(getDiskonHandler))
	http.HandleFunc("/buyVoucher", corsMiddleware(buyVoucherHandler))

	// Endpoint homepage, subkategori, pesan
	http.HandleFunc("/homepage", corsMiddleware(getHomepage))
	http.HandleFunc("/subkategori", corsMiddleware(getSubkategori))
	http.HandleFunc("/pesan", corsMiddleware(createPesanan))

	fmt.Println("Server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}