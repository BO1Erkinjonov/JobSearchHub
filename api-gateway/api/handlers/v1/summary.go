package v1

import (
	"api-gateway/api/models"
	pbs "api-gateway/genproto/client-service"
	"api-gateway/internal/pkg/config"
	l "api-gateway/internal/pkg/logger"
	jwt "api-gateway/internal/pkg/tokens"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// CreateSummary ...
// @Summary CreateSummary
// @Description CreateSummary - Api for creating client
// @Tags summary
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param CreateSummary body models.Summary true "create summary model"
// @Success 200 {object} models.SummaryResponse
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/create/summary [post]
func (h *handlerV1) CreateSummary(c *gin.Context) {
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
	var body models.Summary
	err = c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.log.Error("failed to parse body", zap.Error(err))
		return
	}

	resp, err := h.serviceManager.ClientService().CreateSummary(ctx, &pbs.Summary{
		OwnerId:   cast.ToString(claims["sub"]),
		Skills:    body.Skills,
		Bio:       body.Bio,
		Languages: body.Languages,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create summary", zap.Error(err))
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// GetAllSummaryByOwnerId ...
// @Summary GetAllSummaryByOwnerId
// @Description GetAllSummaryByOwnerId - Api for get summary
// @Tags summary
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param page query int true "page"
// @Param limit query int true "limit"
// @Success 200 {object} models.ListSummary
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/get/all/summary/owner [get]
func (h *handlerV1) GetAllSummaryByOwnerId(c *gin.Context) {
	page := cast.ToInt(c.DefaultQuery("page", "1"))
	limit := cast.ToInt(c.DefaultQuery("limit", "10"))

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

	resp, err := h.serviceManager.ClientService().GetAllSummary(ctx, &pbs.GetAllRequestSummary{
		Page:  int32(page),
		Limit: int32(limit),
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
	c.JSON(http.StatusOK, resp)
}

// UpdateSummary ...
// @Summary UpdateSummary
// @Description UpdateSummary - Api for update summary
// @Tags summary
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param UpdateSummary body models.SummaryUpdateRequest true "update summary"
// @Success 200 {object} models.Summary
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/update/summary [put]
func (h *handlerV1) UpdateSummary(c *gin.Context) {
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
	var body models.SummaryUpdateRequest
	err = c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.log.Error("failed to parse body", zap.Error(err))
		return
	}

	resp, err := h.serviceManager.ClientService().UpdateSummary(ctx, &pbs.Summary{
		OwnerId:   cast.ToString(claims["sub"]),
		Id:        body.Id,
		Skills:    body.Skills,
		Bio:       body.Bio,
		Languages: body.Languages,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update summary", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, resp)
}

// DeleteSummary ...
// @Summary DeleteSummary
// @Description DeleteClient - Api for delete summary
// @Tags summary
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param summary_id query int true "summary_id"
// @Success 200 {object} bool
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/delete/summary [delete]
func (h *handlerV1) DeleteSummary(c *gin.Context) {
	id := cast.ToInt(c.Query("summary_id"))

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

	resp, err := h.serviceManager.ClientService().DeleteSummary(ctx, &pbs.GetRequestSummary{
		Id:      int32(id),
		OwnerId: cast.ToString(claims["sub"]),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get summary", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, resp)
}
