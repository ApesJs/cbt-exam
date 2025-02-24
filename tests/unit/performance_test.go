package unit

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	examv1 "github.com/ApesJs/cbt-exam/api/proto/exam/v1"
	questionv1 "github.com/ApesJs/cbt-exam/api/proto/question/v1"
	scoringv1 "github.com/ApesJs/cbt-exam/api/proto/scoring/v1"
	sessionv1 "github.com/ApesJs/cbt-exam/api/proto/session/v1"
)

// ResponseTimeMetrics menyimpan metrik response time untuk tiap API
type ResponseTimeMetrics struct {
	Name        string
	Count       int
	Total       time.Duration
	Min         time.Duration
	Max         time.Duration
	Average     time.Duration
	Percentiles map[int]time.Duration // p50, p90, p95, p99
	Samples     []time.Duration       // untuk menghitung persentil
}

// TestResponseTimeWithExistingData menguji response time untuk setiap API utama menggunakan data yang sudah ada
func TestResponseTimeWithExistingData(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping response time test in short mode")
	}

	// Konfigurasi test
	config := struct {
		ExamServiceAddr     string
		QuestionServiceAddr string
		SessionServiceAddr  string
		ScoringServiceAddr  string
		NumRequests         int
		Concurrency         int
		ExistingExamID      string // Menggunakan ID ujian yang sudah ada
	}{
		ExamServiceAddr:     "localhost:50051",
		QuestionServiceAddr: "localhost:50052",
		SessionServiceAddr:  "localhost:50053",
		ScoringServiceAddr:  "localhost:50054",
		NumRequests:         100,                                    // Jumlah request per API
		Concurrency:         10,                                     // Jumlah concurrent request
		ExistingExamID:      "33333333-3333-3333-3333-333333333333", // ID ujian yang sudah ada di database
	}

	// Buat koneksi ke service
	examClient, questionClient, sessionClient, scoringClient, connections, err := createServiceClients(config.ExamServiceAddr,
		config.QuestionServiceAddr, config.SessionServiceAddr, config.ScoringServiceAddr)
	if err != nil {
		t.Fatalf("Failed to create service clients: %v", err)
	}
	defer func() {
		for _, conn := range connections {
			conn.Close()
		}
	}()

	// Verifikasi bahwa ujian yang digunakan sudah aktif
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exam, err := examClient.GetExam(ctx, &examv1.GetExamRequest{
		Id: config.ExistingExamID,
	})
	if err != nil {
		t.Fatalf("Failed to get exam: %v", err)
	}

	if exam.Status.State != examv1.ExamState_EXAM_STATE_ACTIVE {
		// Jika ujian tidak aktif, aktivasi dulu
		t.Logf("Activating exam %s...", config.ExistingExamID)
		_, err := examClient.ActivateExam(ctx, &examv1.ActivateExamRequest{
			Id:       config.ExistingExamID,
			ClassIds: exam.ClassIds,
		})
		if err != nil {
			t.Fatalf("Failed to activate exam: %v", err)
		}
	}

	// Metrik response time untuk setiap API
	metrics := make(map[string]*ResponseTimeMetrics)
	metrics["StartSession"] = &ResponseTimeMetrics{
		Name:        "StartSession",
		Percentiles: make(map[int]time.Duration),
		Samples:     make([]time.Duration, 0, config.NumRequests),
	}
	metrics["GetExamQuestions"] = &ResponseTimeMetrics{
		Name:        "GetExamQuestions",
		Percentiles: make(map[int]time.Duration),
		Samples:     make([]time.Duration, 0, config.NumRequests),
	}
	metrics["SubmitAnswer"] = &ResponseTimeMetrics{
		Name:        "SubmitAnswer",
		Percentiles: make(map[int]time.Duration),
		Samples:     make([]time.Duration, 0, config.NumRequests),
	}
	metrics["FinishSession"] = &ResponseTimeMetrics{
		Name:        "FinishSession",
		Percentiles: make(map[int]time.Duration),
		Samples:     make([]time.Duration, 0, config.NumRequests),
	}
	metrics["CalculateScore"] = &ResponseTimeMetrics{
		Name:        "CalculateScore",
		Percentiles: make(map[int]time.Duration),
		Samples:     make([]time.Duration, 0, config.NumRequests),
	}

	// Test StartSession API
	t.Log("Testing StartSession API response time...")
	testStartSessionResponseTime(t, sessionClient, config.ExistingExamID, config.NumRequests, config.Concurrency, metrics["StartSession"])

	// Test GetExamQuestions API
	t.Log("Testing GetExamQuestions API response time...")
	testGetExamQuestionsResponseTime(t, questionClient, config.ExistingExamID, config.NumRequests, config.Concurrency, metrics["GetExamQuestions"])

	// Untuk SubmitAnswer, FinishSession, dan CalculateScore, perlu dibuat sesi terlebih dahulu
	t.Log("Creating sessions for answer submission tests...")
	sessionIDs := createTestSessions(t, sessionClient, config.ExistingExamID, config.NumRequests)

	// Test SubmitAnswer API
	t.Log("Testing SubmitAnswer API response time...")
	testSubmitAnswerResponseTime(t, sessionClient, questionClient, sessionIDs, config.ExistingExamID, config.NumRequests, config.Concurrency, metrics["SubmitAnswer"])

	// Test FinishSession API
	t.Log("Testing FinishSession API response time...")
	testFinishSessionResponseTime(t, sessionClient, sessionIDs, config.NumRequests, config.Concurrency, metrics["FinishSession"])

	// Test CalculateScore API
	t.Log("Testing CalculateScore API response time...")
	testCalculateScoreResponseTime(t, scoringClient, sessionIDs, config.NumRequests, config.Concurrency, metrics["CalculateScore"])

	// Tampilkan hasil dan validasi
	t.Log("\nResponse Time Test Results:")
	for _, m := range metrics {
		// Hitung persentil
		calculatePercentiles(m)

		t.Logf("API: %s", m.Name)
		t.Logf("  Total Requests:    %d", m.Count)
		t.Logf("  Min Response Time: %v", m.Min)
		t.Logf("  Max Response Time: %v", m.Max)
		t.Logf("  Avg Response Time: %v", m.Average)
		t.Logf("  P50 Response Time: %v", m.Percentiles[50])
		t.Logf("  P90 Response Time: %v", m.Percentiles[90])
		t.Logf("  P95 Response Time: %v", m.Percentiles[95])
		t.Logf("  P99 Response Time: %v", m.Percentiles[99])

		// Validasi response time terhadap SLA
		validateResponseTime(t, m)
	}
}

