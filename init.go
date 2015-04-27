package bartender 

import (
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
	"github.com/codegangsta/envy/lib"
	"net/http"
	"time"
	"fmt"
	"os"
	"io"
)

func (b *bartender) Init(configPath string) {
	var err error
	
	// Initialize the cli app
	b.app = cli.NewApp()
	
	// load the configuration
	b.config, err = loadConfig(configPath)
	
	if err != nil {
		fmt.Fprintf(os.Stderr,"Error while loading config: %s\n", err)
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
	
	// Initialize folders for web app
	gopath := envy.MustGet("GOPATH")
	
	// Get the working directory
	wd, _ := os.Getwd()
	
	// Copy framework views to the working directory
	// if the directory does not already exist
	if !dirExists(wd + "/views") {
		fmt.Fprintf(os.Stdout, "Copying /views to working directory")
		err = copyFile(gopath + "/src/github.com/ralph327/bartender/views", wd + "/views")
		
		if err != nil {
			fmt.Fprintf(os.Stderr,"Error while copying /views: %s\n", err)
			os.Exit(1)
		}
	}
	
	// Copy framework views to the working directory
	// if the directory does not already exist
	if !dirExists(wd + "/public") {
		fmt.Fprintf(os.Stdout, "Copying /public to working directory")
		err = copyFile(gopath + "/src/github.com/ralph327/bartender/public", wd + "/public")
		
		if err != nil {
			fmt.Fprintf(os.Stderr,"Error while copying /public: %s\n", err)
			os.Exit(1)
		}
	}
}

func dirExists(dir string) bool {
	// check if dir exists
	src,err := os.Stat(dir)
	
	
	if err == nil {
		if !src.IsDir() {
			fmt.Fprintf(os.Stderr, "%s is not a directory", dir)
			os.Exit(1)
		}
		return true
	}
	
	// Check if err is due to not exist
	if os.IsNotExist(err) { 
		return false
	}
	
	// Default to false with errors
	return false
}

 func copyFile(source string, dest string) (err error) {
     sourcefile, err := os.Open(source)
     if err != nil {
         return err
     }

     defer sourcefile.Close()

     destfile, err := os.Create(dest)
     if err != nil {
         return err
     }

     defer destfile.Close()

     _, err = io.Copy(destfile, sourcefile)
     if err == nil {
         sourceinfo, err := os.Stat(source)
         if err != nil {
             err = os.Chmod(dest, sourceinfo.Mode())
         }
     }

     return
 }

// Recursively copy source to dest
 func copyDir(source string, dest string) (err error) {

     // get properties of source dir
     sourceinfo, err := os.Stat(source)
     if err != nil {
         return err
     }

     // create dest dir
     err = os.MkdirAll(dest, sourceinfo.Mode())
     if err != nil {
         return err
     }

     directory, _ := os.Open(source)

     objects, err := directory.Readdir(-1)

     for _, obj := range objects {

         sourcefilepointer := source + "/" + obj.Name()

         destinationfilepointer := dest + "/" + obj.Name()


         if obj.IsDir() {
             // create sub-directories - recursively
             err = copyDir(sourcefilepointer, destinationfilepointer)
             if err != nil {
                 fmt.Println(err)
             }
         } else {
             // perform copy
             err = copyFile(sourcefilepointer, destinationfilepointer)
             if err != nil {
                 fmt.Println(err)
             }
         }

     }
     return
 }
