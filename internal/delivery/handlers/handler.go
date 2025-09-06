package handlers

import (
	"net/http"
	"strconv"
	"strings"
	httpdto "subcalc/internal/delivery/http"
	"subcalc/internal/domain"
	"subcalc/internal/repository"
	"subcalc/internal/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Handler struct {
	usecase usecase.SubscriptionUsecase
	log     *zap.SugaredLogger
}

func NewHandler(u usecase.SubscriptionUsecase, log *zap.SugaredLogger) *Handler {
	return &Handler{usecase: u, log: log}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		s := api.Group("/subscriptions")
		{
			s.POST("", h.Create)
			s.GET("", h.List)
			s.GET("/sum", h.Sum)
			s.GET("/:id", h.GetByID)
			s.PUT("/:id", h.Update)
			s.DELETE("/:id", h.Delete)
		}
	}
}

// Create godoc
// @Summary Create subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body httpdto.CreateSubscriptionRequest true "subscription"
// @Success 201 {object} domain.Subscription
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subscriptions [post]
func (h *Handler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var req httpdto.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warnf("invalid create body: %v", err)
		RespondError(c, http.StatusBadRequest, "invalid_payload", "invalid request body", map[string]string{"body": err.Error()})
		return
	}
	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid_field", "user_id must be a UUID", map[string]string{"user_id": "invalid uuid"})
		return
	}
	start, err := parseMonthYear(req.StartDate)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid_field", "start_date must be in format MM-YYYY", map[string]string{"start_date": "expected MM-YYYY"})
		return
	}
	var endPtr *time.Time
	if req.EndDate != nil {
		t, err := parseMonthYear(*req.EndDate)
		if err != nil {
			RespondError(c, http.StatusBadRequest, "invalid_field", "end_date must be in format MM-YYYY", map[string]string{"end_date": "expected MM-YYYY"})
			return
		}
		if t.Before(start) {
			RespondError(c, http.StatusBadRequest, "invalid_field", "end_date must be equal or after start_date", map[string]string{"end_date": "must be >= start_date"})
			return
		}
		endPtr = &t
	}

	req.ServiceName = strings.TrimSpace(req.ServiceName)
	if req.ServiceName == "" || len(req.ServiceName) > 255 {
		RespondError(c, http.StatusBadRequest, "invalid_field", "service_name is required and must be <=255 chars", map[string]string{"service_name": "required, max 255 chars"})
		return
	}

	sub := &domain.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      uid,
		StartDate:   start,
		EndDate:     endPtr,
	}
	if err := h.usecase.Create(ctx, sub); err != nil {
		h.log.Errorf("create failed: %v", err)
		RespondError(c, http.StatusInternalServerError, "internal_error", "create failed", nil)
		return
	}
	c.JSON(http.StatusCreated, sub)
}

// List godoc
// @Summary List subscriptions
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "user uuid"
// @Param service_name query string false "service name"
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Success 200 {array} domain.Subscription
// @Header 200 {string} X-Total-Count "Total number of subscriptions matching the filter"
// @Failure 500 {object} ErrorResponse
// @Router /api/subscriptions [get]
func (h *Handler) List(c *gin.Context) {
	ctx := c.Request.Context()

	var filter repository.SubscriptionFilter
	if uidStr := c.Query("user_id"); uidStr != "" {
		uid, err := uuid.Parse(uidStr)
		if err == nil {
			filter.UserID = &uid
		}
	}
	if s := c.Query("service_name"); s != "" {
		filter.ServiceName = &s
	}

	const (
		defaultLimit = 100
		maxLimit     = 1000
	)
	limit := defaultLimit
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			if v > 0 {
				limit = v
			}
		}
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	offset := 0
	if o := c.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}
	filter.Limit = limit
	filter.Offset = offset

	if fromStr := c.Query("from"); fromStr != "" {
		if t, err := parseMonthYear(fromStr); err == nil {
			filter.From = &t
		} else {
			RespondError(c, http.StatusBadRequest, "invalid_field", "from must be MM-YYYY", map[string]string{"from": "expected MM-YYYY"})
			return
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		if t, err := parseMonthYear(toStr); err == nil {
			filter.To = &t
		} else {
			RespondError(c, http.StatusBadRequest, "invalid_field", "to must be MM-YYYY", map[string]string{"to": "expected MM-YYYY"})
			return
		}
	}

	total, err := h.usecase.Count(ctx, filter)
	if err != nil {
		h.log.Errorf("count failed: %v", err)
		RespondError(c, http.StatusInternalServerError, "internal_error", "list failed", nil)
		return
	}
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))

	subs, err := h.usecase.List(ctx, filter)
	if err != nil {
		h.log.Errorf("list failed: %v", err)
		RespondError(c, http.StatusInternalServerError, "internal_error", "list failed", nil)
		return
	}
	c.JSON(http.StatusOK, subs)
}