// Menghitung persentil dari sampel response time
func calculatePercentiles(metrics *ResponseTimeMetrics) {
	if len(metrics.Samples) == 0 {
		return
	}

	// Sort samples
	sort.Slice(metrics.Samples, func(i, j int) bool {
		return metrics.Samples[i] < metrics.Samples[j]
	})

	// Find min and max
	metrics.Min = metrics.Samples[0]
	metrics.Max = metrics.Samples[len(metrics.Samples)-1]

	// Calculate percentiles
	getPercentile := func(percent int) time.Duration {
		idx := int(float64(percent) / 100.0 * float64(len(metrics.Samples)-1))
		return metrics.Samples[idx]
	}

	metrics.Percentiles[50] = getPercentile(50)
	metrics.Percentiles[90] = getPercentile(90)
	metrics.Percentiles[95] = getPercentile(95)
	metrics.Percentiles[99] = getPercentile(99)
}

// Validasi response time terhadap SLA
func validateResponseTime(t *testing.T, metrics *ResponseTimeMetrics) {
	// Definisi SLA berdasarkan API
	slaLimits := map[string]time.Duration{
		"StartSession":     100 * time.Millisecond,
		"GetExamQuestions": 150 * time.Millisecond,
		"SubmitAnswer":     50 * time.Millisecond,
		"FinishSession":    100 * time.Millisecond,
		"CalculateScore":   200 * time.Millisecond,
	}

	// Validasi P95 response time
	slaLimit, exists := slaLimits[metrics.Name]
	if exists {
		assert.LessOrEqual(t, metrics.Percentiles[95], slaLimit,
			"P95 response time for %s exceeds SLA limit of %v", metrics.Name, slaLimit)
	}
}

