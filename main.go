package main

//int add(int a,int b) {
//	return a+b;
//}
import "C"
import (
	"fmt"
	"time"
)

type stu struct {
	Name string `json:"name"`
	Age int
}

//func (s *stu) String() string {
//	return "my name is "+s.Name
//}

func (s *stu) Info() string  {
	return s.Name
}
type People interface {
	Info() string
}
func main() {
	fmt.Println(time.Duration(10)*time.Second)
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