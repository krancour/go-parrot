package ardrone3

import log "github.com/Sirupsen/logrus"

// TODO: Implement this
func (f *feature) numberOfSatelliteChanged(args []interface{}) error {
	log.Debugf("the number of satellites changed: %v", args)
	return nil
}
