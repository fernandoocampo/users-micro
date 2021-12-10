package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // just to load drivers
)

// Parameters contains postgres information
type Parameters struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     int
}

// NewPostgresClient creates a new postgresql client.
func NewPostgresClient(parameters Parameters) (*sql.DB, error) {
	psqlInfo := buildPostgresqlConnection(parameters)
	pgsqlconn, err := connectToDatabase(psqlInfo)
	if err != nil {
		return nil, err
	}

	return pgsqlconn, nil
}

func buildPostgresqlConnection(dbParameters Parameters) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbParameters.Host,
		dbParameters.Port,
		dbParameters.User,
		dbParameters.Password,
		dbParameters.DBName,
	)
}

// connectToDatabase creates a connection to postgresql database based on given client parameters.
func connectToDatabase(psqlInfo string) (*sql.DB, error) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	// ensure connection calling ping method
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
