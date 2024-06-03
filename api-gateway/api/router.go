package api

import (
	// "github.com/casbin/casbin/v2"

	"api-gateway/api/middleware/casbin"
	"time"

	v1 "api-gateway/api/handlers/v1"
	"api-gateway/api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	grpcClients "api-gateway/internal/infastructure/grpc_service_client"
	"api-gateway/internal/pkg/config"
	"api-gateway/internal/usecase/event"
)

type RouteOption struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	Service        grpcClients.ServiceClient
	BrokerProducer event.BrokerProducer
}

// @title Bobur Erkinjonov
// @version 1.7
// @host localhost:1212

// NewRoute
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func NewRoute(option RouteOption) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	HandlerV1 := v1.New(&v1.HandlerV1Config{
		Config:         option.Config,
		Logger:         option.Logger,
		ContextTimeout: option.ContextTimeout,
		Service:        option.Service,
		BrokerProducer: option.BrokerProducer,
	})

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowBrowserExtensions = true
	corsConfig.AllowMethods = []string{"*"}
	router.Use(cors.New(corsConfig))

	router.Use(middleware.Tracing())

	router.Use(casbin.NewAuthorizer())

	api := router.Group("/v1")

	// Registration
	api.POST("/register", HandlerV1.Register)
	api.POST("/Verification", HandlerV1.Verification)
	api.POST("/login", HandlerV1.LogIn)

	// Client
	api.GET("/get/client", HandlerV1.GetClientById)
	api.PUT("/update/client", HandlerV1.UpdateClient)
	api.DELETE("/delete/client", HandlerV1.DeleteClient)

	// Job
	api.POST("/create/job", HandlerV1.CreateJob)
	api.GET("/get/jobs/ownerId", HandlerV1.GetJobsByOwnerId)
	api.GET("/get/jobs", HandlerV1.GetAllJobs)
	api.PUT("/update/job", HandlerV1.UpdateJob)
	api.DELETE("/delete/job", HandlerV1.DeleteJob)

	// Summary
	api.POST("/create/summary", HandlerV1.CreateSummary)
	api.GET("/get/all/summary/owner", HandlerV1.GetAllSummaryByOwnerId)
	api.DELETE("/delete/summary", HandlerV1.DeleteSummary)
	api.PUT("/update/summary", HandlerV1.UpdateSummary)

	// Request
	api.POST("/create/request", HandlerV1.CreateRequest)
	api.GET("/get/all/request", HandlerV1.GetAllRequest)
	api.PUT("/update/request", HandlerV1.UpdateRequest)

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
