package ardrone3

import log "github.com/Sirupsen/logrus"

// TODO: Implement this
func (f *feature) positionChanged(args []interface{}) error {
	log.Debugf("the position changed: %v", args)
	return nil
}

// TODO: Implement this
func (f *feature) speedChanged(args []interface{}) error {
	log.Debugf("the speed changed: %v", args)
	return nil
}

// TODO: Implement this
func (f *feature) attitudeChanged(args []interface{}) error {
	log.Debugf("the attitude changed: %v", args)
	return nil
}

// TODO: Implement this
func (f *feature) altitudeChanged(args []interface{}) error {
	log.Debugf("the altitude changed: %v", args)
	return nil
}

// TODO: Implement this
func (f *feature) gpsLocationChanged(args []interface{}) error {
	log.Debugf("the gps location changed: %v", args)
	return nil
}
