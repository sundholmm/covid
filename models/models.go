package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
)

// Record struct for a single record
type Record struct {
	Date string `json:"dateRep" validate:"required"`
	YearWeek string `json:"year_week" validate:"required"`
	CasesWeekly *int `json:"cases_weekly,omitempty"`
	DeathsWeekly *int `json:"deaths_weekly,omitempty"`
	Country string `json:"countriesAndTerritories" validate:"required"`
	GeoID string `json:"geoId" validate:"required"`
	CountryCode string `json:"countryterritoryCode"`
	Population int `json:"popData2019"`
	Continent string `json:"continentExp" validate:"required"`
	NotificationRate string `json:"notification_rate_per_100000_population_14-days,omitempty"`
}

// RequestError represents an error with an associated HTTP status code, time of the event and an error string
type RequestError struct {
	StatusCode int
	TimeStamp time.Time
	Err string
}

// Error allows RequestError to satisfy the error interface
func (r RequestError) Error() string {
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d - %s", r.TimeStamp.Year(), r.TimeStamp.Month(), r.TimeStamp.Day(), r.TimeStamp.Hour(), r.TimeStamp.Minute(), r.TimeStamp.Second(), r.Err)
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
	OrderASC = "asc"
	OrderDesc = "desc"
	OrderByDate = "date"
	OrderByCasesWeekly = "cases_weekly"
	OrderByDeathsWeekly = "deaths_weekly"
	OrderByCountry = "country"
	OrderByPopulation = "population"
	OrderByNotificationRate = "notification_rate"
)

// QueryParams struct for a single record
type QueryParams struct {
	Country string
	OrderBy string
	Order string
}

// ValidateQueryParams validates all the GET parameters against the defined constants
// returns bool indicating the valid status and string with the invalid value
func (queryParams QueryParams) ValidateQueryParams(ctx context.Context, db *sql.DB) (error) {

	if queryParams.Country != "" {
		countries, err := getAllCountries(ctx, db,)
		if err != nil {
			return RequestError{500, time.Now(), err.Error()}
		}
		if !stringInSlice(queryParams.Country, countries) {
			return RequestError{404, time.Now(), "Query parameter country not found"}
		}
	}

	if	queryParams.OrderBy != OrderByDate && queryParams.OrderBy != OrderByCasesWeekly &&
		queryParams.OrderBy != OrderByDeathsWeekly && queryParams.OrderBy != OrderByCountry &&
		queryParams.OrderBy != OrderByPopulation && queryParams.OrderBy != "" {
			return RequestError{400, time.Now(), fmt.Sprintf("Invalid query parameter orderBy: %s", queryParams.OrderBy)}
	}

	if queryParams.Order != OrderASC && queryParams.Order != OrderDesc && queryParams.Order != "" {
		return RequestError{400, time.Now(), fmt.Sprintf("Invalid query parameter order: %s", queryParams.Order)}
	}

	return nil

}

// stringInSlice is an utility function for checking if an array contains a string
func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

// getQueryString returns an SQL string conditionally based on the search parameters
func (queryParams QueryParams) getQueryString() (string) {

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
	
	rows, err := db.QueryContext(ctx, "SELECT DISTINCT LOWER(\"country\") FROM record;");
    if err != nil {
        return nil, err
	}

    defer rows.Close()

    var countries []string

    for rows.Next() {

        var country string

        err := rows.Scan(&country)
        if err != nil {
            return nil, err
        }

		countries = append(countries, country)

	}

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return countries, nil

}

// GetAllRecords returns all records from the database
func GetAllRecords(ctx context.Context, db *sql.DB, queryParams *QueryParams) ([]Record, error) {

	rows, err := db.QueryContext(ctx, "SELECT " +
	"\"date\", \"year_week\", \"cases_weekly\", \"deaths_weekly\", \"country\", " +
	"\"geo_id\", \"country_code\", \"population\", \"continent\", \"notification_rate\" " +
	"FROM \"record\"" + queryParams.getQueryString() + ";")

    if err != nil {
        return nil, err
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
            return nil, err
        }

		records = append(records, record)

	}

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return records, nil
}

// SaveSingleRecord saves single record into the database
func SaveSingleRecord(ctx context.Context, db *sql.DB, record Record, index int) (error) {

	sql := fmt.Sprintf("INSERT INTO \"record\" ( " +
	"\"date\", \"year_week\", \"cases_weekly\", \"deaths_weekly\", \"country\", " +
	"\"geo_id\", \"country_code\", \"population\", \"continent\", \"notification_rate\" " +
	") VALUES ( '%s', '%s', %d, %d, '%s', '%s', '%s', %d, '%s', '%s' );",
	record.Date, record.YearWeek, *record.CasesWeekly,
	*record.DeathsWeekly, record.Country, record.GeoID, record.CountryCode,
	record.Population, record.Continent, record.NotificationRate)

	_, err := db.ExecContext(ctx, sql)
	if err != nil {
		return err
	}

	log.Printf("Record number #%d save successful", index)

	return nil

}

// SaveMultipleRecords saves multiple records into the database
func SaveMultipleRecords(ctx context.Context, db *sql.DB, records []Record) (error) {

	log.Printf("Number of records to save: #%d", len(records))
	for index, record := range records {
		err := SaveSingleRecord(ctx, db, record, index)
		if err != nil {
			return err
		}
	}

	return nil

}