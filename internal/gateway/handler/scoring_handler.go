package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	scoringv1 "github.com/ApesJs/cbt-exam/api/proto/scoring/v1"
	"github.com/ApesJs/cbt-exam/pkg/client"
)

type ScoringHandler struct {
	client *client.ServiceClient
}

func NewScoringHandler(client *client.ServiceClient) *ScoringHandler {
	return &ScoringHandler{
		client: client,
	}
}

func (h *ScoringHandler) CalculateScore(c *gin.Context) {
	var req scoringv1.CalculateScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	score, err := h.client.CalculateScore(c.Request.Context(), &req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		switch st.Code() {
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": st.Message()})
		case codes.AlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": st.Message()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		}
		return
	}

	c.JSON(http.StatusCreated, score)
}

func (h *ScoringHandler) GetScore(c *gin.Context) {
	id := c.Param("id")
	score, err := h.client.GetScore(c.Request.Context(), &scoringv1.GetScoreRequest{
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

	c.JSON(http.StatusOK, score)
}

func (h *ScoringHandler) ListScores(c *gin.Context) {
	examID := c.Param("examId")
	pageSize := 10 // default page size
	if size := c.Query("pageSize"); size != "" {
		if s, err := strconv.Atoi(size); err == nil {
			pageSize = s
		}
	}
	pageToken := c.Query("pageToken")

	req := &scoringv1.ListScoresRequest{
		ExamId:    examID,
		PageSize:  int32(pageSize),
		PageToken: pageToken,
	}

	resp, err := h.client.ListScores(c.Request.Context(), req)
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

	c.JSON(http.StatusOK, resp)
}
