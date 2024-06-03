package handler

import (
	"admin_api_gateway/api/models"
	"admin_api_gateway/entity"
	l "admin_api_gateway/internal/pkg/logger"
	"admin_api_gateway/mock_data/jobs_service"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
	"time"
)

func CreateRequest(c *gin.Context) {
	var body models.RequestListResp
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		l.Error(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	resp, err := jobs_service.NewMockClientServiceClient().CreateRequests(ctx, &entity.Request{
		JobId:           body.JobId,
		ClientId:        body.ClientId,
		SummaryId:       body.SummaryId,
		StatusResp:      body.StatusResp,
		DescriptionResp: body.DescriptionResp,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		l.Error(err)
		return
	}
	c.JSON(http.StatusCreated, models.RequestResp{
		JobId:     resp.JobId,
		ClientId:  resp.ClientId,
		SummaryId: resp.SummaryId,
	})
}

func GetAllRequest(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	page := cast.ToInt(c.Query("page"))
	limit := cast.ToInt(c.Query("limit"))
	clientId := c.Query("clientId")
	resp, err := jobs_service.NewMockClientServiceClient().GetAllRequest(ctx, &entity.GetAllReq{
		Page:  int32(page),
		Limit: int32(limit),
		Field: "client_id",
		Value: clientId,
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

func DeleteRequest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	clientId := c.Query("ClientId")
	jobId := c.Query("JobId")

	_, err := jobs_service.NewMockClientServiceClient().DeleteRequest(ctx, &entity.GetRequestReq{
		ClientId: clientId,
		JobId:    jobId,
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
