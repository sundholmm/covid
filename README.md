# sundholmm/covid

Go REST API with underlying PostgreSQL database for handling COVID-19 records.

## API

Currently the application serves the paths below:

### HTTP POST /api/v1/record

Save a record to the database.

```
{
    "dateRep" : "28/12/2020",
    "year_week" : "2020-52",
    "cases_weekly" : 1994,
    "deaths_weekly" : 88,
    "countriesAndTerritories" : "Afghanistan",
    "geoId" : "AF",
    "countryterritoryCode" : "AFG",
    "popData2019" : 38041757,
    "continentExp" : "Asia",
    "notification_rate_per_100000_population_14-days" : "7.19"
}
```

### HTTP POST /api/v1/records

Save multiple records to the database.

```
{
    "dateRep" : "28/12/2020",
    "year_week" : "2020-52",
    "cases_weekly" : 1994,
    "deaths_weekly" : 88,
    "countriesAndTerritories" : "Afghanistan",
    "geoId" : "AF",
    "countryterritoryCode" : "AFG",
    "popData2019" : 38041757,
    "continentExp" : "Asia",
    "notification_rate_per_100000_population_14-days" : "7.19"
},
{
    "dateRep" : "21/12/2020",
    "year_week" : "2020-51",
    "cases_weekly" : 740,
    "deaths_weekly" : 111,
    "countriesAndTerritories" : "Afghanistan",
    "geoId" : "AF",
    "countryterritoryCode" : "AFG",
    "popData2019" : 38041757,
    "continentExp" : "Asia",
    "notification_rate_per_100000_population_14-days" : "6.56"
}
```

### HTTP GET /api/v1/records

Get records from the database.

Currently supported HTTP GET query parameters are as listed:

- country
- orderBy
- order

For order the values are:

```
OrderASC = "ASC"
OrderDesc = "DESC"
```

For orderBy the values are:

```
OrderByDate = "date"
OrderByCasesWeekly = "cases_weekly"
OrderByDeathsWeekly = "deaths_weekly"
OrderByCountry = "country"
OrderByPopulation = "population"
OrderByNotificationRate = "notification_rate"
```

## Environmental variables

For initializing the application be sure these variables are in place:

```
"DB_HOST" // localhost
"DB_PORT" // 5432
"DB_USER" // postgres
"DB_USER_PASSWORD" // mysecretpassword
"DB_NAME" // covid
```

## Database setup

[Docker postgres](https://hub.docker.com/_/postgres)

## Migrations

DDL migrations (.sql files) should be placed under migrations\ddl and they'll be run when the application starts
Prefix the migrations files as 001, 002, 003, ..., 999 and add a description with underscores

Example:

```
001_create_record_table.sql
```

## Data source

---

"The downloadable data file was updated daily to 14 December 2020 using the latest available public data on COVID-19. Each row/entry contains the number of new cases reported per day and per country. After 14 December 2020, ECDC shifted to weekly data collection." [European Centre for Disease Prevention and Control](https://www.ecdc.europa.eu/en/publications-data/download-todays-data-geographic-distribution-covid-19-cases-worldwide)

---

[COVID-19 Coronavirus data](https://data.europa.eu/euodp/en/data/dataset/covid-19-coronavirus-data)

Can be applied to the database by running the application and then executing

```
$ python utils\migrate.py
```
