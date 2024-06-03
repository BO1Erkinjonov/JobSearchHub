package v1

import (
	_ "admin_api_gateway/api/docs"
	pbu "admin_api_gateway/genproto/client-service"
	l "admin_api_gateway/internal/pkg/logger"
	token "admin_api_gateway/internal/pkg/tokens"
	//token "api-gateway/internal/pkg/tokens"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

// LogIn ...
// @Summary LogIn
// @Description LogIn - Api for registering users
// @Tags Register
// @Accept json
// @Produce json
// @Param password query string true "Password"
// @Param email query string true "Email"
// @Success 200 {object} models.AccessToken
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/login/ [post]
func (h *handlerV1) LogIn(c *gin.Context) {
	password := c.Query("password")
	email := c.Query("email")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	user, err := h.serviceManager.ClientService().Exists(ctx, &pbu.EmailRequest{
		Email: email,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if user.DeletedAt != "0001-01-01 00:00:00 +0000 UTC" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "such user does not exist",
		})
		h.log.Error("the user has already been deleted", l.Error(err))
		return
	}
	if password != user.Password {

		c.JSON(http.StatusBadRequest, "password error")
		return
	}
	h.jwthandler = token.JWTHandler{
		Sub:     user.Id,
		Role:    "admin",
		SignKey: h.cfg.Token.Secret,
		Timout:  h.cfg.Token.AccessTTL,
	}
	access, refresh, err := h.jwthandler.GenerateAuthJWT()
	if err != nil {
		log.Fatalln(err)
	}
	_, err = h.serviceManager.ClientService().UpdateClient(ctx, &pbu.Client{
		Id:           user.Id,
		Role:         "admin",
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Password:     user.Password,
		RefreshToken: refresh,
	})
	c.JSON(http.StatusOK, access)
}
