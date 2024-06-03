package handler

import (
	"api-gateway/api/models"
	"api-gateway/entity"
	l "api-gateway/internal/pkg/logger"
	"api-gateway/mock_data/jobs_service"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"net/http"
	"time"
)

func CreateJob(c *gin.Context) {
	var body models.JobsOwner

	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		l.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	resp, err := jobs_service.NewMockClientServiceClient().CreateJob(ctx, &entity.Job{
		Id:          uuid.NewString(),
		Owner_id:    body.OwnerId,
		Title:       body.Title,
		Description: body.Description,
		Response:    0,
	})
	c.JSON(http.StatusOK, resp)
}

func GetJobsByOwnerId(c *gin.Context) {
	page := cast.ToInt(c.DefaultQuery("page", "1"))
	limit := cast.ToInt(c.DefaultQuery("limit", "10"))
	clientId := c.Query("ClientId")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	resp, err := jobs_service.NewMockClientServiceClient().GetAllJobs(ctx, &entity.GetAll{
		Page:  int32(page),
		Limit: int32(limit),
		Field: "owner_id",
		Value: clientId,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		l.Error(err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func DeleteJob(c *gin.Context) {
	id := c.Query("id")
	isActive := c.GetBool("is_active")
	isHardDelete := c.GetBool("is_hard")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	_, err := jobs_service.NewMockClientServiceClient().DeleteJob(ctx, &entity.DelReq{
		Id:            id,
		IsActive:      isActive,
		IsHardDeleted: isHardDelete,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		l.Error(err)
		return
	}
	c.JSON(http.StatusOK, true)
}
