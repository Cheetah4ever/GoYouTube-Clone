package server

import (
	"fmt"
	//"log"
	"net"
	//"bufseekio"
	//"github.com/abema/go-mp4"
	//"os"
	//"tag"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// expand all boxes
/*
func readMp4file(file []byte) {
	_, err := mp4.ReadBoxStructure(file, func(h *mp4.ReadHandle) (interface{}, error) {
		fmt.Println("depth", len(h.Path))

		// Box Type (e.g. "mdhd", "tfdt", "mdat")
		fmt.Println("type", h.BoxInfo.Type.String())

		// Box Size
		fmt.Println("size", h.BoxInfo.Size)

		if h.BoxInfo.IsSupportedType() {
			// Payload
			box, _, err := h.ReadPayload()
			if err != nil {
				return nil, err
			}
			str, err := mp4.Stringify(box, h.BoxInfo.Context)
			if err != nil {
				return nil, err
			}
			fmt.Println("payload", str)

			// Expands children
			return h.Expand()
		}
		return nil, nil
	})
}

func Udp() {
	// listen to incoming udp packets
	conn, err := net.ListenPacket("udp4", "localhost:1053")
	fmt.Println("Echo UPD Server is online")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	dat, err := os.ReadFile("goyoutube/z_client/206294_tiny.mp4")
	//os.Stdout.Write(dat)
	//check(err)
	// cache block size   : 128KBytes
	// cache block history: 4
	//dat := bufseekio.NewReadSeeker("goyoutube/z_client/206294_tiny.mp4", 128 * 1024, 4)
	m, err := tag.ReadFrom("goyoutube/z_client/206294_tiny.mp4")
	if err != nil {
		log.Fatal(err)
	}
	log.Print(m.Format()) // The detected format.
	log.Print(m.Title())  // The title of the track (see Metadata interface for more details).
	buf := make([]byte, 1024)
	_, addr, err := conn.ReadFrom(buf)
	log.Print(addr)
	check(err)
	go serve(conn, addr, dat)
	/*
		for {
			buf := make([]byte, 1024)
			_, addr, err := conn.ReadFrom(buf)
			check(err)
			go serve(conn, addr, dat)
		}
*/

func serve(pc net.PacketConn, addr net.Addr, buf []byte) {
	// 0 - 1: ID
	// 2: QR(1): Opcode(4)
	buf[2] |= 0x80 // Set QR bit
	pc.WriteTo(buf, addr)
	fmt.Printf("%s\n", buf)
}
