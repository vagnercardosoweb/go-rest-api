package tests

import (
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/redis"
	"os"
	"path/filepath"
	"runtime"
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
	container, err := testcontainers.GenericContainer(t.Ctx, testcontainers.GenericContainerRequest{
		Started: true,
		ContainerRequest: testcontainers.ContainerRequest{
			Name:         fmt.Sprintf("postgres-test-%s", schema),
			Image:        "bitnami/postgresql:16",
			ExposedPorts: []string{port},
			WaitingFor:   wait.ForLog(waitForLogPg),
			Env: map[string]string{
				"POSTGRESQL_PASSWORD": testValue,
				"POSTGRESQL_USERNAME": testValue,
				"POSTGRESQL_DATABASE": testValue,
			},
		},
	})

	t.Require().Nil(err)
	t.PgContainer = container

	host, err := container.Host(t.Ctx)
	t.Require().Nil(err)
	_ = os.Setenv("DB_HOST", host)

	mappedPort, err := container.MappedPort(t.Ctx, nat.Port(port))
	t.Require().Nil(err)
	_ = os.Setenv("DB_PORT", mappedPort.Port())

	t.PgClient = postgres.NewFromEnv(t.Ctx, t.Logger)
	_, err = t.PgClient.Exec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s";`, schema))
	t.Require().Nil(err)

	driver, err := migratepostgres.WithInstance(t.PgClient.GetDb(), &migratepostgres.Config{
		MigrationsTable: "migrations",
		DatabaseName:    testValue,
		SchemaName:      schema,
	})
	t.Require().Nil(err)

	_, file, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(file), "..")

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+basePath+"/migrations",
		"postgres",
		driver,
	)
	t.Require().Nil(err)

	err = m.Up()
	t.Require().Nil(err)
}

func (t *ContainerTestSuite) createContainerRedis() {
	password := "test"
	_ = os.Setenv("REDIS_PASSWORD", password)
	port := "6379/tcp"

	container, err := testcontainers.GenericContainer(t.Ctx, testcontainers.GenericContainerRequest{
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
		},
	})

	t.Require().Nil(err)
	t.RedisContainer = container

	host, err := container.Host(t.Ctx)
	t.Require().Nil(err)
	_ = os.Setenv("REDIS_HOST", host)

	mappedPort, err := container.MappedPort(t.Ctx, nat.Port(port))
	t.Require().Nil(err)
	_ = os.Setenv("REDIS_PORT", mappedPort.Port())

	t.RedisClient = redis.NewFromEnv(t.Ctx, t.Logger)
}

func (t *ContainerTestSuite) SetupSuite() {
	t.GlobalTestSuite.SetupSuite()
	t.createContainerPostgres()
	t.createContainerRedis()
}

func (t *ContainerTestSuite) TearDownSuite() {
	//_ = t.PgClient.Close()
	//_ = t.PgContainer.Terminate(t.Ctx)

	//_ = t.RedisClient.Close()
	//_ = t.RedisContainer.Terminate(t.Ctx)
}
