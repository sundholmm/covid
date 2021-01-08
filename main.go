package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

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
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pw := os.Getenv("DB_USER_PASSWORD")
	name := os.Getenv("DB_NAME")

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
	http.HandleFunc("/api/v1/records", env.recordsHandler)

	// Listen on port :8080
	log.Println("Server listening on port :8080")
	http.ListenAndServe(":8080", nil)

}

// recordsHandler handles the request types on /api/v1/records path
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

// getRecordsHandler handles the HTTP GET requests on path /api/v1/records
func getRecordsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	log.Printf("getRecordsHandler called by %s", r.Host)
	defer log.Printf("getRecordsHandler called by %s ended", r.Host)

	ctx := r.Context()

	country := r.URL.Query().Get("country")
	orderBy := r.URL.Query().Get("orderBy")
	order := r.URL.Query().Get("order")

	// Uppercase and lowercase letters are treated as equivalent (case-insensitive)
	queryParams := models.QueryParams{
		Country: strings.ToLower(country),
		OrderBy: strings.ToLower(orderBy),
		Order:   strings.ToLower(order),
	}

	err := queryParams.ValidateQueryParams(ctx, db)
	if err != nil {
		log.Println(err)
		r := err.(models.RequestError)
		http.Error(w, r.Error(), r.StatusCode)
		return
	}

	recordsDTO, err := models.GetRecords(ctx, db, &queryParams)

	if err != nil {
		log.Println(err)
		r := err.(models.RequestError)
		http.Error(w, r.Error(), r.StatusCode)
		return
	}

	if recordsDTO.Meta.RecordAmount == 0 {
		http.Error(w, "No records found!", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recordsDTO)

}

// postRecordsHandler handles the HTTP POST requests on path /api/v1/records
func postRecordsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	log.Printf("postRecordsHandler called by %s", r.Host)
	defer log.Printf("postRecordsHandler called by %s ended", r.Host)

	ctx := r.Context()

	var recordDTO models.RecordDTO

	err := json.NewDecoder(r.Body).Decode(&recordDTO)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(400), 400)
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(recordDTO.Records))

	for index := range recordDTO.Records {
		go func(currentRecord models.Record, currenIndex int) {
			defer wg.Done()
			err = currentRecord.Validate()
			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(400), 400)
				return
			}
			log.Printf("Record #%d validate successful", currenIndex)
		}(recordDTO.Records[index], index)
	}

	wg.Wait()

	err = models.SaveRecords(ctx, db, recordDTO)
	if err != nil {
		log.Println(err)
		r := err.(models.RequestError)
		http.Error(w, r.Error(), r.StatusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(http.StatusText(201)))

}
