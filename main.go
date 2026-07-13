package main

import (
	"flag"
	"log"

	"goyoutube/server"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "HTTP server address")
	video := flag.String("video", "client/206294_tiny.mp4", "path to the MP4 video")
	flag.Parse()

	if err := server.RunHTTP(*addr, *video); err != nil {
		log.Fatal(err)
	}
}
