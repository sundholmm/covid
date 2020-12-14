package models

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
)

// Record struct for a single record
type Record struct {
	Date string `json:"dateRep" validate:"required"`
	Day string `json:"day" validate:"required"`
	Month string `json:"month" validate:"required"`
	Year string `json:"year" validate:"required"`
	Cases *int `json:"cases,omitempty"`
	Deaths *int `json:"deaths,omitempty"`
	Country string `json:"countriesAndTerritories" validate:"required"`
	GeoID string `json:"geoId" validate:"required"`
	CountryCode string `json:"countryterritoryCode" validate:"required"`
	Population *int `json:"popData2019,omitempty"`
	Continent string `json:"continentExp" validate:"required"`
	Cumulative string `json:"Cumulative_number_for_14_days_of_COVID-19_cases_per_100000,omitempty"`
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

// GetAllRecords returns all records from the database
func GetAllRecords(db *sql.DB) ([]Record, error) {

	rows, err := db.Query("SELECT " +
	"\"date\", \"day\", \"month\", \"year\", \"cases\", \"deaths\", \"country\", " +
	"\"geo_id\", \"country_code\", \"population\", \"continent\", \"cumulative\" " +
	"FROM \"record\"")

    if err != nil {
        return nil, err
	}

    defer rows.Close()

    var records []Record

    for rows.Next() {

        var record Record

        err := rows.Scan(
			&record.Date,
			&record.Day,
			&record.Month,
			&record.Year,
			&record.Cases,
			&record.Deaths,
			&record.Country,
			&record.GeoID,
			&record.CountryCode,
			&record.Population,
			&record.Continent,
			&record.Cumulative,
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
func SaveSingleRecord(db *sql.DB, record Record, index int) (error) {

	sql := fmt.Sprintf("INSERT INTO \"record\" ( " +
	"\"date\", \"day\", \"month\", \"year\", \"cases\", \"deaths\", \"country\", " +
	"\"geo_id\", \"country_code\", \"population\", \"continent\", \"cumulative\" " +
	") VALUES ( '%s', '%s', '%s', '%s', %d, %d, '%s', '%s', '%s', %d, '%s', '%s' )",
	record.Date, record.Day, record.Month, record.Year, *record.Cases,
	*record.Deaths, record.Country, record.GeoID, record.CountryCode,
	*record.Population, record.Continent, record.Cumulative)

	_, err := db.Exec(sql)
	if err != nil {
		return err
	}

	log.Printf("Record number #%d save successful", index)

	return nil

}

// SaveMultipleRecords saves multiple records into the database
func SaveMultipleRecords(db *sql.DB, records []Record) (error) {

	log.Printf("Number of records to save: #%d", len(records))
	for index, record := range records {
		err := SaveSingleRecord(db, record, index)
		if err != nil {
			return err
		}
	}

	return nil

}