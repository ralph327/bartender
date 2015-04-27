package bartender

import (
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/suapapa/go_sass"
	"log"
	"time"
)

type Bartender interface {
	Init(string) 
	Start([]string) 
}

func NewBartender(configPath string) Bartender {
	b := new(bartender)
	b.Init(configPath)
	return b
}

type bartender struct {
	server 	  *gin.Engine
	database 	  *gorm.DB
	sc 		  sass.Compiler
	config 	  *config
	startTime   time.Time
	logger      *log.Logger
	buildError  error
	app 		  *cli.App
	debug	  bool
	initiated   bool
}



func (b *bartender) Start(args []string) {	
	// Run mainAction
	b.app.Run(args)
}
