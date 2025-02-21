package unit

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	examv1 "github.com/ApesJs/cbt-exam/api/proto/exam/v1"
	questionv1 "github.com/ApesJs/cbt-exam/api/proto/question/v1"
	scoringv1 "github.com/ApesJs/cbt-exam/api/proto/scoring/v1"
	sessionv1 "github.com/ApesJs/cbt-exam/api/proto/session/v1"
)

// TestConfig berisi konfigurasi untuk pengujian performa
type TestConfig struct {
	ExamServiceAddr     string
	QuestionServiceAddr string
	SessionServiceAddr  string
	ScoringServiceAddr  string
	ConcurrentUsers     int
	RequestsPerUser     int
}

// PerformanceResult menyimpan hasil pengujian performa
type PerformanceResult struct {
	TotalRequests      int
	SuccessfulRequests int
	FailedRequests     int
	TotalDuration      time.Duration
	AverageLatency     time.Duration
	MinLatency         time.Duration
	MaxLatency         time.Duration
	RequestsPerSecond  float64
}

// createServiceClients membuat koneksi ke semua service
func createServiceClients(config TestConfig) (
	examv1.ExamServiceClient,
	questionv1.QuestionServiceClient,
	sessionv1.SessionServiceClient,
	scoringv1.ScoringServiceClient,
	[]*grpc.ClientConn,
	error) {

	// Buat koneksi ke ExamService
	examConn, err := grpc.Dial(config.ExamServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to connect to exam service: %v", err)
	}

	// Buat koneksi ke QuestionService
	questionConn, err := grpc.Dial(config.QuestionServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		examConn.Close()
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to connect to question service: %v", err)
	}

	// Buat koneksi ke SessionService
	sessionConn, err := grpc.Dial(config.SessionServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		examConn.Close()
		questionConn.Close()
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to connect to session service: %v", err)
	}

	// Buat koneksi ke ScoringService
	scoringConn, err := grpc.Dial(config.ScoringServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		examConn.Close()
		questionConn.Close()
		sessionConn.Close()
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to connect to scoring service: %v", err)
	}

	// Buat client untuk setiap service
	examClient := examv1.NewExamServiceClient(examConn)
	questionClient := questionv1.NewQuestionServiceClient(questionConn)
	sessionClient := sessionv1.NewSessionServiceClient(sessionConn)
	scoringClient := scoringv1.NewScoringServiceClient(scoringConn)

	connections := []*grpc.ClientConn{examConn, questionConn, sessionConn, scoringConn}

	return examClient, questionClient, sessionClient, scoringClient, connections, nil
}

// closeConnections menutup semua koneksi gRPC
func closeConnections(connections []*grpc.ClientConn) {
	for _, conn := range connections {
		conn.Close()
	}
}

