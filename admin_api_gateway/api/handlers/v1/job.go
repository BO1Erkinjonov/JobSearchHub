package v1

import (
	"admin_api_gateway/api/models"
	pbu "admin_api_gateway/genproto/client-service"
	job "admin_api_gateway/genproto/jobs-service"
	l "admin_api_gateway/internal/pkg/logger"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// CreateJob ...
// @Summary CreateJob
// @Description CreateJob - Api for createJob
// @Tags job
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param CreateJob body models.JobReq true "create job"
// @Success 200 {object} models.JobResponse
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/create/job [post]
func (h *handlerV1) CreateJob(c *gin.Context) {
	var body models.JobReq

	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to parse body", zap.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	resp, err := h.serviceManager.JobService().CreateJob(ctx, &job.Job{
		Id:          uuid.NewString(),
		OwnerId:     body.ClientId,
		Title:       body.Title,
		Description: body.Description,
		Responses:   0,
	})
	c.JSON(http.StatusOK, resp)
}

// GetJobsByOwnerId ...
// @Summary GetJobsByOwnerId
// @Description GetJobsByOwnerId - Api for get jobs
// @Tags job
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Param ClientId query string true "ClientId"
// @Success 200 {object} models.ListJobs
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/get/jobs/ownerId [get]
func (h *handlerV1) GetJobsByOwnerId(c *gin.Context) {
	page := cast.ToInt(c.DefaultQuery("page", "1"))
	limit := cast.ToInt(c.DefaultQuery("limit", "10"))
	clientId := c.Query("ClientId")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	resp, err := h.serviceManager.JobService().GetAllJobs(ctx, &job.GetAll{
		Page:  int32(page),
		Limit: int32(limit),
		Field: "owner_id",
		Value: clientId,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get client", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetAllJobs ...
// @Summary GetAllJobs
// @Description GetAllJobs - Api for get jobs
// @Tags job
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Param title query string false "title"
// @Success 200 {object} models.ListJobs
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/get/jobs [get]
func (h *handlerV1) GetAllJobs(c *gin.Context) {
	page := cast.ToInt(c.DefaultQuery("page", "1"))
	limit := cast.ToInt(c.DefaultQuery("limit", "10"))
	title := c.DefaultQuery("title", "")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	reqTitle := ""

	if title != "" {
		reqTitle = "title"
	}
	resp, err := h.serviceManager.JobService().GetAllJobs(ctx, &job.GetAll{
		Page:  int32(page),
		Limit: int32(limit),
		Field: reqTitle,
		Value: title,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get client", l.Error(err))
		return
	}

	var listJobs models.ListJobs
	for _, i := range resp.Jobs {
		respOwner, err := h.serviceManager.ClientService().GetClientById(ctx, &pbu.GetRequest{
			ClientId: i.OwnerId,
			IsActive: false,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to get client", l.Error(err))
			return
		}

		listJobs.Jobs = append(listJobs.Jobs, models.JobsOwner{
			Id:          i.Id,
			OwnerId:     i.OwnerId,
			Title:       i.Title,
			Description: i.Description,
			Responses:   i.Responses,
			CreatedAt:   i.CreatedAt,
			UpdatedAt:   i.UpdatedAt,
			DeletedAt:   i.DeletedAt,
			Owners: models.ResponseClient{
				Id:        respOwner.Id,
				Role:      respOwner.Role,
				FirstName: respOwner.FirstName,
				LastName:  respOwner.LastName,
				Email:     respOwner.Email,
				Password:  respOwner.Password,
				CreatedAt: respOwner.CreatedAt,
				UpdatedAt: respOwner.UpdatedAt,
				DeletedAt: respOwner.DeletedAt,
			},
		})
	}

	c.JSON(http.StatusOK, listJobs)
}

// UpdateJob ...
// @Summary UpdateJob
// @Description UpdateJob - Api for update job
// @Tags job
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param UpdateJob body models.JobUpdateReq true "update job"
// @Success 200 {object} models.JobResponse
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/update/job [put]
func (h *handlerV1) UpdateJob(c *gin.Context) {
	var body models.JobUpdateReq

	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to parse body", zap.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()
	resp, err := h.serviceManager.JobService().UpdateJob(ctx, &job.Job{
		OwnerId:     body.ClientId,
		Id:          body.JobId,
		Title:       body.Title,
		Description: body.Description,
		Responses:   body.Responses,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get client", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, resp)
}

// DeleteJob ...
// @Summary DeleteJob
// @Description DeleteJob - Api for delete job
// @Tags job
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request query models.JobIdClientIdIsActiveReq true "delete job"
// @Success 200 {object} bool
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/delete/job [delete]
func (h *handlerV1) DeleteJob(c *gin.Context) {
	id := c.Query("id")
	isActive := c.GetBool("is_active")
	isHardDelete := c.GetBool("is_hard")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	resp, err := h.serviceManager.JobService().DeleteJob(ctx, &job.DelReq{
		Id:            id,
		IsActive:      isActive,
		IsHardDeleted: isHardDelete,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get client", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, resp)
}
