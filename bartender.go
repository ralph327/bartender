package bartender

import (
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"time"
)

type Bartender interface {
	Init() 
	Start() 
}

func NewBartender(configPath string) Bartender {
	b := new(bartender)
	b.Init(configPath)
	return b
}

type bartender struct {
	server 	 *gin.Engine
	database 	 *gorm.DB
	config 	 *config
	startTime  time.Time
	logger     *log.Logger
	immediate  bool
	buildError error
	app 		 *cli.App
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
	b.app.Name = "genever"
	b.app.Usage = "A live reload utility for Go web applications."
	b.app.Action = b.mainAction
	b.app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port,p",
			Value: 3000,
			Usage: "port for the proxy server",
		},
		cli.IntFlag{
			Name:  "appPort,a",
			Value: 3001,
			Usage: "port for the Go web server",
		},
		cli.StringFlag{
			Name:  "bin,b",
			Value: "genever",
			Usage: "name of generated binary file",
		},
		cli.StringFlag{
			Name:  "path,t",
			Value: ".",
			Usage: "Path to watch files from",
		},
		cli.BoolFlag{
			Name:  "immediate,i",
			Usage: "run the server immediately after it's built",
		},
		cli.BoolFlag{
			Name:  "godep,g",
			Usage: "use godep when building",
		},
	}
	b.app.Commands = []cli.Command{
		{
			Name:      "run",
			ShortName: "r",
			Usage:     "Run the gin proxy in the current working directory",
			Action:    b.mainAction,
		},
		{
			Name:      "env",
			ShortName: "e",
			Usage:     "Display environment variables set by the .env file",
			Action:    b.envAction,
		},
	}	

	b.startTime  = time.Now()
	b.logger     = log.New(os.Stdout, "[genever] ", 0)
	b.immediate  = false
}

func (b *bartender) Start() {
	b.server.Run(":8989")
}
