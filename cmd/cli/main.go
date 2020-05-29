package main

import (
	"fmt"
	"github.com/djschaap/sqs-to-hec"
	"os"
	"regexp"
)

var (
	buildDt string
	commit  string
	version string
)

func main() {
	printVersion()

	srcQueue := os.Getenv("SRC_QUEUE")
	hasSrcQueue, _ := regexp.MatchString(`^https`, srcQueue)
	if !hasSrcQueue {
		fmt.Println("ERROR: SRC_QUEUE must be set")
		os.Exit(1)
	}

	hecUrl := os.Getenv("HEC_URL")
	hecToken := os.Getenv("HEC_TOKEN")

	hecConfig := sqstohec.HecConfig{
		Token: hecToken,
		Url:   hecUrl,
	}
	sqsConfig := sqstohec.SqsConfig{
		Url: srcQueue,
	}

	app := sqstohec.New(
		hecConfig,
		sqsConfig,
	)

	app.RunOnce()
}

func printVersion() {
	fmt.Println("sqs-to-hec cli  Version:",
		version, " Commit:", commit,
		" Built at:", buildDt)
}
