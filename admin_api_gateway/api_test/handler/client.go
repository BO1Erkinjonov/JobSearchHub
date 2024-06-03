package handler

import (
	"admin_api_gateway/api/models"
	"admin_api_gateway/entity"
	l "admin_api_gateway/internal/pkg/logger"
	"admin_api_gateway/mock_data/client_service"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"net/http"
	"time"
)

func CreateClient(c *gin.Context) {
	var body models.Client
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		l.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	ClientId := uuid.NewString()

	resp, err := client_service.NewMockClientServiceClient().CreateClient(ctx, &entity.Client{
		Id:           ClientId,
		Role:         body.Role,
		FirstName:    body.FirstName,
		LastName:     body.LastName,
		Email:        body.Email,
		Password:     body.Password,
		RefreshToken: "refresh token",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		l.Error(err)
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func GetClientById(c *gin.Context) {
	id := c.Query("id")
	IsActive := cast.ToBool(c.Query("is_active"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()
	resp, err := client_service.NewMockClientServiceClient().GetClientById(ctx, &entity.GetRequest{
		ClientId: id,
		IsActive: IsActive,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		l.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func GetClientList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	page := cast.ToInt(c.Query("page"))
	limit := cast.ToInt(c.Query("limit"))

	Field := c.Query("Field")
	Value := c.Query("Value")

	resp, err := client_service.NewMockClientServiceClient().GetAllClients(ctx, &entity.GetAllRequest{
		Page:  int32(page),
		Limit: int32(limit),
		Field: Field,
		Value: Value,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		l.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func DeleteClient(c *gin.Context) {
	clientId := c.Query("id")
	isActive := cast.ToBool(c.Query("is_active"))
	isHardDelete := cast.ToBool(c.Query("is_hard"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()
	_, err := client_service.NewMockClientServiceClient().DeleteClient(ctx, &entity.DeleteReq{
		ClientId:      clientId,
		IsActive:      isActive,
		IsHardDeleted: isHardDelete,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		l.Error(err)
		return
	}
	c.JSON(http.StatusOK, true)
}
