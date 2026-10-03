package ticket

import (
	"context"
	"database/sql"
	dto "konsera-backend/internal/DTO/ticket"
	"konsera-backend/internal/models"
	"time"

	"github.com/google/uuid"
)

func scanAccessTicket(s interface{ Scan(...any) error }) (*models.Ticket, error) {
	x := &models.Ticket{}
	e := s.Scan(&x.ID, &x.BookingItemID, &x.OwnerUserID, &x.TicketCode, &x.QRToken, &x.QRTokenUpdatedAt, &x.Status, &x.WatermarkUserHash, &x.IssuedAt)
	return x, e
}

const accessTicketCols = `id,booking_item_id,owner_user_id,ticket_code,qr_token,qr_token_updated_at,status,watermark_user_hash,issued_at`

func (r *TicketRepository) Tickets(ctx context.Context, user uuid.UUID) ([]*models.Ticket, error) {
	rows, e := r.db.QueryContext(ctx, `SELECT `+accessTicketCols+` FROM tickets WHERE owner_user_id=$1 ORDER BY issued_at DESC`, user)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []*models.Ticket{}
	for rows.Next() {
		x, e := scanAccessTicket(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *TicketRepository) Ticket(ctx context.Context, user, id uuid.UUID) (*models.Ticket, error) {
	return scanAccessTicket(r.db.QueryRowContext(ctx, `SELECT `+accessTicketCols+` FROM tickets WHERE owner_user_id=$1 AND id=$2`, user, id))
}
func (r *TicketRepository) UpdateTicket(ctx context.Context, user, id uuid.UUID, q *dto.UpdateTicketRequest) (*models.Ticket, error) {
	return scanAccessTicket(r.db.QueryRowContext(ctx, `UPDATE tickets SET status=COALESCE($3,status),qr_token_updated_at=NOW() WHERE owner_user_id=$1 AND id=$2 RETURNING `+accessTicketCols, user, id, q.Status))
}
func scanTransfer(s interface{ Scan(...any) error }) (*models.TicketTransfer, error) {
	x := &models.TicketTransfer{}
	e := s.Scan(&x.ID, &x.TicketID, &x.FromUserID, &x.ToUserID, &x.Status, &x.RequestedAt, &x.ResolvedAt)
	return x, e
}

const transferCols = `id,ticket_id,from_user_id,to_user_id,status,requested_at,resolved_at`

func (r *TicketRepository) CreateTransfer(ctx context.Context, x *models.TicketTransfer) error {
	created, err := scanTransfer(r.db.QueryRowContext(ctx, `
		INSERT INTO ticket_transfers(ticket_id,from_user_id,to_user_id)
		SELECT id,owner_user_id,$2
		FROM tickets
		WHERE id=$1 AND owner_user_id=$3 AND status IN ('issued','transferred')
		RETURNING `+transferCols, x.TicketID, x.ToUserID, x.FromUserID))
	if err != nil {
		return err
	}
	*x = *created
	return nil
}
func (r *TicketRepository) Transfers(ctx context.Context, user uuid.UUID) ([]*models.TicketTransfer, error) {
	rows, e := r.db.QueryContext(ctx, `SELECT `+transferCols+` FROM ticket_transfers WHERE from_user_id=$1 OR to_user_id=$1 ORDER BY requested_at DESC`, user)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []*models.TicketTransfer{}
	for rows.Next() {
		x, e := scanTransfer(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *TicketRepository) Transfer(ctx context.Context, user, id uuid.UUID) (*models.TicketTransfer, error) {
	return scanTransfer(r.db.QueryRowContext(ctx, `SELECT `+transferCols+` FROM ticket_transfers WHERE id=$1 AND (from_user_id=$2 OR to_user_id=$2)`, id, user))
}
func (r *TicketRepository) ResolveTransfer(ctx context.Context, user, id uuid.UUID, status string) (*models.TicketTransfer, error) {
	tx, e := r.db.BeginTx(ctx, nil)
	if e != nil {
		return nil, e
	}
	defer tx.Rollback()
	x, e := scanTransfer(tx.QueryRowContext(ctx, `SELECT `+transferCols+` FROM ticket_transfers WHERE id=$1 AND (to_user_id=$2 OR from_user_id=$2) FOR UPDATE`, id, user))
	if e != nil {
		return nil, e
	}
	if status == "accepted" {
		if x.ToUserID != user {
			return nil, sql.ErrNoRows
		}
		if _, e = tx.ExecContext(ctx, `UPDATE tickets SET owner_user_id=$2,status='transferred' WHERE id=$1 AND owner_user_id=$3`, x.TicketID, x.ToUserID, x.FromUserID); e != nil {
			return nil, e
		}
	}
	if _, e = tx.ExecContext(ctx, `UPDATE ticket_transfers SET status=$2,resolved_at=NOW() WHERE id=$1 AND status='pending'`, id, status); e != nil {
		return nil, e
	}
	x.Status = status
	now := time.Now()
	x.ResolvedAt = &now
	if e = tx.Commit(); e != nil {
		return nil, e
	}
	return x, nil
}
func scanCheckIn(s interface{ Scan(...any) error }) (*models.CheckIn, error) {
	x := &models.CheckIn{}
	e := s.Scan(&x.ID, &x.TicketID, &x.EventSessionID, &x.ScannedBy, &x.Result, &x.DeviceID, &x.ScannedAt, &x.IsOfflineSync)
	return x, e
}

const checkInCols = `id,ticket_id,event_session_id,scanned_by,result,device_id,scanned_at,is_offline_sync`

func (r *TicketRepository) CreateCheckIn(ctx context.Context, x *models.CheckIn) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var status string
	var ticketSession uuid.UUID
	err = tx.QueryRowContext(ctx, `
		SELECT t.status, bi.event_session_id
		FROM tickets t
		JOIN booking_items bi ON bi.id=t.booking_item_id
		WHERE t.id=$1
		FOR UPDATE
	`, x.TicketID).Scan(&status, &ticketSession)
	if err != nil {
		return err
	}

	var gateOpen, gateClose time.Time
	err = tx.QueryRowContext(ctx, `SELECT gate_open_at,gate_close_at FROM event_sessions WHERE id=$1`, x.EventSessionID).Scan(&gateOpen, &gateClose)
	if err != nil {
		return err
	}

	result := "success"
	now := time.Now()
	if ticketSession != x.EventSessionID {
		result = "invalid"
	} else if status == "checked_in" {
		result = "duplicate"
	} else if (status != "issued" && status != "transferred") || now.Before(gateOpen) {
		result = "not_yet_open"
	} else if now.After(gateClose) {
		result = "closed"
	} else if _, err = tx.ExecContext(ctx, `UPDATE tickets SET status='checked_in' WHERE id=$1 AND status='issued'`, x.TicketID); err != nil {
		return err
	}

	created, err := scanCheckIn(tx.QueryRowContext(ctx, `INSERT INTO check_ins(ticket_id,event_session_id,scanned_by,result,device_id,scanned_at,is_offline_sync) VALUES($1,$2,$3,$4,$5,COALESCE($6,NOW()),$7) RETURNING `+checkInCols, x.TicketID, x.EventSessionID, x.ScannedBy, result, x.DeviceID, x.ScannedAt, x.IsOfflineSync))
	if err != nil {
		return err
	}
	*x = *created
	return tx.Commit()
}
func (r *TicketRepository) CheckIns(ctx context.Context, ticketID *uuid.UUID) ([]*models.CheckIn, error) {
	rows, e := r.db.QueryContext(ctx, `SELECT `+checkInCols+` FROM check_ins WHERE ($1::uuid IS NULL OR ticket_id=$1) ORDER BY scanned_at DESC`, ticketID)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []*models.CheckIn{}
	for rows.Next() {
		x, e := scanCheckIn(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
