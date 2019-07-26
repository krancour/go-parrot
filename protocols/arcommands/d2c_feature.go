package arcommands

// D2CFeature ...
// TODO: Document this
type D2CFeature interface {
	FeatureID() uint8
	FeatureName() string
	D2CClasses() []D2CClass
}
