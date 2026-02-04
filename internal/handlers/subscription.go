package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/Estriper0/subscription_service/internal/handlers/dto"
	"github.com/Estriper0/subscription_service/internal/service"
	"github.com/Estriper0/subscription_service/internal/service/domain"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	subscriptionService ISubscriptionService
	validate            *validator.Validate
}

type ISubscriptionService interface {
	Create(ctx context.Context, subscription *domain.SubscriptionCreate) (int, error)
	GetByUser(ctx context.Context, user_id uuid.UUID) ([]*domain.Subscription, error)
	GetById(ctx context.Context, id int) (*domain.Subscription, error)
	DeleteById(ctx context.Context, id int) (*domain.Subscription, error)
	Update(ctx context.Context, data *domain.SubscriptionUpdate) (*domain.Subscription, error)
	GetPriceByFilter(ctx context.Context, user_id *uuid.UUID, service_name *string, startDate, endDate string) (int, error)
}

func NewSubscriptionHandler(g *gin.RouterGroup, subscriptionService ISubscriptionService, validate *validator.Validate) {
	r := &SubscriptionHandler{
		subscriptionService: subscriptionService,
		validate:            validate,
	}

	g.POST("/", r.Add)
	g.GET("/", r.GetById)
	g.DELETE("/", r.DeleteById)
	g.PATCH("/", r.Update)
	g.POST("/price", r.GetPriceByFilter)
	g.GET("/user", r.GetByUser)
}

func (h *SubscriptionHandler) Add(c *gin.Context) {
	var req dto.SubscriptionCreateRequest

	if err := c.Bind(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, ErrStatusBadRequest, err)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respondWithError(c, http.StatusBadRequest, ErrStatusBadRequest, err)
		return
	}

	UUID, _ := uuid.Parse(req.UserId)

	id, err := h.subscriptionService.Create(c.Request.Context(), &domain.SubscriptionCreate{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserId:      UUID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	})
	if err != nil {
		if errors.Is(err, service.ErrIncorrectTime) {
			respondWithError(c, http.StatusBadRequest, ErrStatusBadRequest, err)
			return
		}
		respondWithError(c, http.StatusInternalServerError, ErrStatusInternal, err)
		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"id": id,
		},
	)
}

func (h *SubscriptionHandler) GetByUser(c *gin.Context) {
	userId, ok := c.GetQuery("user_id")
	if !ok {
		respondWithError(c, http.StatusBadRequest, ErrStatusBadRequest, errors.New("invalid query params"))
		return
	}

	parseUUID, err := uuid.Parse(userId)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, ErrStatusBadRequest, errors.New("incorrect uuid"))
		return
	}

	subscriptions, err := h.subscriptionService.GetByUser(c.Request.Context(), parseUUID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, ErrStatusInternal, err)
		return
	}

	var res []dto.Subscription
	for _, s := range subscriptions {
		k := dto.Subscription{
			Id:          s.Id,
			ServiceName: s.ServiceName,
			Price:       s.Price,
			UserId:      s.UserId,
			StartDate:   s.StartDate,
			EndDate:     s.EndDate,
		}
		res = append(res, k)
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"subscriptions": res,
		},
	)
}

func (h *SubscriptionHandler) GetById(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		respondWithError(c, http.StatusBadRequest, ErrStatusBadRequest, errors.New("invalid query params"))
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil || idInt < 0 {
		respondWithError(c, http.StatusBadRequest, ErrStatusBadRequest, errors.New("id must be a non-negative integer"))
		return
	}

	subscription, err := h.subscriptionService.GetById(c.Request.Context(), idInt)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			respondWithError(c, http.StatusNotFound, ErrStatusNotFound, err)
			return
		}
		respondWithError(c, http.StatusInternalServerError, ErrStatusInternal, err)
		return
	}

	res := dto.Subscription{
		Id:          subscription.Id,
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		UserId:      subscription.UserId,
		StartDate:   subscription.StartDate,
		EndDate:     subscription.EndDate,
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"subscription": res,
		},
	)
}

func (h *SubscriptionHandler) DeleteById(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		respondWithError(c, http.StatusBadRequest, ErrStatusBadRequest, errors.New("invalid query params"))
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil || idInt < 0 {
		respondWithError(c, http.StatusBadRequest, ErrStatusBadRequest, errors.New("id must be a non-negative integer"))
		return
	}

	subscription, err := h.subscriptionService.DeleteById(c.Request.Context(), idInt)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			respondWithError(c, http.StatusNotFound, ErrStatusNotFound, err)
			return
		}
		respondWithError(c, http.StatusInternalServerError, ErrStatusInternal, err)
		return
	}

	res := dto.Subscription{
		Id:          subscription.Id,
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		UserId:      subscription.UserId,
		StartDate:   subscription.StartDate,
		EndDate:     subscription.EndDate,
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"subscription": res,
		},
	)
}

func (h *SubscriptionHandler) Update(c *gin.Context) {
	var req dto.SubscriptionUpdateRequest

	if err := c.Bind(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, ErrStatusBadRequest, err)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respondWithError(c, http.StatusBadRequest, ErrStatusBadRequest, err)
		return
	}

	subscription, err := h.subscriptionService.Update(c.Request.Context(), &domain.SubscriptionUpdate{
		Id:          req.Id,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	})
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			respondWithError(c, http.StatusNotFound, ErrStatusNotFound, err)
			return
		}
		respondWithError(c, http.StatusInternalServerError, ErrStatusInternal, err)
		return
	}

	res := dto.Subscription{
		Id:          subscription.Id,
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		UserId:      subscription.UserId,
		StartDate:   subscription.StartDate,
		EndDate:     subscription.EndDate,
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"subscription": res,
		},
	)
}

func (h *SubscriptionHandler) GetPriceByFilter(c *gin.Context) {
	var req dto.SubscriptionFilterRequest

	if err := c.Bind(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, ErrStatusBadRequest, err)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respondWithError(c, http.StatusBadRequest, ErrStatusBadRequest, err)
		return
	}

	var n *string
	serviceName, ok := c.GetQuery("service_name")
	if ok {
		n = &serviceName
	}

	var u *uuid.UUID
	userId, ok := c.GetQuery("user_id")
	if ok {
		parseUUID, err := uuid.Parse(userId)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, ErrStatusBadRequest, errors.New("incorrect uuid"))
			return
		}
		u = &parseUUID
	}

	price, err := h.subscriptionService.GetPriceByFilter(c.Request.Context(), u, n, req.StartDate, req.EndDate)
	if err != nil {
		if errors.Is(err, service.ErrIncorrectTime) {
			respondWithError(c, http.StatusBadRequest, ErrStatusBadRequest, err)
			return
		}
		respondWithError(c, http.StatusInternalServerError, ErrStatusInternal, err)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"price": price,
		},
	)
}
