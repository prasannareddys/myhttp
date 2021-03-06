package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/prasanna-eyewa/myhttp/external"

	"github.com/prasanna-eyewa/myhttp/api"

	"github.com/prasanna-eyewa/myhttp/utils"
)

const defaultParallelProcesses = 10

func main() {
	parallelProcess := flag.Int("parallel", defaultParallelProcesses, "no of concurrent processes")
	flag.Parse()
	urls := flag.Args()

	wg := new(sync.WaitGroup)
	wg.Add(len(urls))
	concurrentProcess := make(chan bool, *parallelProcess)
	for _, reqUrl := range urls {
		concurrentProcess <- true
		go func(reqUrl string) {
			defer func() {
				wg.Done()
				<-concurrentProcess
			}()
			apiClient := api.NewApiClient(reqUrl, external.GetClient())
			callURLResponse, err := apiClient.CallURL()
			if err != nil {
				log.Printf("Call failed for url \"%s\". Error %s.", reqUrl, err.Error())
			} else {
				fmt.Println(fmt.Sprintf("%s %x", callURLResponse.RequestURL, utils.GetHash(callURLResponse.ResponseBody)))
			}
		}(reqUrl)
	}
	wg.Wait()
}
