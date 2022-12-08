package routes

import (
	POST "user-onboarding/controllers/POST"

	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/module/apmgin"
)

func v1Routes(route *gin.RouterGroup) {

	router := gin.New()
	router.Use(apmgin.Middleware(router))

	v1Routes := route.Group("/v1")
	{
		v1Routes.POST("/login", POST.UserDetails)
		v1Routes.POST("/createUser", POST.UserDetails)
		v1Routes.POST("/fetchUserName", POST.FetchUser)
	}
}
