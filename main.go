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

// Struct untuk Testimoni
type Testimoni struct {
	IdTrPemesanan string `json:"idTrPemesanan"`
	Tgl           string `json:"tgl"`
	Teks          string `json:"teks"`
	Rating        int    `json:"rating"`
}

// Struct untuk Diskon, Voucher, Promo
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
	Status  bool         `json:"status"`
	Message string       `json:"message"`
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

func main() {
	pgConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	conn, err := sql.Open("postgres", pgConnStr)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	db = conn
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	fmt.Println("Connected to the PostgreSQL database")

	// Endpoint existing
	http.HandleFunc("/login", corsMiddleware(checkLogin))
	http.HandleFunc("/register", corsMiddleware(register))
	http.HandleFunc("/getUser", corsMiddleware(getUser))
	http.HandleFunc("/updateUser", corsMiddleware(updateUser))

	// Endpoint baru untuk testimoni
	http.HandleFunc("/createTestimoni", corsMiddleware(createTestimoniHandler))
	http.HandleFunc("/getTestimoni", corsMiddleware(getTestimoniHandler))
	http.HandleFunc("/deleteTestimoni", corsMiddleware(deleteTestimoniHandler))

	// Endpoint untuk diskon & pembelian voucher
	http.HandleFunc("/getDiskon", corsMiddleware(getDiskonHandler))
	http.HandleFunc("/buyVoucher", corsMiddleware(buyVoucherHandler))

	fmt.Println("Server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// ------------------------------------------------------
// Bagian Testimoni
// ------------------------------------------------------
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
		return fmt.Errorf("Anda bukan pelanggan yang memesan jasa ini.")
	}

	selesai, err := IsPesananSelesai(db, pemesananID)
	if err != nil {
		return err
	}
	if !selesai {
		return fmt.Errorf("Pesanan belum selesai, tidak dapat memberikan testimoni.")
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
		return fmt.Errorf("Anda bukan pelanggan yang memesan jasa ini, tidak dapat menghapus testimoni.")
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

// Handler create testimoni
func createTestimoniHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
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
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = DeleteTestimoni(db, req.UserID, req.PemesananID, req.Tgl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("Testimoni berhasil dihapus"))
}

// ------------------------------------------------------
// Bagian Diskon & Voucher
// ------------------------------------------------------
func getDiskonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Ambil data voucher
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

	// Ambil data promo
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
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body BuyVoucherRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ambil data voucher
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
			Message: "Voucher tidak ditemukan",
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

	// Asumsi Id MyPay sudah diketahui (cek dari data di DB Anda)
	myPayId := "e2ae7f92-eefb-47a7-aa1b-c7d157ab94d7"

	var tglAwal = time.Now()
	var tglAkhir = tglAwal.AddDate(0, 0, jmlHari)

	if body.MetodeBayarId != myPayId {
		// Metode bayar bukan MyPay, langsung berhasil
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

	// Jika metode bayar adalah MyPay, cek saldo
	var saldo float64
	err = db.QueryRow(`SELECT SaldoMyPay FROM "user" WHERE Id = $1`, body.UserID).Scan(&saldo)
	if err == sql.ErrNoRows {
		response := BuyVoucherResponse{
			Status:  false,
			Message: "User tidak ditemukan",
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
		// Saldo tidak cukup
		response := BuyVoucherResponse{
			Status:  false,
			Message: "Saldo MyPay tidak cukup",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Saldo cukup, update saldo user
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

	// Insert ke TR_PEMBELIAN_VOUCHER
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

// ------------------------------------------------------
// Bagian User (Login, Register, GetUser, UpdateUser)
// ------------------------------------------------------
func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body RegisterRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var userId string
	err = db.QueryRow(`INSERT INTO "user" VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING Id`,
		uuid.New(), body.Nama, body.JenisKelamin, body.NoHP, body.Pwd, body.TglLahir, body.Alamat, 0.0).Scan(&userId)
	if err == sql.ErrNoRows {
		response := &RegisterResponseBody{
			Status:  false,
			Message: "Invalid Credential on user",
		}

		json.NewEncoder(w).Encode(response)
		return
	} else if err != nil {
		response := &RegisterResponseBody{
			Status:  false,
			Message: err.Error(),
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	if body.Role == 0 {
		err = db.QueryRow(`INSERT INTO PELANGGAN VALUES ($1, $2) RETURNING Id`, userId, "Basic").Scan(&userId)
		if err == sql.ErrNoRows {
			response := &RegisterResponseBody{
				Status:  false,
				Message: "Invalid Credential on pelanggan",
			}

			json.NewEncoder(w).Encode(response)
			return
		} else if err != nil {
			response := &RegisterResponseBody{
				Status:  false,
				Message: err.Error(),
			}

			json.NewEncoder(w).Encode(response)
			return
		}
	} else {
		err = db.QueryRow(`INSERT INTO PEKERJA VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING Id`,
			userId,
			body.NamaBank,
			body.NomorRekening,
			body.NPWP,
			body.LinkFoto,
			body.Rating,
			body.JmlPsnananSelesai).Scan(&userId)

		if err == sql.ErrNoRows {
			response := &RegisterResponseBody{
				Status:  false,
				Message: "Invalid Credential on pekerja",
			}

			json.NewEncoder(w).Encode(response)
			return
		} else if err != nil {
			response := &RegisterResponseBody{
				Status:  false,
				Message: err.Error(),
			}

			json.NewEncoder(w).Encode(response)
			return
		}
	}

	response := &RegisterResponseBody{
		Status:  true,
		Message: "User berhasil dibuat",
	}

	json.NewEncoder(w).Encode(response)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body UpdateUserRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var oldValue UpdateUserRequestBody
	err = db.QueryRow(`SELECT Nama, JenisKelamin, NoHP, TglLahir, Alamat FROM "user" WHERE Id = $1`, body.User).
		Scan(
			&oldValue.Nama,
			&oldValue.JenisKelamin,
			&oldValue.NoHP,
			&oldValue.TglLahir,
			&oldValue.Alamat)

	if err == sql.ErrNoRows {
		response := &UpdateUserResponseBody{
			Status:  false,
			Message: "Invalid Credential on user",
		}

		json.NewEncoder(w).Encode(response)
		return
	} else if err != nil {
		response := &UpdateUserResponseBody{
			Status:  false,
			Message: err.Error(),
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	var current_user_id string
	err = db.QueryRow(`UPDATE "user" SET Nama = $1, JenisKelamin = $2, TglLahir = $3, Alamat = $4 WHERE Id = $5 Returning Id`,
		oldValue.Nama,
		oldValue.JenisKelamin,
		oldValue.TglLahir,
		oldValue.Alamat,
		body.User).Scan(&current_user_id)

	if body.NoHP != oldValue.NoHP {
		err = db.QueryRow(`UPDATE "user" SET NoHP = $1 WHERE Id = $2 Returning Id`,
			body.NoHP,
			body.User).Scan(&current_user_id)
	}

	if err == sql.ErrNoRows {
		response := &UpdateUserResponseBody{
			Status:  false,
			Message: "Invalid update on user",
		}

		json.NewEncoder(w).Encode(response)
		return
	} else if err != nil {
		response := &UpdateUserResponseBody{
			Status:  false,
			Message: err.Error() + " User",
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	if body.Role == 1 {
		err = db.QueryRow(`SELECT NPWP, LinkFoto, NamaBank, NomorRekening FROM PEKERJA WHERE Id = $1`, body.User).
			Scan(
				&oldValue.NPWP,
				&oldValue.LinkFoto,
				&oldValue.NamaBank,
				&oldValue.NomorRekening)
		if err == sql.ErrNoRows {
			response := &UpdateUserResponseBody{
				Status:  false,
				Message: "Invalid Credential on pekerja",
			}

			json.NewEncoder(w).Encode(response)
			return
		} else if err != nil {
			response := &UpdateUserResponseBody{
				Status:  false,
				Message: err.Error() + " Update",
			}

			json.NewEncoder(w).Encode(response)
			return
		}

		err = db.QueryRow(`UPDATE PEKERJA SET 
        NPWP = $1, 
        LinkFoto = $2 
        WHERE Id = $3 Returning Id`,
			oldValue.NPWP,
			oldValue.LinkFoto,
			body.User).Scan(&current_user_id)
		if err == sql.ErrNoRows {
			response := &UpdateUserResponseBody{
				Status:  false,
				Message: "Invalid Credential on pekerja",
			}

			json.NewEncoder(w).Encode(response)
			return
		} else if err != nil {
			response := &UpdateUserResponseBody{
				Status:  false,
				Message: err.Error() + " Pekerja",
			}

			json.NewEncoder(w).Encode(response)
			return
		}

		if body.NomorRekening != oldValue.NomorRekening && body.NamaBank != oldValue.NamaBank {
			err = db.QueryRow(`UPDATE PEKERJA SET 
			NamaBank = $1, 
			NomorRekening = $2 
			WHERE Id = $3 Returning Id`,
				body.NamaBank,
				body.NomorRekening,
				body.User).Scan(&current_user_id)
			if err == sql.ErrNoRows {
				response := &UpdateUserResponseBody{
					Status:  false,
					Message: "Invalid Credential on pekerja",
				}

				json.NewEncoder(w).Encode(response)
				return
			} else if err != nil {
				response := &UpdateUserResponseBody{
					Status:  false,
					Message: err.Error() + " Update",
				}

				json.NewEncoder(w).Encode(response)
				return
			}
		} else if body.NamaBank != oldValue.NamaBank {
			err = db.QueryRow(`UPDATE PEKERJA SET 
			NamaBank = $1
			WHERE Id = $2 Returning Id`,
				body.NamaBank,
				body.User).Scan(&current_user_id)
			if err == sql.ErrNoRows {
				response := &UpdateUserResponseBody{
					Status:  false,
					Message: "Invalid Credential on pekerja",
				}

				json.NewEncoder(w).Encode(response)
				return
			} else if err != nil {
				response := &UpdateUserResponseBody{
					Status:  false,
					Message: err.Error() + " Bank Name",
				}

				json.NewEncoder(w).Encode(response)
				return
			}
		} else if body.NomorRekening != oldValue.NomorRekening {
			err = db.QueryRow(`UPDATE PEKERJA SET 
			NomorRekening = $1
			WHERE Id = $2 Returning Id`,
				body.NomorRekening,
				body.User).Scan(&current_user_id)
			if err == sql.ErrNoRows {
				response := &UpdateUserResponseBody{
					Status:  false,
					Message: "Invalid Credential on pekerja",
				}

				json.NewEncoder(w).Encode(response)
				return
			} else if err != nil {
				response := &UpdateUserResponseBody{
					Status:  false,
					Message: err.Error() + " Rekening",
				}

				json.NewEncoder(w).Encode(response)
				return
			}
		}

	}

	response := &RegisterResponseBody{
		Status:  true,
		Message: fmt.Sprintf("User dengan id %s berhasil di update", current_user_id),
	}

	json.NewEncoder(w).Encode(response)
}

func checkLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body LoginRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var userID string
	var name string
	var role int

	err = db.QueryRow(`SELECT Id, Nama FROM "user" WHERE NoHP = $1 AND Pwd = $2`, body.NoHP, body.Pwd).Scan(&userID, &name)
	if err == sql.ErrNoRows {
		response := &LoginResponseBody{
			Status:  false,
			UserId:  userID,
			Name:    name,
			Role:    role,
			Message: "Invalid Credential",
		}

		json.NewEncoder(w).Encode(response)
		return
	} else if err != nil {
		response := &LoginResponseBody{
			Status:  false,
			UserId:  userID,
			Name:    name,
			Role:    role,
			Message: err.Error(),
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	db.QueryRow(`SELECT 1 FROM PELANGGAN WHERE Id = $1`, userID).Scan(&role)
	response := &LoginResponseBody{
		Status:  true,
		UserId:  userID,
		Name:    name,
		Role:    role,
		Message: "Success",
	}

	json.NewEncoder(w).Encode(response)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body GetUserRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var response GetUserResponseBody
	err = db.QueryRow(`SELECT Nama, JenisKelamin, NoHP, Pwd, TglLahir, Alamat, SaldoMyPay FROM "user" WHERE Id = $1`, body.User).Scan(
		&response.Nama,
		&response.JenisKelamin,
		&response.NoHP,
		&response.Pwd,
		&response.TglLahir,
		&response.Alamat,
		&response.SaldoMyPay)

	if err == sql.ErrNoRows {
		response := &GetUserResponseBody{
			Status:  false,
			Message: "Invalid Credential",
		}

		json.NewEncoder(w).Encode(response)
		return
	} else if err != nil {
		response := &GetUserResponseBody{
			Status:  false,
			Message: err.Error(),
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	response.Status = true
	response.Message = "Berhasil mendapatkan data"

	if body.Role == 0 {
		db.QueryRow(`SELECT Level FROM PELANGGAN WHERE Id = $1`, body.User).Scan(&response.Level)
		json.NewEncoder(w).Encode(response)
	} else {
		db.QueryRow(`SELECT NamaBank, NomorRekening, NPWP, LinkFoto, Rating, JmlPsnananSelesai FROM PEKERJA WHERE Id = $1`, body.User).Scan(
			&response.NamaBank,
			&response.NomorRekening,
			&response.NPWP,
			&response.LinkFoto,
			&response.Rating,
			&response.JmlPsnananSelesai)
		json.NewEncoder(w).Encode(response)
	}
}
