package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"

	"github.com/openapp-dev/openapp/pkg/apiserver/handler"
	"github.com/openapp-dev/openapp/pkg/generated/clientset/versioned"
	"github.com/openapp-dev/openapp/pkg/utils"
)

func NewGinContextWithClientLister(k8sClient kubernetes.Interface,
	openappClient versioned.Interface,
	openappHelper *utils.OpenAPPHelper) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(utils.OpenAPPHelperKey, openappHelper)
		c.Next()
	}
}

func NewOpenAPPServerRouter(k8sClient kubernetes.Interface,
	openappClient versioned.Interface,
	openappHelper *utils.OpenAPPHelper) *gin.Engine {
	router := gin.New()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"*"}
	corsHandler := cors.New(config)

	router.Use(corsHandler)
	router.Use(NewGinContextWithClientLister(k8sClient, openappClient, openappHelper))

	// middleware
	cfg, err := openappHelper.ConfigMapLister.ConfigMaps(utils.SystemNamespace).Get(utils.SystemConfigMap)
	if err != nil {
		panic(err)
	}
	router.Use(jwtAuth([]byte(cfg.Data["password"])))

	initAPPRouter(router, corsHandler)
	initPublicServiceRouter(router, corsHandler)
	initConfigRouter(router, corsHandler)
	initVersionRouter(router, corsHandler)
	initLoginRouter(router, corsHandler)

	return router
}

func initAPPRouter(router *gin.Engine, corsHandler gin.HandlerFunc) {
	appGroup := router.Group("/api/v1/apps")
	appGroup.GET("/templates", handler.ListAllAppTemplatesHandler)
	appGroup.GET("/templates/:templateName", handler.GetAppTemplateHandler)

	appGroup.GET("/instances", handler.ListAllAppInstancesHandler)
	appGroup.GET("/instances/:instanceName", handler.GetAppInstanceHandler)
	appGroup.POST("/instances/:instanceName", handler.CreateOrUpdateAppInstanceHandler)
	appGroup.DELETE("/instances/:instanceName", handler.DeleteAppInstanceHandler)
	appGroup.GET("/instances/:instanceName/log", handler.AppInstanceLoggingHandler)
	appGroup.Use(corsHandler)
}

func initPublicServiceRouter(router *gin.Engine, corsHandler gin.HandlerFunc) {
	publicServiceGroup := router.Group("/api/v1/publicservices")
	publicServiceGroup.GET("/templates", handler.ListAllPublicServiceTemplatesHandler)
	publicServiceGroup.GET("/templates/:templateName", handler.GetPublicServiceTemplateHandler)

	publicServiceGroup.GET("/instances", handler.ListAllPublicServiceInstancesHandler)
	publicServiceGroup.GET("/instances/:instanceName", handler.GetPublicServiceInstanceHandler)
	publicServiceGroup.POST("/instances/:instanceName", handler.CreateOrUpdatePublicServiceInstanceHandler)
	publicServiceGroup.DELETE("/instances/:instanceName", handler.DeletePublicServiceInstanceHandler)
	publicServiceGroup.GET("/instances/:instanceName/log", handler.PublicServiceInstanceLoggingHandler)
	publicServiceGroup.Use(corsHandler)
}

func initConfigRouter(router *gin.Engine, corsHandler gin.HandlerFunc) {
	configGroup := router.Group("/api/v1/config")
	configGroup.GET("", handler.GetConfigHandler)
	configGroup.POST("", handler.UpdateConfigHandler)
	configGroup.Use(corsHandler)
}

func initVersionRouter(router *gin.Engine, corsHandler gin.HandlerFunc) {
	versionGroup := router.Group("/version")
	versionGroup.GET("", handler.GetOpenAPPVersionHandler)
	versionGroup.Use(corsHandler)
}

func initLoginRouter(router *gin.Engine, corsHandler gin.HandlerFunc) {
	loginGroup := router.Group("/login")
	loginGroup.POST("", handler.LoginHandler)
	loginGroup.Use(corsHandler)
}

func jwtAuth(secret []byte) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Authorization token is required",
			})
			ctx.Abort()
			return
		}

		jwt := utils.NewJWT(secret)
		_, err := jwt.ParseToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": err.Error(),
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
