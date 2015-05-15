package bartender 

import (
     "github.com/codegangsta/cli"
     "github.com/gin-gonic/gin"
     "github.com/codegangsta/envy/lib"
     "github.com/ralph327/go_sass"
     "net/http"
     "html/template"
     "time"
     "fmt"
     "os"
     "strings"
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
    
     b.startTime  = time.Now()
     b.initiated  = false
	
     b.sass = new(sass.Compiler)
		
     // Initialize folders for web app
     gopath := envy.MustGet("GOPATH")
	
     // Get the working directory
     b.wd, _ = os.Getwd()
	
     // Set the directories that need to be created
     dirs := [...]string{"/views", "/public", "/models", "/controllers"}
	
     // Copy framework directories to the working directory
     // do not copy a directory if it exists
     for _, dir := range dirs {
	  if !dirExists(b.wd + dir) {
	       fmt.Fprintf(os.Stdout, "Copying %s to working directory\n", dir)
	       err = copyDir(gopath + "/src/github.com/ralph327/bartender" + dir, b.wd + dir)
		    
		if err != nil {
			fmt.Fprintf(os.Stderr,"Error while copying %s: %s\n", dir, err)
			os.Exit(1)
		}
	  }
     }
     
     // Initialize gin
     b.server = gin.Default()
     b.server.GET("/", func(c *gin.Context) {
	  c.String(http.StatusOK, b.config.Put)
     })
     
     // Define Routes
     b.server.GET("/", func(c *gin.Context) {
          host_split := strings.Split(c.Request.Host, "."+b.config.DomainName)
	  subdomain := host_split[0]
          data := gin.H{"subdomain": subdomain}
          data["title"] = b.config.SiteName
          data["url"] = b.config.DomainName
                    
          c.HTML(http.StatusOK, "base", data)
     })
     
     // Serve Static files
     b.server.Static("/"+b.config.DomainName+"/style/css","style/css")
     
     // Define templates
     html := template.Must(template.ParseGlob("style/tmpl/*"))
     b.server.SetHTMLTemplate(html)
}
