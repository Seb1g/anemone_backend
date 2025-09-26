package notes_repository;
import (
	"database/sql"
	"github.com/joho/godotenv"
	"log"
	"os"
	_ "github.com/lib/pq"
)

func Database_conn () (*sql.DB, error) {
	error_dotenv := godotenv.Load();

	if error_dotenv != nil {
		log.Fatalf("Error load dotenv file: %v", error_dotenv);
		return nil, error_dotenv
	};

	connStr := os.Getenv("conn_str");

	db, error_db_connect := sql.Open("postgres", connStr);

	if error_db_connect != nil {
		log.Fatalf("Error connect to database: %v", error_db_connect);
		return nil, error_db_connect
	};

	return db, nil
};