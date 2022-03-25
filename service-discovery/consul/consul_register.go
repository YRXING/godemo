package consul

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/health/grpc_health_v1"
	"time"
)

// 根据需求进行服务定制
type ConsulService struct {
	IP string
	Port int
	Tag []string
	Name string
}

func RegisterService(consulAddress string, service *ConsulService) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulAddress
	client, err := api.NewClient(consulConfig)
	if err != nil {
		log.Errorf("New consul client err \n:%v",err)
		return
	}
	agent := client.Agent()

	reg := &api.AgentServiceRegistration{
		ID: fmt.Sprintf("%v-%v-%v",service.Name,service.IP,service.Port), //服务节点名称
		Name: service.Name,	//服务名称
		Tags: service.Tag,
		Port: service.Port,
		Address: service.IP,
		Check: &api.AgentServiceCheck{ //健康检查
			//HTTP: fmt.Sprintf("http://%s:%d",service.IP,service.Port),
			Interval: (time.Duration(10)*time.Second).String(),  //检查间隔
			GRPC: fmt.Sprintf("%v:%v/%v", service.IP, service.Port, service.Name), // grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
			DeregisterCriticalServiceAfter: (time.Duration(1)*time.Minute).String(),   //注销时间，相当于过期时间
		},
	}
	log.Printf("register to %v\n",consulAddress)
	if err := agent.ServiceRegister(reg);err!=nil{
		log.Printf("service register error %v",err)
		return
	}
}

func DeregisterService(consulAddress,serviceName string)  {
	cfg := api.DefaultConfig()
	cfg.Address = consulAddress
	client, err := api.NewClient(cfg)
	if err != nil {
		log.Errorf("New consul client err :%v",err)
		return
	}
	log.Printf("deregister from %v\n",consulAddress)
	if err := client.Agent().ServiceDeregister(serviceName);err != nil{
		log.Printf("service deregister error: %v",err)
		return
	}
}

func FilterServiceByID(consulAddress,filterId string) map[string]*api.AgentService {
	cfg := api.DefaultConfig()
	cfg.Address = consulAddress
	client, err := api.NewClient(cfg)
	if err != nil {
		log.Errorf("New consul client err :%v",err)
		return nil
	}
	data ,err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service ==%s", filterId))
	if err != nil {
		panic(err)
	}
	return data
}

func FindServiceByID(consulAddress,serviceID string) *api.AgentService {
	cfg := api.DefaultConfig()
	cfg.Address = consulAddress
	client, err := api.NewClient(cfg)
	if err != nil {
		log.Errorf("New consul client err :%v",err)
		return nil
	}

	service,_,err := client.Agent().Service(serviceID,nil)
	if err != nil {
		return nil
	}
	return service
}

//HealthImpl 定义一个空结构体用来进行健康检查，专门用于gRPC服务的检查
//HealthImpl 实现了HealthServer 这个接口
type HealthImpl struct{}

// Check 实现健康检查接口，这里直接返回健康状态，这里也可以有更复杂的健康检查策略，比如根据服务器负载来返回
func (h *HealthImpl) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	log.Println("health checking")
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch 实现健康检查接口,这个没用，只是为了让HealthImpl实现RegisterHealthServer内部的interface接口
func (h *HealthImpl) Watch(req *grpc_health_v1.HealthCheckRequest, w grpc_health_v1.Health_WatchServer) error {
	return nil
}

/*
当GRPC服务服务启动的时候要添加：
grpc_health_v1.RegisterHealthServer(s, &HealthImpl{})//比普通的grpc开启多了这一步
 */