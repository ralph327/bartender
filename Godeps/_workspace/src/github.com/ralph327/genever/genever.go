package genever
/*
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

var (
	startTime  = time.Now()
	logger     = log.New(os.Stdout, "[genever] ", 0)
	immediate  = false
	buildError error
)

func MainAction(c *cli.Context) {
	port := c.GlobalInt("port")
	appPort := strconv.Itoa(c.GlobalInt("appPort"))
	immediate = c.GlobalBool("immediate")

	// Bootstrap the environment
	envy.Bootstrap()

	// Set the PORT env
	os.Setenv("PORT", appPort)

	wd, err := os.Getwd()
	if err != nil {
		logger.Fatal(err)
	}

	builder := genever.NewBuilder(c.GlobalString("path"), c.GlobalString("bin"), c.GlobalBool("godep"))
	runner := genever.NewRunner(filepath.Join(wd, builder.Binary()), c.Args()...)
	runner.SetWriter(os.Stdout)
	proxy := genever.NewProxy(builder, runner)

	config := &genever.Config{
		Port:    port,
		ProxyTo: "http://localhost:" + appPort,
	}

	err = proxy.Run(config)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Printf("listening on port %d\n", port)

	shutdown(runner)

	// build right now
	build(builder, runner, logger)

	// scan for changes
	scanChanges(c.GlobalString("path"), func(path string) {
		runner.Kill()
		build(builder, runner, logger)
	})
}

func EnvAction(c *cli.Context) {
	// Bootstrap the environment
	env, err := envy.Bootstrap()
	if err != nil {
		logger.Fatalln(err)
	}

	for k, v := range env {
		fmt.Printf("%s: %s\n", k, v)
	}

}

func build(builder genever.Builder, runner genever.Runner, logger *log.Logger) {
	err := builder.Build()
	if err != nil {
		buildError = err
		logger.Println("ERROR! Build failed.")
		fmt.Println(builder.Errors())
	} else {
		// print success only if there were errors before
		if buildError != nil {
			logger.Println("Build Successful")
		}
		buildError = nil
		if immediate {
			runner.Run()
		}
	}

	time.Sleep(100 * time.Millisecond)
}

type scanCallback func(path string)

func scanChanges(watchPath string, cb scanCallback) {
	for {
		filepath.Walk(watchPath, func(path string, info os.FileInfo, err error) error {
			if path == ".git" {
				return filepath.SkipDir
			}

			// ignore hidden files
			if filepath.Base(path)[0] == '.' {
				return nil
			}

			if (filepath.Ext(path) == ".go" || filepath.Ext(path) == ".tmpl" || filepath.Ext(path) == ".css" || filepath.Ext(path) == ".json") && info.ModTime().After(startTime) {
				cb(path)
				startTime = time.Now()
				return errors.New("done")
			}

			return nil
		})
		time.Sleep(500 * time.Millisecond)
	}
}

func shutdown(runner genever.Runner) {
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
*/