// TestExamSessionPerformance menguji performa saat banyak siswa mengakses ujian secara bersamaan
func TestExamSessionPerformance(t *testing.T) {
	// Skip test jika running di CI environment atau dengan -short flag
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Konfigurasi pengujian
	config := TestConfig{
		ExamServiceAddr:     "localhost:50051",
		QuestionServiceAddr: "localhost:50052",
		SessionServiceAddr:  "localhost:50053",
		ScoringServiceAddr:  "localhost:50054",
		ConcurrentUsers:     100, // Simulasikan 100 siswa secara bersamaan
		RequestsPerUser:     10,  // Setiap siswa melakukan 10 request
	}

	// Buat koneksi ke service
	examClient, questionClient, sessionClient, scoringClient, connections, err := createServiceClients(config)
	if err != nil {
		t.Fatalf("Failed to create service clients: %v", err)
	}
	defer closeConnections(connections)

	// Persiapkan ujian untuk testing
	examID, err := prepareExam(t, examClient, questionClient)
	if err != nil {
		t.Fatalf("Failed to prepare exam: %v", err)
	}

	// Aktifkan ujian
	_, err = examClient.ActivateExam(context.Background(), &examv1.ActivateExamRequest{
		Id:       examID,
		ClassIds: []string{"class-1"},
	})
	if err != nil {
		t.Fatalf("Failed to activate exam: %v", err)
	}

	// Channel untuk mengumpulkan hasil
	resultChan := make(chan time.Duration, config.ConcurrentUsers*config.RequestsPerUser)
	errorChan := make(chan error, config.ConcurrentUsers*config.RequestsPerUser)

	// Tunggu semua goroutine selesai
	var wg sync.WaitGroup

	// Catat waktu mulai
	startTime := time.Now()

	// Simulasikan banyak siswa mengakses ujian secara bersamaan
	for i := 0; i < config.ConcurrentUsers; i++ {
		wg.Add(1)
		go func(studentID int) {
			defer wg.Done()

			// Mulai sesi ujian
			sessionID, err := startExamSession(context.Background(), sessionClient, examID, fmt.Sprintf("student-%d", studentID))
			if err != nil {
				errorChan <- fmt.Errorf("student %d failed to start session: %v", studentID, err)
				return
			}

			// Ambil pertanyaan ujian
			questions, err := getExamQuestions(context.Background(), questionClient, examID)
			if err != nil {
				errorChan <- fmt.Errorf("student %d failed to get questions: %v", studentID, err)
				return
			}

			// Jawab pertanyaan
			for j, question := range questions {
				if j >= config.RequestsPerUser {
					break
				}

				// Catat waktu sebelum request
				beforeRequest := time.Now()

				// Pilih jawaban (untuk simulasi)
				choiceIndex := j % len(question.Choices)
				selectedChoice := string('A' + choiceIndex)

				// Kirim jawaban
				_, err := sessionClient.SubmitAnswer(context.Background(), &sessionv1.SubmitAnswerRequest{
					SessionId:      sessionID,
					QuestionId:     question.Id,
					SelectedChoice: selectedChoice,
				})

				// Catat waktu setelah request
				latency := time.Since(beforeRequest)
				resultChan <- latency

				if err != nil {
					errorChan <- fmt.Errorf("student %d failed to submit answer: %v", studentID, err)
				}

				// Tambahkan jeda kecil untuk simulasi perilaku realistis
				time.Sleep(time.Millisecond * 10)
			}

			// Selesaikan sesi ujian
			_, err = sessionClient.FinishSession(context.Background(), &sessionv1.FinishSessionRequest{
				Id: sessionID,
			})
			if err != nil {
				errorChan <- fmt.Errorf("student %d failed to finish session: %v", studentID, err)
				return
			}

			// Hitung skor
			_, err = scoringClient.CalculateScore(context.Background(), &scoringv1.CalculateScoreRequest{
				SessionId: sessionID,
			})
			if err != nil {
				errorChan <- fmt.Errorf("student %d failed to calculate score: %v", studentID, err)
			}

		}(i)
	}

	// Tunggu semua goroutine selesai
	wg.Wait()
	close(resultChan)
	close(errorChan)

	// Hitung total durasi pengujian
	totalDuration := time.Since(startTime)

	// Aggregate results
	var successCount, errorCount int
	var totalLatency, minLatency, maxLatency time.Duration

	// Set minLatency ke nilai maksimum awal
	minLatency = time.Hour

	for latency := range resultChan {
		successCount++
		totalLatency += latency

		if latency < minLatency {
			minLatency = latency
		}
		if latency > maxLatency {
			maxLatency = latency
		}
	}

	// Hitung jumlah error
	for range errorChan {
		errorCount++
	}

	// Hitung metrik
	totalRequests := successCount + errorCount
	averageLatency := totalLatency / time.Duration(successCount)
	requestsPerSecond := float64(successCount) / totalDuration.Seconds()

	// Tampilkan hasil
	result := PerformanceResult{
		TotalRequests:      totalRequests,
		SuccessfulRequests: successCount,
		FailedRequests:     errorCount,
		TotalDuration:      totalDuration,
		AverageLatency:     averageLatency,
		MinLatency:         minLatency,
		MaxLatency:         maxLatency,
		RequestsPerSecond:  requestsPerSecond,
	}

	t.Logf("Performance Test Results:")
	t.Logf("  Total Requests:       %d", result.TotalRequests)
	t.Logf("  Successful Requests:  %d", result.SuccessfulRequests)
	t.Logf("  Failed Requests:      %d", result.FailedRequests)
	t.Logf("  Total Duration:       %v", result.TotalDuration)
	t.Logf("  Average Latency:      %v", result.AverageLatency)
	t.Logf("  Min Latency:          %v", result.MinLatency)
	t.Logf("  Max Latency:          %v", result.MaxLatency)
	t.Logf("  Requests Per Second:  %.2f", result.RequestsPerSecond)

	// Assertions untuk memastikan performa memenuhi ekspektasi
	successRate := float64(result.SuccessfulRequests) / float64(result.TotalRequests) * 100
	assert.GreaterOrEqual(t, successRate, 95.0, "Success rate should be at least 95%%")
	assert.LessOrEqual(t, result.AverageLatency, 100*time.Millisecond, "Average latency should be less than 100ms")
	assert.GreaterOrEqual(t, result.RequestsPerSecond, 100.0, "Should handle at least 100 requests per second")
}

