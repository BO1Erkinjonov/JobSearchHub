package v1

import (
	"admin_api_gateway/api/models"
	pbu "admin_api_gateway/genproto/client-service"
	l "admin_api_gateway/internal/pkg/logger"
	token "admin_api_gateway/internal/pkg/tokens"
	"context"
	"log"
	"net/http"
	"time"
)

// CreateClient ...
// @Summary CreateClient
// @Description CreateClient - Api for creating client
// @Tags client
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param CreateClient body models.Client true "create client"
// @Success 200 {object} models.ResponseClient
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/create/client [post]
func (h *handlerV1) CreateClient(c *gin.Context) {
	var body models.Client
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.log.Error("failed to bind body", zap.Any("body", body))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	ClientId := uuid.NewString()
	h.jwthandler = token.JWTHandler{
		Sub:     ClientId,
		Role:    "client",
		SignKey: h.cfg.Token.Secret,
		Timout:  h.cfg.Token.AccessTTL,
	}

	_, refresh, err := h.jwthandler.GenerateAuthJWT()
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := h.serviceManager.ClientService().CreateClient(ctx, &pbu.Client{
		Id:           ClientId,
		Role:         body.Role,
		FirstName:    body.FirstName,
		LastName:     body.LastName,
		Email:        body.Email,
		Password:     body.Password,
		RefreshToken: refresh,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		h.log.Error("failed to create client", zap.Any("client", body))
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// GetClientById ...
// @Summary GetClientById
// @Description GetClientById - Api for getting client
// @Tags client
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request query models.IdIsActive true "request"
// @Success 200 {object} models.ResponseClient
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/get/client [get]
func (h *handlerV1) GetClientById(c *gin.Context) {
	id := c.Query("id")
	IsActive := cast.ToBool(c.Query("is_active"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()
	resp, err := h.serviceManager.ClientService().GetClientById(ctx, &pbu.GetRequest{
		ClientId: id,
		IsActive: IsActive,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get client", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetClientList ...
// @Summary GetClientList
// @Description GetClientList - Api for get clients
// @Tags client
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Param Field query string false "Field"
// @Param Value query string false "Value"
// @Success 200 {object} models.ListClients
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/get/all/client [get]
func (h *handlerV1) GetClientList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	page := cast.ToInt(c.Query("page"))
	limit := cast.ToInt(c.Query("limit"))

	Field := c.Query("Field")
	Value := c.Query("Value")

	resp, err := h.serviceManager.ClientService().GetAllClients(ctx, &pbu.GetAllRequest{
		Page:  int32(page),
		Limit: int32(limit),
		Field: Field,
		Value: Value,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get clients", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateClient ...
// @Summary UpdateClient
// @Description UpdateClient - Api for updating client
// @Tags client
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param UpdateClient body models.ReqClient true "updateModel"
// @Success 200 {object} models.ResponseClient
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/update/client [put]
func (h *handlerV1) UpdateClient(c *gin.Context) {

	var body models.ReqClient

	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update client", l.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	h.jwthandler = token.JWTHandler{
		Sub:     body.ClientId,
		Role:    "client",
		SignKey: h.cfg.Token.Secret,
		Timout:  h.cfg.Token.AccessTTL,
	}

	_, refresh, err := h.jwthandler.GenerateAuthJWT()
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := h.serviceManager.ClientService().UpdateClient(ctx, &pbu.Client{
		Id:           body.ClientId,
		Role:         body.Role,
		FirstName:    body.FirstName,
		LastName:     body.LastName,
		Email:        body.Email,
		Password:     body.Password,
		RefreshToken: refresh,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update client", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, resp)
}

// DeleteClient ...
// @Summary DeleteClient
// @Description DeleteClient - Api for deleteClient client
// @Tags client
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Param request query models.IdIsActiveHard true "delete client"
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/delete/client [delete]
func (h *handlerV1) DeleteClient(c *gin.Context) {
	clientId := c.Query("id")
	isActive := cast.ToBool(c.Query("is_active"))
	isHardDelete := cast.ToBool(c.Query("is_hard"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()
	resp, err := h.serviceManager.ClientService().DeleteClient(ctx, &pbu.DeleteReq{
		ClientId:      clientId,
		IsActive:      isActive,
		IsHardDeleted: isHardDelete,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get client", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, resp.Status)
}
