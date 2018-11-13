package ardrone3

import log "github.com/Sirupsen/logrus"

// TODO: Implement this
func (f *feature) orientation(args []interface{}) error {
	log.Debugf("the camera orientation changed: %v", args)
	return nil
}
