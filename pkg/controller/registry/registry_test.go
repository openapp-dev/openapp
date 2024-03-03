package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloneOpenAPPRegistry(t *testing.T) {
	registryList := []string{"https://github.com/openapp-dev/openapp-registry@main"}

	err := CloneOpenAPPRegistry(registryList)

	assert.NoError(t, err)
}
