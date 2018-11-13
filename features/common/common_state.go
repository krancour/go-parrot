package common

import log "github.com/Sirupsen/logrus"

// TODO: Implement this
func (f *feature) batteryStateChanged(args []interface{}) error {
	log.Debugf("the battery state changed: %v", args)
	return nil
}

// TODO: Implement this
func (f *feature) wifiSignalChanged(args []interface{}) error {
	log.Debugf("the wifi signal changed: %v", args)
	return nil
}
