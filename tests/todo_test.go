package tests

import (
	// "database/sql"
	"database/sql"
	"docker-test/todo"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
)

var Conn *sql.DB

func TestMain(m *testing.M) {

	// Create a new pool of connections
	pool, err := dockertest.NewPool("myPool")

	// Connect to docker on local machine
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Spawning a new MySQL docker container. We pass desired credentials for accessing it.
	resource, err := pool.Run("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=secret"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	connected := false
	// Try connecting for 200secs
	for i := 0; i < 20; i++ {
		// Try establishing MySQL connection.
		Conn, err = sql.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql?parseTime=true", resource.GetPort("3306/tcp")))
		if err != nil {
			panic(err)
		}

		err = Conn.Ping()
		if err != nil {
			// connection established success
			connected = true
			break
		}

		// Sleep for 10 sec and try again.
	}

	if !connected {
		fmt.Println("Couldnt connect to SQL")
		pool.Purge(resource)
	}

	// Run our unit test cases
	code := m.Run()

	// Purge our docker containers
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestSomething(t *testing.T) {
	assert := assert.New(t)
	listRepo := todo.NewListRepository(Conn)
	err := listRepo.Init()
	assert.NoError(err, "error while initializing repo")

	ID := uuid.New().String()
	name := "itemName 1"

	// Save Item
	err = listRepo.SaveItem(ID, name)
	assert.NoError(err, "error while saving item")

	// Find Item
	foundName, err := listRepo.FindItem(ID)
	assert.NoError(err, "error while saving item")
	assert.Equal(foundName, name)

	// Delete Item
	err = listRepo.DeleteItem(ID)
	assert.NoError(err, "error while saving item")
	foundName, err = listRepo.FindItem(ID)
	assert.Error(err, "delete unsuccessful")
}

func TestE2E(t *testing.T) {
	assert := assert.New(t)
	client := &http.Client{}

	req, err := http.NewRequest("GET", "localhost:8080/save?item=test122", nil)
	assert.NoError(err, "error while initializing request")

	res, err := client.Do(req)
	assert.NoError(err, "error while making request")
	assert.Equal(res.StatusCode, 200)
	... // Other E2E requests
}