// TestConcurrentExamSessions menguji kemampuan sistem menangani banyak sesi ujian secara bersamaan
func TestConcurrentExamSessions(t *testing.T) {
	// Skip test jika running di CI environment atau dengan -short flag
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Konfigurasi pengujian
	config := TestConfig{
		ExamServiceAddr:     "localhost:50051",
		QuestionServiceAddr: "localhost:50052",
		SessionServiceAddr:  "localhost:50053",
		ScoringServiceAddr:  "localhost:50054",
		ConcurrentUsers:     1000, // Simulasikan 1000 siswa secara bersamaan
		RequestsPerUser:     1,    // Setiap siswa memulai 1 sesi
	}

	// Buat koneksi ke service
	examClient, questionClient, sessionClient, _, connections, err := createServiceClients(config)
	if err != nil {
		t.Fatalf("Failed to create service clients: %v", err)
	}
	defer closeConnections(connections)

	// Persiapkan ujian untuk testing
	examID, err := prepareExam(t, examClient, questionClient)
	if err != nil {
		t.Fatalf("Failed to prepare exam: %v", err)
	}

	// Aktifkan ujian
	_, err = examClient.ActivateExam(context.Background(), &examv1.ActivateExamRequest{
		Id:       examID,
		ClassIds: []string{"class-1"},
	})
	if err != nil {
		t.Fatalf("Failed to activate exam: %v", err)
	}

	// Channel untuk mengumpulkan hasil
	sessionCreationTimes := make(chan time.Duration, config.ConcurrentUsers)
	errorChan := make(chan error, config.ConcurrentUsers)

	// Tunggu semua goroutine selesai
	var wg sync.WaitGroup

	// Catat waktu mulai
	startTime := time.Now()

	// Simulasikan banyak siswa memulai sesi ujian secara bersamaan
	for i := 0; i < config.ConcurrentUsers; i++ {
		wg.Add(1)
		go func(studentID int) {
			defer wg.Done()

			// Catat waktu sebelum request
			beforeRequest := time.Now()

			// Mulai sesi ujian
			_, err := sessionClient.StartSession(context.Background(), &sessionv1.StartSessionRequest{
				ExamId:    examID,
				StudentId: fmt.Sprintf("student-%d", studentID),
			})

			// Catat waktu setelah request
			latency := time.Since(beforeRequest)
			sessionCreationTimes <- latency

			if err != nil {
				errorChan <- err
			}
		}(i)
	}

	// Tunggu semua goroutine selesai
	wg.Wait()
	close(sessionCreationTimes)
	close(errorChan)

	// Hitung total durasi pengujian
	totalDuration := time.Since(startTime)

	// Aggregate results
	var successCount, errorCount int
	var totalLatency, minLatency, maxLatency time.Duration

	// Set minLatency ke nilai maksimum awal
	minLatency = time.Hour

	for latency := range sessionCreationTimes {
		successCount++
		totalLatency += latency

		if latency < minLatency {
			minLatency = latency
		}
		if latency > maxLatency {
			maxLatency = latency
		}
	}

	// Hitung jumlah error
	for range errorChan {
		errorCount++
	}

	// Hitung metrik
	totalRequests := successCount + errorCount
	averageLatency := totalLatency / time.Duration(successCount)
	requestsPerSecond := float64(successCount) / totalDuration.Seconds()

	// Tampilkan hasil
	t.Logf("Concurrent Sessions Test Results:")
	t.Logf("  Total Session Requests:       %d", totalRequests)
	t.Logf("  Successful Sessions Started:  %d", successCount)
	t.Logf("  Failed Session Requests:      %d", errorCount)
	t.Logf("  Total Duration:               %v", totalDuration)
	t.Logf("  Average Session Creation Time: %v", averageLatency)
	t.Logf("  Min Session Creation Time:     %v", minLatency)
	t.Logf("  Max Session Creation Time:     %v", maxLatency)
	t.Logf("  Sessions Started Per Second:   %.2f", requestsPerSecond)

	// Assertions untuk memastikan performa memenuhi ekspektasi
	successRate := float64(successCount) / float64(totalRequests) * 100
	assert.GreaterOrEqual(t, successRate, 95.0, "Session creation success rate should be at least 95%%")
	assert.LessOrEqual(t, averageLatency, 200*time.Millisecond, "Average session creation time should be less than 200ms")
	assert.GreaterOrEqual(t, requestsPerSecond, 50.0, "Should handle at least 50 session creations per second")
}

