package main

import (
	"os"
	"os/exec"
	"log"
	"flag"
	"time"
	"fmt"
	"bufio"
	"strings"
	"sync"
)

func main() {
	var (
		rate     int
		inflight int
		wg       sync.WaitGroup
	)

	start := time.Now()

	flag.IntVar(&rate, "rate", 1, "максимальное кол-во запусков команды в секунду")
	flag.IntVar(&inflight, "inflight", 1, "максимальное кол-во параллельно запущенных команд")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s --rate <N> --inflight <P> <command...>\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "  <command...>: команда для запуска, {} в команде заменяется на строчку из stdin.\n")
	}
	flag.Parse()

	argCmd := strings.Join(flag.Args(), " ")

	limiter := time.Tick(time.Second * 1)
	jobs := make(chan string)

	for w := 1; w <= inflight; w++ {
		go func() {
			for job := range jobs {
				worker(job, &wg)
			}
		}()
	}

	reader := bufio.NewReader(os.Stdin)
	i := 0
	for {
		i++

		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		jobs <- strings.Replace(argCmd, "{}", strings.Trim(line, "\n"), -1)

		if i%rate == 0 {
			<-limiter
		}
	}
	close(jobs)

	wg.Wait()

	fmt.Println("Elapsed time", time.Since(start))
}

func worker(job string, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	cmd := exec.Command("bash", "-c", job)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}
