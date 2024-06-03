package handler

import (
	_ "api-gateway/api/docs"
	"api-gateway/entity"
	"api-gateway/mock_data/client_service"
	"fmt"

	//token "api-gateway/internal/pkg/tokens"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func LogIn(c *gin.Context) {
	password := "Mock password"
	email := c.Query("email")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	user, err := client_service.NewMockClientServiceClient().Exists(ctx, email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Println(user.Password, password, "<================================")
	if password != user.Password {
		c.JSON(http.StatusBadRequest, "password error")
		return
	}

	_, err = client_service.NewMockClientServiceClient().UpdateClient(ctx, &entity.Client{
		Id:           user.Id,
		Role:         "admin",
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Password:     user.Password,
		RefreshToken: "refresh token",
	})
	c.JSON(http.StatusOK, "access token")
}
