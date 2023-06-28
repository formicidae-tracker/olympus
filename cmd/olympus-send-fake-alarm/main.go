package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/formicidae-tracker/olympus/pkg/api"
	"github.com/jessevdk/go-flags"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Options struct {
	Warning bool   `short:"w" long:"warning" description:"send a warning instead of an emergency"`
	Host    string `short:"h" long:"host" description:"olympus host to connect to" default:"localhost"`
	Port    int    `short:"p" long:"port" description:"gRPC port to use" default:"3001"`

	Args struct {
		Identification string
		Description    string
	} `positional-args:"yes" required:"yes"`
}

func main() {
	log.SetFlags(0)
	if err := execute(); err != nil {
		log.Fatalf("%s", err)
	}
}

func (o Options) BuildAlarmUpdate() *api.AlarmUpdate {
	update := &api.AlarmUpdate{
		Identification: o.Args.Identification,
		Description:    o.Args.Description,
		Status:         api.AlarmStatus_ON,
		Level:          api.AlarmLevel_EMERGENCY,
		Time:           timestamppb.New(time.Now()),
	}

	if o.Warning == true {
		update.Level = api.AlarmLevel_WARNING
	}

	return update

}

func (o Options) Dial() (*grpc.ClientConn, error) {
	return grpc.Dial(
		fmt.Sprintf("%s:%d", o.Host, o.Port),
		grpc.WithInsecure(),
	)
}

func execute() error {
	opts := Options{}
	if _, err := flags.Parse(&opts); err != nil {
		return nil
	}

	conn, err := opts.Dial()
	if err != nil {
		return fmt.Errorf("could not connect to host: %w", err)
	}
	defer conn.Close()

	update := opts.BuildAlarmUpdate()

	client := api.NewOlympusClient(conn)

	_, err = client.SendAlarm(context.Background(), update)

	if err != nil {
		return err
	}

	fmt.Printf("sent to %s:%d %s\n", opts.Host, opts.Port, update)

	return nil
}
