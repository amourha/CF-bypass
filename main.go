package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"./http"
	"./output"
	"./providers"
)

const (
	currentVersion = "1.00"
)

func run(config *providers.Config, domains []string) {
	var allProviders []providers.Provider

	for _, providerName := range config.Providers {
		switch providerName {
		case "viewdns":
			viewDNSObj := providers.NewViewDns(config)
			allProviders = append(allProviders, viewDNSObj)
		case "shodan":
			shodanObj := providers.NewShodan(config)
			allProviders = append(allProviders, shodanObj)
		default:
			fmt.Fprintf(os.Stderr, "Error: %s is not a valid provider\n", providerName)
		}
	}

	resultsChannel := make(chan string)
	writewg := &sync.WaitGroup{}
	writewg.Add(1)

	// Start writer thread
	go func() {
		defer writewg.Done()
		err := output.Write(resultsChannel, config.Output)
		if nil != err {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()

	for _, domain := range domains {
		// Run all providers in parallel
		wg := &sync.WaitGroup{}
		wg.Add(len(allProviders))

		for _, provider := range allProviders {
			go func(p providers.Provider) {
				defer wg.Done()
				err := p.BypassCF(domain, resultsChannel)
				if nil != err && config.Verbose {
					fmt.Fprintln(os.Stderr, err)
				}
			}(provider)
		}
		// Wait for providers to finish
		wg.Wait()
	}

	// Close results channel so the writer can return
	close(resultsChannel)
	// Wait for writer to finish
	writewg.Wait()
	os.Exit(0)
}

func main() {
	verbose := flag.Bool("v", false, "enable verbose mode")
	useProviders := flag.String("providers", "viewdns", "providers to try")
	version := flag.Bool("version", false, "show cf-bypass version")
	maxRetries := flag.Uint("retries", 5, "amount of retries for http client")
	output := flag.String("o", "", "filename to write results to")
	shodanAPIKey := flag.String("shodan-key", "", "Your Shodan API key. Required for the shodan provider")
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

	config := providers.Config{
		Verbose:      *verbose,
		Output:       out,
		Providers:    strings.Split(*useProviders, ","),
		Client:       http.NewHTTPClient(*maxRetries),
		ShodanAPIKey: *shodanAPIKey,
	}

	run(&config, domains)
}
