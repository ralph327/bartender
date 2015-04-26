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
	"strconv"
	"syscall"
	"time"
)

func (b *bartender) mainAction(c *cli.Context) {
	port := c.GlobalInt("port")
	appPort := strconv.Itoa(c.GlobalInt("appPort"))
	b.immediate = c.GlobalBool("immediate")

	// Bootstrap the environment
	envy.Bootstrap()

	// Set the PORT env
	os.Setenv("PORT", appPort)

	wd, err := os.Getwd()
	if err != nil {
		b.logger.Fatal(err)
	}
	
	b.logger.Println(c.GlobalString("path"), c.GlobalString("bin"), c.GlobalBool("godep"))

	b.logger.Println("before builder")
	builder := genever.NewBuilder(c.GlobalString("path"), c.GlobalString("bin"), c.GlobalBool("godep"))
	b.logger.Println("before runner")
	runner := genever.NewRunner(filepath.Join(wd, builder.Binary()), c.Args()...)
	b.logger.Println("before writer")
	runner.SetWriter(os.Stdout)
	var proxy *genever.Proxy
	if b.proxyOn == false {
		b.logger.Println("before proxy")
		proxy = genever.NewProxy(builder, runner)
		b.logger.Println("after proxy")
	}
	
	config := &genever.Config{
		Port:    port,
		ProxyTo: "http://localhost:" + appPort,
	}

	b.logger.Println("before proxy run")
	if b.proxyOn == false {
		err = proxy.Run(config)
		b.logger.Println("after proxy run")
		if err != nil {
			b.logger.Fatal(err)
		}
	}
	
	b.logger.Printf("listening on port %d\n", port)

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
		if b.immediate {
			runner.Run()
		}
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
