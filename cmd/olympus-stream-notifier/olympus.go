package main

import (
	"fmt"
	"net/rpc"
)

func CallRPC(method string, args interface{}) error {
	c, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", opts.RPCAddress, opts.RPCPort))
	if err != nil {
		return nil
	}
	unused := 0
	return c.Call(method, args, &unused)
}

type LetoTrackingRegister struct {
	Host, URL string
}

func RegisterOlympus(hostname, URL string) error {
	return CallRPC("Olympus.RegisterTracker", LetoTrackingRegister{Host: hostname, URL: URL})
}

func UnregisterOlympus(hostname string) error {
	return CallRPC("Olympus.UnregisterTracker", hostname)
}
