package handler

import (
	"admin_api_gateway/api/models"
	"admin_api_gateway/entity"
	l "admin_api_gateway/internal/pkg/logger"
	"admin_api_gateway/mock_data/client_service"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
	"time"
)

func CreateSummary(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()
	var body models.SummaryResponse
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		l.Error(err)
		return
	}

	resp, err := client_service.NewMockClientServiceClient().CreateSummary(ctx, &entity.Summary{
		OwnerId:   body.OwnerId,
		Skills:    body.Skills,
		Bio:       body.Bio,
		Languages: body.Languages,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		l.Error(err)
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func GetAllSummaryByOwnerId(c *gin.Context) {
	page := cast.ToInt(c.DefaultQuery("page", "1"))
	limit := cast.ToInt(c.DefaultQuery("limit", "10"))
	field := c.Query("Field")
	value := c.Query("Value")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	resp, err := client_service.NewMockClientServiceClient().GetAllSummary(ctx, &entity.GetAllRequestSummary{
		Page:  int32(page),
		Limit: int32(limit),
		Field: field,
		Value: value,
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

func DeleteSummary(c *gin.Context) {
	id := cast.ToInt(c.Query("summary_id"))
	ownerId := c.Query("owner_id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	_, err := client_service.NewMockClientServiceClient().DeleteSummary(ctx, &entity.GetRequestSummary{
		Id:      int32(id),
		OwnerId: ownerId,
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
