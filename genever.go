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
	b.logger.Println("initating genever")
	
	// Check and set environment
	b.config.Environment = c.String("env")
	
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
	
	b.initiated = true
}

func (b *bartender) checkMode(c *cli.Context) {
	// running with no commands or flags
	if len(c.Args()) == 1 {
		
	} else if len(c.Args()) > 1 {
		if c.IsSet("env") || c.IsSet("e") {
			var tempEnv string
			tempEnv = c.String("env")
			if tempEnv == "" {
				tempEnv = c.String("e")
			}
		}
	}
}

// runs genever
func (b *bartender) mainAction(c *cli.Context) {
	
	// Run initiate if haven't
	if b.initiated == false {
		b.initAction(c)
	}
	
	// Bootstrap the environment
	envy.Bootstrap()
	
	// check the run mode
	b.checkMode(c)
	
	b.logger.Println(c.String("env"))
	
	// Get working directory
	wd, err := os.Getwd()
	if err != nil {
		b.logger.Fatal(err)
	}
	
	// Set builder and runner
	b.logger.Println("before builder")
	b.logger.Println("Running builder: ", c.GlobalString("path"), c.GlobalString("bin"), c.GlobalBool("godep"))
	builder := genever.NewBuilder(c.GlobalString("path"), c.GlobalString("bin"), c.GlobalBool("godep"))
	b.logger.Println("before runner")
	b.logger.Println("Running runner:", filepath.Join(wd, builder.Binary()), "c")
	runner := genever.NewRunner(filepath.Join(wd, builder.Binary()), "c")
	b.logger.Println("before writer")
	runner.SetWriter(os.Stdout)
	
	// Run proxy if not child process and not in production
	proxyPort := b.config.ProxyPort
	appPort := b.config.AppPort
	
	// Set the PORT env
	os.Setenv("PORT", appPort)
	
	var proxy *genever.Proxy
	if b.proxyOn == false {
		b.logger.Println("before proxy")
		proxy = genever.NewProxy(builder, runner)
		b.logger.Println("after proxy")
	}
	
	config := &genever.Config{
		Port:    proxyPort,
		ProxyTo: "http://localhost:" + appPort,
	}

	b.logger.Println("before proxy run")
	if b.proxyOn == false {
		err = proxy.Run(config)
		b.logger.Println("after proxy run")
		if err != nil {
			b.logger.Println("Fatal error after proxy run")
			b.logger.Fatal(err)
		}else{
			b.proxyOn = true
		}
	}
	
	b.logger.Printf("listening on port %s\n", proxyPort)

	b.logger.Println("before shutdown")
	b.shutdown(runner)
	b.logger.Println("after shutdown")

	// build right now
	b.logger.Println("before build")
	b.build(builder, runner, b.logger)
	b.logger.Println("after build")

	// scan for changes
	b.scanChanges(c.GlobalString("path"), func(path string) {
		err := runner.Kill()
		
		b.logger.Println("Kill err:", err)
		
		b.logger.Println("build after kill")
		b.build(builder, runner, b.logger)
		b.logger.Println("after build after kill")
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
	err := builder.Build()
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
		b.logger.Println("RUNNING")
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

			if (filepath.Ext(path) == ".go" || filepath.Ext(path) == ".tmpl" || filepath.Ext(path) == ".css" || filepath.Ext(path) == ".json") && info.ModTime().After(b.startTime) {
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
