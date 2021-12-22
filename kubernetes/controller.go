package main

import (
	"fmt"
	"reflect"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	corev1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

const (
	TunnelServiceNs = "kube-system"

	TunnelServiceName = "x-tunnel-server-svc"

	TunnelServiceType = "LoadBalance"

	TunnelServiceControllerThreadiness = 1
)

// TunnelService is a controller that monitor x-yurt-tunnel-service ip changes
type TunnelServiceController struct {
	client    kubernetes.Interface
	informer  corev1.ServiceInformer
	workqueue workqueue.RateLimitingInterface
	ipChanged chan<- struct{}
}

// NewTunnelServiceController creates a new TunnelServiceController.
func NewTunnelServiceController(clientset kubernetes.Interface, svcinformer corev1.ServiceInformer, ipChanged chan<- struct{}) *TunnelServiceController {
	wq := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	tsc := &TunnelServiceController{
		client:    clientset,
		workqueue: wq,
		ipChanged: ipChanged,
	}
	svcinformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: tsc.updateTunnelService,
	})
	tsc.informer = svcinformer

	return tsc
}

// Run starts the TunnelServiceController.
func (tsc *TunnelServiceController) Run(workers int, stopCh <-chan struct{}) {
	defer runtime.HandleCrash()
	defer tsc.workqueue.ShutDown()

	klog.InfoS("Starting controller", "controller", "TunnelServiceController")
	defer klog.InfoS("Shutting down controller", "controller", "TunnelServiceController")

	if !cache.WaitForCacheSync(stopCh,
		tsc.informer.Informer().HasSynced) {
		klog.Error("sync svc timeout")
		return
	}

	for i := 0; i < workers; i++ {
		go wait.Until(tsc.worker, time.Second, stopCh)
	}

	<-stopCh
}

func (tsc *TunnelServiceController) worker() {
	for tsc.processNextItem() {
	}
}

func (tsc *TunnelServiceController) processNextItem() bool {
	key, quit := tsc.workqueue.Get()
	if quit {
		return false
	}

	metaName, ok := key.(string)
	if !ok {
		tsc.workqueue.Forget(key)
		runtime.HandleError(
			fmt.Errorf("expected string in workqueue but got %#v", key))
		return true
	}

	defer tsc.workqueue.Done(metaName)

	namespace, name, err := cache.SplitMetaNamespaceKey(metaName)
	if err != nil {
		tsc.workqueue.Forget(key)
		runtime.HandleError(
			fmt.Errorf("failed to split meta namespace cache key"))
		return true
	}

	_, err = tsc.informer.Lister().Services(namespace).Get(name)
	if err != nil {
		runtime.HandleError(err)
		if !errors.IsNotFound(err) {
			tsc.workqueue.AddRateLimited(key)
		}
		return true
	}

	return true
}

func (tsc *TunnelServiceController) updateTunnelService(old, cur interface{}) {
	osvc := old.(*v1.Service)
	nsvc := cur.(*v1.Service)
	if !isTunnelService(nsvc) {
		klog.Infof("service %v is not %s service", nsvc, TunnelServiceName)
		return
	}
	if reflect.DeepEqual(osvc.Status.LoadBalancer.Ingress, nsvc.Status.LoadBalancer.Ingress) {
		return
	}
	klog.Infof("the service load balancer's ip has changed")

	tsc.enqueueSvc(nsvc)
}

func (tsc *TunnelServiceController) enqueueSvc(svc *v1.Service) {

	key, err := cache.MetaNamespaceKeyFunc(svc)
	if err != nil {
		runtime.HandleError(fmt.Errorf("couldn't get key for object %#v: %v", svc, err))
		return
	}

	tsc.workqueue.AddRateLimited(key)
}
