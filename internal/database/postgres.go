package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/say8hi/go-jwt-api/internal/models"
)

var db *sql.DB

func Init() {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			fmt.Println("Successfully connected to database")
			return
		}
		fmt.Printf("Couldn't connect to databse: %v. Retry in 1 second.\n", err)
		time.Sleep(1 * time.Second)
	}
}

func CloseConnection() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Successfully closed connection to database")
	}
}

func CreateTables() {
	tables := [2]string{
    `
      CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
    `,
		`
        CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            username TEXT NOT NULL UNIQUE,
            email TEXT,
            refresh_hash TEXT,
            ip TEXT
        );
    `,
	}

	for _, sql := range tables {
		_, err := db.Exec(sql)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Table Users
func CreateUser(request_user models.CreateUserRequest) (models.UserInDatabase, error) {
	var user models.UserInDatabase
	err := db.QueryRow("INSERT INTO users (username, email, ip) VALUES($1, $2, $3) RETURNING *",
		request_user.Username, request_user.Email, request_user.IP).Scan(&user.ID, &user.Username, &user.Email, &user.RefreshHash, &user.IP)

	if err != nil {
		return models.UserInDatabase{}, err
	}

	return user, nil
}

func GetUserByID(userID string) (models.UserInDatabase, error) {
	var user models.UserInDatabase

	err := db.QueryRow("SELECT * FROM users WHERE id=$1",
		userID).Scan(&user.ID, &user.Username, &user.Email, &user.RefreshHash, &user.IP)

	if err != nil {
		return models.UserInDatabase{}, err
	}

	return user, nil
}


func UpdateUser(userID string, updateReq models.UserUpdateRequest) error {
	var setParts []string
	var args []interface{}
	var argIndex int = 1

	if updateReq.Username != "" {
		setParts = append(setParts, fmt.Sprintf("username = $%d", argIndex))
		args = append(args, updateReq.Username)
		argIndex++
	}
	if updateReq.IP != "" {
		setParts = append(setParts, fmt.Sprintf("ip = $%d", argIndex))
		args = append(args, updateReq.IP)
		argIndex++
	}
  if updateReq.RefreshHash != "" {
		setParts = append(setParts, fmt.Sprintf("refresh_hash = $%d", argIndex))
		args = append(args, updateReq.RefreshHash)
		argIndex++
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	setClause := strings.Join(setParts, ", ")
	queryString := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", setClause, argIndex)
	args = append(args, userID)

	_, err := db.Exec(queryString, args...)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}
