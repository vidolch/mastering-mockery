package products

import "fmt"

type Database interface {
	GetProductById(id string) (Product, error)
}

var dbProducts = map[string]Product{
	"id1": {
		ID:   "id1",
		Name: "Macbook",
	},
	"id2": {
		ID:   "id2",
		Name: "Playstation",
	},
	"id3": {
		ID:   "id3",
		Name: "Office chair",
	},
}

type DatabaseImpl struct{}

func NewDatabase() *DatabaseImpl {
	return &DatabaseImpl{}
}

func (d DatabaseImpl) GetProductById(id string) (Product, error) {
	fmt.Println("getting " + id + "sup")
	return dbProducts[id], nil
}
