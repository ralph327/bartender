package bartender

import (
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"fmt"
	"time"
	"net/http"
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
	config 	  *config
	startTime   time.Time
	logger      *log.Logger
	buildError  error
	app 		  *cli.App
	debug	  bool
	initiated   bool
}

func (b *bartender) Init(configPath string) {
	var err error
	
	// Initialize the cli app
	b.app = cli.NewApp()
	
	// load the configuration
	b.config, err = loadConfig(configPath)
	
	if err != nil {
		fmt.Fprint(os.Stderr,"Error while loading config: %s\n",err)
		os.Exit(1)
	}
	
	// Set application data
	b.app.Name = b.config.SiteName
	b.app.Usage = "Bartender is a framework with a hot reload utility baked in. This is disabled on production environments."
	b.app.Action = b.mainAction
	b.app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "env,e",
			Value: "d",
			Usage: "environment to run server under",
		},
		cli.StringFlag{
			Name:  "debugging,d",
			Value: "t",
			Usage: "Output debugging info or not",
		},
		cli.StringFlag{
			Name:  "proxyPort,p",
			Value: "9000",
			Usage: "port for the proxy server",
		},
		cli.StringFlag{
			Name:  "appPort,a",
			Value: "9001",
			Usage: "port for the Go web server",
			EnvVar: "PORT",
		},
		cli.StringFlag{
			Name:  "bin,b",
			Value: b.config.SiteName,
			Usage: "name of generated binary file",
		},
		cli.StringFlag{
			Name:  "path,t",
			Value: ".",
			Usage: "Path to watch files from",
		},
		cli.BoolFlag{
			Name:  "godep,g",
			Usage: "use godep when building",
		},
	}
	b.app.Commands = []cli.Command{
		{
			Name:	 "init",
			ShortName: "i",
			Usage:	 "Initiate bartender",
			Action:	 b.initAction,
		},
		{
			Name:      "run",
			ShortName: "r",
			Usage:     "Run the genever proxy in the current working directory",
			Action:    b.mainAction,
		},
		{
			Name:      "env",
			ShortName: "e",
			Usage:     "Display environment variables set by the .env file",
			Action:    b.envAction,
		},
	}	
	b.server = gin.Default()
	b.server.GET("/", func(c *gin.Context) {
        c.String(http.StatusOK, b.config.Put)
    })
	b.startTime  = time.Now()
	b.initiated  = false
}

func (b *bartender) Start(args []string) {	
	// Run mainAction
	b.app.Run(args)
}
