package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kripst/krosovka/inventory_service/config"
	"github.com/kripst/krosovka/inventory_service/internal/model"
	"go.uber.org/zap"
)

type PostgresStorageImpl struct {
	pool  *pgxpool.Pool //postgres
	log *zap.Logger
	ctx context.Context
	sq  squirrel.StatementBuilderType
}

func NewPostgresStorageImpl(storageConfig *config.StorageConfig, log *zap.Logger, ctx context.Context) (*PostgresStorageImpl, error) {
	poolConfig, err := pgxpool.ParseConfig(storageConfig.DSN())
    if err != nil {
        return nil, err
    }
    
    poolConfig.MaxConns = int32(storageConfig.PoolMax)
    // Здесь можно установить другие параметры конфигурации, например:
    // poolConfig.MinConns = int32(storageConfig.PoolMin)
    // poolConfig.MaxConnLifetime = storageConfig.MaxConnLifetime
    // poolConfig.MaxConnIdleTime = storageConfig.MaxConnIdleTime
    
    pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
    if err != nil {
        return nil, err
    }
	
	return &PostgresStorageImpl{
		pool: pool,
		log: log,
		sq: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}, nil
}

func (s *PostgresStorageImpl) Close() error {
	//TODO graceful shd wait for all conns, actions. select
	timeLimit := 60 * time.Second
	select {
	case <- s.ctx.Done():
		s.log.Info("ctx done, close pool")
	case <- time.After(timeLimit):
		s.log.Info("time limit ", zap.Any("time limit", timeLimit))
		return fmt.Errorf("time limit done")
	}
	s.pool.Close()
	return nil
}

func (s *PostgresStorageImpl) CreateSneakers(ctx context.Context, sneakers []*model.Sneaker) error {
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

func (s *PostgresStorageImpl) UpdateSneakers(ctx context.Context, sneakers []*model.Sneaker) error {
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

func (s *PostgresStorageImpl) DeleteSneakers(ctx context.Context, itemIDs []int32) error {
	// Проверяем, не отменен ли контекст перед началом
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context canceled before starting: %w", err)
	}

    query := fmt.Sprintf(`
        UPDATE %s 
        SET %s = CURRENT_TIMESTAMP 
        WHERE %s = ANY($1) AND %s IS NULL`,
        SneakersTable,
        SneakersDeletedAt,  // Поле для мягкого удаления
        SneakersID,
        SneakersDeletedAt,  // Проверка что запись еще не удалена
    )

    // Выполняем запрос с массивом ID
    result, err := s.pool.Exec(ctx, query, itemIDs)
    if err != nil {
        return fmt.Errorf("failed to soft delete items: %w", err)
    }

    // Проверяем что действительно обновили записи
    if result.RowsAffected() == 0 {
        return fmt.Errorf("items not found or already deleted")
    }

    return nil
}