package consul

// RegisterToConsul 调用RegisterService向consul中注册
func RegisterToConsul() {
	RegisterService("192.168.53.205:8500", &ConsulService{
		Name: "helloworld",
		Tag:  []string{"helloworld", "gopher"},
		IP:   "192.168.53.205",
		Port: 50051,
	})
}

func gRPC_Server(){
	/*
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterGopherServer(s, &server{})
		grpc_health_v1.RegisterHealthServer(s, &HealthImpl{})
		RegisterToConsul()
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	*/
}