package item

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"task-api/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type Controller struct {
	Service Service
}

func NewController(db *gorm.DB) Controller {
	return Controller{
		Service: NewService(db),
	}
}

type ApiError struct {
	Field  string
	Reason string
}

func msgForTag(tag, param string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "gt":
		return fmt.Sprintf("Number must greater than %v", param)
	case "gte":
		return fmt.Sprintf("Number must greater than or equal %v", param)
	}
	return ""
}
func getValidationErrors(err error) []ApiError {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ApiError, len(ve))
		for i, fe := range ve {
			out[i] = ApiError{fe.Field(), msgForTag(fe.Tag(), fe.Param())}
		}
		return out
	}
	return nil
}

func (controller Controller) CreateItem(ctx *gin.Context) {
	// Bind
	var request model.RequestItem
	if err := ctx.Bind(&request); err != nil {
		fmt.Println("Validation error:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": getValidationErrors(err),
		})
		return
	}
	// Create item
	item, err := controller.Service.Create(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}
	// Response
	ctx.JSON(http.StatusCreated, gin.H{
		"data": item,
	})
}

func (controller Controller) FindItems(ctx *gin.Context) {
	// Bind query parameters
	var (
		request model.RequestFindItem
	)
	if err := ctx.BindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}
	// Find
	items, err := controller.Service.Find(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": items,
	})
}

func (controller Controller) UpdateItemStatus(ctx *gin.Context) {
	// Bind
	var (
		request model.RequestUpdateItem
	)
	if err := ctx.Bind(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}
	// Path param
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	// Update status
	item, err := controller.Service.UpdateStatus(uint(id), request.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": item,
	})
}

func (controller *Controller) GetItemByID(c *gin.Context) {
	// Get the ID from the URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID",
		})
		return
	}

	// Convert id to uint since FindByID expects a uint
	itemID := uint(id)

	// Fetch the item by ID using the service
	item, err := controller.Service.FindItemByID(itemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Item not found",
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (controller *Controller) UpdateItem(c *gin.Context) {
	// Get the ID from the URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID",
		})
		return
	}

	// Fetch the existing item from the database
	existingItem, err := controller.Service.FindItemByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Item not found",
		})
		return
	}

	// Bind the incoming JSON to the model
	var item model.Item
	if err := c.BindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	// Keep the existing status if not provided in the incoming request
	if item.Status == "" {
		item.Status = existingItem.Status
	}

		// Merge the existing item's fields with the new data
		existingItem.Title = item.Title
		existingItem.Amount = item.Amount
		existingItem.Quantity = item.Quantity
		existingItem.Status = item.Status // This will either be unchanged or updated based on the incoming request

	// Update the item using the service
	updatedItem, err := controller.Service.UpdateItem(uint(id), item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to update item",
		})
		return
	}

	c.JSON(http.StatusOK, updatedItem)
}

func (controller *Controller) DeleteItem(c *gin.Context) {
	// Get the ID from the URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID",
		})
		return
	}

	// Delete the item using the service
	err = controller.Service.DeleteItem(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to delete item",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item deleted successfully",
	})
}
