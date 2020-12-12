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
		http.Error(w, http.StatusText(501), 501)
	}

}

// postRecordHandler handles the HTTP POST requests on path /record
func postRecordHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var record models.Record

	err := json.NewDecoder(r.Body).Decode(&record)
    if err != nil {
        log.Println(err)
		http.Error(w, http.StatusText(400), 400)
		return
	}

    err = models.SaveSingleRecord(db, record)
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
			getRecordsHandler(w, env.db)
		case http.MethodPost:
			postRecordsHandler(w, r, env.db)
		default:
			http.Error(w, http.StatusText(501), 501)
	}

}

// getRecordsHandler handles the HTTP GET requests on path /records
func getRecordsHandler(w http.ResponseWriter, db *sql.DB) {

	records, err := models.GetAllRecords(db)

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

	var records []models.Record

	err := json.NewDecoder(r.Body).Decode(&records)
    if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(400), 400)
		return
	}

    err = models.SaveMultipleRecords(db, records)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
	w.Write([]byte(http.StatusText(201)))

}