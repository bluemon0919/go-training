package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

var ql []string = []string{
	"hoge",
	"fuga",
	"the",
	"time",
	"returned",
	"by",
	"time",
	"now",
	"contains",
	"a",
	"monotonic",
	"clock",
	"reading",
	"if",
	"time",
	"t",
	"has",
	"a",
	"monotonic",
	"clock",
	"reading",
	"t",
	"add",
	"adds",
	"the",
	"same",
}

func main() {
	ok, ng := run()

	fmt.Println("OK:", ok)
	fmt.Println("NG:", ng)
}

func run() (int, int) {
	ch := input(os.Stdin)

	bc := context.Background()
	t := 10 * time.Second
	ctx, cancel := context.WithTimeout(bc, t)
	defer cancel()

	countOK := 0
	countNG := 0

	for _, q := range ql {
		fmt.Println("$ ", q)
		select {
		case v1 := <-ch:
			result := judge(v1, q)
			if "OK" == result {
				countOK++
			} else {
				countNG++
			}
			fmt.Println(result)
		case <-ctx.Done():
			fmt.Println("Timeout!!")
			return countOK, countNG
		}
	}
	return countOK, countNG
}

func judge(s string, q string) string {
	result := "NG"
	if s == q {
		result = "OK"
	}
	return result
}

func input(r io.Reader) <-chan string {
	ch := make(chan string)
	go func() {
		s := bufio.NewScanner(r)
		for s.Scan() {
			ch <- s.Text()
		}
		close(ch)
	}()
	return ch
}
