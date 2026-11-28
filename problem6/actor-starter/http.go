package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateReq struct {
	ID    string `json:"id" binding:"required"`
	Quota int    `json:"quota" binding:"gt=0"`
}

type UpdateReq struct {
	Quota int `json:"quota" binding:"gt=0"`
}

type GrabReq struct {
	UserID string `json:"user_id" binding:"required"`
}

type ActivityResp struct {
	ID            string `json:"id"`
	Quota         int    `json:"quota"`
	Remaining     int    `json:"remaining"`
	EnrolledCount int    `json:"enrolled_count"`
}

type GrabResp struct {
	Code      string `json:"code"`
	Remaining int    `json:"remaining"`
}

var (
	ErrExists        = errors.New("ACTIVITY_EXISTS")
	ErrNotFound      = errors.New("NOT_FOUND")
	ErrRepeat        = errors.New("ALREADY_GRABBED")
	ErrSoldOut       = errors.New("SOLD_OUT")
	ErrQuotaTooSmall = errors.New("QUOTA_TOO_SMALL")
	ErrBusy          = errors.New("BUSY")
)

type Service interface {
	Create(ctx context.Context, id string, quota int) (ActivityResp, error)
	Grab(ctx context.Context, id, userID string) (remaining int, err error)
	Status(ctx context.Context, id string) (ActivityResp, error)
	Update(ctx context.Context, id string, quota int) (ActivityResp, error)
	Delete(ctx context.Context, id string) error
	Shutdown()
}

func NewApp() (http.Handler, func()) {
	svc := newService()
	r := gin.New()
	r.POST("/activities", func(c *gin.Context) {
		var req CreateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST"})
			return
		}
		activity, err := svc.Create(c.Request.Context(), req.ID, req.Quota)
		respond(c, http.StatusCreated, activity, err)
	})
	r.POST("/activities/:id/grab", func(c *gin.Context) {
		var req GrabReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST"})
			return
		}
		remaining, err := svc.Grab(c.Request.Context(), c.Param("id"), req.UserID)
		respond(c, http.StatusOK, GrabResp{Code: "OK", Remaining: remaining}, err)
	})
	r.GET("/activities/:id", func(c *gin.Context) {
		activity, err := svc.Status(c.Request.Context(), c.Param("id"))
		respond(c, http.StatusOK, activity, err)
	})
	r.PATCH("/activities/:id", func(c *gin.Context) {
		var req UpdateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST"})
			return
		}
		activity, err := svc.Update(c.Request.Context(), c.Param("id"), req.Quota)
		respond(c, http.StatusOK, activity, err)
	})
	r.DELETE("/activities/:id", func(c *gin.Context) {
		if err := svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
			respondError(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	})
	return r, svc.Shutdown
}

func respond(c *gin.Context, status int, body any, err error) {
	if err == nil {
		c.JSON(status, body)
		return
	}
	respondError(c, err)
}

func respondError(c *gin.Context, err error) {
	for _, failure := range httpErrors {
		if errors.Is(err, failure.err) {
			c.JSON(failure.status, gin.H{"code": failure.err.Error()})
			return
		}
	}
	if c.Request.Context().Err() != nil {
		return // 客户端已结束等待。
	}
	c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_ERROR"})
}

var httpErrors = []struct {
	err    error
	status int
}{
	{ErrExists, http.StatusConflict},
	{ErrQuotaTooSmall, http.StatusConflict},
	{ErrNotFound, http.StatusNotFound},
	{ErrRepeat, http.StatusConflict},
	{ErrSoldOut, http.StatusGone},
	{ErrBusy, http.StatusServiceUnavailable},
}
