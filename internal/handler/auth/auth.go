package auth

import (
	"book_loop/internal/model"
	"book_loop/internal/service/auth"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)
const errorKey = "error"
type handler struct {
	authService auth.Service
	logger *slog.Logger
}

type Handler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Refresh(c *gin.Context)
}

func New(authService auth.Service, logger *slog.Logger) Handler {
	return &handler{
		authService: authService,
		logger: logger}
}

func (h *handler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{errorKey: err.Error()})
		return
	}
	err := h.authService.Register(c.Request.Context(), req.Login, req.Pass, req.Role)
	if err != nil {
		h.logger.Error("Register: ", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{errorKey: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "ok"})
}
func (h *handler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{errorKey: err.Error()})
		return
	}
	user, err := h.authService.Login(c.Request.Context(), req.Login, req.Pass)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound), errors.Is(err, model.ErrInvalid):
			c.JSON(http.StatusUnauthorized, gin.H{errorKey: "invalid credentials"})
		default:
			h.logger.Error("Login", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{errorKey: "internal error"})
		}
	}
	c.JSON(http.StatusOK, user)
}

func (h *handler) Refresh(c *gin.Context) {
	var req model.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Refresh: ", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{errorKey: err.Error()})
		return
	}
	token, err := h.authService.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		h.logger.Error("Refresh: ", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{errorKey: err.Error()})
		return
	}
	c.JSON(http.StatusOK, token)
}