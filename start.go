package ginger

import (
	"github.com/gin-gonic/gin"
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
}

//默认就是DEBUG。如果设置为FALSE，则ReleaseMode
func (g *Ginger) Debug(isDebug bool) *Ginger {
	if !isDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	return g
}

//加载路由
func (g *Ginger) LoadRouter(router func(router *gin.Engine) gin.IRoutes) *Ginger {
	g.App = gin.Default()
	g.App.Use(Cors())
	router(g.App)
	return g
}

//开启服务
func (g *Ginger) Start(addr string) {
	if g.App == nil {
		log.Fatalf("No Router")
	}
	g.App.Run(addr)
}
