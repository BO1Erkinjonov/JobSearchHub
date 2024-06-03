package v1

import (
	_ "api-gateway/api/docs"
	"api-gateway/api/models"
	pbu "api-gateway/genproto/client-service"
	l "api-gateway/internal/pkg/logger"
	token "api-gateway/internal/pkg/tokens"

	//token "api-gateway/internal/pkg/tokens"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

// Register ...
// @Summary Register
// @Description Register - Api for registering users
// @Tags Register
// @Accept json
// @Produce json
// @Param Register body models.ReqClient true "createRegisterModel"
// @Success 200 {object} bool
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/register/ [post]err
func (h *handlerV1) Register(c *gin.Context) {
	var (
		body models.RespClient
	)

	err := c.ShouldBindJSON(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Failed to bind json", l.Error(err))
		return
	}
	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	exists_email, err := h.serviceManager.ClientService().CheckUniques(ctx, &pbu.CheckUniquesRequest{
		Field: "email",
		Value: body.Email,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to check email uniques")
		return
	}

	if exists_email.IsExist {
		c.JSON(http.StatusConflict, gin.H{
			"error": "This email already in use, please use another email address",
		})
		h.log.Error("failed to check email uniques", l.Error(err))
		return
	}

	body.Id = uuid.New().String()

	rdb := redis.NewClient(&redis.Options{
		Addr: "redisdb:6379",
	})
	byteDate, err := json.Marshal(&body)
	if err != nil {
		log.Fatalln(err)
	}

	err = rdb.Set(context.Background(), "email_"+body.Email, byteDate, 0).Err()
	if err != nil {
		log.Fatalln(err)
	}
	code := rand.Int() % 1000000
	err = rdb.SetEx(context.Background(), body.Email, code, time.Minute*1).Err()
	if err != nil {
		log.Fatalln(err)
	}
	codeResp := strconv.Itoa(code)

	auth := smtp.PlainAuth("", "boburerkinzonov@gmail.com", "llqmgbilccvhltfd", "smtp.gmail.com")
	err = smtp.SendMail("smtp.gmail.com:587", auth, "boburerkinzonov@gmail.com", []string{body.Email}, []byte(codeResp))
	if err != nil {
		log.Fatalln(err)
	}
	c.JSON(http.StatusOK, true)
}

// Verification ...
// @Summary Verification
// @Description Verification - Api for registering users
// @Tags Register
// @Accept json
// @Produce json
// @Param code query int true "Code"
// @Param email query string true "Email"
// @Success 200 {object} models.AccessToken
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/Verification/ [post]
func (h *handlerV1) Verification(c *gin.Context) {
	codeRegis := c.Query("code")
	emailRegis := c.Query("email")

	rdb := redis.NewClient(&redis.Options{
		Addr: "redisdb:6379",
	})

	respCode, err := rdb.Get(context.Background(), emailRegis).Result()
	if err != nil {
		log.Fatalln(err)
	}
	var code int
	if err := json.Unmarshal([]byte(respCode), &code); err != nil {
		log.Fatalln(err)
	}

	code1, err := strconv.Atoi(codeRegis)
	if err != nil {
		log.Fatalln(err)
	}

	if code != code1 {
		c.JSON(http.StatusBadRequest, false)
	} else {

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
		defer cancel()
		err = rdb.Del(context.Background(), emailRegis).Err()
		if err != nil {
			log.Fatalln(err)
		}
		var regis models.RespClient
		respUser, err := rdb.Get(context.Background(), "email_"+emailRegis).Result()
		if err != nil {
			log.Fatalln(err)
		}
		err = rdb.Del(context.Background(), "email_"+emailRegis).Err()
		if err != nil {
			log.Fatalln(err)
		}
		if err := json.Unmarshal([]byte(respUser), &regis); err != nil {
			log.Fatalln(err)
		}

		h.jwthandler = token.JWTHandler{
			Sub:     regis.Id,
			Role:    "client",
			SignKey: h.cfg.Token.Secret,
			Timout:  h.cfg.Token.AccessTTL,
		}
		//pr := kafka.NewProducerInit(h.cfg, h.log)
		//err = pr.ProduceUser(ctx, "user.create", regis)
		//if err != nil {
		//	fmt.Println("qwe")
		//	return
		//}
		access, refresh, err := h.jwthandler.GenerateAuthJWT()
		if err != nil {
			log.Fatalln(err)
		}
		_, err = h.serviceManager.ClientService().CreateClient(ctx, &pbu.Client{
			Id:           regis.Id,
			Role:         "client",
			FirstName:    regis.FirstName,
			LastName:     regis.LastName,
			Email:        regis.Email,
			Password:     regis.Password,
			RefreshToken: refresh,
		})
		if err != nil {
			log.Fatalln(err)
		}

		c.JSON(http.StatusOK, access)
	}
}

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
		Role:    "client",
		SignKey: h.cfg.Token.Secret,
		Timout:  h.cfg.Token.AccessTTL,
	}
	access, refresh, err := h.jwthandler.GenerateAuthJWT()
	if err != nil {
		log.Fatalln(err)
	}
	_, err = h.serviceManager.ClientService().UpdateClient(ctx, &pbu.Client{
		Id:           user.Id,
		Role:         user.Role,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Password:     user.Password,
		RefreshToken: refresh,
	})
	c.JSON(http.StatusOK, access)
}
