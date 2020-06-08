package ginger

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"log"
	"sync"
)

var IGinger *Ginger
var once sync.Once

//单例模式获取实例
func GetIns() *Ginger {
	once.Do(func() {
		IGinger = new(Ginger)
	})
	return IGinger
}

type Ginger struct {
	App *gin.Engine
	//下面的变量不希望出现在用户的可选参数中，误导用户，所以用小写
	openSwagger bool
	hasRouter   bool
	routers     func(router *gin.Engine) gin.IRoutes
}

//默认就是DEBUG。如果设置为FALSE，则ReleaseMode
func (g *Ginger) Debug(isDebug bool) *Ginger {
	if !isDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	return g
}

//开启swagger
func (g *Ginger) Swagger(openSwagger bool) *Ginger {
	if openSwagger {
		g.openSwagger = true
	}
	return g
}

//加载路由
func (g *Ginger) LoadRouter(router func(router *gin.Engine) gin.IRoutes) *Ginger {
	g.hasRouter = true
	g.routers = router
	return g
}

//开启服务
func (g *Ginger) Start(addr string) {
	if g.hasRouter {
		g.App = gin.Default()
		g.App.Use(Cors())
	}
	if g.App == nil || !g.hasRouter {
		log.Fatalf("No Router")
	}
	cRoute := g.routers(g.App)
	if g.openSwagger {
		cRoute.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	g.App.Run(addr)
}
