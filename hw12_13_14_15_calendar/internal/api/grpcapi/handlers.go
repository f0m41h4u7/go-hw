package grpcapi

import (
	"context"
	"strconv"

	"github.com/araddon/dateparse"
	g "github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/api/grpcspec"
	in "github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
)

func convertGrpc(ev *g.Event) in.Event {
	return in.Event{
		UUID:     ev.Uuid,
		Title:    ev.Title,
		Start:    ev.Start.String(),
		End:      ev.End.String(),
		OwnerID:  ev.Ownerid,
		NotifyIn: string(ev.Notifyin),
	}
}

func convertInternal(ev in.Event) (*g.Event, error) {
	stTime, err := dateparse.ParseAny(ev.Start)
	if err != nil {
		return nil, err
	}
	endTime, err := dateparse.ParseAny(ev.End)
	if err != nil {
		return nil, err
	}
	notif, err := strconv.Atoi(ev.NotifyIn)
	if err != nil {
		return nil, err
	}

	st, err := ptypes.TimestampProto(stTime)
	if err != nil {
		return nil, err
	}
	end, err := ptypes.TimestampProto(endTime)
	if err != nil {
		return nil, err
	}
	return &g.Event{
		Uuid:     ev.UUID,
		Title:    ev.Title,
		Start:    st,
		End:      end,
		Ownerid:  ev.OwnerID,
		Notifyin: int64(notif),
	}, nil
}

func convertUpdate(upd *g.UpdateRequest) (in.Event, string) {
	return convertGrpc(upd.Event), upd.Uuid
}

func (s *Server) Create(ctx context.Context, req *g.CreateRequest) (*g.CreateResponse, error) {
	id, err := app.CreateEvent(convertGrpc(req.Event))
	if err != nil {
		return nil, err
	}
	return &g.CreateResponse{
		Uuid: id.String(),
	}, nil
}

func (s *Server) Update(ctx context.Context, req *g.UpdateRequest) (*g.UpdateResponse, error) {
	ev, id := convertUpdate(req)
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return &g.UpdateResponse{}, app.UpdateEvent(ev, uuid)
}

func (s *Server) Delete(ctx context.Context, req *g.DeleteRequest) (*g.DeleteResponse, error) {
	uuid, err := uuid.Parse(req.Uuid)
	if err != nil {
		return nil, err
	}
	return &g.DeleteResponse{}, app.DeleteEvent(uuid)
}

func (s *Server) GetForDay(ctx context.Context, req *g.GetRequest) (*g.GetResponse, error) {
	date, err := ptypes.Timestamp(req.Date)
	if err != nil {
		return nil, err
	}
	evs, err := app.GetEventsForDay(date)
	if err != nil {
		return nil, err
	}
	res := g.GetResponse{
		Event: make([]*g.Event, len(evs)),
	}
	for _, ev := range evs {
		tmp, err := convertInternal(ev)
		if err != nil {
			return nil, err
		}
		res.Event = append(res.Event, tmp)
	}
	return &res, nil
}

func (s *Server) GetForWeek(ctx context.Context, d *g.GetRequest) (*g.GetResponse, error) {
	date, err := ptypes.Timestamp(d.Date)
	if err != nil {
		return nil, err
	}
	evs, err := app.GetEventsForWeek(date)
	if err != nil {
		return nil, err
	}
	res := g.GetResponse{
		Event: make([]*g.Event, len(evs)),
	}
	for _, ev := range evs {
		tmp, err := convertInternal(ev)
		if err != nil {
			return nil, err
		}
		res.Event = append(res.Event, tmp)
	}
	return &res, nil
}

func (s *Server) GetForMonth(ctx context.Context, d *g.GetRequest) (*g.GetResponse, error) {
	date, err := ptypes.Timestamp(d.Date)
	if err != nil {
		return nil, err
	}
	evs, err := app.GetEventsForMonth(date)
	if err != nil {
		return nil, err
	}
	res := g.GetResponse{
		Event: make([]*g.Event, len(evs)),
	}
	for _, ev := range evs {
		tmp, err := convertInternal(ev)
		if err != nil {
			return nil, err
		}
		res.Event = append(res.Event, tmp)
	}
	return &res, nil
}
