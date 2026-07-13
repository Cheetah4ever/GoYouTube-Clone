package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"goyoutube/ffmpeg_handler"
)

const legacyTCPOutputPath = "client/out.mp4"

// RunTCP starts the original TCP video experiment.
func RunTCP() {

	listener, err := net.Listen("tcp", "localhost:4200")
	if err != nil {
		log.Fatal(err)
	}
	//defer listener.Close()

	conn, err := listener.Accept()
	fmt.Println("Go Server")
	if err != nil {
		log.Fatal(err)
	}
	for {
		func(c net.Conn) {
			buf := make([]byte, 1024)
			n, err := c.Read(buf)
			fmt.Fprintf(os.Stdout, "%s\n", buf[:n])
			if err != nil {
				log.Fatal(err)
			}
			if bytes.Equal([]byte("Video"), buf[:n]) {
				pr, pw := io.Pipe()

				go ffmpeghandler.TranscodeSampleToMPEGTS(pw)
				vidbuf, err := io.ReadAll(pr)
				if err != nil {
					fmt.Println(err)
				}
				ffmpeghandler.RemuxMPEGTSFile(vidbuf, legacyTCPOutputPath)
				_, e := conn.Write(vidbuf)
				if e != nil {
					fmt.Println(e)
				}
			} else {
				//os.Stdout.Write(buf[:n])
				_, e := conn.Write(buf[:n])
				if e != nil {
					fmt.Println(e)
				}
			}
		}(conn)
	}
}
