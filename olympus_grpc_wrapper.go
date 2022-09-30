package main

import (
	"fmt"
	"io"

	"github.com/formicidae-tracker/olympus/proto"
)

type OlympusGRPCWrapper Olympus

type serverStream[UpStream any, DownStream any] interface {
	Recv() (*UpStream, error)
	Send(*DownStream) error
}

func readAll[UpStream any, DownStream any](
	stream serverStream[UpStream, DownStream],
	messages chan<- *UpStream,
	errors chan<- error) {

	defer close(messages)
	defer close(errors)
	for {
		m, err := stream.Recv()
		if err != nil {
			errors <- err
		} else {
			messages <- m
		}
	}
}

func serveLoop[UpStream any, DownStream any](
	stream serverStream[UpStream, DownStream],
	handleMessage func(*UpStream) error,
	finish *<-chan struct{}) error {

	messages := make(chan *UpStream)
	errors := make(chan error)

	go readAll(stream, messages, errors)

	for {
		select {
		case <-(*finish):
			return nil
		case err := <-errors:
			if err == io.EOF {
				return nil
			}
			return err
		case m := <-messages:
			err := handleMessage(m)
			if err != nil {
				return err
			}
		}
	}
}

func (o *OlympusGRPCWrapper) Zone(stream proto.Olympus_ZoneServer) (err error) {
	var z ZoneLogger = nil
	var finish <-chan struct{} = nil

	defer func() {
		if z == nil {
			return
		}
		graceful := err != nil
		(*Olympus)(o).UnregisterZone(z.ZoneIdentifier(), graceful)
	}()

	handleMessage := func(m *proto.ZoneUpStream) error {
		if z == nil {
			if m.Declaration == nil {
				return fmt.Errorf("first message of stream must contain ZoneDeclaration")
			}
			var err error
			z, finish, err = (*Olympus)(o).RegisterZone(m.Declaration)
			if err != nil {
				return err
			}
		}

		if m.Target != nil {
			z.PushTarget(m.Target)
		}
		if len(m.Alarms) > 0 {
			z.PushAlarms(m.Alarms)
		}
		if len(m.Reports) > 0 {
			z.PushReports(m.Reports)
		}
		return nil
	}

	return serveLoop[proto.ZoneUpStream, proto.ZoneDownStream](stream, handleMessage, &finish)
}

func (o *OlympusGRPCWrapper) Tracking(stream proto.Olympus_TrackingServer) (err error) {
	var t TrackingLogger = nil
	var finish <-chan struct{} = nil
	hostname := ""
	defer func() {
		if t == nil {
			return
		}
		graceful := err != nil
		(*Olympus)(o).UnregisterTracker(hostname, graceful)
	}()

	handleMessage := func(m *proto.TrackingUpStream) error {
		if t == nil {
			if m.Declaration == nil {
				return fmt.Errorf("first message of stream must contain TrackingDeclaration")
			}
			var err error
			t, finish, err = (*Olympus)(o).RegisterTracker(m.Declaration)
			hostname = m.Declaration.Hostname
			if err != nil {
				return err
			}
		}
		return nil
	}

	return serveLoop[proto.TrackingUpStream, proto.TrackingDownStream](stream, handleMessage, &finish)
}
