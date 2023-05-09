package main

import "fmt"

func bkdrhash(s string) int {
	seed := 131
	hash := 0
	for i := 0; i < len(s); i++ {
		hash = hash*seed + int(s[i]);
	}
	return hash & 0x7FFFFFFF
}
func main() {
	fmt.Println(bkdrhash("rold=db;namespace=defaut"))
}

func libp2p()  {
	//addr ,_ :=ma.NewMultiaddr("/ip4/127.0.0.1/tcp/12000/ipfs/QmddTrQXhA9AkCpXPTkcY7e22NK73TwkUms3a44DhTKJTD")
	//pid, err := addr.ValueForProtocol(ma.P_IPFS)
	//fmt.Println(pid)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//peerid, err := peer.Decode(pid)
	//fmt.Println(peerid)
	//targetPeerAddr, _ := ma.NewMultiaddr(
	//	fmt.Sprintf("/ipfs/%s", peer.Encode(peerid)))
	//fmt.Println(targetPeerAddr)
	//targetAddr := addr.Decapsulate(targetPeerAddr)
	//fmt.Println(targetAddr)
}