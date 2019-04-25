package main

import (
	"github.com/andreymgn/RSOI-post/pkg/post"
	"github.com/andreymgn/RSOI/pkg/tracer"
)

func runPost(port int, connString, jaegerAddr string) error {
	tracer, closer, err := tracer.NewTracer("post", jaegerAddr)
	if err != nil {
		return err
	}

	defer closer.Close()

	server, err := post.NewServer(connString)
	if err != nil {
		return err
	}

	return server.Start(port, tracer)
}
