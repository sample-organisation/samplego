package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
  "strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
)

var db *sql.DB

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

  pgopts := &dockertest.RunOptions{
    Name:       "testing-pg",
    Repository: "postgres",
    Tag:        "11",
    Env: []string{
      "POSTGRES_PASSWORD=postgrespassword",
      "POSTGRES_DB=postgres",
    },
    ExposedPorts: []string{"5432/tcp"},
  }
  pg, err := pool.RunWithOptions(pgopts)
  if err != nil {
    t.Fatalf("Could not start resource: %s", err)
  }
  var db *sql.DB
	if err = pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", fmt.Sprintf("postgres://postgres:postgrespassword@%s:%s/%s?sslmode=disable", "localhost", pg.GetPort("5432/tcp"), "postgres"))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		t.Fatal(err)
	}

	// // pulls an image, creates a container based on it and runs it
	// resource, err := pool.Run("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=secret"})
	// if err != nil {
	// 	log.Fatalf("Could not start resource: %s", err)
	// }

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	// if err := pool.Retry(func() error {
	// 	var err error
  //   p := resource.GetPort("3306/tcp")
  //   dHostFull := os.Getenv("DOCKER_HOST")
  //   dHostWithPort := strings.Replace(dHostFull, "tcp://", "", 1)
  //   dHost := strings.Replace(dHostWithPort, ":2376", "", 1)
  //   dHost = "localhost"
  //   fmt.Printf("p = %+v DOCKER_HOST = %s", p, dHost)
	// 	db, err = sql.Open("mysql", fmt.Sprintf("root:secret@(%s:%s)/mysql", dHost, p))
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return db.Ping()
	// }); err != nil {
	// 	log.Fatalf("Could not connect to docker: %s", err)
	// }
  //
	// code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
  if err = pool.Purge(pg); err != nil {
    t.Fatalf("Could not purge resource: %s", err)
  }

	os.Exit(code)
}
