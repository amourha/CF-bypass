package main

import (
	"fmt"
	"flag"
	"os"
	"strings"
	"log"
	"io"
	"bufio"
	"./providers"
	"./http"
)

const (
	currentVersion = "1.00"
)

func run(config *providers.Config, domains []string) {
}

func main() {
	verbose := flag.Bool("v", false, "enable verbose mode")
	useProviders := flag.String("providers", "provider1,provider2", "providers to try")
	version := flag.Bool("version", false, "show cf-bypass version")
	maxRetries := flag.Uint("retries", 5, "amount of retries for http client")
	output := flag.String("o", "", "filename to write results to")
	jsonOut := flag.Bool("json", false, "write output as json")
	flag.Parse()

	if *version {
		fmt.Printf("CF-bypass %v\n", currentVersion)
		os.Exit(0)
	}

	var out io.Writer

	if "" == *output {
		out = os.Stdout
	} else {
		fp, err := os.OpenFile(*output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if nil != err {
			log.Fatalf("Could not open output file: %v\n", err)
		}
		defer fp.Close()
		out = fp
	}

	var domains []string

	if flag.NArg() > 0 {
		domains = flag.Args()
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			domains = append(domains, scanner.Text())
		}
	}

	config := providers.Config {
		Verbose:           *verbose,
		Output:             out,
		JSON:              *jsonOut,
		Providers: strings.Split(*useProviders, ","),
		Client: http.NewHTTPClient(*maxRetries),
	}

	run(&config, domains)
}