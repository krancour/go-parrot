package arnetwork

// The arnetwork package implements the ARNetwork protocol and is responsible
// for a buffer abstraction that is utilized by other, higher level protocols.
//
// The buffer abstraction intelligently handles frame packing/unpacking,
// acknowledgement of frame receipt, retries, timeouts, etc. as defined by the
// ARNetwork protocol.
//
// This package is agnostic of underlying transport mechanisms; relying instead
// on the ARNetworkAL network abstraction found in the arnetworkal package.
