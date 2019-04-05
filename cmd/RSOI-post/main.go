package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	conn := os.Getenv("CONN")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		fmt.Println("PORT parse error")
		return
	}

	jaegerAddr := os.Getenv("JAEGER-ADDR")

	fmt.Printf("running post service on port %d\n", port)
	err = runPost(port, conn, jaegerAddr)

	if err != nil {
		fmt.Printf("finished with error %v", err)
	}
}
