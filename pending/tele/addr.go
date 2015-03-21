// Copyright 2013 Tumblr, Inc.
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package tele

import (
	"encoding/gob"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/gocircuit/alef/errors"
	"github.com/gocircuit/alef/peer"
)

// Addr maintains a single unique instance for each addr.
// Addr object uniqueness is required by the n.Addr interface.
type Addr struct {
	ID  peer.Id
	PID int
	TCP *net.TCPAddr
}

func init() {
	gob.Register(&Addr{})
}

func MustParseNetAddr(s string) net.Addr {
	addr, err := ParseNetAddr(s)
	if err != nil {
		panic(err)
	}
	return addr
}

func ParseNetAddr(s string) (net.Addr, error) {
	return net.ResolveTCPAddr("tcp", s)
}

func NewNetAddr(id peer.Id, pid int, addr net.Addr) *Addr {
	return &Addr{ID: id, PID: pid, TCP: addr.(*net.TCPAddr)}
}

func NewAddr(id peer.Id, pid int, hostport string) (peer.Addr, error) {
	a, err := net.ResolveTCPAddr("tcp", hostport)
	if err != nil {
		return nil, err
	}
	return &Addr{ID: id, PID: pid, TCP: a}, nil
}

func (a *Addr) Id() peer.Id {
	return a.ID
}

func (a *Addr) String() string {
	u := url.URL{
		Scheme: peer.Scheme,
		Host:   sanitizeTCP(a.TCP),
		Path:   "/" + strconv.Itoa(a.PID) + "/" + a.ID.String(),
	}
	return u.String()
}

// circuit://123.3.45.0:3456/2345/R1122334455667788
func ParseAddr(s string) (*Addr, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	if u.Scheme != peer.Scheme {
		return nil, errors.NewError("worker address URL scheme mismatch")
	}
	// Net address
	naddr, err := ParseNetAddr(u.Host)
	if err != nil {
		return nil, err
	}
	// Parse path
	parts := strings.Split(u.Path, "/")
	if len(parts) != 3 {
		return nil, errors.NewError(fmt.Sprintf("parse path: %#v", parts))
	}
	if parts[0] != "" {
		return nil, errors.NewError("must start with slash")
	}
	// PID
	pid, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}
	// Worker ID
	id, err := peer.ParseId(parts[2])
	if err != nil {
		return nil, err
	}
	return &Addr{ID: id, PID: pid, TCP: naddr.(*net.TCPAddr)}, nil
}

func sanitizeTCP(a *net.TCPAddr) string {
	if len(a.IP) == 0 {
		return "noaddr"
	}
	return a.String()
}
