package grpcapi

import (
	"context"

	"github.com/araddon/dateparse"
	g "github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/api/grpcspec"
	in "github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal"
	"github.com/google/uuid"
)

func convertGrpc(ev *g.Event) in.Event {
	return in.Event{
		UUID:     ev.Uuid,
		Title:    ev.Title,
		Start:    ev.Start,
		End:      ev.End,
		OwnerID:  ev.Ownerid,
		NotifyIn: ev.Notifyin,
	}
}

func convertInternal(ev in.Event) *g.Event {
	return &g.Event{
		Uuid:     ev.UUID,
		Title:    ev.Title,
		Start:    ev.Start,
		End:      ev.End,
		Ownerid:  ev.OwnerID,
		Notifyin: ev.NotifyIn,
	}
}

func convertUpdate(upd *g.UpdateRequest) (in.Event, string) {
	return convertGrpc(upd.Event), upd.Id
}

func (s *Server) Create(ctx context.Context, grpcEv *g.Event) (*g.Id, error) {
	id, err := app.CreateEvent(convertGrpc(grpcEv))
	if err != nil {
		return nil, err
	}
	return &g.Id{
		Id: id.String(),
	}, nil
}

func (s *Server) Update(ctx context.Context, upd *g.UpdateRequest) (*g.Empty, error) {
	ev, id := convertUpdate(upd)
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return &g.Empty{}, app.UpdateEvent(ev, uuid)
}

func (s *Server) Delete(ctx context.Context, id *g.Id) (*g.Empty, error) {
	uuid, err := uuid.Parse(id.Id)
	if err != nil {
		return nil, err
	}
	return &g.Empty{}, app.DeleteEvent(uuid)
}

func (s *Server) GetForDay(ctx context.Context, d *g.Date) (*g.EventList, error) {
	date, err := dateparse.ParseAny(d.Date)
	if err != nil {
		return nil, err
	}
	evs, err := app.GetEventsForDay(date)
	if err != nil {
		return nil, err
	}
	res := g.EventList{
		Event: []*g.Event{},
	}
	for _, ev := range evs {
		res.Event = append(res.Event, convertInternal(ev))
	}
	return &res, nil
}

func (s *Server) GetForWeek(ctx context.Context, d *g.Date) (*g.EventList, error) {
	date, err := dateparse.ParseAny(d.Date)
	if err != nil {
		return nil, err
	}
	evs, err := app.GetEventsForWeek(date)
	if err != nil {
		return nil, err
	}
	res := g.EventList{
		Event: []*g.Event{},
	}
	for _, ev := range evs {
		res.Event = append(res.Event, convertInternal(ev))
	}
	return &res, nil
}

func (s *Server) GetForMonth(ctx context.Context, d *g.Date) (*g.EventList, error) {
	date, err := dateparse.ParseAny(d.Date)
	if err != nil {
		return nil, err
	}
	evs, err := app.GetEventsForMonth(date)
	if err != nil {
		return nil, err
	}
	res := g.EventList{
		Event: []*g.Event{},
	}
	for _, ev := range evs {
		res.Event = append(res.Event, convertInternal(ev))
	}
	return &res, nil
}
