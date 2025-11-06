package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func parseLocationCount(locationCode string, body []byte) (int, error) {
	// yuck
	pattern := fmt.Sprintf(`'%s' : \{\s*'capacity' : \d+,\s*'count' : (\d+),`, locationCode)
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return -1, err
	}

	matches := regex.FindStringSubmatch(string(body))
	if len(matches) > 1 {
		count, err := strconv.Atoi(matches[1])
		if err != nil {
			return -1, err
		}
		return count, nil
	}

	return -1, fmt.Errorf("no count found for %s", locationCode)
}

func getOccupancy() (int, int, error) {
	resp, err := http.Get("https://portal.rockgympro.com/portal/public/8490bc5e774d0034d09420df23a224b9/occupancy?=&iframeid=occupancyCounter&fId=")
	if err != nil {
		return -1, -1, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, -1, err
	}

	ptmCount, ptmErr := parseLocationCount("PTM", body)
	scbCount, scbErr := parseLocationCount("SCB", body)

	if ptmErr != nil && scbErr != nil {
		return -1, -1, errors.New("no count found in response")
	}

	return ptmCount, scbCount, nil
}

func initDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS occupancy (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME NOT NULL,
		location TEXT NOT NULL,
		count INTEGER NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func saveOccupancy(db *sql.DB, location string, count int) error {
	insertSQL := `INSERT INTO occupancy (timestamp, location, count) VALUES (?, ?, ?)`
	_, err := db.Exec(insertSQL, time.Now(), location, count)
	return err
}

func main() {
	db, err := initDB("gym_occupancy.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	ptmCount, scbCount, err := getOccupancy()
	if err != nil {
		log.Fatalf("Failed to get occupancy: %v", err)
	}

	fmt.Printf("Portsmouth occupancy: %d\n", ptmCount)
	fmt.Printf("Scarborough occupancy: %d\n", scbCount)

	if err := saveOccupancy(db, "PTM", ptmCount); err != nil {
		log.Printf("Failed to save PTM occupancy: %v", err)
	}

	if err := saveOccupancy(db, "SCB", scbCount); err != nil {
		log.Printf("Failed to save SCB occupancy: %v", err)
	}

	fmt.Println("Data saved successfully to gym_occupancy.db")
}
