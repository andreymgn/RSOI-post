package main

import (
	"log"
	"os"
	"strconv"
)

func main() {
	conn := os.Getenv("CONN")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Println("PORT parse error")
		return
	}

	jaegerAddr := os.Getenv("JAEGER-ADDR")

	log.Printf("running post service on port %d\n", port)
	err = runPost(port, conn, jaegerAddr)

	if err != nil {
		log.Printf("finished with error %v", err)
	}
}
