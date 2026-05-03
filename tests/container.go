package tests

import (
	"fmt"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/moby/moby/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/redis"
)

type ContainerTestSuite struct {
	GlobalTestSuite
	RedisClient    *redis.Client
	RedisContainer testcontainers.Container
	PgClient       *postgres.Client
	PgContainer    testcontainers.Container
}

var postgresContainer testcontainers.Container
var redisContainer testcontainers.Container

func (t *ContainerTestSuite) createContainerPostgres() {
	testValue := "test"
	_ = os.Setenv("DB_NAME", testValue)

	schema := uuid.NewString()
	_ = os.Setenv("DB_SCHEMA", schema)

	_ = os.Setenv("DB_PASSWORD", testValue)
	_ = os.Setenv("DB_CONN_MAX_OPEN", "100")
	_ = os.Setenv("DB_USERNAME", testValue)

	port := "5432/tcp"

	if postgresContainer == nil {
		var err error

		postgresContainer, err = testcontainers.GenericContainer(t.Ctx, testcontainers.GenericContainerRequest{
			Started: true,
			Reuse:   true,
			ContainerRequest: testcontainers.ContainerRequest{
				ExposedPorts: []string{port},
				Image:        "postgres:18-alpine",
				Name:         fmt.Sprintf("postgres-test-%s", schema),
				WaitingFor:   wait.ForListeningPort(port),
				Env: map[string]string{
					"POSTGRES_USER":     testValue,
					"POSTGRES_PASSWORD": testValue,
					"POSTGRES_DB":       testValue,
				},
				HostConfigModifier: func(config *container.HostConfig) {
					config.AutoRemove = true
				},
			},
		})

		t.Require().Nil(err)
	}

	t.PgContainer = postgresContainer

	host, err := postgresContainer.Host(t.Ctx)
	t.Require().Nil(err)
	_ = os.Setenv("DB_HOST", host)

	mappedPort, err := postgresContainer.MappedPort(t.Ctx, port)
	t.Require().Nil(err)
	_ = os.Setenv("DB_PORT", mappedPort.Port())

	t.PgClient = postgres.FromEnv(t.Ctx, t.Logger)
}

func (t *ContainerTestSuite) createContainerRedis() {
	password := "test"
	_ = os.Setenv("REDIS_PASSWORD", password)
	port := "6379/tcp"

	if redisContainer == nil {
		var err error

		redisContainer, err = testcontainers.GenericContainer(t.Ctx, testcontainers.GenericContainerRequest{
			Reuse:   true,
			Started: true,
			ContainerRequest: testcontainers.ContainerRequest{
				Name:         fmt.Sprintf("redis-test-%s", uuid.NewString()),
				Image:        "bitnami/redis:latest",
				ExposedPorts: []string{port},
				WaitingFor:   wait.ForListeningPort(port),
				Env: map[string]string{
					"ALLOW_EMPTY_PASSWORD": "no",
					"REDIS_PASSWORD":       password,
				},
				HostConfigModifier: func(config *container.HostConfig) {
					config.AutoRemove = true
				},
			},
		})

		t.Require().Nil(err)
	}

	t.RedisContainer = redisContainer

	host, err := redisContainer.Host(t.Ctx)
	t.Require().Nil(err)
	_ = os.Setenv("REDIS_HOST", host)

	mappedPort, err := redisContainer.MappedPort(t.Ctx, port)
	t.Require().Nil(err)
	_ = os.Setenv("REDIS_PORT", mappedPort.Port())

	t.RedisClient = redis.FromEnv(t.Ctx)
}

func (t *ContainerTestSuite) SetupSuite() {
	t.GlobalTestSuite.SetupSuite()
	t.createContainerPostgres()
	t.createContainerRedis()
}

func (t *ContainerTestSuite) TearDownSuite() {
	// t.Require().Nil(t.PgClient.Close())
	// _ = t.PgContainer.Terminate(t.Ctx)

	// t.Require().Nil(t.RedisClient.Close())
	// _ = t.RedisContainer.Terminate(t.Ctx)
}
