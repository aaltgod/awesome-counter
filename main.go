package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
)

func main() {

	var (
		reader = bufio.NewReader(os.Stdin)
		wg = &sync.WaitGroup{}
		quotaCh = make(chan struct{}, 1)
		total int
	)

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			log.Println(err)
			break
		} else if err == io.EOF {
				break
		}

		line = strings.Replace(line, "\n", "", -1)

		for runtime.NumGoroutine() >= 5 {
		}

		wg.Add(1)

		go func(url string, wg *sync.WaitGroup, quotaCh chan struct{}) {
			defer wg.Done()

			quotaCh <- struct{}{}

			resp, err := http.Get(url)
			if err != nil {
				log.Println(err, "\n[URL]", url)
			}
			defer resp.Body.Close()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println(err, "\n[URL]", url)
			}

			count := strings.Count(string(data), "Go")
			fmt.Printf("Count for %s: %d\n", url, count)

			total += count
			<- quotaCh
		}(line, wg, quotaCh)
	}

	wg.Wait()

	fmt.Printf("Total: %d\n", total)
}