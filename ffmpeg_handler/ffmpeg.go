package ffmpeg_handler

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	//"github.com/quic-go/quic-go"
)

func serveBufferedRawVideoTCP() {
	log.Print("Opening TCP Sockets")
	listener, err := net.Listen("tcp", "localhost:4200")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	conn, err := listener.Accept()
	fmt.Println("Go Server")
	if err != nil {
		log.Fatal(err)
	}

	pr, pw := io.Pipe()

	log.Print("About to read the video")

	go ffmpeg.Input("goyoutube/z_client/206294_tiny.mp4").
		Output("pipe:", ffmpeg.KwArgs{"format": "rawvideo", "pix_fmt": "rgb24"}).
		WithOutput(pw).
		Run()

	vbuf, err := io.ReadAll(pr)

	if err != nil {
		fmt.Println(err)
	}

	for {
		func(c net.Conn) {
			buf := make([]byte, 4096)
			n, err := c.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s", buf[:n])
			//os.Stdout.Write(buf[:n])
			if bytes.Equal([]byte("Send"), buf[:n]) {
				_, e := conn.Write(vbuf)

				if e != nil {
					fmt.Println(e)
				}
			} else {
				_, e := conn.Write(vbuf)
				if e != nil {
					fmt.Println(e)
				}
			}
		}(conn)
	}
}

func decodeRawVideo() {
	log.Print("ffmpeg")
	pr, pw := io.Pipe()
	ffmpeg.Input("/Users/vgupta/projects/Go Practice/goyoutube/z_client/206294_tiny.mp4").
		Output("pipe:", ffmpeg.KwArgs{"format": "rawvideo", "pix_fmt": "rgb24"}).
		WithOutput(pw).
		ErrorToStdOut().
		Run()
		//Users/vgupta/projects/Go Practice/goyoutube/z_client/206294_tiny.mp4
	_, err := io.ReadAll(pr)
	//log.Print(vbuf)
	if err != nil {
		fmt.Println(err)
	}
}

func streamRawVideoTCP() {

	log.Print("Opening TCP Sockets")
	listener, err := net.Listen("tcp", "localhost:4200")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	conn, err := listener.Accept()
	fmt.Println("Go Server")
	if err != nil {
		log.Fatal(err)
	}
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()
		_ = ffmpeg.Input("goyoutube/z_client/206294_tiny.mp4").
			Output("pipe:", ffmpeg.KwArgs{"format": "rawvideo", "pix_fmt": "rgb24"}).
			WithOutput(pw).
			Run()
	}()

	buf := make([]byte, 4096)
	for {
		n, err := pr.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println(err)
			break
		}
		conn.Write(buf[:n])
	}
}

func streamH264MP4TCP() {
	listener, err := net.Listen("tcp", "localhost:4200")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("Waiting for connection...")
	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Println("Client connected")

	pr, pw := io.Pipe()

	// Spawn ffmpeg process
	go func() {
		out := bytes.NewBuffer(nil)
		defer pw.Close()
		err := ffmpeg.Input("/Users/vgupta/projects/Go Practice/goyoutube/z_client/206294_tiny.mp4").
			Output("pipe:1", ffmpeg.KwArgs{
				"vcodec": "libx264",
				"preset": "fast",
				"f":      "mp4",
			}).
			WithOutput(pw).
			WithErrorOutput(out).
			OverWriteOutput().
			Run()
		fmt.Println(out.String())
		if err != nil {
			log.Printf("ffmpeg error: %v", err)
		}
		log.Println("ffmpeg finished")
	}()

	buf := make([]byte, 32*1024) // stream in chunks
	for {
		n, err := pr.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("read error: %v", err)
			break
		}
		_, err = conn.Write(buf[:n])
		if err != nil {
			log.Printf("write error: %v", err)
			break
		}
	}
	log.Println("stream finished")

}

func streamInputVideoTCP() {
	listener, err := net.Listen("tcp", "localhost:4200")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("Waiting for connection...")
	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Println("Client connected")

	pr, pw := io.Pipe()

	// Spawn ffmpeg process
	go func() {
		out := bytes.NewBuffer(nil)
		defer pw.Close()
		err := ffmpeg.Input("/Users/vgupta/projects/Go Practice/goyoutube/z_client/206294_tiny.mp4").
			Output("pipe:1").
			WithOutput(pw).
			WithErrorOutput(out).
			Run()
		fmt.Println(out.String())
		if err != nil {
			log.Printf("ffmpeg error: %v", err)
		}
		log.Println("ffmpeg finished")
	}()

	buf := make([]byte, 32*1024) // stream in chunks
	for {
		n, err := pr.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("read error: %v", err)
			break
		}
		_, err = conn.Write(buf[:n])
		if err != nil {
			log.Printf("write error: %v", err)
			break
		}
	}
	log.Println("stream finished")

}

func inspectMediaStreams() {

	split := ffmpeg.Input("/Users/vgupta/projects/Go Practice/goyoutube/z_client/206294_tiny.mp4").
		Split()
	fmt.Println(split)
	split0, split1 := split.Get("0"), split.Get("1")
	fmt.Println(split0)
	fmt.Println(split1)
	vid := ffmpeg.Input("/Users/vgupta/projects/Go Practice/goyoutube/z_client/206294_tiny.mp4").
		Video().Output("pipe:0", ffmpeg.KwArgs{"t": "20", `f`: `mp4`, `vcodec`: `rawvideo`})
	fmt.Println(vid)
	vid.Run()
	aud := ffmpeg.Input("/Users/vgupta/projects/Go Practice/goyoutube/z_client/206294_tiny.mp4").
		Audio()
	fmt.Println(aud)
}

