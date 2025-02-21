package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	examv1 "github.com/ApesJs/cbt-exam/api/proto/exam/v1"
	"github.com/ApesJs/cbt-exam/pkg/client"
)

type ExamHandler struct {
	client *client.ServiceClient
}

func NewExamHandler(client *client.ServiceClient) *ExamHandler {
	return &ExamHandler{
		client: client,
	}
}

func (h *ExamHandler) CreateExam(c *gin.Context) {
	var req examv1.CreateExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exam, err := h.client.CreateExam(c.Request.Context(), &req)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		switch st.Code() {
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{"error": st.Message()})
		case codes.AlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": st.Message()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": st.Message()})
		}
		return
	}

	c.JSON(http.StatusCreated, exam)
}

func (h *ExamHandler) GetExam(c *gin.Context) {
	id := c.Param("id")
	exam, err := h.client.GetExam(c.Request.Context(), id)
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

	c.JSON(http.StatusOK, exam)
}

func (h *ExamHandler) UpdateExam(c *gin.Context) {
	id := c.Param("id")
	var exam examv1.Exam
	if err := c.ShouldBindJSON(&exam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := &examv1.UpdateExamRequest{
		Id:   id,
		Exam: &exam,
	}

	updatedExam, err := h.client.UpdateExam(c.Request.Context(), req)
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

	c.JSON(http.StatusOK, updatedExam)
}

func (h *ExamHandler) DeleteExam(c *gin.Context) {
	id := c.Param("id")
	err := h.client.DeleteExam(c.Request.Context(), id)
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

func (h *ExamHandler) ListExams(c *gin.Context) {
	teacherID := c.Query("teacherId")
	pageSize := 10 // default page size
	pageToken := c.Query("pageToken")

	req := &examv1.ListExamsRequest{
		TeacherId: teacherID,
		PageSize:  int32(pageSize),
		PageToken: pageToken,
	}

	resp, err := h.client.ListExams(c.Request.Context(), req)
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

func (h *ExamHandler) ActivateExam(c *gin.Context) {
	id := c.Param("id")
	var req examv1.ActivateExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Id = id

	exam, err := h.client.ActivateExam(c.Request.Context(), &req)
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

	c.JSON(http.StatusOK, exam)
}

func (h *ExamHandler) DeactivateExam(c *gin.Context) {
	id := c.Param("id")
	exam, err := h.client.DeactivateExam(c.Request.Context(), &examv1.DeactivateExamRequest{
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

	c.JSON(http.StatusOK, exam)
}
