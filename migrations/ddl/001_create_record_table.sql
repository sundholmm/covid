CREATE SEQUENCE IF NOT EXISTS record_id_seq;

CREATE TABLE IF NOT EXISTS record (
	"id" integer PRIMARY KEY NOT NULL DEFAULT nextval('record_id_seq'),
	"date" varchar(256) NULL,
	"year_week" varchar(256) NULL,
	"cases_weekly" int NULL,
	"deaths_weekly" int NULL,
	"country" varchar(256) NULL,
	"geo_id" varchar(256) NULL,
	"country_code" varchar(256) NULL,
	"population" int NULL,
	"continent" varchar(256) NULL,
	"notification_rate" varchar(256) NULL
);

CREATE INDEX IF NOT EXISTS record_index ON record ("id", "date", "year_week", "cases_weekly", "deaths_weekly", "country", "geo_id",
"country_code", "population", "continent", "notification_rate");