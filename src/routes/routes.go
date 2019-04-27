package routes

import (
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/slack-bot-4all/slack-bot/cmd/config"
	"github.com/slack-bot-4all/slack-bot/cmd/model"
	"github.com/slack-bot-4all/slack-bot/cmd/repository"
	"github.com/slack-bot-4all/slack-bot/cmd/resource"
)

// ErrorMap : struct
type ErrorMap struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// GetRoutes : function to map all security and routes permissions
func GetRoutes() *gin.Engine {
	r := gin.Default()

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("jeremyBOT-github"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: "user",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					"user": v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &model.User{
				Username: claims["user"].(string),
			}
		},
		Authenticator: Authenticator,
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// Auth Group
	{
		authGroup := r.Group("/auth")
		authGroup.POST("/login", authMiddleware.LoginHandler)
		authGroup.GET("/refresh_token", authMiddleware.RefreshHandler)
	}

	// v1 Group
	v1 := r.Group("/v1")
	v1.Use(authMiddleware.MiddlewareFunc())

	// Users Group
	// {
	// 	usersGroup := v1.Group("/users")

	// 	usersGroup.GET("/", resource.ListUsers)
	// }

	// Ranchers Group
	{
		ranchersGroup := v1.Group("/ranchers")

		ranchersGroup.GET("/", resource.ListRancher)
		ranchersGroup.POST("/", resource.AddRancher)
	}

	return r
}

// Authenticator ::
func Authenticator(c *gin.Context) (interface{}, error) {
	var userLogin model.User
	if err := c.ShouldBind(&userLogin); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	password := userLogin.Password
	if err := repository.FindUserByUsername(&userLogin); err != nil {
		return nil, jwt.ErrFailedAuthentication
	} else {
		var hash config.Hash
		if err := hash.Compare(userLogin.Password, password); err != nil {
			return nil, jwt.ErrFailedAuthentication
		}
		return &model.User{
			Username: userLogin.Username,
		}, nil
	}

}
