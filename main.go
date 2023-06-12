package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	faker "github.com/bxcodec/faker/v3"
	_ "github.com/lib/pq"
)

var (
	dbPool *sql.DB
	config Config
	delay  float64
)

type Config struct {
	Database struct {
		DSN string `json:"dsn"`
	} `json:"database"`
	TestRun    int `json:"test_run"`
	RPS        int `json:"rps"`
	MaxRetry   int `json:"max_retry"`
	DelayRetry int `json:"delay_retry"`
}

func init() {

	err := error(nil)
	config, err = loadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	delay = 1000 / float64(config.RPS)
	log.Printf("delay : %v", delay)

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
	log.Printf("Config : %+v", config)
	return config, err
}

func main() {
	createDatabase(dbPool, "main")
	createTableIfNotExists(dbPool)
	writeData()

}

func writeData() {
	log.Printf("================\n Start Write Data Simulation")
	startTime := time.Now()
	log.Println("Start time: ", startTime)
	for i := 0; i < config.TestRun; i++ {
		generateDummyData(dbPool, 1)
		time.Sleep(time.Duration(delay) * time.Millisecond)

	}
	duration := time.Since(startTime)
	log.Printf("End time: %v \n", time.Now())
	log.Printf("Duration: %s\n", duration)
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
	log.Println("Table busy is ready")
}

func generateDummyData(db *sql.DB, n int) {
	stmt := `INSERT INTO busy (description, status) VALUES ($1, $2)`

	for i := 0; i < n; i++ {

		data := faker.Email()
		retryCount := 1
		lastError := time.Now()
		for {
			_, err := db.Exec(stmt, data, "idle")
			if err != nil {
				if retryCount == 1 {
					lastError = time.Now()
				}
				if retryCount >= config.MaxRetry {
					log.Fatalf("Failed to insert data: %s. Error: %v\n", data, err)
				}

				log.Printf("Failed to insert: %s. Error: %v. Retrying (%d/%d)...\n", data, err, retryCount, config.MaxRetry)
				retryCount++
				time.Sleep(time.Duration(config.DelayRetry) * time.Second)

			} else {
				if retryCount > 1 {
					downTime := time.Since(lastError).Milliseconds()
					log.Printf("DownTime: %dms\n", downTime)
				}
				log.Printf("Insert: %s success\n", data)
				break
			}
		}
	}

}

type BusyRow struct {
	ID          int
	Description string
	Status      string
	Time        string
}

func readBusyTable(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM busy ORDER BY id desc limit 1")
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
