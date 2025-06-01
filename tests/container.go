package tests

import (
	"fmt"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/redis"
)

const (
	waitForLogPg    = "database system is ready to accept connections"
	waitForLogRedis = "Ready to accept connections tcp"
)

type ContainerTestSuite struct {
	GlobalTestSuite
	RedisClient    *redis.Client
	RedisContainer testcontainers.Container
	PgClient       *postgres.Client
	PgContainer    testcontainers.Container
}

func (t *ContainerTestSuite) createContainerPostgres() {
	testValue := "test"
	_ = os.Setenv("DB_NAME", testValue)

	schema := uuid.NewString()
	_ = os.Setenv("DB_SCHEMA", schema)

	_ = os.Setenv("DB_PASSWORD", testValue)
	_ = os.Setenv("DB_MAX_OPEN_CONN", "100")
	_ = os.Setenv("DB_USERNAME", testValue)

	port := "5432/tcp"
	postgresContainer, err := testcontainers.GenericContainer(t.Ctx, testcontainers.GenericContainerRequest{
		Started: true,
		ContainerRequest: testcontainers.ContainerRequest{
			Name:         fmt.Sprintf("postgres-test-%s", schema),
			Image:        "bitnami/postgresql:16",
			ExposedPorts: []string{port},
			WaitingFor:   wait.ForLog(waitForLogPg),
			Env: map[string]string{
				"POSTGRESQL_PASSWORD":                 testValue,
				"POSTGRESQL_USERNAME":                 testValue,
				"POSTGRESQL_REPLICATION_USE_PASSFILE": "no",
				"POSTGRESQL_DATABASE":                 testValue,
			},
			HostConfigModifier: func(config *container.HostConfig) {
				config.AutoRemove = true
			},
		},
	})

	t.Require().Nil(err)
	t.PgContainer = postgresContainer

	host, err := postgresContainer.Host(t.Ctx)
	t.Require().Nil(err)
	_ = os.Setenv("DB_HOST", host)

	mappedPort, err := postgresContainer.MappedPort(t.Ctx, nat.Port(port))
	t.Require().Nil(err)
	_ = os.Setenv("DB_PORT", mappedPort.Port())

	t.PgClient = postgres.FromEnv(t.Ctx, t.Logger)
}

func (t *ContainerTestSuite) createContainerRedis() {
	password := "test"
	_ = os.Setenv("REDIS_PASSWORD", password)
	port := "6379/tcp"

	redisContainer, err := testcontainers.GenericContainer(t.Ctx, testcontainers.GenericContainerRequest{
		Started: true,
		ContainerRequest: testcontainers.ContainerRequest{
			Name:         fmt.Sprintf("redis-test-%s", uuid.NewString()),
			Image:        "bitnami/redis:latest",
			ExposedPorts: []string{port},
			WaitingFor:   wait.ForLog(waitForLogRedis),
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
	t.RedisContainer = redisContainer

	host, err := redisContainer.Host(t.Ctx)
	t.Require().Nil(err)
	_ = os.Setenv("REDIS_HOST", host)

	mappedPort, err := redisContainer.MappedPort(t.Ctx, nat.Port(port))
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
	t.Require().Nil(t.PgClient.Close())
	_ = t.PgContainer.Terminate(t.Ctx)

	t.Require().Nil(t.RedisClient.Close())
	_ = t.RedisContainer.Terminate(t.Ctx)
}
