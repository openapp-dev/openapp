package utils

import (
	"regexp"
	"testing"
)

func TestGetLocalServerIPAddress(t *testing.T) {
	ip, err := GetLocalServerIPAddress()
	if err != nil {
		t.Errorf("GetLocalServerIPAddress() failed, err: %v", err)
	}
	ipRegex := regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)
	if !ipRegex.MatchString(ip) {
		t.Errorf("GetLocalServerIPAddress() failed, ip: %s", ip)
	}
}
