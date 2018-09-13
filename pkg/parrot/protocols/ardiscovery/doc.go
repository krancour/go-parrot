package ardiscovery

// The ardiscovery package partially implements the ARDiscovery protocol, which
// is responsible for both device discovery and connection negotiation. This
// partial implementation only provides an abstraction for connection
// negotiation.
//
// Two implementations of this abstraction are possible, with one being UDP/IP
// based (implemented) and the other being BLE (Bluetooth) based (currently not
// implemented).
//
// Device discovery is not implemented due to the following assumptions:
//
//   1. Users of a Wi-Fi-enabled device have pre-connected to the WLAN provided
//   by their device.
//
//   2. Wi-Fi-enabled devices are available at a well-known IP address.
//
// To help maintain traceability to relevant sections of the Parrot developer
// documentation, the UDP/IP based implementation of this abstraction will be
// abbreviated as "wifi".
//
// // TODO: Create BLE based implementation.