func transcodeH265MP4ToBuffer() {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input("/Users/vgupta/projects/Go Practice/goyoutube/z_client/206294_tiny.mp4").
		Output("pipe:", ffmpeg.KwArgs{"c:v": "libx265", "f": "mp4"}).
		WithOutput(buf).
		OverWriteOutput().
		Run()

	fmt.Println(err)
	fmt.Println(buf)

}

func TranscodeSampleToMPEGTS(pw *io.PipeWriter) {

	defer pw.Close()
	err := ffmpeg.Input("/Users/vgupta/projects/Go Practice/goyoutube/client/206294_tiny.mp4").
		Output("pipe:1", ffmpeg.KwArgs{
			"c:v": "libx265",
			"f":   "mpegts",
		}).
		OverWriteOutput().
		WithOutput(pw).
		Run()

	fmt.Println(err)

}

func transcodeH265MatroskaToStdout() {
	err := ffmpeg.Input("/Users/vgupta/projects/Go Practice/goyoutube/z_client/206294_tiny.mp4").
		Output("pipe:", ffmpeg.KwArgs{
			"c:v": "libx265",
			"f":   "matroska", // MKV supports H.265 cleanly
		}).
		OverWriteOutput().
		WithOutput(os.Stdout).
		Run()
	fmt.Println(err)

}

func transcodeFragmentedMP4ToStdout() {
	err := ffmpeg.Input("/Users/vgupta/projects/Go Practice/goyoutube/z_client/206294_tiny.mp4").
		Output("pipe:", ffmpeg.KwArgs{
			"c:v":      "libx265",
			"f":        "mp4",
			"movflags": "frag_keyframe+empty_moov+default_base_moof",
		}).
		OverWriteOutput().
		WithOutput(os.Stdout).
		Run()
	fmt.Println(err)
}

func transcodeH265MPEGTSStdout() {
	err := ffmpeg.Input("./sample_data/in1.mp4").
		Output("pipe:", ffmpeg.KwArgs{
			"c:v":    "libx265",
			"tag:v":  "hvc1", // ensure TS muxer recognizes HEVC
			"f":      "mpegts",
			"vcodec": "libx265",
		}).
		OverWriteOutput().
		WithOutput(os.Stdout).
		Run()
	fmt.Println(err)
}

func serveVideoHTTPExperiment() {
	http.HandleFunc("/video", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "video.mp4")
	})

	log.Println("Serving Video over TCP (HTTP/1.1) on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func registerHTTP3VideoExperiment() {
	http.HandleFunc("/video", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "video.mp4")
	})
	//create ssl: openssl req -x509 -newkey rsa:4096 -nodes -keyout key.pem -out cert.pem -days 365  -subj "/CN=localhost"
	/*
		addr := ":4433"
		certFile := "cert.pem"
		keyFile := "key.pem"
		log.Println("Serving Video over QUIC (HTTP/3) on")
		log.Fatal(http.ListenAndServeQUIC(addr, certFile, keyFile, nil))
	*/
}

func generateHLS() {
	//python3 -m http.server 8080 -d output/
	err := ffmpeg.Input("./sample_data/in1.mp4").
		Output("output/playlist.m3u8", ffmpeg.KwArgs{
			"c:v":                  "libx264",
			"c:a":                  "aac",
			"f":                    "hls", // HLS output format
			"hls_time":             4,     // 4-second segments
			"hls_list_size":        0,     // include all segments in manifest
			"hls_segment_filename": "output/segment_%03d.ts",
		}).
		OverWriteOutput().
		Run()
	if err != nil {
		panic(err)
	}
}

func RemuxMPEGTSFile(tsData []byte, filename string) {
	cmd := exec.Command("ffmpeg",
		"-f", "mpegts", "-i", "pipe:0",
		"-c", "copy",
		filename,
	)
	stdin, _ := cmd.StdinPipe()
	cmd.Start()

	// Write TS bytes:
	stdin.Write(tsData)
	stdin.Close()
	cmd.Wait()

}

func transcodeMPEGTSExampleToMP3() {

	// Example MPEG-TS data — in practice, this comes from a network or file.
	// For demo, replace with real TS bytes.
	tsData := []byte{ /* ... raw MPEG-TS bytes ... */ }

	// Set up FFmpeg command:
	cmd := exec.Command("ffmpeg",
		"-f", "mpegts", // input format
		"-i", "pipe:0", // read from stdin
		"-acodec", "libmp3lame", // encode to MP3
		"-ar", "44100", // sample rate
		"-ab", "192k", // bitrate
		"output.mp3", // output file
	)

	// Attach stdin to send TS bytes.
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("stdin error: %v", err)
	}

	// Optional: capture FFmpeg logs if needed
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Start FFmpeg process
	if err := cmd.Start(); err != nil {
		log.Fatalf("ffmpeg start failed: %v", err)
	}

	// Feed your MPEG-TS bytes
	if _, err := stdin.Write(tsData); err != nil {
		log.Fatalf("writing TS failed: %v", err)
	}
	stdin.Close()

	// Wait for FFmpeg to finish
	if err := cmd.Wait(); err != nil {
		log.Fatalf("ffmpeg error: %v\n%s", err, stderr.String())
	}

	log.Println("Done: output.mp3 created")
}

func transcodeMPEGTSBytesToMP4(tsData []byte) {
	cmd := exec.Command("ffmpeg",
		"-f", "mpegts", "-i", "pipe:0",
		"-c:v", "libx264", "-preset", "veryfast", "-pix_fmt", "yuv420p",
		"-c:a", "aac", "-ar", "44100", "-b:a", "128k",
		"-f", "mp4",
		"pipe:1",
	)
	stdin, _ := cmd.StdinPipe()
	cmd.Start()

	// Write TS bytes:
	stdin.Write(tsData)
	stdin.Close()
	cmd.Wait()

}
