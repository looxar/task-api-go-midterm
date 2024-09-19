package model

import "task-api/internal/constant"

type RequestItem struct {
	Title string `binding:"required"`
	// Add      float64 `binding:"required,gte=10"`
	Amount   uint `binding:"required"`
	Quantity uint `binding:"required"`
}
type RequestFindItem struct {
	// Statuses []constant.ItemStatus `form:"status[]"`
	Statuses constant.ItemStatus `form:"status"`
}

type RequestUpdateItem struct {
	Status constant.ItemStatus
}

type RequestLogin struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}
