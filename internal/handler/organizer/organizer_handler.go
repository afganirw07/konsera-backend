package organizer

import (
	"database/sql"
	"errors"
	dto "konsera-backend/internal/DTO/organizer"
	"konsera-backend/internal/helpers"
	service "konsera-backend/internal/services/organizer"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrganizerHandler struct{ service *service.OrganizerService }

func NewOrganizerHandler(service *service.OrganizerService) *OrganizerHandler {
	return &OrganizerHandler{service: service}
}

func respondOrganizerError(c *gin.Context, err error) {
	status := http.StatusBadRequest
	if errors.Is(err, sql.ErrNoRows) || err.Error() == "organizer not found" {
		status = http.StatusNotFound
	}
	helpers.Error(c, status, err.Error(), nil)
}

// @Summary Create organizer
// @Description Create an organizer profile for a user
// @Tags Organizers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizer body organizer.CreateOrganizerRequest true "Organizer data"
// @Success 201 {object} models.Organizer
// @Failure 400 {object} helpers.APIResponse
// @Failure 401 {object} helpers.APIResponse
// @Failure 403 {object} helpers.APIResponse
// @Router /organizers [post]
func (h *OrganizerHandler) Create(c *gin.Context) {
	var req dto.CreateOrganizerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Error(c, 400, "Invalid request", err.Error())
		return
	}
	result, err := h.service.Create(c, &req)
	if err != nil {
		respondOrganizerError(c, err)
		return
	}
	helpers.Success(c, http.StatusCreated, "Organizer created successfully", result)
}

// @Summary List organizers
// @Description Get all active organizers
// @Tags Organizers
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Organizer
// @Failure 401 {object} helpers.APIResponse
// @Router /organizers [get]
func (h *OrganizerHandler) List(c *gin.Context) {
	result, err := h.service.List(c)
	if err != nil {
		respondOrganizerError(c, err)
		return
	}
	helpers.Success(c, http.StatusOK, "Organizers retrieved successfully", result)
}

// @Summary Get organizer
// @Tags Organizers
// @Produce json
// @Security BearerAuth
// @Param organizer_id path string true "Organizer ID"
// @Success 200 {object} models.Organizer
// @Failure 401 {object} helpers.APIResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /organizers/{organizer_id} [get]
func (h *OrganizerHandler) Get(c *gin.Context) {
	result, err := h.service.Get(c, c.Param("organizer_id"))
	if err != nil {
		respondOrganizerError(c, err)
		return
	}
	helpers.Success(c, http.StatusOK, "Organizer retrieved successfully", result)
}

// @Summary Update organizer
// @Tags Organizers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizer_id path string true "Organizer ID"
// @Param organizer body organizer.UpdateOrganizerRequest true "Organizer changes"
// @Success 200 {object} models.Organizer
// @Failure 400 {object} helpers.APIResponse
// @Failure 401 {object} helpers.APIResponse
// @Failure 403 {object} helpers.APIResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /organizers/{organizer_id} [put]
func (h *OrganizerHandler) Update(c *gin.Context) {
	var req dto.UpdateOrganizerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Error(c, 400, "Invalid request", err.Error())
		return
	}
	result, err := h.service.Update(c, c.Param("organizer_id"), &req)
	if err != nil {
		respondOrganizerError(c, err)
		return
	}
	helpers.Success(c, http.StatusOK, "Organizer updated successfully", result)
}

// @Summary Delete organizer
// @Tags Organizers
// @Produce json
// @Security BearerAuth
// @Param organizer_id path string true "Organizer ID"
// @Success 204 {object} helpers.APIResponse
// @Failure 401 {object} helpers.APIResponse
// @Failure 403 {object} helpers.APIResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /organizers/{organizer_id} [delete]
func (h *OrganizerHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c, c.Param("organizer_id")); err != nil {
		respondOrganizerError(c, err)
		return
	}
	helpers.Message(c, http.StatusNoContent, "Organizer deleted successfully")
}

