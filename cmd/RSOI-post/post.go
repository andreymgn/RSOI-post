package main

import (
	"log"

	"github.com/andreymgn/RSOI/pkg/tracer"
	"github.com/andreymgn/RSOI-post/pkg/post"
)

const (
	PostAppID     = "PostAPI"
	PostAppSecret = "0JDt37eVLP0VcEJB"
)

func runPost(port int, connString, jaegerAddr, redisAddr, redisPassword string, redisDB int) error {
	tracer, err := tracer.NewTracer("post", jaegerAddr)
	if err != nil {
		log.Fatal(err)
	}

	knownKeys := map[string]string{PostAppID: PostAppSecret}

	server, err := post.NewServer(connString, redisAddr, redisPassword, redisDB, knownKeys)
	if err != nil {
		log.Fatal(err)
	}

	return server.Start(port, tracer)
}
