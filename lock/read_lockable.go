package lock

type ReadLockable interface {
	// RLock blocks until a read lock is obtained. This permits callers to procede
	// with querying any or all related attributes of some stateful object without
	// worry that some attributes will be overwritten as others are read. i.e. It
	// permits the possibility of taking an atomic snapshop of a stateful object.
	// Note that use of this function is not obligatory for applications that do
	// not require such guarantees. Callers MUST call RUnlock() or else the
	// stateful object will never resume updating.
	RLock()
	// RUnlock releases a read lock on the state. See RLock().
	RUnlock()
}
