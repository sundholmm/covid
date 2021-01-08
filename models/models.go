package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
)

// RecordDTO struct for multiple records and metadata
type RecordDTO struct {
	Records []Record `json:"records" validate:"required"`
	Meta    Metadata `json:"metadata,omitempty"`
}

// Metadata struct for RecordDTO metadata
type Metadata struct {
	RecordAmount int `json:"record_amount" validate:"required"`
}

// Record struct for a single record
type Record struct {
	Date             string `json:"dateRep" validate:"required"`
	YearWeek         string `json:"year_week" validate:"required"`
	CasesWeekly      *int   `json:"cases_weekly,omitempty"`
	DeathsWeekly     *int   `json:"deaths_weekly,omitempty"`
	Country          string `json:"countriesAndTerritories" validate:"required"`
	GeoID            string `json:"geoId" validate:"required"`
	CountryCode      string `json:"countryterritoryCode"`
	Population       int    `json:"popData2019"`
	Continent        string `json:"continentExp" validate:"required"`
	NotificationRate string `json:"notification_rate_per_100000_population_14-days,omitempty"`
}

// RequestError represents an error with an associated HTTP status code, time of the event and an error string
type RequestError struct {
	StatusCode int
	TimeStamp  time.Time
	Err        string
}

// Error allows RequestError to satisfy the error interface
func (r RequestError) Error() string {
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d - HTTP STATUS %d - %s", r.TimeStamp.Year(), r.TimeStamp.Month(), r.TimeStamp.Day(), r.TimeStamp.Hour(), r.TimeStamp.Minute(), r.TimeStamp.Second(), r.StatusCode, r.Err)
}

// Use a single instance of Validate
var validate *validator.Validate

// Validate validates the record struct
func (record Record) Validate() error {
	validate = validator.New()
	err := validate.Struct(record)
	if err != nil {
		return err
	}
	return nil
}

// Constants for defining allowed GET parameter values
const (
	OrderASC                = "asc"
	OrderDesc               = "desc"
	OrderByDate             = "date"
	OrderByCasesWeekly      = "cases_weekly"
	OrderByDeathsWeekly     = "deaths_weekly"
	OrderByCountry          = "country"
	OrderByPopulation       = "population"
	OrderByNotificationRate = "notification_rate"
)

// QueryParams struct for a single record
type QueryParams struct {
	Country string
	OrderBy string
	Order   string
}

// ValidateQueryParams validates all the GET parameters against the defined constants
// returns bool indicating the valid status and string with the invalid value
func (queryParams QueryParams) ValidateQueryParams(ctx context.Context, db *sql.DB) error {

	if queryParams.Country != "" {
		countries, err := getAllCountries(ctx, db)
		if err != nil {
			return err
		}
		if !stringInSlice(queryParams.Country, countries) {
			return RequestError{404, time.Now(), fmt.Sprintf("Query parameter country value \"%s\" not found", queryParams.Country)}
		}
	}

	if (queryParams.OrderBy == "" && queryParams.Order != "") || queryParams.Order == "" && queryParams.OrderBy != "" {
		return RequestError{400, time.Now(), "Query parameters order and orderBy must both be included"}
	}

	if queryParams.OrderBy != OrderByDate && queryParams.OrderBy != OrderByCasesWeekly &&
		queryParams.OrderBy != OrderByDeathsWeekly && queryParams.OrderBy != OrderByCountry &&
		queryParams.OrderBy != OrderByPopulation && queryParams.OrderBy != "" {
		return RequestError{400, time.Now(), fmt.Sprintf("Invalid query parameter orderBy value \"%s\"", queryParams.OrderBy)}
	}

	if queryParams.Order != OrderASC && queryParams.Order != OrderDesc && queryParams.Order != "" {
		return RequestError{400, time.Now(), fmt.Sprintf("Invalid query parameter order value \"%s\"", queryParams.Order)}
	}

	return nil

}

