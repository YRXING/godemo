package main

import (
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"log"
)

func main() {
	addr ,_ :=ma.NewMultiaddr("/ip4/127.0.0.1/tcp/12000/ipfs/QmddTrQXhA9AkCpXPTkcY7e22NK73TwkUms3a44DhTKJTD")
	pid, err := addr.ValueForProtocol(ma.P_IPFS)
	fmt.Println(pid)
	if err != nil {
		log.Fatalln(err)
	}

	peerid, err := peer.Decode(pid)
	fmt.Println(peerid)
	targetPeerAddr, _ := ma.NewMultiaddr(
		fmt.Sprintf("/ipfs/%s", peer.Encode(peerid)))
	fmt.Println(targetPeerAddr)
	targetAddr := addr.Decapsulate(targetPeerAddr)
	fmt.Println(targetAddr)
}
