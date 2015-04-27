package bartender

import (
	"errors"
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/codegangsta/envy/lib"
	"github.com/ralph327/genever"

	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// Checks and sets flags
func (b *bartender) initAction(c *cli.Context) {
	
	// Check and set environment
	b.config.Environment = c.String("env")
	
	// Initiate logger print out
	switch b.config.Environment {
		case "production", "p", "development", "d":
			b.logger = log.New(os.Stdout, "[genever] ", 0)
		case "child", "c" :
			b.logger = log.New(os.Stdout, "[genever-child] ", 0)
	}
	
	b.logger.Println("initating genever")
	
	// Ensure proxy and app ports are set
	var tempPort string
	
	tempPort = c.String("appPort")
	// Override config file to use flag
	if tempPort != "" {
		b.config.AppPort = tempPort
	}
	// If still empty, set to default
	if b.config.AppPort == "" {		
		b.config.AppPort = "9001"
	}
	
	tempPort = c.String("proxyPort")
	// Override config file to use flag
	if tempPort != "" {
		b.config.ProxyPort = tempPort
	}
	// If still empty, set to default
	if b.config.ProxyPort == "" {
		b.config.ProxyPort = "9000"
	}
	
	// Set debugging flag, on by default
	switch b.config.Debugging {
		case "","t", "true":
			b.config.Debugging = "true"
			b.debug = true
		case "false", "f":
			b.config.Debugging = "false"
			b.debug = false
	}
	
	// Override config with command line args
	if c.IsSet("debugging") {
		switch c.GlobalString("debugging") {
			case "t", "true" :
				b.config.Debugging = "true"
				b.debug = true
			case "f", "false" :
				b.config.Debugging = "false"
				b.debug = false
		}
	}
	
	b.initiated = true
}

// runs genever
func (b *bartender) mainAction(c *cli.Context) {
	
	// Run initiate if haven't
	if b.initiated == false {
		b.initAction(c)
	}
	
	// Bootstrap the environment
	envy.Bootstrap()
	
	// If the environment is set to production or child
	switch b.config.Environment {
		case "production", "p", "child", "c":
			b.server.Run(":" + b.config.AppPort)
	}
	
	// Get working directory
	wd, err := os.Getwd()
	if err != nil {
		b.logger.Fatal(err)
	}
	
	// Set builder and runner
	if b.debug {
		b.logger.Println("Running NewBuilder: ", c.GlobalString("path"), c.GlobalString("bin")+"_child", c.GlobalBool("godep"))
	}
	builder := genever.NewBuilder(c.GlobalString("path"), c.GlobalString("bin")+"_child", c.GlobalBool("godep"))
	
	if b.debug {
		b.logger.Println("Running NewRunner:", filepath.Join(wd, builder.Binary()), "-e", "c", "-d", "f")
	}
	runner := genever.NewRunner(filepath.Join(wd, builder.Binary()), "-e", "c", "-d", "f")
	
	runner.SetWriter(os.Stdout)
	
	// Run proxy for development environment
	var proxy *genever.Proxy
	if c.String("env") == "development" || c.String("env") == "d" || c.String("env") == "dev" {
		
		// Run proxy if not child process and not in production
		proxyPort := b.config.ProxyPort
		appPort := b.config.AppPort
		
		// Set the PORT env
		os.Setenv("PORT", appPort)
		proxy = genever.NewProxy(builder, runner)

		config := &genever.Config{
			Port:    proxyPort,
			ProxyTo: "http://localhost:" + appPort,
		}

		err = proxy.Run(config)

		if err != nil {
			if b.debug {
				b.logger.Println("Fatal error on proxy run")
				b.logger.Fatal(err)
			}
		}
		
		if b.debug {
			b.logger.Printf("listening on port %s\n", proxyPort)
		}
	}

	b.shutdown(runner)

	// build right now
	b.build(builder, runner, b.logger)

	// scan for changes
	b.scanChanges(c.GlobalString("path"), func(path string) {
		err := runner.Kill()
		
		if b.debug {
			if err != nil {
				b.logger.Println("Kill err:", err)
			}
		}
		b.build(builder, runner, b.logger)
	})
}

// reads .env
func (b *bartender) envAction(c *cli.Context) {
	// Bootstrap the environment
	env, err := envy.Bootstrap()
	if err != nil {
		b.logger.Fatalln(err)
	}

	for k, v := range env {
		fmt.Printf("%s: %s\n", k, v)
	}

}

func (b *bartender) build(builder genever.Builder, runner genever.Runner, logger *log.Logger) {
	
	// Build the binary
	err := builder.Build()
	
	// Scan and compile scss files
     b.sc.CompileFolder("views/sass/", "public/css")

	if err != nil {
		b.buildError = err
		b.logger.Println("ERROR! Build failed.")
		fmt.Println(builder.Errors())
	} else {
		// print success only if there were errors before
		if b.buildError != nil {
			b.logger.Println("Build Successful")
		}
		b.buildError = nil
		
		// run the server
		runner.Run()
	}

	time.Sleep(100 * time.Millisecond)
}

type scanCallback func(path string)

func (b *bartender) scanChanges(watchPath string, cb scanCallback) {
	b.logger.Println("scanning for changes")
	for {
		filepath.Walk(watchPath, func(path string, info os.FileInfo, err error) error {
			if path == ".git" {
				return filepath.SkipDir
			}

			// ignore hidden files
			if filepath.Base(path)[0] == '.' {
				return nil
			}

			if (filepath.Ext(path) == ".go" || filepath.Ext(path) == ".tmpl" || filepath.Ext(path) == ".scss" || filepath.Ext(path) == ".json") && info.ModTime().After(b.startTime) {
				cb(path)
				b.startTime = time.Now()
				return errors.New("done")
			}

			return nil
		})
		time.Sleep(500 * time.Millisecond)
	}
}

func (b *bartender) shutdown(runner genever.Runner) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-c
		log.Println("Got signal: ", s)
		err := runner.Kill()
		if err != nil {
			log.Print("Error killing: ", err)
		}
		os.Exit(1)
	}()
}
