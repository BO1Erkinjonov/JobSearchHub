package v1

import (
	"api-gateway/api/models"
	pbu "api-gateway/genproto/client-service"
	"api-gateway/internal/pkg/config"
	l "api-gateway/internal/pkg/logger"
	jwt "api-gateway/internal/pkg/tokens"
	token "api-gateway/internal/pkg/tokens"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"log"
	"net/http"
	"time"
)

// GetClientById ...
// @Summary GetClientById
// @Description GetClientById - Api for getting client
// @Tags client
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} models.ResponseClient
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/get/client [get]
func (h *handlerV1) GetClientById(c *gin.Context) {
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
	resp, err := h.serviceManager.ClientService().GetClientById(ctx, &pbu.GetRequest{
		ClientId: cast.ToString(claims["sub"]),
		IsActive: false,
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

	tok := c.GetHeader("Authorization")
	claims, err := jwt.ExtractClaim(tok, []byte(config.Token()))

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
		Sub:     cast.ToString(claims["sub"]),
		Role:    "client",
		SignKey: h.cfg.Token.Secret,
		Timout:  h.cfg.Token.AccessTTL,
	}

	_, refresh, err := h.jwthandler.GenerateAuthJWT()
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := h.serviceManager.ClientService().UpdateClient(ctx, &pbu.Client{
		Id:           cast.ToString(claims["sub"]),
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
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/delete/client [delete]
func (h *handlerV1) DeleteClient(c *gin.Context) {
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
	resp, err := h.serviceManager.ClientService().DeleteClient(ctx, &pbu.DeleteReq{
		ClientId:      cast.ToString(claims["sub"]),
		IsActive:      false,
		IsHardDeleted: false,
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
