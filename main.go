package main

import (
	"log"
	"net"

	__ "URL_Shorten/proto"
	shorten "URL_Shorten/source"

	"google.golang.org/grpc"
)

func main() {
	s := grpc.NewServer()
	srv := &shorten.GRPCServer{}
	__.RegisterShortenServer(s, srv)

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
