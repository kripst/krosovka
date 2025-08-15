package postgres_test

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/kripst/krosovka/inventory_service/internal/model"
	"go.uber.org/zap"
)

func (s *MockPostgresStorageImpl) Close() error {
	s.pool.Close()
	return nil
}

func (s *MockPostgresStorageImpl) CreateSneakers(ctx context.Context, sneakers []*model.Sneaker) error {
    if len(sneakers) == 0 {
        return nil
    }

	 if err := ctx.Err(); err != nil {
        return fmt.Errorf("context canceled before starting transaction: %w", err)
    }

    tx, err := s.pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    // SQL запрос с включением ID
    query := fmt.Sprintf(`
        INSERT INTO %s (
            %s, %s, %s, %s, %s, %s, %s, %s
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8
        )`,
        SneakersTable,
        SneakersID,
        SneakersArticle,
        SneakersName,
        SneakersDescription,
        SneakersPrice,
        SneakersSize,
        SneakersBrand,
        SneakersProductionAddress,
    )

    batch := &pgx.Batch{}
    
    for _, sneaker := range sneakers {
		if err := ctx.Err(); err != nil {
            return fmt.Errorf("context canceled during batch preparation: %w", err)
        }

        if sneaker.Price <= float64(0) {
            s.log.Error("Price belong or eq zero", zap.Float64("Price", sneaker.Price), zap.String("ID", SneakersID))
            return fmt.Errorf("Price belong or eq zero")
        }
        
        batch.Queue(query,
            sneaker.ID,
            sneaker.Article,
            sneaker.SneakerName,
            sneaker.SneakerDescription,
            sneaker.Price,
            sneaker.Size,
            sneaker.Brand,
            sneaker.ProductionAddress,
        )
    }

	if err := ctx.Err(); err != nil {
        return fmt.Errorf("context canceled before sending batch: %w", err)
    }


    br := tx.SendBatch(ctx, batch)
    if err := br.Close(); err != nil {
        return fmt.Errorf("batch insert failed: %w", err)
    }

	if err := ctx.Err(); err != nil {
        return fmt.Errorf("context canceled before commit: %w", err)
    }


    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("transaction commit failed: %w", err)
    }

    return nil
}

func (s *MockPostgresStorageImpl) UpdateSneakers(sneakers []*model.Sneaker, ctx context.Context) error {
	if len(sneakers) == 0 {
		return nil
	}

	// Проверяем, не отменен ли контекст перед началом
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context canceled before starting: %w", err)
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// SQL запрос для UPDATE с использованием констант
	query := fmt.Sprintf(`
		UPDATE %s SET 
			%s = $1,
			%s = $2,
			%s = $3,
			%s = $4,
			%s = $5,
			%s = $6,
			%s = $7,
		WHERE %s = $8`,
		SneakersTable,
		SneakersArticle,
		SneakersName,
		SneakersDescription,
		SneakersPrice,
		SneakersSize,
		SneakersBrand,
		SneakersProductionAddress,
		SneakersID,
	)

	batch := &pgx.Batch{}

	for _, sneaker := range sneakers {
		// Проверяем контекст перед добавлением в батч
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("context canceled during batch preparation: %w", err)
		}

		batch.Queue(query,
			sneaker.Article,
			sneaker.SneakerName,
			sneaker.SneakerDescription,
			sneaker.Price,
			sneaker.Size,
			sneaker.Brand,
			sneaker.ProductionAddress,
			sneaker.ID,
		)
	}

	// Отправляем batch
	br := tx.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return fmt.Errorf("batch update failed: %w", err)
	}

	// Коммитим транзакцию
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}

func (s *MockPostgresStorageImpl) DeleteSneaker(ctx context.Context, itemID string) error {
    query := fmt.Sprintf(`
        UPDATE %s 
        SET %s = CURRENT_TIMESTAMP 
        WHERE %s = $1 AND %s IS NULL`,
        SneakersTable,
        SneakersDeletedAt,  // Поле для мягкого удаления
        SneakersID,
        SneakersDeletedAt,  // Проверка что запись еще не удалена
    )

    // Выполняем запрос
    result, err := s.pool.Exec(ctx, query, itemID)
    if err != nil {
        return fmt.Errorf("failed to soft delete item: %w", err)
    }

    // Проверяем что действительно обновили запись
    if result.RowsAffected() == 0 {
        return fmt.Errorf("item not found or already deleted")
    }

    return nil
}