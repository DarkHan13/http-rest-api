package sqlstore_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "host=localhost port=5432 user=postgres password=7777 dbname=restapi_test sslmode=disable"
	}

	os.Exit(m.Run())
}
