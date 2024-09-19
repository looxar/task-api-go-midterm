package item

import (
	"task-api/internal/constant"
	"task-api/internal/model"

	"gorm.io/gorm"
)

type Service struct {
	Repository Repository
}

func NewService(db *gorm.DB) Service {
	return Service{
		Repository: NewRepository(db),
	}
}
func (service Service) Create(req model.RequestItem) (model.Item, error) {
	item := model.Item{
		Title:    req.Title,
		Amount:    req.Amount,
		Quantity: req.Quantity,
		Status:   constant.ItemPendingStatus,
	}
	if err := service.Repository.Create(&item); err != nil {
		return model.Item{}, err
	}

	return item, nil
}

func (service Service) Find(query model.RequestFindItem) ([]model.Item, error) {
    return service.Repository.Find(query)
}

func (service Service) UpdateStatus(id uint, status constant.ItemStatus) (model.Item, error) {
    // Find item
    item, err := service.Repository.FindByID(id)
    if err != nil {
        return model.Item{}, err
    }
    // Fill data
    item.Status = status
    // Replace
    if err := service.Repository.Replace(item); err != nil {
        return model.Item{}, err
    }
    return item, nil
}

func (s *Service) FindItemByID(id uint) (model.Item, error) {
	return s.Repository.FindByID(id)
}

func (s *Service) UpdateItem(id uint, updatedItem model.Item) (model.Item, error) {
	// Find the existing item by ID
	item, err := s.Repository.FindByID(id)
	if err != nil {
		return model.Item{}, err
	}

	// Update the item's fields with the new data
	item.Title = updatedItem.Title
	item.Amount = updatedItem.Amount
	item.Quantity = updatedItem.Quantity
	item.Status = updatedItem.Status

	// Save the updated item back to the database
	if err := s.Repository.Save(item); err != nil {
		return model.Item{}, err
	}

	return item, nil
}

func (s *Service) DeleteItem(id uint) error {
	return s.Repository.DeleteByID(id)
}