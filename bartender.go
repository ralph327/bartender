package bartender

import (
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"os"
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
	proxyOn	  bool
	initiated   bool
}

func (b *bartender) Init(configPath string) {
	var err error
	
	// Initialize the cli app
	b.app = cli.NewApp()
	
	// load the configuration
	b.config, err = loadConfig(configPath)
	
	if err != nil {
		b.logger.Println(err)
		os.Exit(1)
	}
	
	// Set application data
	b.app.Name = b.config.SiteName
	b.app.Usage = "Bartender is a framework with a hot reload utility baked in. This is disabled on production environments."
	b.app.Action = b.initAction
	b.app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "env,e",
			Value: "d",
			Usage: "environment to run server under",
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
	b.logger     = log.New(os.Stdout, "[genever] ", 0)
	b.proxyOn    = false
	b.initiated  = false
}

func (b *bartender) Start(args []string) {
	b.logger.Println(args)

	// Run initAction
	if len(args) == 1 {
		b.app.Run(args)
		
		// Ensure that genever will run
		args = append(args, "c")
	}

	b.logger.Println("Env: ", b.config.Environment)

	// Run server based on environment
	switch b.config.Environment {
		case "production", "prod", "p", "child", "c":
			b.server.Run(":8989")
		case "development", "dev", "d":
			b.app.Run(args)
	}
}
