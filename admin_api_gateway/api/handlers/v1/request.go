package v1

import (
	"admin_api_gateway/api/models"
	pbs "admin_api_gateway/genproto/client-service"
	job "admin_api_gateway/genproto/jobs-service"
	l "admin_api_gateway/internal/pkg/logger"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// CreateRequest ...
// @Summary CreateRequest
// @Description CreateRequest - Api for create request
// @Tags requests
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param CreateRequest body models.RequestListResp true "create request"
// @Success 200 {object} models.RequestResp
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/create/request [post]
func (h *handlerV1) CreateRequest(c *gin.Context) {
	var body models.RequestListResp
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.log.Error("Failed to parse request body", zap.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	respSummary, err := h.serviceManager.ClientService().GetAllSummary(ctx, &pbs.GetAllRequestSummary{
		Page:  0,
		Limit: 0,
		Field: "owner_id",
		Value: body.ClientId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get summary", zap.Error(err))
		return
	}
	t := false
	for _, i := range respSummary.Summary {
		if i.Id == body.SummaryId {
			t = true
			break
		}
	}
	if !t {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "summary already exists",
		})
		h.log.Error("failed to update summary", l.Error(errors.New("summary already exists")))
		return
	}

	resp, err := h.serviceManager.JobService().CreateRequests(ctx, &job.Request{
		JobId:     body.JobId,
		ClientId:  body.ClientId,
		SummaryId: body.SummaryId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create request", l.Error(err))
		return
	}
	c.JSON(http.StatusCreated, models.RequestResp{
		JobId:     resp.JobId,
		ClientId:  resp.ClientId,
		SummaryId: resp.SummaryId,
	})
}

// GetAllRequest ...
// @Summary GetAllRequest
// @Description GetAllRequest - Api for get request
// @Tags requests
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Param clientId query string false "clientId"
// @Success 200 {object} models.ListRequest
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/get/all/request [get]
func (h *handlerV1) GetAllRequest(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	page := cast.ToInt(c.Query("page"))
	limit := cast.ToInt(c.Query("limit"))
	clientId := c.Query("clientId")
	resp, err := h.serviceManager.JobService().GetAllRequest(ctx, &job.GetAllReq{
		Page:  int32(page),
		Limit: int32(limit),
		Field: "client_id",
		Value: clientId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get request", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateRequest ...
// @Summary UpdateRequest
// @Description UpdateRequest - Api for update request
// @Tags requests
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param UpdateRequest body models.RequestListResp true "update request"
// @Success 200 {object} models.Request
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/update/request [put]
func (h *handlerV1) UpdateRequest(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	var body models.RequestListResp
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to parse request body", zap.Error(err))
		return
	}
	resp, err := h.serviceManager.JobService().GetAllRequest(ctx, &job.GetAllReq{
		Page:  0,
		Limit: 0,
		Field: "client_id",
		Value: body.ClientId,
	})
	t := false
	for _, i := range resp.Req {
		if i.JobId == body.JobId {
			t = true
			break
		}
	}

	if !t {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "request already exists",
		})
		h.log.Error("failed to update request", l.Error(errors.New("request already exists")))
		return
	}
	respReq, err := h.serviceManager.JobService().UpdateRequest(ctx, &job.Request{
		StatusResp:      body.StatusResp,
		JobId:           body.JobId,
		ClientId:        body.ClientId,
		DescriptionResp: body.DescriptionResp,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update request", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, respReq)
}

// DeleteRequest ...
// @Summary DeleteRequest
// @Description DeleteRequest - Api for delete request
// @Tags requests
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param ClientId query string true "ClientId"
// @Param JobId query string true "JobId"
// @Success 200 {object} bool
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/delete/request [delete]
func (h *handlerV1) DeleteRequest(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	clientId := c.Query("ClientId")
	jobId := c.Query("JobId")

	resp, err := h.serviceManager.JobService().DeleteRequest(ctx, &job.GetRequest{
		ClientId: clientId,
		JobId:    jobId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete request", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, resp.Status)
}
