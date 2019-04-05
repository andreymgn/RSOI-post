package main

import (
	"github.com/andreymgn/RSOI-post/pkg/post"
)

func runPost(port int, connString string) error {
	server, err := post.NewServer(connString)
	if err != nil {
		return err
	}

	return server.Start(port)
}
