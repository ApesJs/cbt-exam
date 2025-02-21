package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sessionv1 "github.com/ApesJs/cbt-exam/api/proto/session/v1"
	"github.com/ApesJs/cbt-exam/pkg/client"
)

type SessionHandler struct {
	client *client.ServiceClient
}

func NewSessionHandler(client *client.ServiceClient) *SessionHandler {
	return &SessionHandler{
		client: client,
	}
}

func (h *SessionHandler) StartSession(c *gin.Context) {
	var req sessionv1.StartSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.client.StartSession(c.Request.Context(), &req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		switch st.Code() {
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": st.Message()})
		case codes.FailedPrecondition:
			c.JSON(http.StatusPreconditionFailed, gin.H{"error": st.Message()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		}
		return
	}

	c.JSON(http.StatusCreated, session)
}

func (h *SessionHandler) GetSession(c *gin.Context) {
	id := c.Param("id")
	session, err := h.client.GetSession(c.Request.Context(), &sessionv1.GetSessionRequest{
		Id: id,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if st.Code() == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": st.Message()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *SessionHandler) SubmitAnswer(c *gin.Context) {
	id := c.Param("id")
	var req sessionv1.SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.SessionId = id

	resp, err := h.client.SubmitAnswer(c.Request.Context(), &req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		switch st.Code() {
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": st.Message()})
		case codes.FailedPrecondition:
			c.JSON(http.StatusPreconditionFailed, gin.H{"error": st.Message()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *SessionHandler) FinishSession(c *gin.Context) {
	id := c.Param("id")
	session, err := h.client.FinishSession(c.Request.Context(), &sessionv1.FinishSessionRequest{
		Id: id,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		switch st.Code() {
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": st.Message()})
		case codes.FailedPrecondition:
			c.JSON(http.StatusPreconditionFailed, gin.H{"error": st.Message()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		}
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *SessionHandler) GetRemainingTime(c *gin.Context) {
	id := c.Param("id")
	time, err := h.client.GetRemainingTime(c.Request.Context(), &sessionv1.GetRemainingTimeRequest{
		SessionId: id,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if st.Code() == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": st.Message()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}

	c.JSON(http.StatusOK, time)
}
