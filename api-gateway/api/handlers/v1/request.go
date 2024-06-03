package v1

import (
	"api-gateway/api/models"
	pbs "api-gateway/genproto/client-service"
	job "api-gateway/genproto/jobs-service"
	"api-gateway/internal/pkg/config"
	l "api-gateway/internal/pkg/logger"
	jwt "api-gateway/internal/pkg/tokens"
	"context"
	"errors"
	"fmt"
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
// @Param CreateRequest body models.Request true "create request"
// @Success 200 {object} models.RequestResp
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/create/request [post]
func (h *handlerV1) CreateRequest(c *gin.Context) {
	var body models.Request
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.log.Error("Failed to parse request body", zap.Error(err))
		return
	}
	tok := c.GetHeader("Authorization")
	claims, err := jwt.ExtractClaim(tok, []byte(config.Token()))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get client", l.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	respSummary, err := h.serviceManager.ClientService().GetAllSummary(ctx, &pbs.GetAllRequestSummary{
		Page:  0,
		Limit: 0,
		Field: "owner_id",
		Value: cast.ToString(claims["sub"]),
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
		fmt.Println(i)
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
		ClientId:  cast.ToString(claims["sub"]),
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
// @Success 200 {object} models.ListRequest
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/get/all/request [get]
func (h *handlerV1) GetAllRequest(c *gin.Context) {
	tok := c.GetHeader("Authorization")
	claims, err := jwt.ExtractClaim(tok, []byte(config.Token()))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get client", l.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	page := cast.ToInt(c.Query("page"))
	limit := cast.ToInt(c.Query("limit"))

	resp, err := h.serviceManager.JobService().GetAllRequest(ctx, &job.GetAllReq{
		Page:  int32(page),
		Limit: int32(limit),
		Field: "client_id",
		Value: cast.ToString(claims["sub"]),
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
// @Param UpdateRequest body models.RequestResponse true "update request"
// @Success 200 {object} models.Request
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/update/request [put]
func (h *handlerV1) UpdateRequest(c *gin.Context) {

	tok := c.GetHeader("Authorization")
	claims, err := jwt.ExtractClaim(tok, []byte(config.Token()))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get client", l.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	var body models.RequestResponse
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
		Value: cast.ToString(claims["sub"]),
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
		ClientId:        cast.ToString(claims["sub"]),
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
