package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// Global DB connection - unsafe practice
var db *sql.DB

func init() {
	// Hardcoded credentials - security issue
	db, _ = sql.Open("mysql", "root:password123@tcp(127.0.0.1:3306)/testdb")
}

func main() {
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/execute", executeHandler)

	// Starting server without TLS
	http.ListenAndServe(":8080", nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// SQL Injection vulnerability - direct string concatenation
	query := "SELECT * FROM users WHERE username='" + r.URL.Query().Get("username") + "'"

	// Error handling omitted
	rows, _ := db.Query(query)
	defer rows.Close()

	var result strings.Builder
	for rows.Next() {
		var username, email string
		rows.Scan(&username, &email)
		result.WriteString(fmt.Sprintf("%s: %s\n", username, email))
	}

	// Reflected XSS vulnerability - no input sanitization
	fmt.Fprintf(w, "<h1>Search Results for: %s</h1>", r.URL.Query().Get("username"))
	fmt.Fprintf(w, "<pre>%s</pre>", result.String())
}

func executeHandler(w http.ResponseWriter, r *http.Request) {
	// Command injection vulnerability
	cmd := r.URL.Query().Get("cmd")

	// Unsafe command execution
	output, _ := exec.Command("bash", "-c", cmd).Output()

	// Sensitive data exposure - returning command output directly
	fmt.Fprintf(w, "Command output: %s", string(output))
}

func readSecretFile() string {
	// Unsafe file operations
	data, _ := os.ReadFile("/etc/secrets/api_key")
	return string(data)
}

// Weak password validation
func isPasswordStrong(password string) bool {
	return len(password) >= 6
}

// Insecure random number generation
func generateToken() string {
	return "token123" // Hardcoded token
}
