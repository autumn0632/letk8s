package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	config, err := NewConfig("./testdata/kube/config")
	// 判断err

	assert.NoError(t, err)

	fmt.Printf("%s\n", config)

}
