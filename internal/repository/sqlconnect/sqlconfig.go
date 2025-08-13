package sqlconnect

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func ConnectDb() (*sql.DB, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("error opening up dot env file")
		return nil, err
	}

	connectionString := os.Getenv("DB_CONNECT")

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to mariadb")
	return db, nil
}
