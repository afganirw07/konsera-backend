package ticket

import (
	"context"
	"fmt"
	dto "konsera-backend/internal/DTO/ticket"
	"konsera-backend/internal/models"
	"time"

	"github.com/google/uuid"
)

func (s *TicketService) Tickets(c context.Context, u uuid.UUID) ([]*models.Ticket, error) {
	return s.repo.Tickets(c, u)
}
func (s *TicketService) TicketByOwner(c context.Context, u, id uuid.UUID) (*models.Ticket, error) {
	return s.repo.Ticket(c, u, id)
}
func (s *TicketService) UpdateOwned(c context.Context, u, id uuid.UUID, q *dto.UpdateTicketRequest) (*models.Ticket, error) {
	return s.repo.UpdateTicket(c, u, id, q)
}
func (s *TicketService) saveTransfer(c context.Context, x *models.TicketTransfer) error {
	return s.repo.CreateTransfer(c, x)
}
func (s *TicketService) Transfers(c context.Context, u uuid.UUID) ([]*models.TicketTransfer, error) {
	return s.repo.Transfers(c, u)
}
func (s *TicketService) ResolveTransferByUser(c context.Context, u, id uuid.UUID, status string) (*models.TicketTransfer, error) {
	return s.repo.ResolveTransfer(c, u, id, status)
}
func (s *TicketService) saveCheckIn(c context.Context, x *models.CheckIn) error {
	return s.repo.CreateCheckIn(c, x)
}
func (s *TicketService) ListCheckIns(c context.Context, id *uuid.UUID) ([]*models.CheckIn, error) {
	return s.repo.CheckIns(c, id)
}
func (s *TicketService) MyTickets(c context.Context, u string) ([]*models.Ticket, error) {
	id, e := parseID(u, "user_id")
	if e != nil {
		return nil, e
	}
	return s.Tickets(c, id)
}
func (s *TicketService) Ticket(c context.Context, u, v string) (*models.Ticket, error) {
	a, e := parseID(u, "user_id")
	if e != nil {
		return nil, e
	}
	b, e := parseID(v, "ticket_id")
	if e != nil {
		return nil, e
	}
	return s.TicketByOwner(c, a, b)
}
func (s *TicketService) UpdateOwnedTicket(c context.Context, u, v string, q *dto.UpdateTicketRequest) (*models.Ticket, error) {
	a, e := parseID(u, "user_id")
	if e != nil {
		return nil, e
	}
	b, e := parseID(v, "ticket_id")
	if e != nil {
		return nil, e
	}
	return s.UpdateOwned(c, a, b, q)
}
func (s *TicketService) CreateTransfer(c context.Context, u, v string, q *dto.CreateTransferRequest) (*models.TicketTransfer, error) {
	from, e := parseID(u, "user_id")
	if e != nil {
		return nil, e
	}
	ticket, e := parseID(v, "ticket_id")
	if e != nil {
		return nil, e
	}
	to, e := parseID(q.ToUserID, "to_user_id")
	if e != nil {
		return nil, e
	}
	if from == to {
		return nil, fmt.Errorf("cannot transfer ticket to yourself")
	}
	x := &models.TicketTransfer{TicketID: ticket, FromUserID: from, ToUserID: to, Status: "pending", RequestedAt: time.Now()}
	e = s.saveTransfer(c, x)
	return x, e
}
func (s *TicketService) MyTransfers(c context.Context, u string) ([]*models.TicketTransfer, error) {
	id, e := parseID(u, "user_id")
	if e != nil {
		return nil, e
	}
	return s.Transfers(c, id)
}
func (s *TicketService) ResolveTransfer(c context.Context, u, v string, q *dto.UpdateTransferRequest) (*models.TicketTransfer, error) {
	user, e := parseID(u, "user_id")
	if e != nil {
		return nil, e
	}
	id, e := parseID(v, "transfer_id")
	if e != nil {
		return nil, e
	}
	return s.ResolveTransferByUser(c, user, id, q.Status)
}
func (s *TicketService) CreateCheckIn(c context.Context, u string, q *dto.CreateCheckInRequest) (*models.CheckIn, error) {
	scanner, e := parseID(u, "user_id")
	if e != nil {
		return nil, e
	}
	ticket, e := parseID(q.TicketID, "ticket_id")
	if e != nil {
		return nil, e
	}
	session, e := parseID(q.EventSessionID, "event_session_id")
	if e != nil {
		return nil, e
	}
	x := &models.CheckIn{TicketID: ticket, EventSessionID: session, ScannedBy: scanner, Result: q.Result, DeviceID: q.DeviceID, IsOfflineSync: q.IsOfflineSync, ScannedAt: time.Now()}
	if q.ScannedAt != nil {
		x.ScannedAt = *q.ScannedAt
	}
	e = s.saveCheckIn(c, x)
	return x, e
}
func (s *TicketService) CheckIns(c context.Context, v string) ([]*models.CheckIn, error) {
	var id *uuid.UUID
	if v != "" {
		x, e := parseID(v, "ticket_id")
		if e != nil {
			return nil, e
		}
		id = &x
	}
	return s.ListCheckIns(c, id)
}
