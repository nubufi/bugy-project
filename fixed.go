package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var (
	db *sql.DB
	wg sync.WaitGroup
)

func main() {
	var err error
	db, err = sql.Open("postgres", "user=postgres dbname=test sslmode=disable")
	// check if the connection failed
	if err != nil {
		log.Fatal(err)
	}
	// Make sure the database is closed when the program exits
	defer db.Close()

	// Make sure the database is reachable
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/create", createUser)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT name FROM users")
	// check if the query failed
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var name string
		// check if the scan failed
		if err := rows.Scan(&name); err != nil {
			http.Error(w, "Failed to scan user", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "User: %s\n", name)
	}

	// check if the iteration failed
	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to iterate over users", http.StatusInternalServerError)
		return
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second) // Simulate a long database operation
	username := r.URL.Query().Get("name")

	if username == "" {
		http.Error(w, "Missing name parameter", http.StatusBadRequest)
		return
	}

	// Use parameterized queries to prevent SQL injection
	_, err := db.Exec("INSERT INTO users (name) VALUES ($1)", username)
	if err != nil {
		fmt.Fprintf(w, "Failed to create user: %v", err)
		return
	}

	fmt.Fprintf(w, "User %s created successfully", username)
}
