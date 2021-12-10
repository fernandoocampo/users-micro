package postgresql_test

import (
	"database/sql"
	"flag"
	"testing"

	"github.com/fernandoocampo/users-micro/internal/adapter/postgresql"
	"github.com/stretchr/testify/assert"
)

var (
	integration = flag.Bool("integration", false, "run integration tests")
	dbport      = flag.Int("port", 5432, "postgresql database port")
	dbuser      = flag.String("user", "postgres", "postgresql database user")
	dbpassword  = flag.String("password", "postgres", "postgresql database password")
	dbname      = flag.String("database", "postgres", "postgresql database name")
	dbhost      = flag.String("host", "localhost", "postgresql database host")
)

func TestPostgresqlCreation(t *testing.T) {
	if !*integration {
		t.Skip("this is an integration test, to execute this test send integration flag to true")
	}

	t.Log("executing integration test to create a postgresql client")

	givenParameters := postgresql.Parameters{
		DBName:   *dbname,
		Host:     *dbhost,
		User:     *dbuser,
		Password: *dbpassword,
		Port:     *dbport,
	}

	client, err := postgresql.NewPostgresClient(givenParameters)
	assert.NoError(t, err)
	if err != nil {
		t.Errorf("unexpected error trying to connect to default database: %s", err)
		t.FailNow()
	}
	client.Close()
}

func createClient(t *testing.T, parameters postgresql.Parameters) *sql.DB {
	t.Helper()
	client, err := postgresql.NewPostgresClient(parameters)
	assert.NoError(t, err)
	if err != nil {
		t.Errorf("unexpected error trying to connect to default database: %s", err)
		t.FailNow()
	}
	return client
}
