package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/credentials/oauth"
	"io/ioutil"
	"log"
	auth "godemo/gRPC/productinfo/authorize"
	pb "godemo/gRPC/productinfo/client/ecommerce"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	address    = "localhost:50051"
	serverName = "www.xing.com"
	crtFile    = "../pki/client.crt"
	keyFile    = "../pki/client.key"
	caFile     = "../pki/ca.crt"
)

func mTLSClient() {
	//cert , err := credentials.NewClientTLSFromFile(crtFile,serverName)
	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("failed to load credentials\n")
	}

	// 创建认证服务器的权威根证书列表
	rootCAPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatalf("could not read ca certificate\n")
	}
	if ok := rootCAPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("failed to append ca certs\n")
	}

	//创建连接选项
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			ServerName:   serverName,
			Certificates: []tls.Certificate{cert},
			RootCAs:      rootCAPool,
		})),
		//grpc.WithTransportCredentials(cert),
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect:%v", err)
	}

	defer conn.Close()
	c := pb.NewProductInfoClient(conn)

	name := "Apple iPhone 12 pro"
	description := `Meet Apple iPhone 12 pro. All-new dual-camera system with Ultra while and Night mode.`
	price := float32(1000.0)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.AddProduct(ctx, &pb.Product{Name: name, Description: description, Price: price})

	if err != nil {
		log.Fatal("Could not add product :%v", err)
	}
	log.Printf("Product ID :%s added successfully", r.Value)

	product, err := c.GetProduct(ctx, &pb.ProductID{Value: r.Value})
	if err != nil {
		log.Fatal("Could not get product:%v", err)
	}
	log.Printf("Product: ", product.String())
}

func basicAuthClient() {
	creds, _ := credentials.NewClientTLSFromFile(caFile, serverName)

	auth := auth.NewBasicAuth("admin","admin")

	opts := []grpc.DialOption{
		//传递凭证
		grpc.WithPerRPCCredentials(auth),
		grpc.WithTransportCredentials(creds),
	}

	conn, _ := grpc.Dial(address,opts...)
	defer conn.Close()
	c := pb.NewProductInfoClient(conn)

	name := "Apple iPhone 12 pro"
	description := `Meet Apple iPhone 12 pro. All-new dual-camera system with Ultra while and Night mode.`
	price := float32(1000.0)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.AddProduct(ctx, &pb.Product{Name: name, Description: description, Price: price})

	if err != nil {
		log.Fatal("Could not add product :%v", err)
	}
	log.Printf("Product ID :%s added successfully", r.Value)
}

func oAuthClient() {
	auth := oauth.NewOauthAccess(fetchToken())

	creds, _ := credentials.NewClientTLSFromFile(caFile,serverName)
	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(auth),
		grpc.WithTransportCredentials(creds),
	}

	conn, _ := grpc.Dial(address,opts...)
	defer conn.Close()
	c := pb.NewProductInfoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	name := "Apple iPhone 12 pro"
	description := `Meet Apple iPhone 12 pro. All-new dual-camera system with Ultra while and Night mode.`
	price := float32(1000.0)

	r, err := c.AddProduct(ctx, &pb.Product{Name: name, Description: description, Price: price})

	if err != nil {
		log.Fatal("Could not add product :%v", err)
	}
	log.Printf("Product ID :%s added successfully", r.Value)
}

//由于没有授权服务器，硬编码任意一个字符串来作为令牌
func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: "some-secret-token",
	}
}

func main() {
	//oAuthClient()
	conn , err := grpc.Dial("localhost:8080",grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	pb.NewProductInfoClient(conn)
	fmt.Println("success")
}
