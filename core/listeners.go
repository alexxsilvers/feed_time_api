package core

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
)

type listeners struct {
	http net.Listener
}

func newListeners(publicPort uint) (*listeners, error) {
	l := &listeners{}

	http, err := net.Listen("tcp", fmt.Sprintf(":%v", publicPort))
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create HTTP listener")
	}
	l.http = http

	return l, nil
}
