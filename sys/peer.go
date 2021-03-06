// Copyright 2015 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2015 Petar Maymounkov <p@gocircuit.org>

package sys

const Scheme = "circuit"

// Addr represents a locator for a unique remote worker.
// The implementing type must be registered with package encoding/gob.
type Addr interface {

	// String returns an equivalent textual representation of the address.
	String() string

	// Id returns a unique identifier for the underlying worker.
	Id() Id
}

// Conn is a connection to a remote endpoint, identified by an Addr.
type Conn interface {

	// The language runtime does not itself utilize timeouts on send/receive
	// operations. Instead, it requires that calls to Send and Receive be blocking
	// until success or irrecoverable failure is reached.
	//
	// The implementation of Conn must monitor the liveness of the remote
	// endpoint using an out-of-band (non-visible to the runtime) method. If
	// the endpoint is considered dead, all pending Send and Receive request must
	// return with non-nil error.
	//
	// A non-nil error returned on any invokation of Send and Receive signals to
	// the runtime that the connection is dead.
	//
	// Such an event triggers various language runtime actions such as, for
	// example, releasing all values exported to that runtime. Therefore, a
	// typical Conn implementation might choose to attempt various physical
	// connectivity recovery methods, before it reports an error on any pending
	// connection. Such implentation strategies are facilitated by the fact
	// that the runtime has no semantic limits on the length of blocking waits.
	// In fact, the runtime has no notion of time altogether.

	// Send/Receive operations must panic on any encoding/decoding errors.
	// Whereas they must return an error for any exernal (network) unexpected
	// conditions. Encoding errors indicate compile-time errors (which would be
	// caught automatically once the system has its own compiler) but might be
	// missed by the bare Go compiler.
	//
	// Send/Receive must be re-entrant.

	// Receive reads the next value from the connection.
	Receive() (interface{}, error)

	// Send writes the given value to the connection.
	Send(interface{}) error

	// Close closes the connection.
	Close() error

	// Addr returns the address of the remote endpoint.
	Addr() Addr
}

// Listener is a device for accepting incoming connections.
type Listener interface {

	// Accept returns the next incoming connection.
	Accept() (Conn, error)

	// Addr returns the address of this endpoint.
	Addr() Addr
}

// Dialer is a device for initating connections to addressed remote endpoints.
type Dialer interface {

	// Dial connects to the endpoint specified by addr and returns a respective connection object.
	Dial(addr Addr) (Conn, error)
}

// Peer cumulatively represents the ability to listen for connections and dial into remote endpoints.
type Peer interface {
	Dialer
	Listener
}
