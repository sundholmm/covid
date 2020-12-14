CREATE SEQUENCE IF NOT EXISTS record_id_seq;

CREATE TABLE IF NOT EXISTS record (
	"id" integer PRIMARY KEY NOT NULL DEFAULT nextval('record_id_seq'),
	"date" varchar(256) NULL,
	"day" varchar(256) NULL,
	"month" varchar(256) NULL,
	"year" varchar(256) NULL,
	"cases" int NULL,
	"deaths" int NULL,
	"country" varchar(256) NULL,
	"geo_id" varchar(256) NULL,
	"country_code" varchar(256) NULL,
	"population" int NULL,
	"continent" varchar(256) NULL,
	"cumulative" varchar(256) NULL
);

CREATE INDEX IF NOT EXISTS record_index ON record ("id", "date", "day", "month", "year", "cases", "deaths", "country", "geo_id",
"country_code", "population", "continent", "cumulative");