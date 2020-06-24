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
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/pkg/config"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/pkg/logger"
	. "github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var (
	cl calendar.Calendar
	st TestStorage

	lis *bufconn.Listener

	date      = time.Date(2023, 3, 11, 9, 0, 0, 0, time.UTC)
	GrpcEvent = grpcspec.Event{
		Title:       "cool event",
		Start:       date.String(),
		End:         date.Add(3 * time.Hour).String(),
		Description: "test",
		Ownerid:     uuid.New().String(),
		Notifyin:    "1h",
	}
	UpdReq = grpcspec.UpdateRequest{
		Id:    uuid.New().String(),
		Event: &GrpcEvent,
	}
	DelID = grpcspec.Id{
		Id: uuid.New().String(),
	}
	GetReq = grpcspec.Date{
		Date: "2023-03-11T00:00:00",
	}
)

const bufSize = 1024 * 1024

func InitTest(cl *calendar.Calendar) {
	lis = bufconn.Listen(bufSize)
	go func() {
		if err := InitServer(cl).Grpc.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestCreate(t *testing.T) {
	_ = config.InitConfig("../../../../tests/testdata/config.json")
	_ = logger.InitLogger()
	st.Err = nil
	cl = calendar.NewCalendar(&st)
	InitTest(&cl)

	t.Run("simple", func(t *testing.T) {
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		require.Nil(t, err)
		defer conn.Close()

		client := grpcspec.NewCalendarClient(conn)
		resp, err := client.Create(ctx, &GrpcEvent)
		require.Nil(t, err)
		_, err = uuid.Parse(resp.Id)
		require.Nil(t, err)
	})

	t.Run("error in db", func(t *testing.T) {
		st = TestStorage{Err: errors.New("some db error")}
		cl = calendar.NewCalendar(&st)
		InitTest(&cl)

		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		require.Nil(t, err)
		defer conn.Close()

		client := grpcspec.NewCalendarClient(conn)
		_, err = client.Create(ctx, &GrpcEvent)
		require.NotNil(t, err)
	})
}

func TestUpdate(t *testing.T) {
	_ = config.InitConfig("../../../../tests/testdata/config.json")
	_ = logger.InitLogger()
	st.Err = nil
	cl = calendar.NewCalendar(&st)
	InitTest(&cl)

	t.Run("simple", func(t *testing.T) {
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		require.Nil(t, err)
		defer conn.Close()

		client := grpcspec.NewCalendarClient(conn)
		_, err = client.Update(ctx, &UpdReq)
		require.Nil(t, err)
	})

	t.Run("not valid uuid", func(t *testing.T) {
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		require.Nil(t, err)
		defer conn.Close()

		UpdReq.Id = "deadbeef"
		client := grpcspec.NewCalendarClient(conn)
		_, err = client.Update(ctx, &UpdReq)
		require.NotNil(t, err)
	})
}

func TestDelete(t *testing.T) {
	_ = config.InitConfig("../../../../tests/testdata/config.json")
	_ = logger.InitLogger()
	InitTest(&cl)

	t.Run("simple", func(t *testing.T) {
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		require.Nil(t, err)
		defer conn.Close()

		client := grpcspec.NewCalendarClient(conn)
		_, err = client.Delete(ctx, &DelID)
		require.Nil(t, err)
	})

	t.Run("not valid uuid", func(t *testing.T) {
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		require.Nil(t, err)
		defer conn.Close()

		DelID.Id = "deadbeef"
		client := grpcspec.NewCalendarClient(conn)
		_, err = client.Delete(ctx, &DelID)
		require.NotNil(t, err)
	})
}

func TestGet(t *testing.T) {
	_ = config.InitConfig("../../../../tests/testdata/config.json")
	_ = logger.InitLogger()
	InitTest(&cl)

	t.Run("simple", func(t *testing.T) {
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		require.Nil(t, err)
		defer conn.Close()

		client := grpcspec.NewCalendarClient(conn)
		_, err = client.GetForDay(ctx, &GetReq)
		require.Nil(t, err)
	})

	t.Run("not valid date", func(t *testing.T) {
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		require.Nil(t, err)
		defer conn.Close()

		GetReq.Date = "tomorrow"
		client := grpcspec.NewCalendarClient(conn)
		_, err = client.GetForDay(ctx, &GetReq)
		require.NotNil(t, err)
	})
}
