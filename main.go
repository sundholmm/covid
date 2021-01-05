package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"covid.sundholm.io/migrations"
	"covid.sundholm.io/models"

	_ "github.com/lib/pq"
)

// Env struct for holding the connection pool
type Env struct {
    db *sql.DB
}

func main() {

	// Load the env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	// Database connection config
	host	:= os.Getenv("DB_HOST")
	port	:= os.Getenv("DB_PORT")
	user	:= os.Getenv("DB_USER")
	pw		:= os.Getenv("DB_USER_PASSWORD")
	name	:= os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
	"password=%s dbname=%s sslmode=disable",
	host, port, user, pw, name)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Database connection successful")

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	env := &Env{db: db}

	// Run database migrations
	migrations.MigrateDB(env.db)

	// Register handlers
	http.HandleFunc("/api/v1/record", env.recordHandler)
	http.HandleFunc("/api/v1/records", env.recordsHandler)

	// Listen on port :8080
	log.Println("Server listening on port :8080")
	http.ListenAndServe(":8080", nil)

}

// recordHandler handles the request types on /record path
func (env *Env) recordHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
		case http.MethodPost:
			postRecordHandler(w, r, env.db)
		default:
			log.Printf("recordHandler called with not yet implemented method %s by %s", r.Method, r.Host)
			http.Error(w, http.StatusText(501), 501)
	}

}

// postRecordHandler handles the HTTP POST requests on path /record
func postRecordHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	log.Printf("postRecordHandler called by %s", r.Host)
	defer log.Printf("postRecordHandler called by %s ended", r.Host)

	ctx := r.Context()

	var record models.Record

	index := 1

	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(400), 400)
		return
	}

	err = record.Validate()
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(400), 400)
		return
	}

	err = models.SaveSingleRecord(ctx, db, record, index)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
	w.Write([]byte(http.StatusText(201)))

}

// recordsHandler handles the request types on /records path
func (env *Env) recordsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
		case http.MethodGet:
			getRecordsHandler(w, r, env.db)
		case http.MethodPost:
			postRecordsHandler(w, r, env.db)
		default:
			log.Printf("recordsHandler called with not yet implemented method %s by %s", r.Method, r.Host)
			http.Error(w, http.StatusText(501), 501)
	}

}

// getRecordsHandler handles the HTTP GET requests on path /records
func getRecordsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	log.Printf("getRecordsHandler called by %s", r.Host)
	defer log.Printf("getRecordsHandler called by %s ended", r.Host)

	ctx := r.Context()

	country := r.URL.Query().Get("country")
	orderBy := r.URL.Query().Get("orderBy")
	order := r.URL.Query().Get("order")

	queryParams := models.QueryParams{
		Country: country,
		OrderBy: orderBy,
		Order: order,
	}

	valid, invalid := queryParams.ValidateQueryParams(ctx, db)
	if valid == false {
		log.Printf("Invalid GET query parameter %s sent by %s", invalid, r.Host)
		http.Error(w, fmt.Sprintf("Invalid GET query parameter value %s", invalid), 400)
		return
	}

	records, err := models.GetAllRecords(ctx, db, &queryParams)

	if err != nil {
        log.Println(err)
        http.Error(w, http.StatusText(500), 500)
        return
	}

	if records == nil {
		http.Error(w, "No records found!", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)

}

// postRecordsHandler handles the HTTP POST requests on path /records
func postRecordsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	log.Printf("postRecordsHandler called by %s", r.Host)
	defer log.Printf("postRecordsHandler called by %s ended", r.Host)

	ctx := r.Context()
	
	var records []models.Record

	err := json.NewDecoder(r.Body).Decode(&records)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(400), 400)
		return
	}

	for index, record := range records {
		err = record.Validate()
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(400), 400)
			return
		}
		log.Printf("Record #%d validate successful", index)
	}

	err = models.SaveMultipleRecords(ctx, db, records)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
	w.Write([]byte(http.StatusText(201)))

}
