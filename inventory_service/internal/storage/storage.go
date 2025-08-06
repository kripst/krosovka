package storage


type Storage interface {
	//TODO any replace on model.item
	CreateItem(item any) error
	UpdateItem(item any) error
	DeleteItem(itemID string) error
	GetItems(itemIDs []string) ([]any, error)
}

type StorageImpl struct {
	pgx
}