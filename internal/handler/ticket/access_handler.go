package ticket

import (
	"database/sql"
	"errors"
	dto "konsera-backend/internal/DTO/ticket"
	"konsera-backend/internal/helpers"

	"github.com/gin-gonic/gin"
)

func accessFail(c *gin.Context, e error) {
	st := 400
	if errors.Is(e, sql.ErrNoRows) {
		st = 404
	}
	helpers.Error(c, st, e.Error(), nil)
}
func accessUser(c *gin.Context) string { v, _ := c.Get("user_id"); return v.(string) }
func (h *TicketHandler) ListOwned(c *gin.Context) {
	x, e := h.s.MyTickets(c, accessUser(c))
	if e != nil {
		accessFail(c, e)
		return
	}
	helpers.Success(c, 200, "Tickets retrieved successfully", x)
}
func (h *TicketHandler) GetOwned(c *gin.Context) {
	x, e := h.s.Ticket(c, accessUser(c), c.Param("ticket_id"))
	if e != nil {
		accessFail(c, e)
		return
	}
	helpers.Success(c, 200, "Ticket retrieved successfully", x)
}
func (h *TicketHandler) UpdateOwned(c *gin.Context) {
	q := dto.UpdateTicketRequest{}
	if e := c.ShouldBindJSON(&q); e != nil {
		helpers.Error(c, 400, "Invalid request", e.Error())
		return
	}
	x, e := h.s.UpdateOwnedTicket(c, accessUser(c), c.Param("ticket_id"), &q)
	if e != nil {
		accessFail(c, e)
		return
	}
	helpers.Success(c, 200, "Ticket updated successfully", x)
}
func (h *TicketHandler) CreateTransfers(c *gin.Context) {
	q := dto.CreateTransferRequest{}
	if e := c.ShouldBindJSON(&q); e != nil {
		helpers.Error(c, 400, "Invalid request", e.Error())
		return
	}
	x, e := h.s.CreateTransfer(c, accessUser(c), c.Param("ticket_id"), &q)
	if e != nil {
		accessFail(c, e)
		return
	}
	helpers.Success(c, 201, "Transfer requested successfully", x)
}
func (h *TicketHandler) Transfers(c *gin.Context) {
	x, e := h.s.MyTransfers(c, accessUser(c))
	if e != nil {
		accessFail(c, e)
		return
	}
	helpers.Success(c, 200, "Transfers retrieved successfully", x)
}
func (h *TicketHandler) ResolveTransfer(c *gin.Context) {
	q := dto.UpdateTransferRequest{}
	if e := c.ShouldBindJSON(&q); e != nil {
		helpers.Error(c, 400, "Invalid request", e.Error())
		return
	}
	x, e := h.s.ResolveTransfer(c, accessUser(c), c.Param("transfer_id"), &q)
	if e != nil {
		accessFail(c, e)
		return
	}
	helpers.Success(c, 200, "Transfer resolved successfully", x)
}
func (h *TicketHandler) CreateCheckIn(c *gin.Context) {
	q := dto.CreateCheckInRequest{}
	if e := c.ShouldBindJSON(&q); e != nil {
		helpers.Error(c, 400, "Invalid request", e.Error())
		return
	}
	x, e := h.s.CreateCheckIn(c, accessUser(c), &q)
	if e != nil {
		accessFail(c, e)
		return
	}
	helpers.Success(c, 201, "Check-in recorded successfully", x)
}
func (h *TicketHandler) ListCheckIns(c *gin.Context) {
	x, e := h.s.CheckIns(c, c.Query("ticket_id"))
	if e != nil {
		accessFail(c, e)
		return
	}
	helpers.Success(c, 200, "Check-ins retrieved successfully", x)
}
