package postgres_test

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kripst/krosovka/inventory_service/internal/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// Тест №1: Успешная вставка нескольких записей (happy path).
func TestCreateSneakers_Success(t *testing.T) {
	// --- Arrange ---
	require := require.New(t)
	ctx := context.Background()
	storage := NewMockPostgresStorageImpl(TestDbPool, zap.NewNop())

	sneakersToCreate := []*model.Sneaker{
		{ID: 1, Article: "ART-001", SneakerName: "Runner Pro", Price: 150.00, Brand: "Nike"},
		{ID: 2, Article: "ART-002", SneakerName: "Classic", Price: 120.50, Brand: "Adidas"},
	}

	// Очистка таблицы после теста
	t.Cleanup(func() {
		_, err := TestDbPool.Exec(ctx, "TRUNCATE TABLE sneakers RESTART IDENTITY CASCADE")
		require.NoError(err)
	})

	// --- Act ---
	err := storage.CreateSneakers(ctx, sneakersToCreate)

	// --- Assert ---
	require.NoError(err)

	// Проверяем, что данные действительно появились в БД
	var count int
	err = TestDbPool.QueryRow(ctx, "SELECT COUNT(*) FROM sneakers").Scan(&count)
	require.NoError(err)
	require.Equal(2, count, "должно быть вставлено 2 записи")
}



// Тест №2: Обработка пустого среза на входе.
func TestCreateSneakers_EmptySlice(t *testing.T) {
	// --- Arrange ---
	require := require.New(t)
	ctx := context.Background()
	storage := NewMockPostgresStorageImpl(TestDbPool, zap.NewNop())
	
	sneakersToCreate := []*model.Sneaker{}
	t.Cleanup(func() {
		_, err := TestDbPool.Exec(ctx, "TRUNCATE TABLE sneakers RESTART IDENTITY CASCADE")
		require.NoError(err)
	})
	
	// --- Act ---
	err := storage.CreateSneakers(ctx, sneakersToCreate)

	// --- Assert ---
	require.NoError(err, "ошибки быть не должно, если срез пустой")
	
	// Убеждаемся, что таблица осталась пустой
	var count int
	err = TestDbPool.QueryRow(ctx, "SELECT COUNT(*) FROM sneakers").Scan(&count)
	require.NoError(err)
	require.Equal(0, count)
}

// Тест №3: Ошибка при вставке дубликата по первичному ключу.
func TestCreateSneakers_Failure_DuplicateKey(t *testing.T) {
	// --- Arrange ---
	require := require.New(t)
	ctx := context.Background()
	storage := NewMockPostgresStorageImpl(TestDbPool, zap.NewNop())

	sneakersToCreate := []*model.Sneaker{
		{ID: 10, Article: "ART-010", SneakerName: "Duplicate Test 1", Price: 99.99, Brand: "Puma"},
		{ID: 10, Article: "ART-011", SneakerName: "Duplicate Test 2", Price: 99.99, Brand: "Puma"},
	}
	t.Cleanup(func() {
		_, err := TestDbPool.Exec(ctx, "TRUNCATE TABLE sneakers RESTART IDENTITY CASCADE")
		require.NoError(err)
	})

	// --- Act ---
	err := storage.CreateSneakers(ctx, sneakersToCreate)

	// --- Assert ---
	require.Error(err, "должна быть ошибка из-за дублирования ключа")
	log.Println(err)

	// Проверяем тип и код ошибки
	var pgErr *pgconn.PgError
	require.True(errors.As(err, &pgErr) && pgErr.Code == "23505", "ожидается ошибка unique_violation (23505)")
	
	// Проверяем, что транзакция была откатана
	var count int
	err = TestDbPool.QueryRow(ctx, "SELECT COUNT(*) FROM sneakers").Scan(&count)
	require.NoError(err)
	require.Equal(0, count, "в таблице не должно быть записей после отката транзакции")
}


// Тест №4: Ошибка при отмененном контексте.
func TestCreateSneakers_Failure_ContextCanceled(t *testing.T) {
	// --- Arrange ---
	require := require.New(t)
	storage := NewMockPostgresStorageImpl(TestDbPool, zap.NewNop())

	// Создаем и сразу отменяем контекст
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	sneakersToCreate := []*model.Sneaker{
		{ID: 20, Article: "ART-020", SneakerName: "Context Test", Price: 50.00, Brand: "Reebok"},
	}

	// --- Act ---
	err := storage.CreateSneakers(ctx, sneakersToCreate)

	// --- Assert ---
	require.Error(err, "должна быть ошибка из-за отмененного контекста")
	require.ErrorIs(err, context.Canceled, "ошибка должна оборачивать context.Canceled")
}