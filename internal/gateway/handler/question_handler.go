package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	questionv1 "github.com/ApesJs/cbt-exam/api/proto/question/v1"
	"github.com/ApesJs/cbt-exam/pkg/client"
)

type QuestionHandler struct {
	client *client.ServiceClient
}

func NewQuestionHandler(client *client.ServiceClient) *QuestionHandler {
	return &QuestionHandler{
		client: client,
	}
}

func (h *QuestionHandler) CreateQuestion(c *gin.Context) {
	var req questionv1.CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question, err := h.client.CreateQuestion(c.Request.Context(), &req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		switch st.Code() {
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{"error": st.Message()})
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": st.Message()})
		case codes.FailedPrecondition:
			c.JSON(http.StatusPreconditionFailed, gin.H{"error": st.Message()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		}
		return
	}

	c.JSON(http.StatusCreated, question)
}

func (h *QuestionHandler) GetQuestion(c *gin.Context) {
	id := c.Param("id")
	question, err := h.client.GetQuestion(c.Request.Context(), &questionv1.GetQuestionRequest{
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

	c.JSON(http.StatusOK, question)
}

func (h *QuestionHandler) ListQuestions(c *gin.Context) {
	examID := c.Query("examId")
	pageSize := 10 // default page size
	if size := c.Query("pageSize"); size != "" {
		if s, err := strconv.Atoi(size); err == nil {
			pageSize = s
		}
	}
	pageToken := c.Query("pageToken")

	req := &questionv1.ListQuestionsRequest{
		ExamId:    examID,
		PageSize:  int32(pageSize),
		PageToken: pageToken,
	}

	resp, err := h.client.ListQuestions(c.Request.Context(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *QuestionHandler) GetExamQuestions(c *gin.Context) {
	examID := c.Param("examId")
	randomize := c.DefaultQuery("randomize", "false") == "true"
	limit := int32(0) // 0 means no limit
	if l := c.Query("limit"); l != "" {
		if lim, err := strconv.Atoi(l); err == nil {
			limit = int32(lim)
		}
	}

	req := &questionv1.GetExamQuestionsRequest{
		ExamId:    examID,
		Randomize: randomize,
		Limit:     limit,
	}

	resp, err := h.client.GetExamQuestions(c.Request.Context(), req)
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

func (h *QuestionHandler) UpdateQuestion(c *gin.Context) {
	id := c.Param("id")
	var question questionv1.Question
	if err := c.ShouldBindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := &questionv1.UpdateQuestionRequest{
		Id:       id,
		Question: &question,
	}

	updatedQuestion, err := h.client.UpdateQuestion(c.Request.Context(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		switch st.Code() {
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": st.Message()})
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{"error": st.Message()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		}
		return
	}

	c.JSON(http.StatusOK, updatedQuestion)
}

func (h *QuestionHandler) DeleteQuestion(c *gin.Context) {
	id := c.Param("id")
	err := h.client.DeleteQuestion(c.Request.Context(), &questionv1.DeleteQuestionRequest{
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

	c.JSON(http.StatusNoContent, nil)
}
