package bebop2

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/features/animation"
	"github.com/krancour/go-parrot/features/ardrone3"
	"github.com/krancour/go-parrot/features/common"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/protocols/arnetwork"
	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/krancour/go-parrot/protocols/arnetworkal/wifi"
	"github.com/pkg/errors"
)

// Controller ...
// TODO: Document this
type Controller interface {
	Common() common.Feature
	ARDrone3() ardrone3.Feature
	Animation() animation.Feature
}

type controller struct {
	common    common.Feature
	ardrone3  ardrone3.Feature
	animation animation.Feature
}

// NewController ...
// TODO: Document this
func NewController() (Controller, error) {
	frameSender, frameReceiver, err := wifi.Connect()
	if err != nil {
		return nil, errors.Wrap(err, "connection error")
	}
	c2dChs, d2cChs, err := arnetwork.NewBuffers(
		frameSender,
		frameReceiver,
		[]arnetwork.C2DBufferConfig{
			// Non ack data (periodic commands for piloting and camera orientation)
			// This buffer transports arcommands
			{
				ID:            10,
				FrameType:     arnetworkal.FrameTypeData,
				Size:          2, // PCMD + camera
				MaxDataSize:   128,
				IsOverwriting: true, // Periodic data; most recent is better
			},

			// Ack data (events, settings, etc.)
			// This buffer transports arcommands
			{
				ID:            11,
				FrameType:     arnetworkal.FrameTypeDataWithAck,
				AckTimeout:    150 * time.Millisecond,
				MaxRetries:    5,
				Size:          20,
				MaxDataSize:   128,
				IsOverwriting: false, // Events should not be dropped
			},

			// Emergency data (emergency commands only)
			// This buffer transports arcommands
			{
				ID:            12,
				FrameType:     arnetworkal.FrameTypeDataWithAck,
				AckTimeout:    150 * time.Millisecond,
				MaxRetries:    -1, // Infinite
				Size:          1,
				MaxDataSize:   128,
				IsOverwriting: false, // Events should not be dropped
			},

			// // TODO: Do something about video streaming?
			// // arstream video acks
			// // This buffer transports arstream data
			// {
			// 	ID:            13,
			// 	FrameType:     arnetworkal.FrameTypeLowLatencyData,
			// 	Size:          1000, // Enough space
			// 	MaxDataSize:   18,   // Size of an ack
			// 	IsOverwriting: true, // New is always better
			// },
		},

		[]arnetwork.D2CBufferConfig{
			// Non ack data (periodic reports from the device)
			// This buffer transports arcommands
			{
				ID:            127,
				FrameType:     arnetworkal.FrameTypeData,
				Size:          20,
				MaxDataSize:   128,
				IsOverwriting: true, // Periodic data: most recent is better
			},

			// Ack data (events, settings, etc.)
			// This buffer transports arcommands
			{
				ID:            126,
				FrameType:     arnetworkal.FrameTypeDataWithAck,
				Size:          256,
				MaxDataSize:   128,
				IsOverwriting: false, // Events should not be dropped
			},

			// // TODO: Do something about video streaming?
			// // arstream video data
			// // This buffer transports arstream data
			// {
			// 	ID:        125,
			// 	FrameType: arnetworkal.FrameTypeLowLatencyData,
			// 	// TODO: According to documentation, size should be set to
			// 	// "arstream_fragment_maximum_number * 2"
			// 	// I think this is supposed to be determined during connection
			// 	// negotiation???
			// 	Size: 1000, // This value is a placeholder!
			// 	// TODO: According to documentation, this should be set to
			// 	// "arstream_fragment_size"
			// 	// I think this is supposed to be determined during connection
			// 	// negotiation???
			// 	MaxDataSize:   256,  // This value is a placeholder!
			// 	IsOverwriting: true, // New is always better
			// },
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "error creating buffer manager")
	}
	c2dCommandClient := arcommands.NewC2DCommandClient(
		c2dChs[11], // ACK
		c2dChs[10], // NON_ACK
		c2dChs[12], // HIGH_PRIO
	)
	commonFeature := common.NewFeature(c2dCommandClient)
	ardrone3Feature := ardrone3.NewFeature(c2dCommandClient)
	d2cCommandServer, err := arcommands.NewD2CCommandServer(
		d2cChs,
		[]arcommands.D2CFeature{
			commonFeature,
			ardrone3Feature,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "error creating d2c command server")
	}
	d2cCommandServer.Start()
	if err := commonFeature.Settings().AllSettings(); err != nil {
		log.Error(err)
	}
	if err := commonFeature.Common().AllStates(); err != nil {
		log.Error(err)
	}
	// TODO: Should we wait until all states are received before returning?
	return &controller{
		common:   commonFeature,
		ardrone3: ardrone3Feature,
	}, nil
}

func (c *controller) Common() common.Feature {
	return c.common
}

func (c *controller) ARDrone3() ardrone3.Feature {
	return c.ardrone3
}

func (c *controller) Animation() animation.Feature {
	return c.animation
}
