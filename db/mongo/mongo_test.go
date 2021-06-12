package mongo

import (
	"database/sql"
	"github.com/ory/dockertest/v3"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mongo", "4.2.9", []string{"MONGO_INITDB_DATABASE=test"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err = sql.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql", resource.GetPort("3306/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
	//pool, err := dockertest.NewPool("")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//resource, err := pool.Run("mongo", "4.2.9", []string{
	//	"MONGO_INITDB_DATABASE=roava",
	//})
	//if err != nil {
	//	log.Fatalf("Could not start resource: %s", err)
	//}
	//
	//mongoDbPort = resource.GetPort("27017/tcp")
	//if err := pool.Retry(func() error {
	//	var err error
	//	connectUrl := fmt.Sprintf("mongodb://localhost:%s", mongoDbPort)
	//	_, _, err = New(connectUrl, "roava")
	//	if err != nil {
	//		return err
	//	}
	//	return nil
	//}); err != nil {
	//	log.Fatalf("Could not connect to docker: %s", err)
	//}
	//code := m.Run()
	//
	//err = pool.Purge(resource)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//os.Exit(code)
}
