package main

import "earnth/enet"

func main() {
	s := enet.NewServer("earnth-v0.1")
	s.Serve()
}
