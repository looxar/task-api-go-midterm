package model

import "task-api/internal/constant"

type Item struct {
	ID       uint                `json:"id" gorm:"primaryKey"`
	Title    string              `json:"title" binding:"required"`
	Amount   uint                `json:"amount"`
	Quantity uint                `json:"quantity"`
	Status   constant.ItemStatus `json:"status"`
}
