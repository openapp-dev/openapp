package types

import "github.com/openapp-dev/openapp/pkg/utils"

type ControllerInterface interface {
	Start()
}

type NewControllerFunc func(openappHelper *utils.OpenAPPHelper) ControllerInterface
