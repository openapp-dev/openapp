package mdns

import (
	"github.com/openapp-dev/openapp/pkg/controller/types"
)

type OpenAPPMDNSController struct {
}

func NewOpenAPPMDNSController() types.ControllerInterface {
	oc := &OpenAPPMDNSController{}
	return oc
}

func (oc *OpenAPPMDNSController) Start() {
	mdnsFunc := func() {
	}

	go mdnsFunc()
}

func (oc *OpenAPPMDNSController) Reconcile(_ interface{}) error {
	return nil
}
