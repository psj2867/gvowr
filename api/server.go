package api

import (
	"gvowr/util"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type videosType = util.SyncMap[string, *videoRoom]
type GvowrServer struct {
	*echo.Echo
	rooms       *videosType
	recommender Recommender
}

func Server() *GvowrServer {
	rooms := util.NewSyncMap[string, *videoRoom]()
	recommender := &MinConnectRecommender{}
	server := &GvowrServer{Echo: echo.New(), rooms: rooms, recommender: recommender}
	recommender.Server = server
	recommender.maxPriority = 3
	setMiddleware(server.Echo)

	server.GET("/health", func(ctx echo.Context) error {
		return ctx.String(200, "ok")
	})
	server.GET("/static/*", StaticEchoServer("/static/"))
	server.setVideoApi()
	server.setNodeApi()
	return server
}
func setMiddleware(server *echo.Echo) {
	server.Use(middleware.CORS())
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.Contains(strings.ToLower(c.Request().Header.Get("Content-Type")), "text/plain") {
				c.Request().Header.Set("Content-Type", "application/json")
			}
			return next(c)
		}
	})
}

func Run() error {
	return http.ListenAndServe(":5050", Server())
}

func makeUuid() string {
	u, _ := uuid.NewUUID()
	return u.String()
}

func jsonFunc(h func(map[string]any, map[string]any) error) func(echo.Context) error {
	return func(c echo.Context) error {
		reqMap := make(map[string]any)
		resMap := make(map[string]any)
		err := c.Bind(&reqMap)
		if err != nil {
			c.Logger().Error(err)
			return err
		}
		err = h(reqMap, resMap)
		if err != nil {
			c.Logger().Error(err)
			return err
		}
		c.JSON(200, resMap)
		return nil
	}
}
