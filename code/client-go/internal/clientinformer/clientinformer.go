package clientinformer

import (
	"log"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func DoInformer(config *rest.Config) {

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	stopCh := make(chan struct{})
	defer close(stopCh)
	//表示每分钟进行一次resync，resync会周期性地执行List操作
	sharedInformers := informers.NewSharedInformerFactory(clientset, time.Second*10)

	// 每种资源都有自己的informer
	informer := sharedInformers.Core().V1().Pods().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			mObj := obj.(v1.Object)
			log.Printf("%s: New Pod Added to Store - %s", mObj.GetNamespace(), mObj.GetName())
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oObj := oldObj.(v1.Object)
			nObj := newObj.(v1.Object)
			log.Printf("%s Pod Updated to %s", oObj.GetName(), nObj.GetName())
		},
		DeleteFunc: func(obj interface{}) {
			mObj := obj.(v1.Object)
			log.Printf("%s: Pod Deleted from Store: %s", mObj.GetNamespace(), mObj.GetName())
		},
	})

	informer.Run(stopCh)

}
