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
	GetByUser(ctx context.Context, userId uuid.UUID) ([]*domain.Subscription, error)
	GetById(ctx context.Context, id int) (*domain.Subscription, error)
	DeleteById(ctx context.Context, id int) (*domain.Subscription, error)
	Update(ctx context.Context, data *domain.SubscriptionUpdate) (*domain.Subscription, error)
	GetPriceByFilter(ctx context.Context, userId *uuid.UUID, serviceName *string, startDate, endDate string) (int, error)
}

func NewSubscriptionHandler(g *gin.RouterGroup, subscriptionService ISubscriptionService, validate *validator.Validate) {
	r := &SubscriptionHandler{
		subscriptionService: subscriptionService,
		validate:            validate,
	}

	g.POST("/", r.Add)
	g.GET("/:id", r.GetById)
	g.DELETE("/:id", r.DeleteById)
	g.PATCH("/:id", r.Update)
	g.GET("/price", r.GetPriceByFilter)
	g.GET("/user/:userId", r.GetByUser)
}

// Add godoc
// @Summary Создать новую подписку
// @Description Создаёт новую подписку для пользователя
// @Tags subscription
// @Accept json
// @Produce json
// @Param request body dto.SubscriptionCreateRequest true "Данные для создания подписки"
// @Router /subscription [post]
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

// GetByUser godoc
// @Summary Получить подписки пользователя
// @Description Возвращает список всех подписок пользователя по его ID
// @Tags subscription
// @Accept json
// @Produce json
// @Param userId path string true "UUID пользователя" format(uuid)
// @Router /subscription/user/{userId} [get]
func (h *SubscriptionHandler) GetByUser(c *gin.Context) {
	userId := c.Param("userId")
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

// GetById godoc
// @Summary Получить подписку по ID
// @Description Возвращает информацию о конкретной подписке по её ID
// @Tags subscription
// @Accept json
// @Produce json
// @Param id path integer true "ID подписки" minimum(0)
// @Router /subscription/{id} [get]
func (h *SubscriptionHandler) GetById(c *gin.Context) {
	id := c.Param("id")
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

// DeleteById godoc
// @Summary Удалить подписку по ID
// @Description Удаляет подписку по её ID и возвращает информацию об удалённой подписке
// @Tags subscription
// @Accept json
// @Produce json
// @Param id path integer true "ID подписки для удаления" minimum(0)
// @Router /subscription/{id} [delete]
func (h *SubscriptionHandler) DeleteById(c *gin.Context) {
	id := c.Param("id")
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

// Update godoc
// @Summary Обновить подписку
// @Description Обновляет информацию о существующей подписке
// @Tags subscription
// @Accept json
// @Produce json
// @Param id path integer true "ID подписки для обновления" minimum(0)
// @Param request body dto.SubscriptionUpdateRequest true "Данные для обновления подписки"
// @Router /subscription/{id} [patch]
func (h *SubscriptionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil || idInt < 0 {
		respondWithError(c, http.StatusBadRequest, ErrStatusBadRequest, errors.New("id must be a non-negative integer"))
		return
	}

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
		Id:          idInt,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	})
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			respondWithError(c, http.StatusNotFound, ErrStatusNotFound, err)
			return
		} else if errors.Is(err, service.ErrIncorrectTime) {
			respondWithError(c, http.StatusBadRequest, ErrStatusBadRequest, err)
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

// GetPriceByFilter godoc
// @Summary Получить сумму подписок по фильтрам
// @Description Рассчитывает общую стоимость подписок по заданным фильтрам (пользователь, сервис, период)
// @Tags subscription
// @Accept json
// @Produce json
// @Param user_id query string false "UUID пользователя для фильтрации" format(uuid)
// @Param service_name query string false "Название сервиса для фильтрации"
// @Param start_date query string true "Начальная дата для подсчета суммы"
// @Param end_date query string true "Конечная дата для подсчета суммы"
// @Router /subscription/price [get]
func (h *SubscriptionHandler) GetPriceByFilter(c *gin.Context) {
	startDate, ok := c.GetQuery("start_date")
	if !ok {
		respondWithError(c, http.StatusBadRequest, ErrStatusBadRequest, errors.New("no start date"))
		return
	}

	endDate, ok := c.GetQuery("end_date")
	if !ok {
		respondWithError(c, http.StatusBadRequest, ErrStatusBadRequest, errors.New("no end date"))
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

	price, err := h.subscriptionService.GetPriceByFilter(c.Request.Context(), u, n, startDate, endDate)
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
