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
	Status              bool      `json:"status"`
	Message             string    `json:"message"`
	User                string    `json:"userid"`
	Role                int       `json:"role"`
	Nama                string    `json:"name"`
	JenisKelamin        string    `json:"sex"`
	NoHP                string    `json:"number"`
	Pwd                 string    `json:"password"`
	TglLahir            time.Time `json:"date"`
	Alamat              string    `json:"address"`
	SaldoMyPay          float64   `json:"saldo"`
	Level               string    `json:"level"`
	NamaBank            string    `json:"bank"`
	NomorRekening       string    `json:"noRek"`
	NPWP                string    `json:"npwp"`
	LinkFoto            string    `json:"link"`
	Rating              float64   `json:"rating"`
	JmlPsnananSelesai   int       `json:"amount"`
	PekerjaKategoriJasa []string  `json:"pekerjakategorijasa"`
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

// BAGIAN MERAH
type MyPayHistory struct {
	ID       string  `json:"id"`
	Tgl      string  `json:"date"`
	Nominal  float64 `json:"nominal"`
	Kategori string  `json:"category"`
}

type MyPayHistoryResponse struct {
	UserID  string         `json:"user_id"`
	History []MyPayHistory `json:"history"`
}

type MyPayTransactionTransfer struct {
	UserID     string  `json:"user_id"`
	KategoriID string  `json:"kategori_id"`
	Nominal    float64 `json:"nominal"`
	ToUserID   string  `json:"to_user_id,omitempty"` // For transfers
}

type MyPayTransactionTopUp struct {
	UserID     string  `json:"userId"`
	Nominal    float64 `json:"nominal"`
	KategoriID int     `json:"kategoriId"`
}

type MyPayKategori struct {
	NamaKategori string `json:"namaKategori"`
}

type GetPesananJasaRequestBody struct {
	User string `json:"user"`
}

type GetPesananJasaResponseBody struct {
	Status  bool          `json:"status"`
	Message string        `json:"message"`
	Pesanan []PesananJasa `json:"pesanan"`
}

type PesananJasa struct {
	NamaJasa   string  `json:"nama_jasa"`
	TotalBiaya float64 `json:"total_biaya"`
}

type MyPayTransactionPay struct {
	UserID     string  `json:"user_id"`
	KategoriID string  `json:"kategori_id"`
	Nominal    float64 `json:"nominal"`
	ToUserID   string  `json:"to_user_id,omitempty"` // For transfers
	BankName   string  `json:"bank_name,omitempty"`  // For withdrawals
	AccountNo  string  `json:"account_no,omitempty"`
}

type GetJobsRequest struct {
	UserID string `json:"user_id"`
}

type JobsData struct {
	Id              string    `json:"id"`
	Kategori        string    `json:"kategori"`
	NamaSubkategori string    `json:"subkategori"`
	TanggalPesan    time.Time `json:"tanggal"`
	NamaPelanggan   string    `json:"nama"`
	Sesi            int       `json:"sesi"`
	Total           float64   `json:"total"`
}

type GetJobsResponse struct {
	Status  bool       `json:"status"`
	Message string     `json:"message"`
	Pesanan []JobsData `json:"pesanan"`
}

type PickJobRequest struct {
	UserID string `json:"user_id"`
	TRID   string `json:"transaksi_pemesanan_jasa_id"`
}

type PickJobResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Id      string `json:"id"`
}

type PekerjaJobRequest struct {
	UserID string `json:"user_id"`
}

type JobsDataDemand struct {
	Id              string    `json:"id"`
	Kategori        string    `json:"kategori"`
	NamaSubkategori string    `json:"subkategori"`
	TanggalPesan    time.Time `json:"tanggal"`
	NamaPelanggan   string    `json:"nama"`
	Sesi            int       `json:"sesi"`
	Total           float64   `json:"total"`
	Status          int       `json:"status"`
}

type PekerjaJobResponse struct {
	Status    bool             `json:"status"`
	Message   string           `json:"message"`
	Pekerjaan []JobsDataDemand `json:"pekerjaan"`
}

type JobUpdateStatusRequest struct {
	TRID string `json:"transaksi_pemesanan_jasa_id"`
}

type JobUpdateStatusResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Id      string `json:"id"`
}


// Struct Testimoni, Diskon, Voucher dari kode terakhir
type Testimoni struct {
	IdTrPemesanan string `json:"idTrPemesanan"`
	Tgl           string `json:"tgl"`
	Teks          string `json:"teks"`
	Rating        int    `json:"rating"`
}

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

	// tambah endpoint disini
	http.HandleFunc("/login", corsMiddleware(checkLogin))
	http.HandleFunc("/register", corsMiddleware(register))
	http.HandleFunc("/getUser", corsMiddleware(getUser))
	http.HandleFunc("/updateUser", corsMiddleware(updateUser))
	http.HandleFunc("/homepage", getHomepage)
	http.HandleFunc("/subkategori", getSubkategori)
	http.HandleFunc("/pesan", createPesanan)

	http.HandleFunc("/mypay/balance", corsMiddleware(getMyPayBalance))
	http.HandleFunc("/mypay/history", corsMiddleware(getMyPayHistory))
	http.HandleFunc("/mypay/topup", corsMiddleware(handleTopUp))
	http.HandleFunc("/mypay/get-category-id", corsMiddleware(GetCategoryIdByName))
	http.HandleFunc("/mypay/getPesananJasa", corsMiddleware(getPesananJasa))
	http.HandleFunc("/mypay/getStatusIdByName", corsMiddleware(GetStatusIdByName))
	http.HandleFunc("/mypay/processPayment", corsMiddleware(ProcessPayment))
	// http.HandleFunc("/mypay/transaction", corsMiddleware(handleMyPayTransaction))
	http.HandleFunc("/pekerja/get-kategori-sub", corsMiddleware(getKategoriFromSub))

	http.HandleFunc("/jobs/available", corsMiddleware(getAvailableJobs))
	http.HandleFunc("/jobs/get-job", corsMiddleware(pickAJob))

	http.HandleFunc("/jobs/job-pekerja-id", corsMiddleware(seePekerjaJob))
	http.HandleFunc("/jobs/job-pekerja-update", corsMiddleware(updatePekerjaJob))


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

		// Call the next handler
		next.ServeHTTP(w, r)
	}
}

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
	err = db.QueryRow(`INSERT INTO "user" VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING Id`, uuid.New(), body.Nama, body.JenisKelamin, body.NoHP, body.Pwd, body.TglLahir, body.Alamat, 0.0).Scan(&userId)
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
			fmt.Println(body.NomorRekening + " " + body.NamaBank + " " + oldValue.NomorRekening + " " + oldValue.NamaBank)
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
			fmt.Println("Masuk-[0]")
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
			fmt.Println("Masuk")
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
		Message: fmt.Sprint("User dengan id $s berhasil di update", current_user_id),
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

		rows, err := db.Query(`SELECT NamaKategori FROM KATEGORI_JASA LEFT JOIN PEKERJA_KATEGORI_JASA 
		ON Id = KategoriJasaId WHERE PekerjaId = $1`, body.User)

		var kategoriList []string
		if err != nil {
			log.Println("Error executing query:", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var namaKategori string
			if err := rows.Scan(&namaKategori); err != nil {
				log.Println("Error scanning row:", err)
				return
			}
			kategoriList = append(kategoriList, namaKategori)
		}
		response.PekerjaKategoriJasa = kategoriList

		json.NewEncoder(w).Encode(response)
	}
}

