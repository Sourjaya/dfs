package main

import (
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/Sourjaya/dfs/p2p"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts

	peerLock sync.Mutex
	peers    map[string]p2p.Peer

	store  *Store
	quitch chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		StorageRoot:       opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitch:         make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
	}
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	s.peers[p.RemoteAddress().String()] = p
	log.Printf("connected with remote %s", p.RemoteAddress())

	return nil
}

func (s *FileServer) bootstrapNetwork() error {
	for _, address := range s.BootstrapNodes {
		if len(address) == 0 {
			continue
		}
		go func(address string) {
			fmt.Println("attempting to connect to remote: ", address)
			if err := s.Transport.Dial(address); err != nil {
				log.Println("Dial error: ", err)
			}
		}(address)
		//s.Transport.Dial()
	}
	return nil
}

func (s *FileServer) Start() error {
	fmt.Println("Listen and Accept from Server")
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
	fmt.Println("Call bootstrapNetwork")
	s.bootstrapNetwork()
	fmt.Println("Call loop")
	s.loop()
	fmt.Println("Return from loop")

	return nil
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) loop() {
	defer func() {
		log.Println("file server stopped due to user quit action")
		s.Transport.Close()
	}()

	for {
		select {
		case msg := <-s.Transport.Consume():
			fmt.Println(msg)
		case <-s.quitch:
			return
		}
	}
}

func (s *FileServer) Store(key string, r io.Reader) error {
	return s.store.Write(key, r)
}
