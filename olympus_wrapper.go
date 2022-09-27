package main

import (
	"fmt"
	"io"

	"github.com/formicidae-tracker/olympus/proto"
)

type OlympusGRPCWrapper Olympus

func readAll(stream proto.Olympus_ZoneServer, messages chan<- *proto.ZoneUpStream, errors chan<- error) {
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

func (o *OlympusGRPCWrapper) Zone(stream proto.Olympus_ZoneServer) error {
	var z ZoneLogger = nil
	messages := make(chan *proto.ZoneUpStream)
	errors := make(chan error)
	var finish <-chan struct{} = nil
	defer func() {
		if z == nil {
			return
		}
		(*Olympus)(o).UnregisterZone(z.ZoneIdentifier())
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

	var err error
	go readAll(stream, messages, errors)
	for {
		select {
		case <-finish:
			return nil
		case err := <-errors:
			if err == io.EOF {
				return nil
			}
			return err
		case m := <-messages:
			err = handleMessage(m)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