// @Summary Submit organizer verification
// @Tags Organizer Verifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizer_id path string true "Organizer ID"
// @Param verification body organizer.CreateOrganizerVerificationRequest true "Verification document"
// @Success 201 {object} models.OrganizerVerification
// @Failure 400 {object} helpers.APIResponse
// @Failure 401 {object} helpers.APIResponse
// @Failure 403 {object} helpers.APIResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /organizers/{organizer_id}/verifications [post]
func (h *OrganizerHandler) CreateVerification(c *gin.Context) {
	var req dto.CreateOrganizerVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Error(c, 400, "Invalid request", err.Error())
		return
	}
	result, err := h.service.CreateVerification(c, c.Param("organizer_id"), &req)
	if err != nil {
		respondOrganizerError(c, err)
		return
	}
	helpers.Success(c, http.StatusCreated, "Verification created successfully", result)
}

// @Summary List organizer verifications
// @Tags Organizer Verifications
// @Produce json
// @Security BearerAuth
// @Param organizer_id path string true "Organizer ID"
// @Success 200 {array} models.OrganizerVerification
// @Failure 401 {object} helpers.APIResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /organizers/{organizer_id}/verifications [get]
func (h *OrganizerHandler) ListVerifications(c *gin.Context) {
	result, err := h.service.ListVerifications(c, c.Param("organizer_id"))
	if err != nil {
		respondOrganizerError(c, err)
		return
	}
	helpers.Success(c, http.StatusOK, "Verifications retrieved successfully", result)
}

// @Summary Get organizer verification
// @Tags Organizer Verifications
// @Produce json
// @Security BearerAuth
// @Param organizer_id path string true "Organizer ID"
// @Param verification_id path string true "Verification ID"
// @Success 200 {object} models.OrganizerVerification
// @Failure 401 {object} helpers.APIResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /organizers/{organizer_id}/verifications/{verification_id} [get]
func (h *OrganizerHandler) GetVerification(c *gin.Context) {
	result, err := h.service.GetVerification(c, c.Param("organizer_id"), c.Param("verification_id"))
	if err != nil {
		respondOrganizerError(c, err)
		return
	}
	helpers.Success(c, http.StatusOK, "Verification retrieved successfully", result)
}

// @Summary Update organizer verification
// @Tags Organizer Verifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param organizer_id path string true "Organizer ID"
// @Param verification_id path string true "Verification ID"
// @Param verification body organizer.UpdateOrganizerVerificationRequest true "Verification changes"
// @Success 200 {object} models.OrganizerVerification
// @Failure 400 {object} helpers.APIResponse
// @Failure 401 {object} helpers.APIResponse
// @Failure 403 {object} helpers.APIResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /organizers/{organizer_id}/verifications/{verification_id} [put]
func (h *OrganizerHandler) UpdateVerification(c *gin.Context) {
	var req dto.UpdateOrganizerVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Error(c, 400, "Invalid request", err.Error())
		return
	}
	result, err := h.service.UpdateVerification(c, c.Param("organizer_id"), c.Param("verification_id"), &req)
	if err != nil {
		respondOrganizerError(c, err)
		return
	}
	helpers.Success(c, http.StatusOK, "Verification updated successfully", result)
}

// @Summary Delete organizer verification
// @Tags Organizer Verifications
// @Produce json
// @Security BearerAuth
// @Param organizer_id path string true "Organizer ID"
// @Param verification_id path string true "Verification ID"
// @Success 204 {object} helpers.APIResponse
// @Failure 401 {object} helpers.APIResponse
// @Failure 403 {object} helpers.APIResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /organizers/{organizer_id}/verifications/{verification_id} [delete]
func (h *OrganizerHandler) DeleteVerification(c *gin.Context) {
	if err := h.service.DeleteVerification(c, c.Param("organizer_id"), c.Param("verification_id")); err != nil {
		respondOrganizerError(c, err)
		return
	}
	helpers.Message(c, http.StatusNoContent, "Verification deleted successfully")
}
