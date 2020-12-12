package migrations

import (
	"database/sql"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

// MigrateDB finds the last run migration
// and runs all the migrations after it in order
func MigrateDB(db *sql.DB) {
	var migrations []string
	var completed []string

	// Create the migration table
	_, err := db.Exec("CREATE SEQUENCE IF NOT EXISTS migration_id_seq; " +
		"CREATE TABLE IF NOT EXISTS \"migration\" " + 
		"( id integer PRIMARY KEY NOT NULL DEFAULT nextval('migration_id_seq'), version varchar(256), summary varchar(256) )")
	if err != nil {
		log.Println(err)
	}

	// Get all the migration files
	files, err := filepath.Glob("./migrations/ddl/*")
	if err != nil {
		log.Fatalln(err)
	}

	sort.Strings(files)

	// Get all the existing migrations
	migrations = readMetadata(db)

	// Run the migrations conditionally
	for _, file := range files {
		filename := filepath.Base(file)

		if !contains(filename[0:3], migrations) {
			log.Printf("Running migration %s", filename)

			sql, err := ioutil.ReadFile("./migrations/ddl/" + filename)
			if err != nil {
				log.Fatalln(err)
			}

			db.Exec(string(sql))
			if err != nil {
				log.Fatalln(err)
			}

			writeMetadata(db, filename)

			completed = append(completed, filename)
			log.Printf("Completed migration %s", filename)
		}
	}

	// Log the final status
	if len(completed) > 0 {
		log.Printf("%d migration(s) completed", len(completed))
	} else {
		log.Println("No migrations to run")
	}

}

// readMetadata reads the metadata from the migration table
func readMetadata(db *sql.DB) []string {
	var migrations []string

	// Get the existing migration versions
	rows, err := db.Query("SELECT \"version\" FROM \"migration\" ORDER BY \"id\" DESC")
	if err != nil {
		log.Fatalln(err)
	}

	defer rows.Close()

	// Append all the versions to migrations array
	for rows.Next() {

		var migration string

		err := rows.Scan(&migration)

		if err != nil {
			log.Fatalln(err)
		}

		migrations = append(migrations, migration)

	}

	return migrations
}

// writeMetadata writes all the metadata of migrations to migration table
func writeMetadata(db *sql.DB, filename string) {

		result, err := db.Exec("INSERT INTO \"migration\" (\"version\", \"summary\") " +
		"VALUES ('" + filename[0:3] + "', '" + strings.TrimSuffix(filename[4:], ".sql") + "')")
		if err != nil {
			log.Printf("Database ERROR %s %s", err, result)
		}

}

// contains checks whether an array of strings contains a string
func contains(s string, a []string) bool {
	for _, k := range a {
		if s == k {
			return true
		}
	}
	return false
}