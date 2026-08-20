package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/NodeFitter/sample/measurements"
	"github.com/NodeFitter/sample/measurements/abstract"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var meter abstract.Imeter

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received normal request")

	// Metrics
	if err := meter.Increment(); err != nil {
		log.Printf("failed to increment request counter: %v", err)
	}

	var number int

	// Get the current number.
	err := db.QueryRow("SELECT number FROM random_number WHERE id = 1").Scan(&number)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate the next random number.
	newNumber := rand.Intn(1000000)

	// Store it.
	_, err = db.Exec(
		"UPDATE random_number SET number = ? WHERE id = 1",
		newNumber,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Returning magic number: ", number)

	fmt.Fprintf(w, "%d\n", number)
}

func pingHandle(w http.ResponseWriter, _ *http.Request) {
	log.Println("Received ping request")

	// Metrics
	if err := meter.Increment(); err != nil {
		log.Printf("failed to increment request counter: %v", err)
		fmt.Fprintf(w, "pong (FAILURE - no counter incremented)")
		return
	}

	fmt.Fprintf(w, "pong (SUCCESS)")
}

func dbBypassCounter(w http.ResponseWriter, _ *http.Request) {
	log.Println("Received bypass request")

	var number int

	// Get the current number.
	err := db.QueryRow("SELECT number FROM random_number WHERE id = 1").Scan(&number)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate the next random number.
	newNumber := rand.Intn(1000000)

	// Store it.
	_, err = db.Exec(
		"UPDATE random_number SET number = ? WHERE id = 1",
		newNumber,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Returning magic number: ", number)

	fmt.Fprintf(w, "%d\n", number)
}

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "test:test@tcp(mariadb:3306)/test"
	}

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	muxWork := http.NewServeMux()
	muxMeter := http.NewServeMux()

	// Init for metrics
	meter = &measurements.Meter{}
	if err := meter.Init(muxMeter); err != nil {
		log.Fatal(err)
	}

	muxWork.HandleFunc("/", handler)
	muxWork.HandleFunc("/pingo", pingHandle)
	muxWork.HandleFunc("/bypassCounter", dbBypassCounter)

	go func() {
		log.Println("listening on :8080")
		log.Fatal(http.ListenAndServe(":8080", muxWork))
	}()

	log.Println("listening on :9090")
	log.Fatal(http.ListenAndServe(":9090", muxMeter))
}
