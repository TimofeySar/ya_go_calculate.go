package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/TimofeySar/ya_go_calculate.go/internal/agent"
)

func main() {
	power, _ := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if power <= 0 {
		power = 1
	}
	fmt.Printf("Agent starting with %d workers\n", power)
	agent.Run(power)
}
