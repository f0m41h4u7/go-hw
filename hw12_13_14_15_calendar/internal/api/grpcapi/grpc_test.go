package grpcapi

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"
	"time"

	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/api/grpcspec"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/logger"
	. "github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/tests"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var (
	cl        calendar.Calendar
	st        TestStorage
	grpcEvent grpcspec.Event
	lis       *bufconn.Listener

	date = time.Date(2023, 3, 11, 9, 0, 0, 0, time.UTC)
)

const bufSize = 1024 * 1024

func initTest(cl *calendar.Calendar) {
	lis = bufconn.Listen(bufSize)
	go func() {
		if err := InitServer(cl).grpc.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestCreate(t *testing.T) {
	_ = config.InitCalendarConfig("../../../../tests/testdata/config.json")
	_ = logger.InitLogger()
	st.Err = nil
	cl = calendar.NewCalendar(&st)
	initTest(&cl)

	start, err := ptypes.TimestampProto(date)
	require.Nil(t, err)
	end, err := ptypes.TimestampProto(date.Add(3 * time.Hour))
	require.Nil(t, err)
	grpcEvent = grpcspec.Event{
		Title:       "cool event",
		Start:       start,
		End:         end,
		Description: "test",
		Ownerid:     uuid.New().String(),
		Notifyin:    time.Hour.Milliseconds(),
	}

	crReq := grpcspec.CreateRequest{
		Event: &grpcEvent,
	}

	t.Run("simple", func(t *testing.T) {
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		require.Nil(t, err)
		defer conn.Close()

		client := grpcspec.NewCalendarClient(conn)
		resp, err := client.Create(ctx, &crReq)
		require.Nil(t, err)
		_, err = uuid.Parse(resp.Uuid)
		require.Nil(t, err)
	})

	t.Run("error in db", func(t *testing.T) {
		st = TestStorage{Err: errors.New("some db error")}
		cl = calendar.NewCalendar(&st)
		initTest(&cl)

		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		require.Nil(t, err)
		defer conn.Close()

		client := grpcspec.NewCalendarClient(conn)
		_, err = client.Create(ctx, &crReq)
		require.NotNil(t, err)
	})
}

func TestUpdate(t *testing.T) {
	_ = config.InitCalendarConfig("../../../../tests/testdata/config.json")
	_ = logger.InitLogger()
	st.Err = nil
	cl = calendar.NewCalendar(&st)
	initTest(&cl)
	updReq := grpcspec.UpdateRequest{
		Uuid:  uuid.New().String(),
		Event: &grpcEvent,
	}

	t.Run("simple", func(t *testing.T) {
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		require.Nil(t, err)
		defer conn.Close()

		client := grpcspec.NewCalendarClient(conn)
		_, err = client.Update(ctx, &updReq)
		require.Nil(t, err)
	})

	t.Run("not valid uuid", func(t *testing.T) {
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		require.Nil(t, err)
		defer conn.Close()

		updReq.Uuid = "deadbeef"
		client := grpcspec.NewCalendarClient(conn)
		_, err = client.Update(ctx, &updReq)
		require.NotNil(t, err)
	})
}

func TestDelete(t *testing.T) {
	_ = config.InitCalendarConfig("../../../../tests/testdata/config.json")
	_ = logger.InitLogger()
	initTest(&cl)
	delReq := grpcspec.DeleteRequest{
		Uuid: uuid.New().String(),
	}

	t.Run("simple", func(t *testing.T) {
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		require.Nil(t, err)
		defer conn.Close()

		client := grpcspec.NewCalendarClient(conn)
		_, err = client.Delete(ctx, &delReq)
		require.Nil(t, err)
	})

	t.Run("not valid uuid", func(t *testing.T) {
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		require.Nil(t, err)
		defer conn.Close()

		delReq.Uuid = "deadbeef"
		client := grpcspec.NewCalendarClient(conn)
		_, err = client.Delete(ctx, &delReq)
		require.NotNil(t, err)
	})
}

func TestGet(t *testing.T) {
	_ = config.InitCalendarConfig("../../../../tests/testdata/config.json")
	_ = logger.InitLogger()
	initTest(&cl)
	dt, err := ptypes.TimestampProto(date)
	require.Nil(t, err)
	getReq := grpcspec.GetRequest{
		Date: dt,
	}

	t.Run("simple", func(t *testing.T) {
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		require.Nil(t, err)
		defer conn.Close()

		client := grpcspec.NewCalendarClient(conn)
		_, err = client.GetForDay(ctx, &getReq)
		require.Nil(t, err)
	})
}
