package bebop2

import (
	"time"

	"github.com/krancour/go-parrot/features/ardrone3"
	"github.com/krancour/go-parrot/features/common"
	"github.com/krancour/go-parrot/protocols/arnetwork"
	"github.com/krancour/go-parrot/protocols/arnetworkal"
	"github.com/krancour/go-parrot/protocols/arnetworkal/wifi"
	"github.com/pkg/errors"
)

// Controller ...
// TODO: Document this
type Controller interface {
	// TODO: Add all supported commands to this interface
	// The first ones we experiment with should be configuration related and
	// not anything to do with flight. We want to be pretty sure that everything
	// works before we try flying!
}

type controller struct {
	common   common.Feature
	ardrone3 ardrone3.Feature
}

// NewController ...
// TODO: Document this
func NewController() (Controller, error) {
	frameSender, frameReceiver, err := wifi.Connect()
	if err != nil {
		return nil, errors.Wrap(err, "connection error")
	}
	bufMan, err := arnetwork.NewBufferManager(
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
	return &controller{
		common:   common.NewFeature(bufMan),
		ardrone3: ardrone3.NewFeature(bufMan),
	}, nil
}
