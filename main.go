package main

import (
	"log"

	"github.com/Sourjaya/dfs/p2p"
)

func OnPeer(peer p2p.Peer) error {
	//peer.Close()
	return nil
}

func makeServer(listenAddress string, nodes ...string) *FileServer {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddress: listenAddress,
		Decoder:       p2p.DefaultDecoder{},
		HandshakeFunc: p2p.NOPHandshakeFunc,
	}
	tcpTransport := p2p.NewTCPTransport(tcpOpts)
	//fmt.Printf("tcpTransport for server at %v is %+v\n", listenAddress, tcpTransport)

	fileServerOpts := FileServerOpts{
		StorageRoot:       listenAddress + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}
	s := NewFileServer(fileServerOpts)

	//fmt.Printf("s for server at %v is %+v\n", listenAddress, s)

	tcpTransport.OnPeer = s.OnPeer

	//fmt.Printf("tcpTransport for server at %v is %+v after setting OnPeer\n ", listenAddress, tcpTransport)

	return s
}

func main() {

	// go func() {
	// 	time.Sleep(time.Second * 3)
	// 	server.Stop()
	// }()
	//fmt.Println("making servers")
	server1 := makeServer(":3000", "")
	server2 := makeServer(":4000", ":3000")
	//fmt.Println("Done making servers")

	go func() {
		//fmt.Println("Starting server 1")
		log.Fatal(server1.Start())
	}()
	//fmt.Println("Starting server 2")
	server2.Start()
	// fmt.Println("Starting server 2")
	// server2.Start()
	// if err := server.Start(); err != nil {
	// 	log.Fatal(err)
	// }
	//select {}

	// go func() {
	// 	for {
	// 		msg := <-tr.Consume()
	// 		fmt.Printf("%+v\n", msg)
	// 	}
	// }()
	// if err := tr.ListenAndAccept(); err != nil {
	// 	log.Fatal(err)
	// }
	//select {}
}