// Sum godoc
// @Summary Sum subscriptions for period
// @Tags subscriptions
// @Produce json
// @Param from query string true "start month-year MM-YYYY"
// @Param to query string true "end month-year MM-YYYY"
// @Param user_id query string false "user uuid"
// @Param service_name query string false "service name"
// @Success 200 {object} httpdto.TotalResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subscriptions/sum [get]
func (h *Handler) Sum(c *gin.Context) {
	ctx := c.Request.Context()

	fromStr := c.Query("from")
	toStr := c.Query("to")
	if fromStr == "" || toStr == "" {
		RespondError(c, http.StatusBadRequest, "invalid_request", "from and to query params required, format MM-YYYY", map[string]string{"from": "required", "to": "required"})
		return
	}
	from, err := parseMonthYear(fromStr)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid_field", "invalid from format, expected MM-YYYY", map[string]string{"from": "expected MM-YYYY"})
		return
	}
	to, err := parseMonthYear(toStr)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid_field", "invalid to format, expected MM-YYYY", map[string]string{"to": "expected MM-YYYY"})
		return
	}
	if from.After(to) {
		RespondError(c, http.StatusBadRequest, "invalid_request", "'from' must be before or equal to 'to'", map[string]string{"from": "must be <= to"})
		return
	}

	var filter repository.SubscriptionFilter
	filter.From = &from
	filter.To = &to
	if uidStr := c.Query("user_id"); uidStr != "" {
		uid, err := uuid.Parse(uidStr)
		if err == nil {
			filter.UserID = &uid
		}
	}
	if s := c.Query("service_name"); s != "" {
		filter.ServiceName = &s
	}

	total, err := h.usecase.SumSubscriptions(ctx, filter)
	if err != nil {
		h.log.Errorf("sum failed: %v", err)
		RespondError(c, http.StatusInternalServerError, "internal_error", "sum failed", nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": total})
}

// GetByID godoc
// @Summary Get subscription by id
// @Tags subscriptions
// @Produce json
// @Param id path string true "subscription id"
// @Success 200 {object} domain.Subscription
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subscriptions/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid_field", "invalid id", map[string]string{"id": "invalid uuid"})
		return
	}
	sub, err := h.usecase.GetByID(ctx, id)
	if err != nil {
		h.log.Errorf("get failed: %v", err)
		RespondError(c, http.StatusInternalServerError, "internal_error", "get failed", nil)
		return
	}
	if sub == nil {
		RespondError(c, http.StatusNotFound, "not_found", "not found", nil)
		return
	}
	c.JSON(http.StatusOK, sub)
}

// Update godoc
// @Summary Update subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "subscription id"
// @Param input body httpdto.UpdateSubscriptionRequest true "update fields"
// @Success 200 {object} domain.Subscription
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/subscriptions/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid_field", "invalid id", map[string]string{"id": "invalid uuid"})
		return
	}
	existing, err := h.usecase.GetByID(ctx, id)
	if err != nil {
		h.log.Errorf("get for update failed: %v", err)
		RespondError(c, http.StatusInternalServerError, "internal_error", "get failed", nil)
		return
	}
	if existing == nil {
		RespondError(c, http.StatusNotFound, "not_found", "not found", nil)
		return
	}

	var req httpdto.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warnf("invalid update body: %v", err)
		RespondError(c, http.StatusBadRequest, "invalid_payload", "invalid request body", map[string]string{"body": err.Error()})
		return
	}

	if req.ServiceName != nil {
		s := strings.TrimSpace(*req.ServiceName)
		if s == "" || len(s) > 255 {
			RespondError(c, http.StatusBadRequest, "invalid_field", "service_name must be non-empty and <=255 chars", map[string]string{"service_name": "required, max 255 chars"})
			return
		}
		existing.ServiceName = s
	}
	if req.Price != nil {
		if *req.Price < 0 {
			RespondError(c, http.StatusBadRequest, "invalid_field", "price must be >= 0", map[string]string{"price": "must be >= 0"})
			return
		}
		existing.Price = *req.Price
	}
	if req.StartDate != nil {
		t, err := parseMonthYear(*req.StartDate)
		if err != nil {
			RespondError(c, http.StatusBadRequest, "invalid_field", "invalid start_date, expected MM-YYYY", map[string]string{"start_date": "expected MM-YYYY"})
			return
		}
		existing.StartDate = t
		if existing.EndDate != nil && existing.EndDate.Before(existing.StartDate) {
			RespondError(c, http.StatusBadRequest, "invalid_field", "existing end_date is before new start_date", map[string]string{"end_date": "must be >= start_date"})
			return
		}
	}
	if req.EndDate != nil {
		if *req.EndDate == "" {
			existing.EndDate = nil
		} else {
			t, err := parseMonthYear(*req.EndDate)
			if err != nil {
				RespondError(c, http.StatusBadRequest, "invalid_field", "invalid end_date, expected MM-YYYY", map[string]string{"end_date": "expected MM-YYYY"})
				return
			}
			if t.Before(existing.StartDate) {
				RespondError(c, http.StatusBadRequest, "invalid_field", "end_date must be equal or after start_date", map[string]string{"end_date": "must be >= start_date"})
				return
			}
			existing.EndDate = &t
		}
	}

	if err := h.usecase.Update(ctx, existing); err != nil {
		h.log.Errorf("update failed: %v", err)
		RespondError(c, http.StatusInternalServerError, "internal_error", "update failed", nil)
		return
	}
	c.JSON(http.StatusOK, existing)
}

// Delete godoc
// @Summary Delete subscription
// @Tags subscriptions
// @Param id path string true "subscription id"
// @Success 204
// @Failure 500 {object} ErrorResponse
// @Router /api/subscriptions/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "invalid_field", "invalid id", map[string]string{"id": "invalid uuid"})
		return
	}
	if err := h.usecase.Delete(ctx, id); err != nil {
		h.log.Errorf("delete failed: %v", err)
		RespondError(c, http.StatusInternalServerError, "internal_error", "delete failed", nil)
		return
	}
	c.Status(http.StatusNoContent)
}