// Helper function untuk mempersiapkan ujian
func prepareExam(t *testing.T, examClient examv1.ExamServiceClient, questionClient questionv1.QuestionServiceClient) (string, error) {
	// Buat ujian
	examResp, err := examClient.CreateExam(context.Background(), &examv1.CreateExamRequest{
		Title:           "Performance Test Exam",
		Subject:         "Performance Testing",
		DurationMinutes: 60,
		TotalQuestions:  20,
		IsRandom:        true,
		TeacherId:       "teacher-performance-test",
		ClassIds:        []string{"class-1"},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create exam: %v", err)
	}

	examID := examResp.Id

	// Buat pertanyaan untuk ujian
	for i := 0; i < 20; i++ {
		choices := []*questionv1.Choice{
			{Text: "Pilihan A"},
			{Text: "Pilihan B"},
			{Text: "Pilihan C"},
			{Text: "Pilihan D"},
		}

		_, err := questionClient.CreateQuestion(context.Background(), &questionv1.CreateQuestionRequest{
			ExamId:        examID,
			QuestionText:  fmt.Sprintf("Pertanyaan performance test %d?", i+1),
			CorrectAnswer: "A", // Untuk simplifikasi, jawaban benar selalu A
			Choices:       choices,
		})
		if err != nil {
			return "", fmt.Errorf("failed to create question %d: %v", i+1, err)
		}
	}

	return examID, nil
}

// Helper function untuk memulai sesi ujian
func startExamSession(ctx context.Context, sessionClient sessionv1.SessionServiceClient, examID, studentID string) (string, error) {
	resp, err := sessionClient.StartSession(ctx, &sessionv1.StartSessionRequest{
		ExamId:    examID,
		StudentId: studentID,
	})
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

// Helper function untuk mendapatkan pertanyaan ujian
func getExamQuestions(ctx context.Context, questionClient questionv1.QuestionServiceClient, examID string) ([]*questionv1.Question, error) {
	resp, err := questionClient.GetExamQuestions(ctx, &questionv1.GetExamQuestionsRequest{
		ExamId:    examID,
		Randomize: true,
		Limit:     20,
	})
	if err != nil {
		return nil, err
	}
	return resp.Questions, nil
}
