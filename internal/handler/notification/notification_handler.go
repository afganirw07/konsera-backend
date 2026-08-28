package notification

import (
	"database/sql"
	"errors"
	"net/http"

	dto "konsera-backend/internal/DTO/notification"
	"konsera-backend/internal/helpers"
	svc "konsera-backend/internal/services/notification"

	"github.com/gin-gonic/gin"
)

type Handler struct{ service *svc.Service }

func NewHandler(service *svc.Service) *Handler { return &Handler{service: service} }
func fail(c *gin.Context, err error) {
	status := http.StatusBadRequest
	if errors.Is(err, sql.ErrNoRows) {
		status = http.StatusNotFound
	}
	helpers.Error(c, status, err.Error(), nil)
}
func bind(c *gin.Context, value any) bool {
	if err := c.ShouldBindJSON(value); err != nil {
		helpers.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return false
	}
	return true
}
func userID(c *gin.Context) string { value, _ := c.Get("user_id"); return value.(string) }

// @Summary List notification templates
// @Tags Notification Templates
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.NotificationTemplate
// @Router /notification-templates [get]
func (h *Handler) ListTemplates(c *gin.Context) {
	items, err := h.service.ListTemplates(c)
	if err != nil {
		fail(c, err)
		return
	}
	helpers.Success(c, 200, "Notification templates retrieved successfully", items)
}

// @Summary Get notification template
// @Tags Notification Templates
// @Produce json
// @Security BearerAuth
// @Param template_id path string true "Template ID"
// @Success 200 {object} models.NotificationTemplate
// @Router /notification-templates/{template_id} [get]
func (h *Handler) GetTemplate(c *gin.Context) {
	item, err := h.service.GetTemplate(c, c.Param("template_id"))
	if err != nil {
		fail(c, err)
		return
	}
	helpers.Success(c, 200, "Notification template retrieved successfully", item)
}

// @Summary Create notification template
// @Tags Notification Templates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param template body notification.CreateNotificationTemplateRequest true "Template"
// @Success 201 {object} models.NotificationTemplate
// @Router /notification-templates [post]
func (h *Handler) CreateTemplate(c *gin.Context) {
	q := dto.CreateNotificationTemplateRequest{}
	if !bind(c, &q) {
		return
	}
	item, err := h.service.CreateTemplate(c, &q)
	if err != nil {
		fail(c, err)
		return
	}
	helpers.Success(c, 201, "Notification template created successfully", item)
}

// @Summary Update notification template
// @Tags Notification Templates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param template_id path string true "Template ID"
// @Param template body notification.UpdateNotificationTemplateRequest true "Template changes"
// @Success 200 {object} models.NotificationTemplate
// @Router /notification-templates/{template_id} [put]
func (h *Handler) UpdateTemplate(c *gin.Context) {
	q := dto.UpdateNotificationTemplateRequest{}
	if !bind(c, &q) {
		return
	}
	item, err := h.service.UpdateTemplate(c, c.Param("template_id"), &q)
	if err != nil {
		fail(c, err)
		return
	}
	helpers.Success(c, 200, "Notification template updated successfully", item)
}

// @Summary Delete notification template
// @Tags Notification Templates
// @Produce json
// @Security BearerAuth
// @Param template_id path string true "Template ID"
// @Success 204 {object} helpers.APIResponse
// @Router /notification-templates/{template_id} [delete]
func (h *Handler) DeleteTemplate(c *gin.Context) {
	if err := h.service.DeleteTemplate(c, c.Param("template_id")); err != nil {
		fail(c, err)
		return
	}
	helpers.Message(c, 204, "Notification template deleted successfully")
}

// @Summary List user notifications
// @Tags Notifications
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Notification
// @Router /notifications [get]
func (h *Handler) List(c *gin.Context) {
	items, err := h.service.List(c, userID(c))
	if err != nil {
		fail(c, err)
		return
	}
	helpers.Success(c, 200, "Notifications retrieved successfully", items)
}

// @Summary Get user notification
// @Tags Notifications
// @Produce json
// @Security BearerAuth
// @Param notification_id path string true "Notification ID"
// @Success 200 {object} models.Notification
// @Router /notifications/{notification_id} [get]
func (h *Handler) Get(c *gin.Context) {
	item, err := h.service.Get(c, userID(c), c.Param("notification_id"))
	if err != nil {
		fail(c, err)
		return
	}
	helpers.Success(c, 200, "Notification retrieved successfully", item)
}

// @Summary Create notification
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param notification body notification.CreateNotificationRequest true "Notification"
// @Success 201 {object} models.Notification
// @Router /notifications [post]
func (h *Handler) Create(c *gin.Context) {
	q := dto.CreateNotificationRequest{}
	if !bind(c, &q) {
		return
	}
	item, err := h.service.Create(c, &q)
	if err != nil {
		fail(c, err)
		return
	}
	helpers.Success(c, 201, "Notification created successfully", item)
}

// @Summary Update notification
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param notification_id path string true "Notification ID"
// @Param notification body notification.UpdateNotificationRequest true "Notification changes"
// @Success 200 {object} models.Notification
// @Router /notifications/{notification_id} [put]
func (h *Handler) Update(c *gin.Context) {
	q := dto.UpdateNotificationRequest{}
	if !bind(c, &q) {
		return
	}
	item, err := h.service.Update(c, userID(c), c.Param("notification_id"), &q)
	if err != nil {
		fail(c, err)
		return
	}
	helpers.Success(c, 200, "Notification updated successfully", item)
}

// @Summary Delete notification
// @Tags Notifications
// @Produce json
// @Security BearerAuth
// @Param notification_id path string true "Notification ID"
// @Success 204 {object} helpers.APIResponse
// @Router /notifications/{notification_id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.Delete(c, userID(c), c.Param("notification_id")); err != nil {
		fail(c, err)
		return
	}
	helpers.Message(c, 204, "Notification deleted successfully")
}