// Get Homepage (all categories and subcategories)
func getHomepage(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT k.id, k.nama, s.id, s.nama
		FROM kategori_jasa k
		LEFT JOIN subkategori_jasa s ON k.id = s.id_kategori`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var data []map[string]interface{}
	for rows.Next() {
		var kategoriID, subkategoriID int
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

// Get Subkategori and Sessions
func getSubkategori(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	rows, err := db.Query(`
		SELECT s.id, s.nama, s.deskripsi, sesi.id, sesi.nama_sesi, sesi.harga
		FROM subkategori_jasa s
		LEFT JOIN sesi_layanan sesi ON s.id = sesi.id_subkategori
		WHERE s.id = $1`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var data []map[string]interface{}
	for rows.Next() {
		var subID, sesiID int
		var subNama, subDeskripsi, sesiNama string
		var harga float64
		if err := rows.Scan(&subID, &subNama, &subDeskripsi, &sesiID, &sesiNama, &harga); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data = append(data, map[string]interface{}{
			"subkategori_id":        subID,
			"subkategori_nama":      subNama,
			"subkategori_deskripsi": subDeskripsi,
			"sesi_id":               sesiID,
			"sesi_nama":             sesiNama,
			"harga":                 harga,
		})
	}
	json.NewEncoder(w).Encode(data)
}

// Create Order
func createPesanan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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

// BAGIAN MERAH
// MyPay model to track user's balance
type MyPay struct {
	ID      int     `json:"id"`
	UserID  int     `json:"user_id"`
	Balance float64 `json:"balance"`
}

// Service Order model
type ServiceOrder struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	WorkerID    int    `json:"worker_id"`
	ServiceName string `json:"service_name"`
	Status      string `json:"status"` // Example: "Looking for Worker", "Worker Assigned", etc.
	CreatedAt   string `json:"created_at"`
	ScheduledAt string `json:"scheduled_at"`
}

// Get MyPay Balance for User
func getMyPayBalance(w http.ResponseWriter, r *http.Request) {
	var request GetUserRequestBody

	// Decode the JSON body to get user ID
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.User == "" {
		http.Error(w, "Missing user ID in request body", http.StatusBadRequest)
		return
	}

	var balance float64
	var noHP string

	// Use the provided user ID to get the balance
	err := db.QueryRow("SELECT SaldoMyPay, NoHP FROM \"user\" WHERE Id = $1", request.User).Scan(&balance, &noHP)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to retrieve balance", http.StatusInternalServerError)
		return
	}

	// Set the response headers and write the balance in JSON format
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"user_id": request.User,
		"balance": balance,
		"no_hp":   noHP,
	}
	json.NewEncoder(w).Encode(response)
}

// Get History Transaction for User
func getMyPayHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // Change to POST if you're sending data in the request body
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body into the GetUserRequestBody struct
	var requestBody GetUserRequestBody
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Check if user ID is present in the request body
	if requestBody.User == "" {
		http.Error(w, "Missing user ID in request body", http.StatusBadRequest)
		return
	}

	// Query the transaction history using the user ID from the request body
	rows, err := db.Query(`
		SELECT t.Id, t.Tgl, t.Nominal, k.Nama 
		FROM TR_MYPAY t
		JOIN KATEGORI_TR_MYPAY k ON t.KategoriId = k.Id
		WHERE t.UserId = $1
		ORDER BY t.Tgl DESC`, requestBody.User)
	if err != nil {
		http.Error(w, "Failed to fetch transaction history", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var history []MyPayHistory
	for rows.Next() {
		var transaction MyPayHistory
		err := rows.Scan(&transaction.ID, &transaction.Tgl, &transaction.Nominal, &transaction.Kategori)
		if err != nil {
			http.Error(w, "Failed to parse transaction history", http.StatusInternalServerError)
			return
		}
		history = append(history, transaction)
	}

	// Respond with the transaction history in the required format
	response := MyPayHistoryResponse{
		UserID:  requestBody.User,
		History: history,
	}
	json.NewEncoder(w).Encode(response)
}

// Handle Top-Up
func handleTopUp(w http.ResponseWriter, r *http.Request) {
	// var request MyPayKategori
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var transaction MyPayTransactionTopUp
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("UPDATE \"user\" SET SaldoMyPay = SaldoMyPay + $1 WHERE Id = $2", transaction.Nominal, transaction.UserID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to top-up", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`INSERT INTO TR_MYPAY (Id, UserId, Tgl, Nominal, KategoriId) VALUES (uuid_generate_v4(), $1, CURRENT_DATE, $2, $3)`,
		transaction.UserID, transaction.Nominal, transaction.KategoriID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to record top-up transaction", http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Top-up successful"})
}

// GetCategoryIdByName fetches the category UUID based on the category name
func GetCategoryIdByName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body into the MyPayKategori struct
	var requestBody MyPayKategori
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Validate the category name
	if requestBody.NamaKategori == "" {
		http.Error(w, "Category name is required", http.StatusBadRequest)
		return
	}

	// Query to get the category UUID based on the category name
	var kategoriId uuid.UUID
	err := db.QueryRow(`
        SELECT Id
        FROM KATEGORI_TR_MYPAY
        WHERE Nama = $1`, requestBody.NamaKategori).Scan(&kategoriId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Category not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch category UUID", http.StatusInternalServerError)
		}
		return
	}

	// Respond with the category UUID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"kategoriId": kategoriId.String(),
	})
}

func getPesananJasa(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body GetPesananJasaRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var response GetPesananJasaResponseBody

	rows, err := db.Query(`
		SELECT 
			k.NamaKategori, 
			tpj.TotalBiaya 
		FROM 
			TR_PEMESANAN_JASA tpj
		JOIN 
			KATEGORI_JASA k ON k.Id = tpj.IdKategoriJasa
		JOIN 
			TR_PEMESANAN_STATUS tps ON tpj.Id = tps.IdTrPemesanan
		WHERE 
			tps.Keterangan = 'Menunggu Pembayaran'
			AND tpj.IdPelanggan = $1`, body.User)
	if err != nil {
		response.Status = false
		response.Message = "Error executing query: " + err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	defer rows.Close()

	var pesanan []PesananJasa
	for rows.Next() {
		var pesananItem PesananJasa
		err := rows.Scan(&pesananItem.NamaJasa, &pesananItem.TotalBiaya)
		if err != nil {
			response.Status = false
			response.Message = "Error scanning row: " + err.Error()
			json.NewEncoder(w).Encode(response)
			return
		}
		pesanan = append(pesanan, pesananItem)
	}

	// Pastikan response selalu mengembalikan pesanan, meskipun kosong
	response.Status = true
	if len(pesanan) == 0 {
		response.Message = "No orders found"
	} else {
		response.Message = "Successfully retrieved orders"
	}
	response.Pesanan = pesanan

	json.NewEncoder(w).Encode(response)
}

func GetStatusIdByName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parsing request body
	var requestBody struct {
		StatusName string `json:"statusName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if requestBody.StatusName == "" {
		http.Error(w, "Status name is required", http.StatusBadRequest)
		return
	}

	// Fetch status ID from STATUS_PESANAN
	var statusId uuid.UUID
	err := db.QueryRow(`SELECT Id FROM STATUS_PESANAN WHERE Status = $1`, requestBody.StatusName).Scan(&statusId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Status not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch status ID", http.StatusInternalServerError)
		}
		return
	}

	// Respond with status ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"statusId": statusId.String(),
	})
}

func ProcessPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parsing request body
	var requestBody struct {
		UserId    string `json:"userId"`    // UUID format
		ServiceId string `json:"serviceId"` // UUID format
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate UUID format for UserId and ServiceId
	if _, err := uuid.Parse(requestBody.UserId); err != nil {
		http.Error(w, "Invalid UserId format", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(requestBody.ServiceId); err != nil {
		http.Error(w, "Invalid ServiceId format", http.StatusBadRequest)
		return
	}

	// Fetch service price
	var servicePrice float64
	err := db.QueryRow(`SELECT TotalBiaya FROM TR_PEMESANAN_JASA WHERE Id = $1`, requestBody.ServiceId).Scan(&servicePrice)
	if err != nil {
		http.Error(w, "Failed to fetch service price", http.StatusInternalServerError)
		return
	}

	// Fetch user balance
	var userBalance float64
	err = db.QueryRow(`SELECT SaldoMyPay FROM "user" WHERE Id = $1`, requestBody.UserId).Scan(&userBalance)
	if err != nil {
		http.Error(w, "Failed to fetch user balance", http.StatusInternalServerError)
		return
	}

	// Validate balance
	if userBalance < servicePrice {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Saldo tidak mencukupi untuk melakukan pembayaran.",
		})
		return
	}

	// Update user balance
	_, err = db.Exec(`UPDATE "user" SET SaldoMyPay = SaldoMyPay - $1 WHERE Id = $2`, servicePrice, requestBody.UserId)
	if err != nil {
		http.Error(w, "Failed to update user balance", http.StatusInternalServerError)
		return
	}

	// Fetch new status ID
	var newStatusId uuid.UUID
	err = db.QueryRow(`SELECT Id FROM STATUS_PESANAN WHERE Status = 'Mencari Pekerja Terdekat'`).Scan(&newStatusId)
	if err != nil {
		http.Error(w, "Failed to fetch new status ID", http.StatusInternalServerError)
		return
	}

	// Update service status
	_, err = db.Exec(`
		UPDATE TR_PEMESANAN_JASA 
		SET IdKategoriJasa = $1, IdDiskon = NULL, IdMetodeBayar = NULL 
		WHERE Id = $2`, newStatusId, requestBody.ServiceId)
	if err != nil {
		http.Error(w, "Failed to update service status", http.StatusInternalServerError)
		return
	}

	// Insert into TR_MYPAY
	var categoryId uuid.UUID
	err = db.QueryRow(`SELECT Id FROM KATEGORI_TR_MYPAY WHERE Nama = 'membayar transaksi jasa'`).Scan(&categoryId)
	if err != nil {
		http.Error(w, "Failed to fetch category ID", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(`INSERT INTO TR_MYPAY (Id, UserId, Tgl, Nominal, KategoriId) VALUES (uuid_generate_v4(), $1, CURRENT_DATE, $2, $3)`,
		requestBody.UserId, servicePrice, categoryId)
	if err != nil {
		http.Error(w, "Failed to insert transaction", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Payment successful",
	})
}

func getAvailableJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body GetJobsRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var response GetJobsResponse
	rows, err := db.Query(`SELECT tj.Id 
	FROM tr_pemesanan_jasa AS tj 
	LEFT JOIN tr_pemesanan_status as ts ON tj.Id = ts.IdTrPemesanan 
	LEFT JOIN status_pesanan as sp ON sp.Id = ts.IdStatus 
	LEFT JOIN SUBKATEGORI_JASA as sj ON sj.Id = tj.IdKategoriJasa
	LEFT JOIN PEKERJA_KATEGORI_JASA as pj ON pj.KategoriJasaId = sj.KategoriJasaId
	WHERE 
	sp.Status LIKE '%Terdekat%' AND 
	pj.PekerjaId = $1
	`, body.UserID)

	var pesananList []JobsData
	if err != nil {
		response := &GetJobsResponse{
			Status:  false,
			Message: err.Error(),
			Pesanan: nil,
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var pemesanan string
		if err := rows.Scan(&pemesanan); err != nil {
			log.Println("Error scanning row:", err)
			return
		}
		var exist int
		db.QueryRow(`SELECT 1 AS data 
		FROM tr_pemesanan_status 
		LEFT JOIN status_pesanan ON status_pesanan.Id = tr_pemesanan_status.IdStatus 
		WHERE tr_pemesanan_status.IdTrPemesanan = $1 
		AND status_pesanan.Status LIKE '%Berangkat%'
		`, pemesanan).Scan(&exist)

		var response_pesan JobsData
		if exist != 1 {
			db.QueryRow(`SELECT 
			TJ.Id, 
			SJ.NamaSubkategori, 
			TJ.TglPemesanan, 
			U.Nama,  
			TJ.Sesi,
			TJ.TotalBiaya,
			KJ.NamaKategori
			FROM TR_PEMESANAN_JASA AS TJ
			LEFT JOIN SUBKATEGORI_JASA AS SJ ON TJ.IdKategoriJasa = SJ.Id
			LEFT JOIN KATEGORI_JASA AS KJ ON KJ.Id = SJ.KategoriJasaId
			LEFT JOIN "user" AS U ON U.Id = TJ.IdPelanggan
			WHERE  
			TJ.Id = $1
			`, pemesanan).Scan(
				&response_pesan.Id,
				&response_pesan.NamaSubkategori,
				&response_pesan.TanggalPesan,
				&response_pesan.NamaPelanggan,
				&response_pesan.Sesi,
				&response_pesan.Total,
				&response_pesan.Kategori,
			)

			pesananList = append(pesananList, response_pesan)
		}
	}

	response.Status = true
	response.Message = "Berhasil mengambil data"
	response.Pesanan = pesananList
	json.NewEncoder(w).Encode(response)
}

func pickAJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body PickJobRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var sesi int
	err = db.QueryRow(`SELECT Sesi FROM TR_PEMESANAN_JASA WHERE Id = $1`, body.TRID).Scan(&sesi)
	if err == sql.ErrNoRows {
		response := &PickJobResponse{
			Status:  false,
			Message: "Invalid Credential",
		}

		json.NewEncoder(w).Encode(response)
		return
	} else if err != nil {
		response := &PickJobResponse{
			Status:  false,
			Message: err.Error(),
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		response := &PickJobResponse{
			Status:  false,
			Message: err.Error(),
		}

		json.NewEncoder(w).Encode(response)
		return
	}
	currentTime := time.Now().In(location)

	date := currentTime.Format("2006-01-02")
	time := currentTime.AddDate(0, 0, sesi).Format("2006-01-02 15:04:05")
	time_status := currentTime.Format("2006-01-02 15:04:05")

	fmt.Println(currentTime, date, time)

	var value string
	err = db.QueryRow(`
	UPDATE TR_PEMESANAN_JASA 
	SET IdPekerja = $1, TglPekerjaan = $2, WaktuPekerjaan = $3 WHERE Id = $4 RETURNING Id`,
		body.UserID, date, time, body.TRID).Scan(&value)

	if err == sql.ErrNoRows {
		response := &PickJobResponse{
			Status:  false,
			Message: "Invalid Credential",
		}

		json.NewEncoder(w).Encode(response)
		return
	} else if err != nil {
		response := &PickJobResponse{
			Status:  false,
			Message: err.Error(),
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	// b8a64f3b-8c66-4e74-b8c4-d62d1b42d33b | Menunggu Pembayaran
	// afac6469-f299-4a56-9eb5-5e7b9f84a6b6 | Mencari Pekerja Terdekat
	// e88a03a5-7de1-4f5d-9d77-1d8149b0aab6 | Menunggu Pekerja Berangkat
	// c5a47e2c-dba6-445e-b98e-29553f74e6a7 | Pekerja tiba di lokasi
	// 2d79e2eb-64f5-4718-bb1e-c9d5e10b3274 | Pelayanan jasa sedang dilakukan
	// a0f51f69-bcb5-45a7-9d55-09c2a15ae4bc | Pesanan selesai
	// 56bb004e-0b0e-4cb8-982b-98eb4f5dc542 | Pesanan dibatal

	err = db.QueryRow(`
	INSERT INTO TR_PEMESANAN_STATUS VALUES ($1, $2, $3) RETURNING IdTrPemesanan`,
		body.TRID, "e88a03a5-7de1-4f5d-9d77-1d8149b0aab6", time_status).Scan(&value)
	if err == sql.ErrNoRows {
		response := &PickJobResponse{
			Status:  false,
			Message: "Invalid Credential",
		}

		json.NewEncoder(w).Encode(response)
		return
	} else if err != nil {
		response := &PickJobResponse{
			Status:  false,
			Message: err.Error() + "Here",
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	response := &PickJobResponse{
		Status:  true,
		Message: "Succes",
		Id:      value,
	}

	json.NewEncoder(w).Encode(response)
}

func seePekerjaJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body PekerjaJobRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	rows, err := db.Query(`SELECT Id FROM TR_PEMESANAN_JASA WHERE IdPekerja = $1`, body.UserID)

	fmt.Println(body.UserID)

	var pekerjaanList []JobsDataDemand
	var response PekerjaJobResponse
	if err != nil {
		response := &PekerjaJobResponse{
			Status:    false,
			Message:   err.Error(),
			Pekerjaan: nil,
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var pemesanan string
		if err := rows.Scan(&pemesanan); err != nil {
			log.Println("Error scanning row:", err)
			return
		}

		var response_pesan JobsDataDemand
		db.QueryRow(`SELECT 
			TJ.Id, 
			SJ.NamaSubkategori, 
			TJ.TglPemesanan, 
			U.Nama,  
			TJ.Sesi,
			TJ.TotalBiaya,
			KJ.NamaKategori
			FROM TR_PEMESANAN_JASA AS TJ
			LEFT JOIN SUBKATEGORI_JASA AS SJ ON TJ.IdKategoriJasa = SJ.Id
			LEFT JOIN KATEGORI_JASA AS KJ ON KJ.Id = SJ.KategoriJasaId
			LEFT JOIN "user" AS U ON U.Id = TJ.IdPelanggan
			WHERE  
			TJ.Id = $1
			`, pemesanan).Scan(
			&response_pesan.Id,
			&response_pesan.NamaSubkategori,
			&response_pesan.TanggalPesan,
			&response_pesan.NamaPelanggan,
			&response_pesan.Sesi,
			&response_pesan.Total,
			&response_pesan.Kategori,
		)

		db.QueryRow(`
		SELECT 
			TJ.Id, 
			SJ.NamaSubkategori, 
			TJ.TglPemesanan, 
			U.Nama AS NamaPelanggan,  
			TJ.Sesi,
			TJ.TotalBiaya,
			KJ.NamaKategori,
			COUNT(TS.IdTrPemesanan) AS JumlahStatus
		FROM 
		    TR_PEMESANAN_JASA AS TJ
			LEFT JOIN SUBKATEGORI_JASA AS SJ ON TJ.IdKategoriJasa = SJ.Id
			LEFT JOIN KATEGORI_JASA AS KJ ON KJ.Id = SJ.KategoriJasaId
			LEFT JOIN "user" AS U ON U.Id = TJ.IdPelanggan
			LEFT JOIN TR_PEMESANAN_STATUS AS TS ON TS.IdTrPemesanan = TJ.Id
		WHERE  
   			 TJ.Id = $1
		GROUP BY 
  		  	TJ.Id, SJ.NamaSubkategori, TJ.TglPemesanan, U.Nama, TJ.Sesi, TJ.TotalBiaya, KJ.NamaKategori;
		`, pemesanan).Scan(
			&response_pesan.Id,
			&response_pesan.NamaSubkategori,
			&response_pesan.TanggalPesan,
			&response_pesan.NamaPelanggan,
			&response_pesan.Sesi,
			&response_pesan.Total,
			&response_pesan.Kategori,
			&response_pesan.Status,
		)
		if response_pesan.Status > 2 {
			pekerjaanList = append(pekerjaanList, response_pesan)
		}
	}

	response.Pekerjaan = pekerjaanList
	response.Status = true
	response.Message = "Berhasil Mendapatkan data"
	json.NewEncoder(w).Encode(response)
}

func updatePekerjaJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body JobUpdateStatusRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// b8a64f3b-8c66-4e74-b8c4-d62d1b42d33b | Menunggu Pembayaran
	// afac6469-f299-4a56-9eb5-5e7b9f84a6b6 | Mencari Pekerja Terdekat
	// e88a03a5-7de1-4f5d-9d77-1d8149b0aab6 | Menunggu Pekerja Berangkat
	// c5a47e2c-dba6-445e-b98e-29553f74e6a7 | Pekerja tiba di lokasi
	// 2d79e2eb-64f5-4718-bb1e-c9d5e10b3274 | Pelayanan jasa sedang dilakukan
	// a0f51f69-bcb5-45a7-9d55-09c2a15ae4bc | Pesanan selesai
	// 56bb004e-0b0e-4cb8-982b-98eb4f5dc542 | Pesanan dibatal

	var status_pesanan = [6]string{
		"b8a64f3b-8c66-4e74-b8c4-d62d1b42d33b",
		"afac6469-f299-4a56-9eb5-5e7b9f84a6b6",
		"e88a03a5-7de1-4f5d-9d77-1d8149b0aab6",
		"c5a47e2c-dba6-445e-b98e-29553f74e6a7",
		"2d79e2eb-64f5-4718-bb1e-c9d5e10b3274",
		"a0f51f69-bcb5-45a7-9d55-09c2a15ae4bc",
	}

	var value int
	err = db.QueryRow(`SELECT COUNT(IdTrPemesanan) AS Jumlah FROM TR_PEMESANAN_STATUS WHERE IdTrPemesanan = $1`, body.TRID).Scan(&value)
	if err == sql.ErrNoRows {
		response := &JobUpdateStatusResponse{
			Status:  false,
			Message: "Invalid Credential",
		}

		json.NewEncoder(w).Encode(response)
		return
	} else if err != nil {
		response := &JobUpdateStatusResponse{
			Status:  false,
			Message: err.Error(),
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		response := &JobUpdateStatusResponse{
			Status:  false,
			Message: err.Error(),
		}

		json.NewEncoder(w).Encode(response)
		return
	}
	currentTime := time.Now().In(location)
	time_status := currentTime.Format("2006-01-02 15:04:05")

	if value == 6 {
		response := &JobUpdateStatusResponse{
			Status:  false,
			Message: "Pekerja telah menyelesaikan pesanan",
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	var id_tr string
	err = db.QueryRow(`
	INSERT INTO TR_PEMESANAN_STATUS VALUES ($1, $2, $3) RETURNING IdTrPemesanan`,
		body.TRID, status_pesanan[value], time_status).Scan(&id_tr)
	if err == sql.ErrNoRows {
		response := &JobUpdateStatusResponse{
			Status:  false,
			Message: "Invalid Credential",
		}

		json.NewEncoder(w).Encode(response)
		return
	} else if err != nil {
		response := &JobUpdateStatusResponse{
			Status:  false,
			Message: err.Error() + "Here",
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	response := &JobUpdateStatusResponse{
		Status:  true,
		Message: "Berhasil Memperbaharui data",
		Id:      id_tr + " " + status_pesanan[value],
	}

	json.NewEncoder(w).Encode(response)
}

type getKategoriFromSubRequest struct {
	Id string `json:"id"`
}

type getKategoriFromSubResponse struct {
	Kategori    []string   `json:"kategori"`
	SubKategori [][]string `json:"subkategori"`
	Status      bool       `json:"status"`
	Message     string     `json:"message"`
}

func getKategoriFromSub(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body getKategoriFromSubRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	rows, err := db.Query(`SELECT 
	Id, NamaKategori FROM KATEGORI_JASA 
	LEFT JOIN PEKERJA_KATEGORI_JASA 
	ON Id = KategoriJasaId WHERE PekerjaId = $1`, body.Id)

	var kategori_list []string
	var sub_kategori_list [][]string
	if err != nil {
		response := &getKategoriFromSubResponse{
			Status:  false,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id_kategori string
		var namaKategori string
		var sub_kategoris []string
		if err := rows.Scan(&id_kategori, &namaKategori); err != nil {
			response := &getKategoriFromSubResponse{
				Status:  false,
				Message: err.Error(),
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		kategori_list = append(kategori_list, namaKategori)

		rows2, err := db.Query(`SELECT
		NamaSubkategori FROM SUBKATEGORI_JASA
		WHERE KategoriJasaId = $1
		`, id_kategori)

		if err != nil {
			response := &getKategoriFromSubResponse{
				Status:  false,
				Message: err.Error(),
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		defer rows2.Close()

		for rows2.Next() {
			var sub_kategori string
			if err := rows2.Scan(&sub_kategori); err != nil {
				response := &getKategoriFromSubResponse{
					Status:  false,
					Message: err.Error(),
				}
				json.NewEncoder(w).Encode(response)
				return
			}
			sub_kategoris = append(sub_kategoris, sub_kategori)
		}

		sub_kategori_list = append(sub_kategori_list, sub_kategoris)
	}

	response := &getKategoriFromSubResponse{
		Kategori:    kategori_list,
		SubKategori: sub_kategori_list,
		Status:      true,
		Message:     "Berhasil Mendapatkan data",
	}

	json.NewEncoder(w).Encode(response)
}

// ------------------------------------------------------
// Bagian Testimoni (Dari kode yang ingin digabungkan)
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
// Bagian Diskon & Voucher (dari kode yang ingin digabungkan)
// ------------------------------------------------------
func getDiskonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body BuyVoucherRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
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
		response := BuyVoucherResponse{
			Status:  false,
			Message: "Saldo MyPay tidak cukup",
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