// Test StartSession API response time
func testStartSessionResponseTime(t *testing.T, sessionClient sessionv1.SessionServiceClient, examID string, numRequests, concurrency int, metrics *ResponseTimeMetrics) {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Channel untuk throttling concurrency
	semaphore := make(chan struct{}, concurrency)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire

		go func(idx int) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release

			// Generate UUID untuk student ID (harus dalam format UUID)
			studentID := uuid.New().String()

			// Measure response time
			start := time.Now()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err := sessionClient.StartSession(ctx, &sessionv1.StartSessionRequest{
				ExamId:    examID,
				StudentId: studentID,
			})

			duration := time.Since(start)

			// Record metrics
			mutex.Lock()
			if err == nil {
				metrics.Count++
				metrics.Total += duration
				metrics.Samples = append(metrics.Samples, duration)
			} else {
				t.Logf("Error starting session %d: %v", idx, err)
			}
			mutex.Unlock()
		}(i)
	}

	wg.Wait()
	if metrics.Count > 0 {
		metrics.Average = metrics.Total / time.Duration(metrics.Count)
	}
}

// Test GetExamQuestions API response time
func testGetExamQuestionsResponseTime(t *testing.T, questionClient questionv1.QuestionServiceClient, examID string, numRequests, concurrency int, metrics *ResponseTimeMetrics) {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Channel untuk throttling concurrency
	semaphore := make(chan struct{}, concurrency)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire

		go func(idx int) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release

			// Measure response time
			start := time.Now()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err := questionClient.GetExamQuestions(ctx, &questionv1.GetExamQuestionsRequest{
				ExamId:    examID,
				Randomize: true,
				Limit:     20,
			})

			duration := time.Since(start)

			// Record metrics
			mutex.Lock()
			if err == nil {
				metrics.Count++
				metrics.Total += duration
				metrics.Samples = append(metrics.Samples, duration)
			} else {
				t.Logf("Error getting exam questions %d: %v", idx, err)
			}
			mutex.Unlock()
		}(i)
	}

	wg.Wait()
	if metrics.Count > 0 {
		metrics.Average = metrics.Total / time.Duration(metrics.Count)
	}
}

// Helper function untuk membuat sesi test
func createTestSessions(t *testing.T, sessionClient sessionv1.SessionServiceClient, examID string, numSessions int) []string {
	sessionIDs := make([]string, 0, numSessions)
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// Batasi jumlah goroutine yang berjalan bersamaan
	semaphore := make(chan struct{}, 10)

	for i := 0; i < numSessions; i++ {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire

		go func(idx int) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release

			studentID := uuid.New().String()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := sessionClient.StartSession(ctx, &sessionv1.StartSessionRequest{
				ExamId:    examID,
				StudentId: studentID,
			})

			if err != nil {
				t.Logf("Error creating test session %d: %v", idx, err)
				return
			}

			mutex.Lock()
			sessionIDs = append(sessionIDs, resp.Id)
			mutex.Unlock()
		}(i)
	}

	wg.Wait()
	return sessionIDs
}

// Test SubmitAnswer API response time
func testSubmitAnswerResponseTime(t *testing.T, sessionClient sessionv1.SessionServiceClient, questionClient questionv1.QuestionServiceClient,
	sessionIDs []string, examID string, numRequests, concurrency int, metrics *ResponseTimeMetrics) {

	if len(sessionIDs) == 0 {
		t.Log("No valid sessions to test SubmitAnswer API")
		return
	}

	// Get questions first
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	questionsResp, err := questionClient.GetExamQuestions(ctx, &questionv1.GetExamQuestionsRequest{
		ExamId:    examID,
		Randomize: false,
		Limit:     20,
	})
	if err != nil {
		t.Fatalf("Failed to get questions for submit answer test: %v", err)
		return
	}

	if len(questionsResp.Questions) == 0 {
		t.Fatalf("No questions available for submit answer test")
		return
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Channel untuk throttling concurrency
	semaphore := make(chan struct{}, concurrency)

	// Limit to the number of available sessions or requested tests
	testCount := min(numRequests, len(sessionIDs))

	for i := 0; i < testCount; i++ {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire

		go func(idx int) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release

			sessionID := sessionIDs[idx%len(sessionIDs)]
			questionIdx := idx % len(questionsResp.Questions)
			questionID := questionsResp.Questions[questionIdx].Id

			// Random choice A, B, C, or D
			choiceIdx := rand.Intn(4)
			selectedChoice := string(rune('A' + choiceIdx))

			// Measure response time
			start := time.Now()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err := sessionClient.SubmitAnswer(ctx, &sessionv1.SubmitAnswerRequest{
				SessionId:      sessionID,
				QuestionId:     questionID,
				SelectedChoice: selectedChoice,
			})

			duration := time.Since(start)

			// Record metrics
			mutex.Lock()
			if err == nil {
				metrics.Count++
				metrics.Total += duration
				metrics.Samples = append(metrics.Samples, duration)
			} else {
				t.Logf("Error submitting answer %d: %v", idx, err)
			}
			mutex.Unlock()
		}(i)
	}

	wg.Wait()
	if metrics.Count > 0 {
		metrics.Average = metrics.Total / time.Duration(metrics.Count)
	}
}

