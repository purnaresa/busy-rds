package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	faker "github.com/bxcodec/faker/v3"
	_ "github.com/lib/pq"
)

var (
	dbPool *sql.DB
)

type Config struct {
	Database struct {
		DSN string `json:"dsn"`
	} `json:"database"`
}

func init() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	dbPool, err = sql.Open("postgres", config.Database.DSN)
	if err != nil {
		log.Fatal(err)
	}
	err = dbPool.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to the database")
}

func loadConfig(filename string) (Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	return config, err
}

func main() {
	createDatabase(dbPool, "main")
	createTableIfNotExists(dbPool)
	generateDummyData(dbPool, 10)
	readBusyTable(dbPool)
}

func createDatabase(db *sql.DB, dbName string) {
	// First, check if the database exists
	rows, err := db.Query(fmt.Sprintf("SELECT datname FROM pg_catalog.pg_database WHERE lower(datname) = lower('%s')", dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// If the database does not exist, create it
	if !rows.Next() {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Database %s created successfully\n", dbName)
	} else {
		fmt.Printf("Database %s already exists\n", dbName)
	}
}

func createTableIfNotExists(db *sql.DB) {

	stmt := `CREATE TABLE IF NOT EXISTS busy(
		id SERIAL PRIMARY KEY,
		description VARCHAR(255),
		status VARCHAR(50),
		time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table busy created successfully")
}

func generateDummyData(db *sql.DB, n int) {
	stmt := `INSERT INTO busy (description, status) VALUES ($1, $2)`

	for i := 0; i < n; i++ {
		_, err := db.Exec(stmt, faker.Sentence(), "idle")
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("%d rows of dummy data inserted into busy table\n", n)
}

type BusyRow struct {
	ID          int
	Description string
	Status      string
	Time        string
}

func readBusyTable(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM busy")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var r BusyRow
		err = rows.Scan(&r.ID, &r.Description, &r.Status, &r.Time)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Description: %s, Status: %s, Time: %s\n", r.ID, r.Description, r.Status, r.Time)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
