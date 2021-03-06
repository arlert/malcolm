package server

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/arlert/malcolm/model"
	"github.com/arlert/malcolm/server/middleware/header"
	"github.com/arlert/malcolm/server/service"
	_ "github.com/arlert/malcolm/utils/loghook"
	"github.com/arlert/malcolm/utils/reqlog"
)

// Load loads the router
func Load(cfg *model.Config) http.Handler {

	logrus.Debugf("\n\nLoad with config:\n %+v\n\n", cfg)

	e := gin.New()
	e.Use(gin.Recovery())

	e.Use(header.NoCache)
	e.Use(header.Secure)
	e.Use(header.Version)
	e.Use(header.Options)
	e.Use(reqlog.ReqLoggerMiddleware(logrus.New(), time.RFC3339, true))

	svc := service.New(cfg)
	//svc := &service.Service{}

	e.GET("ping", svc.GetPing)

	//e.Use(static.Serve("/", utils.Frontend("dist")))

	v1group := e.Group("/v1")
	{
		//-----------------------------------------------------------------
		// pipe
		v1group.POST("/pipe", svc.PostPipe)             // create pipe
		v1group.PATCH("/pipe/:pipeid", svc.PatchPipe)   // update pipe
		v1group.GET("/pipe", svc.GetPipe)               // get pipe
		v1group.GET("/pipe/:pipeid", svc.GetPipe)       // get pipe
		v1group.DELETE("/pipe/:pipeid", svc.DeletePipe) // delete pipe

		//-----------------------------------------------------------------
		// build
		v1group.GET("/build/queue", svc.GetBuildInQueue)   // build queue
		v1group.POST("/pipe/:pipeid/build", svc.PostBuild) // trigger build
		v1group.PATCH("/pipe/:pipeid/build/:buildid", svc.PatchBuild)
		v1group.GET("/pipe/:pipeid/build/:buildid", svc.GetBuild) // get build

		//-----------------------------------------------------------------
		// plugin
		v1group.POST("/plugin", svc.PostPlugin)               // create build
		v1group.PATCH("/plugin/:pluginid", svc.PatchPlugin)   // update build
		v1group.GET("/plugin/:pluginid", svc.GetPlugin)       // get build
		v1group.DELETE("/plugin/:pluginid", svc.DeletePlugin) // delete build

		//-----------------------------------------------------------------
		// log & message
		v1group.GET("/pipe/:pipeid/build/:buildid/log", svc.GetLog)
		v1group.GET("/pipe/:pipeid/build/:buildid/message", svc.GetMessage) // allow build step send message to master

		//-----------------------------------------------------------------
		// hook
		v1group.POST("/hook", svc.GetHook) // hook for git repo
	}

	return e
}