// Test FinishSession API response time
func testFinishSessionResponseTime(t *testing.T, sessionClient sessionv1.SessionServiceClient, sessionIDs []string, numRequests, concurrency int, metrics *ResponseTimeMetrics) {
	if len(sessionIDs) == 0 {
		t.Log("No valid sessions to test FinishSession API")
		return
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Channel untuk throttling concurrency
	semaphore := make(chan struct{}, concurrency)

	// Use only half of the sessions for finish test (we'll need some for calculate score)
	finishCount := min(numRequests, len(sessionIDs)/2)
	if finishCount == 0 {
		t.Log("Not enough sessions to test FinishSession API")
		return
	}

	finishSessionIDs := sessionIDs[:finishCount]

	for i := 0; i < len(finishSessionIDs); i++ {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire

		go func(idx int) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release

			sessionID := finishSessionIDs[idx]

			// Measure response time
			start := time.Now()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err := sessionClient.FinishSession(ctx, &sessionv1.FinishSessionRequest{
				Id: sessionID,
			})

			duration := time.Since(start)

			// Record metrics
			mutex.Lock()
			if err == nil {
				metrics.Count++
				metrics.Total += duration
				metrics.Samples = append(metrics.Samples, duration)
			} else {
				t.Logf("Error finishing session %d: %v", idx, err)
			}
			mutex.Unlock()
		}(i)
	}

	wg.Wait()
	if metrics.Count > 0 {
		metrics.Average = metrics.Total / time.Duration(metrics.Count)
	}
}

// Test CalculateScore API response time
func testCalculateScoreResponseTime(t *testing.T, scoringClient scoringv1.ScoringServiceClient, sessionIDs []string, numRequests, concurrency int, metrics *ResponseTimeMetrics) {
	if len(sessionIDs) == 0 {
		t.Log("No valid sessions to test CalculateScore API")
		return
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Channel untuk throttling concurrency
	semaphore := make(chan struct{}, concurrency)

	// Use the remaining half of sessions for calculate score
	calculateCount := min(numRequests, len(sessionIDs)/2)
	startIdx := len(sessionIDs) - calculateCount
	if startIdx < 0 || calculateCount == 0 {
		t.Log("Not enough sessions to test CalculateScore API")
		return
	}

	calculateSessionIDs := sessionIDs[startIdx:]

	for i := 0; i < len(calculateSessionIDs); i++ {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire

		go func(idx int) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release

			sessionID := calculateSessionIDs[idx]

			// Measure response time
			start := time.Now()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err := scoringClient.CalculateScore(ctx, &scoringv1.CalculateScoreRequest{
				SessionId: sessionID,
			})

			duration := time.Since(start)

			// Record metrics
			mutex.Lock()
			if err == nil {
				metrics.Count++
				metrics.Total += duration
				metrics.Samples = append(metrics.Samples, duration)
			} else {
				t.Logf("Error calculating score %d: %v", idx, err)
			}
			mutex.Unlock()
		}(i)
	}

	wg.Wait()
	if metrics.Count > 0 {
		metrics.Average = metrics.Total / time.Duration(metrics.Count)
	}
}

// Helper untuk menentukan nilai minimum dari dua int
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// createServiceClients membuat koneksi ke semua service
func createServiceClients(examAddr, questionAddr, sessionAddr, scoringAddr string) (
	examv1.ExamServiceClient,
	questionv1.QuestionServiceClient,
	sessionv1.SessionServiceClient,
	scoringv1.ScoringServiceClient,
	[]*grpc.ClientConn,
	error) {

	// Buat koneksi ke ExamService
	examConn, err := grpc.Dial(examAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to connect to exam service: %v", err)
	}

	// Buat koneksi ke QuestionService
	questionConn, err := grpc.Dial(questionAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		examConn.Close()
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to connect to question service: %v", err)
	}

	// Buat koneksi ke SessionService
	sessionConn, err := grpc.Dial(sessionAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		examConn.Close()
		questionConn.Close()
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to connect to session service: %v", err)
	}

	// Buat koneksi ke ScoringService
	scoringConn, err := grpc.Dial(scoringAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
