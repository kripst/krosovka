package model

import (
	"time"

	pb "github.com/kripst/krosovka/inventory_service/proto"
	"github.com/pkg/errors"
)

type Sneaker struct {
	ID                 int32       `json:"id" db:"id"`
	Article            string    `json:"article" db:"article"`
	SneakerName        string    `json:"sneaker_name" db:"sneaker_name"`
	SneakerDescription string    `json:"sneaker_description,omitempty" db:"sneaker_description"`
	Price              float64   `json:"price" db:"price"`
	Size               float64   `json:"size" db:"size"`
	Brand              string    `json:"brand" db:"brand"`
	ProductionAddress  string    `json:"production_address,omitempty" db:"production_address"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt          time.Time `json:"deleted_at" db:"deleted_at"`
}

func (s *Sneaker) FromGrpc(in *pb.Sneaker) error {
	if s == nil {
		return errors.New("nil struct")
	}
    if in == nil {
        return errors.New("nil request")
    }

    s.ID = in.GetSneakerId()
    s.Article = in.GetArticle()
    s.SneakerName = in.GetSneakerName()
    s.SneakerDescription = in.GetSneakerDescription()
    s.Price = in.GetPrice()
    s.Size = float64(in.GetSize())
    s.Brand = in.GetBrand()
    s.ProductionAddress = in.GetProductionAddress()

    return nil
}