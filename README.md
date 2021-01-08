# sundholmm/covid

Go REST API with underlying PostgreSQL database for handling COVID-19 records.

- SQL Database Transactions
- Multithreading
- Pure Go

## Installation

```
$ git clone https://github.com/sundholmm/covid.git
$ cd covid
$ go build
$ ./covid.sundholm.io
```

## API

The endpoints are compatible between each other and as such the JSON objects are directly forwardable.
Currently the application serves the paths below:

### HTTP POST /api/v1/records

Save an array of records to the database.

```
{
   "records":[
      {
         "dateRep":"14/12/2020",
         "year_week":"2020-50",
         "cases_weekly":3179,
         "deaths_weekly":46,
         "countriesAndTerritories":"Finland",
         "geoId":"FI",
         "countryterritoryCode":"FIN",
         "popData2019":5517919,
         "continentExp":"Europe",
         "notification_rate_per_100000_population_14-days":"112.02"
      },
      {
         "dateRep":"07/12/2020",
         "year_week":"2020-49",
         "cases_weekly":3002,
         "deaths_weekly":22,
         "countriesAndTerritories":"Finland",
         "geoId":"FI",
         "countryterritoryCode":"FIN",
         "popData2019":5517919,
         "continentExp":"Europe",
         "notification_rate_per_100000_population_14-days":"108.59"
      }
   ]
}
```

### HTTP GET /api/v1/records

Get records from the database.

```
{
   "records":[
      {
         "dateRep":"14/12/2020",
         "year_week":"2020-50",
         "cases_weekly":3179,
         "deaths_weekly":46,
         "countriesAndTerritories":"Finland",
         "geoId":"FI",
         "countryterritoryCode":"FIN",
         "popData2019":5517919,
         "continentExp":"Europe",
         "notification_rate_per_100000_population_14-days":"112.02"
      },
      {
         "dateRep":"07/12/2020",
         "year_week":"2020-49",
         "cases_weekly":3002,
         "deaths_weekly":22,
         "countriesAndTerritories":"Finland",
         "geoId":"FI",
         "countryterritoryCode":"FIN",
         "popData2019":5517919,
         "continentExp":"Europe",
         "notification_rate_per_100000_population_14-days":"108.59"
      }
   ],
   "metadata":{
      "record_amount":2
   }
}
```

Currently supported HTTP GET query parameters are as listed:

- country
- orderBy
- order

For country there are currently 214 countries in the dataset with spaces as underscores:

```
"United_Kingdom"
"Australia"
"United_Arab_Emirates"
```

For order the values are:

```
"ASC"
"DESC"
```

For orderBy the values are:

```
"date"
"cases_weekly"
"deaths_weekly"
"country"
"population"
"notification_rate"
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

[Docker postgres](https://hub.docker.com/_/postgres) ([or for Windows without Docker](https://www.enterprisedb.com/downloads/postgres-postgresql-downloads))

## Migrations

DDL migrations (.sql files) should be placed under migrations\ddl and they'll be run when the application starts.
Prefix the migrations files as 001, 002, 003, ..., 999 and add a description with underscores.

Example:

```
001_create_record_table.sql
```

## Data source

---

"The downloadable data file was updated daily to 14 December 2020 using the latest available public data on COVID-19. Each row/entry contains the number of new cases reported per day and per country. After 14 December 2020, ECDC shifted to weekly data collection." [European Centre for Disease Prevention and Control](https://www.ecdc.europa.eu/en/publications-data/download-todays-data-geographic-distribution-covid-19-cases-worldwide)

---

[COVID-19 Coronavirus data](https://data.europa.eu/euodp/en/data/dataset/covid-19-coronavirus-data)

Can be applied to the database by running the application and then executing:

```
$ python utils\migrate.py
```
