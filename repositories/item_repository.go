package repositories

import (
	"basic-go/models"
	"basic-go/msg"
	"errors"

	"gorm.io/gorm"
)

type IItemRepository interface {
	FindAll() (*[]models.Item, error)
	FindById(itemId uint) (*models.Item, error)
	Create(newItem models.Item) (*models.Item, error)
	Update(updateItem models.Item) (*models.Item, error)
	Delete(itemId uint) error
}

type ItemMemoryRepository struct {
	items []models.Item
}

// 引数アイテムはDBモックなので本来はない。
// IFで定義したitems構造体をレスポンスする
// このitems構造体のメソッドが実装の実体であり、リポジトリの利用者であるserviceは
// メソッドの引数と戻り値だけを知ることになる
func NewItemMemoryRepository(items []models.Item) IItemRepository {
	return &ItemMemoryRepository{items: items}
}

func (r *ItemMemoryRepository) FindAll() (*[]models.Item, error) {
	return &r.items, nil
}

func (r *ItemMemoryRepository) FindById(itemId uint) (*models.Item, error) {
	for _, v := range r.items {
		if v.ID == itemId {
			return &v, nil
		}
	}
	return nil, errors.New("Item not found")
}

func (r *ItemMemoryRepository) Create(newItem models.Item) (*models.Item, error) {
	newItem.ID = uint(len(r.items) + 1)
	r.items = append(r.items, newItem)
	return &newItem, nil
}

func (r *ItemMemoryRepository) Update(updateItem models.Item) (*models.Item, error) {
	for i, v := range r.items {
		if v.ID == updateItem.ID {
			r.items[i] = updateItem
			return &r.items[i], nil
		}
	}
	return nil, errors.New(msg.UE)
}

func (r *ItemMemoryRepository) Delete(itemId uint) error {
	for i, v := range r.items {
		if v.ID == itemId {
			r.items = append(r.items[:i], r.items[i+1:]...)
			return nil
		}
	}
	return errors.New("Item not found")
}

type ItemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) IItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) FindAll() (*[]models.Item, error) {
	var items []models.Item
	if err := r.db.Find(&items).Error; err != nil {
		return nil, err
	}
	return &items, nil
}

func (r *ItemRepository) FindById(itemId uint) (*models.Item, error) {

	return nil, nil
}

func (r *ItemRepository) Create(newItem models.Item) (*models.Item, error) {
	result := r.db.Create(&newItem)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newItem, nil
}

func (r *ItemRepository) Update(updateItem models.Item) (*models.Item, error) {
	return nil, nil
}

func (r *ItemRepository) Delete(itemId uint) error {
	return nil
}
