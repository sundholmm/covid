package models

import (
	"database/sql"
	"fmt"
)

// Record struct for a single record
type Record struct {
	Date string `json:"date"`
	Day string `json:"day"`
	Month string `json:"month"`
	Year string `json:"year"`
	Cases int `json:"cases"`
	Deaths int `json:"deaths"`
	Country string `json:"country"`
	GeoID string `json:"geoId"`
	CountryCode string `json:"countryCode"`
	Population int `json:"population"`
	Continent string `json:"continent"`
	Cumulative string `json:"cumulative"`
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
func SaveSingleRecord(db *sql.DB, record Record) (error) {

	sql := fmt.Sprintf("INSERT INTO \"record\" ( " +
	"\"date\", \"day\", \"month\", \"year\", \"cases\", \"deaths\", \"country\", " +
	"\"geo_id\", \"country_code\", \"population\", \"continent\", \"cumulative\" " +
	") VALUES ( '%s', '%s', '%s', '%s', %d, %d, '%s', '%s', '%s', %d, '%s', '%s' )",
	record.Date, record.Day, record.Month, record.Year, record.Cases,
	record.Deaths, record.Country, record.GeoID, record.CountryCode,
	record.Population, record.Continent, record.Cumulative)

	_, err := db.Exec(sql)
	if err != nil {
		return err
	}

	return nil

}

// SaveMultipleRecords saves multiple records into the database
func SaveMultipleRecords(db *sql.DB, records []Record) (error) {

	for _, record := range records {
		err := SaveSingleRecord(db, record)
		if err != nil {
			return err
		}
	}

	return nil

}