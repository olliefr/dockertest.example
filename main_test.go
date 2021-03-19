package main

import (
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

const (
	postgresDefaultVersionTag = "13.0"
)

func TestDockerTest(t *testing.T) {

	postgresVersionTag := os.Getenv("PGTAG")
	if postgresVersionTag == "" {
		postgresVersionTag = postgresDefaultVersionTag
		log.Printf("Postgres Docker image tag (PGTAG) not set, defaulting to %s", postgresDefaultVersionTag)
	}

	postgresPassword := os.Getenv("PGPASSWORD")
	if postgresPassword == "" {
		log.Fatal("Panic: Postgres password (PGPASSWORD) not set")
	}

	postgresDatabase := os.Getenv("PGDATABASE")
	if postgresDatabase == "" {
		log.Fatal("Panic: Postgres database (PGDATABASE) not set")
	}

	log.Print("Creating a Docker pool...")
	pool, err := dockertest.NewPool("")

	if err != nil {
		log.Fatalf("Panic: Failed to create the Docker pool: %s", err)
	}

	log.Print("Starting a container...")
	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        postgresVersionTag,
			Env: []string{
				"POSTGRES_PASSWORD=" + postgresPassword,
				"POSTGRES_DB=" + postgresDatabase,
			},
		},
		func(hc *docker.HostConfig) {
			hc.AutoRemove = true
			hc.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		},
	)

	if err != nil {
		log.Fatalf("Panic: Failed to start the container: %s", err)
	}

	t.Cleanup(func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("Panic: Failed to remove a container or a volume: %s", err)
		}
	})
	resource.Expire(60) // Tell Docker to hard kill the container in X seconds

	log.Println("OK, I suppose...")
}
