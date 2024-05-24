package containers

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func Postgres(ctx context.Context, name string) (*postgres.PostgresContainer, error) {
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:         "kanthorlabs-common-postgres",
			Image:        "postgres:16-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     os.Getenv("TEST_CONTAINER_POSTGRES_USER"),
				"POSTGRES_PASSWORD": os.Getenv("TEST_CONTAINER_POSTGRES_PASSWORD"),
				"POSTGRES_DB":       os.Getenv("TEST_CONTAINER_POSTGRES_USER"),
			},
			Cmd: []string{"postgres", "-c", "fsync=off"},
			WaitingFor: wait.
				ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5 * time.Second),
		},
		Started: true,
		Reuse:   true,
		Logger:  &Logger{},
	}
	if name != "" {
		req.ContainerRequest.Name = name
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, err
	}

	return &postgres.PostgresContainer{Container: container}, nil
}

func PostgresConnectionString(ctx context.Context, container *postgres.PostgresContainer) (string, error) {
	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return "", err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return "", err
	}

	user := os.Getenv("TEST_CONTAINER_POSTGRES_USER")
	pass := os.Getenv("TEST_CONTAINER_POSTGRES_PASSWORD")

	uri := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, pass, net.JoinHostPort(host, port.Port()), user)
	return uri, nil
}
