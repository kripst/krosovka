package postgres_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kripst/krosovka/inventory_service/migrate"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

var TestDbPool *pgxpool.Pool

type MockPostgresStorageImpl struct {
	pool  *pgxpool.Pool //postgres
	log *zap.Logger
	sq  squirrel.StatementBuilderType
}

func NewMockPostgresStorageImpl(pool *pgxpool.Pool, log *zap.Logger) *MockPostgresStorageImpl {
	return &MockPostgresStorageImpl{
		pool: pool,
		log: log,
		sq: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// Эта функция будет вызываться перед запуском тестов в пакете
func TestMain(m *testing.M) {
	ctx := context.Background()

	// 1. Создаем запрос на контейнер PostgreSQL
	// Мы указываем образ, имя пользователя, пароль и имя БД
	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("could not start postgres container: %s", err)
	}

	// Обязательно останавливаем и удаляем контейнер после всех тестов
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	// 2. Получаем динамический connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("could not get connection string: %s", err)
	}
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("could not connect to database: %s", err)
	}
	TestDbPool = pool

    log.Println(connStr)
    // Здесь вы должны накатить миграции на тестовую БД
    if err := migrate.RunMigrations(connStr); err != nil {
		log.Fatalf("could not run migrations: %s", err)
	}

	// Теперь можно запускать тесты
	exitCode := m.Run()
    
    // Выходим с кодом завершения тестов
	os.Exit(exitCode)
}

