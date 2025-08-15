package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/kripst/krosovka/inventory_service/internal/model"
)

type SneakerFilters struct {
	Brand    string
	Name     string
	MinPrice float64
	MaxPrice float64
	Size     float32
}

type Pagination struct {
	Limit  int
	Offset int
}

// GetSneakers получает кроссовки с фильтрацией и пагинацией.
func (r *PostgresStorageImpl) GetSneakers(ctx context.Context, filter SneakerFilters, pagination Pagination) ([]model.Sneaker, error) {
	// Начинаем строить запрос, выбирая все поля
	queryBuilder := r.sq.Select("*").From(SneakersTable)

	// Последовательно применяем фильтры с помощью вспомогательных методов
	queryBuilder = r.applyBrandFilter(queryBuilder, filter.Brand)
	queryBuilder = r.applyNameFilter(queryBuilder, filter.Name)
	queryBuilder = r.applyPriceFilter(queryBuilder, filter.MinPrice, filter.MaxPrice)
	queryBuilder = r.applySizeFilter(queryBuilder, filter.Size)

	// Добавляем сортировку для стабильной пагинации
	queryBuilder = queryBuilder.OrderBy(SneakersCreatedAt + "DESC")

	// Применяем пагинацию
	if pagination.Limit > 0 {
		queryBuilder = queryBuilder.Limit(uint64(pagination.Limit))
	}
	if pagination.Offset > 0 {
		queryBuilder = queryBuilder.Offset(uint64(pagination.Offset))
	}

	// Генерируем финальный SQL и слайс аргументов
	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("ошибка при построении SQL-запроса: %w", err)
	}

	// Выполняем запрос к базе данных
	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса к БД: %w", err)
	}
	defer rows.Close()

	// Сканируем все полученные строки в срез структур Sneaker
	// pgx.RowToStructByName автоматически сопоставит snake_case поля из БД с CamelCase полями Go
	sneakers, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Sneaker])
	if err != nil {
		return nil, fmt.Errorf("ошибка при сканировании результатов: %w", err)
	}

	return sneakers, nil
}

// --- Вспомогательные методы для фильтрации ---

func (r *PostgresStorageImpl) applyBrandFilter(builder squirrel.SelectBuilder, brand string) squirrel.SelectBuilder {
	if brand != "" {
		return builder.Where(squirrel.Eq{SneakersBrand: brand})
	}
	return builder
}

func (r *PostgresStorageImpl) applyNameFilter(builder squirrel.SelectBuilder, name string) squirrel.SelectBuilder {
	if name != "" {
		// Используем ILIKE для регистронезависимого поиска (стандарт для PostgreSQL)
		return builder.Where(SneakersName+ " ILIKE ?", "%"+name+"%")
	}
	return builder
}

func (r *PostgresStorageImpl) applyPriceFilter(builder squirrel.SelectBuilder, minPrice, maxPrice float64) squirrel.SelectBuilder {
	if minPrice > 0 {
		builder = builder.Where(squirrel.GtOrEq{SneakersPrice: minPrice})
	}
	if maxPrice > 0 {
		builder = builder.Where(squirrel.LtOrEq{SneakersPrice: maxPrice})
	}
	return builder
}

func (r *PostgresStorageImpl) applySizeFilter(builder squirrel.SelectBuilder, size float32) squirrel.SelectBuilder {
	if size > 0 {
		return builder.Where(squirrel.Eq{SneakersSize: size})
	}
	return builder
}