// stringInSlice returns a boolean based on wether an array contains a string
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// getQueryString returns an SQL string conditionally based on the search parameters
func (queryParams QueryParams) getQueryString() string {

	var whereCountry string
	var orderBy string

	if queryParams.Country != "" {
		whereCountry = fmt.Sprintf(" WHERE LOWER(\"country\")='%s' ", queryParams.Country)
	}

	if queryParams.OrderBy != "" && queryParams.Order != "" {
		orderBy = fmt.Sprintf(" ORDER BY \"%s\" %s", queryParams.OrderBy, queryParams.Order)
	}

	return whereCountry + orderBy

}

// getAllCountries returns all countries from the database that have records saved
func getAllCountries(ctx context.Context, db *sql.DB) ([]string, error) {

	requestError := RequestError{500, time.Now(), http.StatusText(500)}

	rows, err := db.QueryContext(ctx, "SELECT DISTINCT LOWER(\"country\") FROM record;")
	if err != nil {
		return nil, requestError
	}

	defer rows.Close()

	var countries []string

	for rows.Next() {

		var country string

		err := rows.Scan(&country)
		if err != nil {
			return nil, requestError
		}

		countries = append(countries, country)

	}

	if err = rows.Err(); err != nil {
		return nil, requestError
	}

	return countries, nil

}

// GetRecords returns all records from the database
func GetRecords(ctx context.Context, db *sql.DB, queryParams *QueryParams) (*RecordDTO, error) {

	requestError := RequestError{500, time.Now(), http.StatusText(500)}

	rows, err := db.QueryContext(ctx, "SELECT "+
		"\"date\", \"year_week\", \"cases_weekly\", \"deaths_weekly\", \"country\", "+
		"\"geo_id\", \"country_code\", \"population\", \"continent\", \"notification_rate\" "+
		"FROM \"record\""+queryParams.getQueryString()+";")

	if err != nil {
		return nil, requestError
	}

	defer rows.Close()

	var records []Record

	for rows.Next() {

		var record Record

		err := rows.Scan(
			&record.Date,
			&record.YearWeek,
			&record.CasesWeekly,
			&record.DeathsWeekly,
			&record.Country,
			&record.GeoID,
			&record.CountryCode,
			&record.Population,
			&record.Continent,
			&record.NotificationRate,
		)

		if err != nil {
			return nil, requestError
		}

		records = append(records, record)

	}

	if err = rows.Err(); err != nil {
		return nil, requestError
	}

	recordDTO := RecordDTO{
		Records: records,
		Meta: Metadata{
			RecordAmount: len(records),
		},
	}

	return &recordDTO, nil
}

// SaveRecords saves an array of records into the database
func SaveRecords(ctx context.Context, db *sql.DB, recordDTO RecordDTO) error {

	requestError := RequestError{500, time.Now(), http.StatusText(500)}

	log.Printf("Number of records to save: %d", len(recordDTO.Records))

	txn, err := db.Begin()
	if err != nil {
		return requestError
	}

	stmt, err := txn.PrepareContext(ctx, pq.CopyIn("record", "date", "year_week", "cases_weekly", "deaths_weekly", "country", "geo_id", "country_code", "population", "continent", "notification_rate"))
	if err != nil {
		return requestError
	}

	for _, record := range recordDTO.Records {
		_, err = stmt.ExecContext(ctx, record.Date, record.YearWeek, record.CasesWeekly, record.DeathsWeekly, record.Country, record.GeoID, record.CountryCode, record.Population, record.Continent, record.NotificationRate)
		if err != nil {
			return requestError
		}
	}

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return requestError
	}

	err = stmt.Close()
	if err != nil {
		return requestError
	}

	err = txn.Commit()
	if err != nil {
		return requestError
	}

	log.Printf("Successfully saved %d records", len(recordDTO.Records))

	return nil

}
