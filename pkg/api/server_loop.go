package api

import (
	"context"
	"io"
)

type HandleFunc[Up, Down any] func(*Up) (*Down, error)

type ServerStream[Up, Down any] interface {
	Recv() (*Up, error)
	Send(*Down) error
	Context() context.Context
}

type recvResult[Up any] struct {
	Message *Up
	Error   error
}

func readAll[Up, Down any](ch chan<- recvResult[Up],
	s ServerStream[Up, Down]) {

	defer close(ch)

	for {
		m, err := s.Recv()
		select {
		case ch <- recvResult[Up]{m, err}:
		case <-s.Context().Done():
		}

		if err != nil {
			return
		}
	}
}

func ServerLoop[Up, Down any](
	ctx context.Context,
	s ServerStream[Up, Down],
	handler HandleFunc[Up, Down]) error {

	recv := make(chan recvResult[Up])
	go func() {
		readAll[Up, Down](recv, s)
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-s.Context().Done():
			return nil
		case m, ok := <-recv:
			if ok == false {
				return nil
			}
			if m.Error != nil {
				if m.Error == io.EOF || m.Error == context.Canceled {
					return nil
				}
				return m.Error
			}
			resp, err := handler(m.Message)
			if err != nil {
				return err
			}
			err = s.Send(resp)
			if err != nil {
				return err
			}
		}
	}
}
