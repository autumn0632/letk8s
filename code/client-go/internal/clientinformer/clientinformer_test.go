package clientinformer

import (
	"testing"

	"autumn.io/client-go/internal/config"
	"k8s.io/client-go/rest"
)

var kubeConfig *rest.Config

func init() {

	c, err := config.NewConfig("../config/testdata/kube/config")
	if err != nil {
		panic(err)
	}
	kubeConfig = c

}

func TestDoInformer(t *testing.T) {
	DoInformer(kubeConfig)
}
