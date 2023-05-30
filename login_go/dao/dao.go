package dao

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	db, err := sql.Open("mysql", "root:123456@(127.0.0.1:3306)/pwdatabase?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

func CreateUserTable(db *sql.DB) {
	// Create a new table
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME,
		PRIMARY KEY (id)
	);`

	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}
}

func CreateUser(db *sql.DB, username string, password string) {
	// Insert a new user
	createdAt := time.Now()

	result, err := db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`, username, password, createdAt)
	if err != nil {
		// log.Fatal(err)
	}

	id, err := result.LastInsertId()
	fmt.Println(id)
}

func CheckUserPassword(db *sql.DB, username string, password string) (bool, error) {
	// Query a single user
	var (
		id        int
		username1 string
		password1 string
		createdAt time.Time
	)

	query := "SELECT id, username, password, created_at FROM users WHERE username = ?"
	if err := db.QueryRow(query, username).Scan(&id, &username1, &password1, &createdAt); err != nil {
		return false, err // database error
	}

	if username1 == username && password1 == password {
		return true, nil
	} else {
		return false, nil
	}
}

func CheckUserExist(db *sql.DB, username string) (bool, error) {
	var count int

	query := "SELECT COUNT(*) FROM users WHERE username = ?"
	err := db.QueryRow(query, username).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // user not found
		}
		return false, err // database error
	}

	return count > 0, nil
}

func Test() {

}
