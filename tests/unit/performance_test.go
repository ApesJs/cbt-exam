package unit

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	examv1 "github.com/ApesJs/cbt-exam/api/proto/exam/v1"
	questionv1 "github.com/ApesJs/cbt-exam/api/proto/question/v1"
	scoringv1 "github.com/ApesJs/cbt-exam/api/proto/scoring/v1"
	sessionv1 "github.com/ApesJs/cbt-exam/api/proto/session/v1"
)

// TestConfig menyimpan konfigurasi untuk pengujian performa
type TestConfig struct {
	// Alamat service
	ExamServiceAddr     string
	QuestionServiceAddr string
	SessionServiceAddr  string
	ScoringServiceAddr  string

	// Parameter pengujian
	NumUsers            int           // Jumlah pengguna simulasi
	RampUpSeconds       int           // Periode ramp-up dalam detik
	TestDurationMinutes int           // Durasi pengujian dalam menit
	ExistingExamID      string        // ID ujian yang sudah ada
	ConcurrentUsers     int           // Jumlah maksimum pengguna concurrent
	RequestTimeout      time.Duration // Timeout untuk setiap request

	// Simulasi perilaku pengguna
	ThinkTimeMin  int // Waktu minimum berpikir (detik)
	ThinkTimeMax  int // Waktu maksimum berpikir (detik)
	AnswerTimeMin int // Waktu minimum menjawab soal (detik)
	AnswerTimeMax int // Waktu maksimum menjawab soal (detik)

	// Simulasi error
	NetworkFailureRate float32 // Kemungkinan kegagalan jaringan (0-1)
	SlowResponseRate   float32 // Kemungkinan respons lambat (0-1)

	// Output
	ExportMetrics    bool   // Ekspor metrik ke CSV
	OutputDir        string // Direktori output
	MonitorResources bool   // Monitor penggunaan resource
}

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
	ErrorCount  int                   // jumlah error
	mutex       sync.Mutex            // untuk thread safety
}

// ResourceMetrics menyimpan metrik penggunaan resource
type ResourceMetrics struct {
	Timestamp  time.Time
	CPUUsage   float64
	MemoryUsed uint64
	MemoryFree uint64
	GoRoutines int
}

// TestRealisticUserJourney mensimulasikan penggunaan aplikasi seperti kondisi sebenarnya
func TestRealisticUserJourney(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping realistic user journey test in short mode")
	}

	// Konfigurasi test
	config := TestConfig{
		ExamServiceAddr:     "localhost:50051",
		QuestionServiceAddr: "localhost:50052",
		SessionServiceAddr:  "localhost:50053",
		ScoringServiceAddr:  "localhost:50054",
		NumUsers:            100,
		RampUpSeconds:       60,
		TestDurationMinutes: 15,
		ExistingExamID:      "33333333-3333-3333-3333-333333333333",
		ConcurrentUsers:     30,
		RequestTimeout:      5 * time.Second,
		ThinkTimeMin:        5,
		ThinkTimeMax:        15,
		AnswerTimeMin:       20,
		AnswerTimeMax:       90,
		NetworkFailureRate:  0.01,
		SlowResponseRate:    0.05,
		ExportMetrics:       true,
		OutputDir:           "./test_results",
		MonitorResources:    true,
	}

	// Buat output directory jika belum ada
	if config.ExportMetrics || config.MonitorResources {
		os.MkdirAll(config.OutputDir, 0755)
	}

	// Buat koneksi ke service
	examClient, questionClient, sessionClient, scoringClient, connections, err := createServiceClients(
		config.ExamServiceAddr,
		config.QuestionServiceAddr,
		config.SessionServiceAddr,
		config.ScoringServiceAddr,
	)
	if err != nil {
		t.Fatalf("Failed to create service clients: %v", err)
	}
	defer func() {
		for _, conn := range connections {
			conn.Close()
		}
	}()

	// Verifikasi bahwa ujian yang digunakan sudah aktif
	examActive, examErr := verifyAndActivateExam(examClient, config.ExistingExamID)
	if examErr != nil {
		t.Fatalf("Failed to verify exam: %v", examErr)
	}
	if !examActive {
		t.Fatalf("Exam %s is not active and could not be activated", config.ExistingExamID)
	}

	// Metrik response time untuk setiap API
	metrics := make(map[string]*ResponseTimeMetrics)
	metrics["StartSession"] = &ResponseTimeMetrics{
		Name:        "StartSession",
		Percentiles: make(map[int]time.Duration),
		Samples:     make([]time.Duration, 0, config.NumUsers),
	}
	metrics["GetExamQuestions"] = &ResponseTimeMetrics{
		Name:        "GetExamQuestions",
		Percentiles: make(map[int]time.Duration),
		Samples:     make([]time.Duration, 0, config.NumUsers),
	}
	metrics["SubmitAnswer"] = &ResponseTimeMetrics{
		Name:        "SubmitAnswer",
		Percentiles: make(map[int]time.Duration),
		Samples:     make([]time.Duration, 0, config.NumUsers*5), // Estimasi 5 jawaban per user
	}
	metrics["FinishSession"] = &ResponseTimeMetrics{
		Name:        "FinishSession",
		Percentiles: make(map[int]time.Duration),
		Samples:     make([]time.Duration, 0, config.NumUsers),
	}
	metrics["CalculateScore"] = &ResponseTimeMetrics{
		Name:        "CalculateScore",
		Percentiles: make(map[int]time.Duration),
		Samples:     make([]time.Duration, 0, config.NumUsers),
	}
	metrics["TotalUserFlow"] = &ResponseTimeMetrics{
		Name:        "TotalUserFlow",
		Percentiles: make(map[int]time.Duration),
		Samples:     make([]time.Duration, 0, config.NumUsers),
	}

	// Channel untuk memantau resource
	var resourceMetrics []ResourceMetrics
	resourceChan := make(chan struct{})

	// Mulai monitoring resource jika diaktifkan
	if config.MonitorResources {
		go monitorResources(&resourceMetrics, resourceChan, 5*time.Second)
	}

	// Waktu mulai dan selesai pengujian
	startTime := time.Now()
	endTime := startTime.Add(time.Duration(config.TestDurationMinutes) * time.Minute)

	t.Logf("Starting realistic user journey test at %s", startTime.Format(time.RFC3339))
	t.Logf("Test will run until %s", endTime.Format(time.RFC3339))
	t.Logf("Simulating %d users with ramp-up period of %d seconds", config.NumUsers, config.RampUpSeconds)

	// Semaphore untuk membatasi jumlah pengguna bersamaan
	userSemaphore := make(chan struct{}, config.ConcurrentUsers)

	// WaitGroup untuk menunggu semua goroutine selesai
	var wg sync.WaitGroup

	// Interval waktu antar user login untuk mencapai ramp-up yang diinginkan
	userInterval := time.Duration(config.RampUpSeconds) * time.Second / time.Duration(config.NumUsers)

	// Jalankan simulasi untuk setiap pengguna
	for i := 0; i < config.NumUsers; i++ {
		wg.Add(1)

		go func(userID int) {
			defer wg.Done()

			// Simulasikan waktu login yang berbeda untuk setiap pengguna (ramp-up)
			loginDelay := time.Duration(userID) * userInterval
			time.Sleep(loginDelay)

			// Cek apakah sudah melebihi waktu pengujian
			if time.Now().After(endTime) {
				t.Logf("User %d skipped - test duration exceeded", userID)
				return
			}

			// Acquire semaphore untuk membatasi konkurensi
			userSemaphore <- struct{}{}
			defer func() { <-userSemaphore }()

			// Jalankan simulasi alur ujian lengkap untuk pengguna ini
			t.Logf("User %d starting exam flow", userID)
			userStartTime := time.Now()

			simulateCompleteExamFlow(
				examClient, questionClient, sessionClient, scoringClient,
				config, userID, metrics, endTime,
			)

			// Catat total waktu alur pengguna
			userDuration := time.Since(userStartTime)
			recordMetric(metrics["TotalUserFlow"], userDuration, nil)

			t.Logf("User %d completed exam flow in %v", userID, userDuration)
		}(i)
	}

	// Tunggu semua simulasi pengguna selesai
	wg.Wait()

	// Hentikan monitoring resource
	if config.MonitorResources {
		close(resourceChan)
	}

	// Hitung waktu pengujian
	testDuration := time.Since(startTime)

	// Hitung dan tampilkan hasil metrik
	t.Log("\nPerformance Test Results:")
	t.Logf("Total test duration: %v", testDuration)
	t.Logf("Total users simulated: %d", config.NumUsers)

	// Tabel hasil per API
	for _, m := range metrics {
		// Hitung persentil
		calculatePercentiles(m)

		if m.Count == 0 {
			t.Logf("API: %s - No data collected", m.Name)
			continue
		}

		t.Logf("API: %s", m.Name)
		t.Logf("  Total Requests:    %d", m.Count)
		t.Logf("  Error Rate:        %.2f%%", float64(m.ErrorCount)/float64(m.Count)*100)
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

	// Jika monitoring resource diaktifkan, tampilkan ringkasan
	if config.MonitorResources && len(resourceMetrics) > 0 {
		var totalCPU float64
		var peakMem uint64

		for _, rm := range resourceMetrics {
			totalCPU += rm.CPUUsage
			if rm.MemoryUsed > peakMem {
				peakMem = rm.MemoryUsed
			}
		}

		avgCPU := totalCPU / float64(len(resourceMetrics))
		t.Logf("\nResource Usage Summary:")
		t.Logf("  Average CPU Usage: %.2f%%", avgCPU)
		t.Logf("  Peak Memory Usage: %d MB", peakMem/1024/1024)
		t.Logf("  Sample Count: %d", len(resourceMetrics))
	}

	// Ekspor hasil ke CSV jika diaktifkan
	if config.ExportMetrics {
		exportResultsToCSV(metrics, config.OutputDir)

		// Jika monitoring resource diaktifkan, ekspor data resource juga
		if config.MonitorResources {
			exportResourceMetricsToCSV(resourceMetrics, config.OutputDir)
		}

		t.Logf("Test results exported to %s directory", config.OutputDir)
	}
}

// Simulasikan alur lengkap ujian untuk satu pengguna
func simulateCompleteExamFlow(
	examClient examv1.ExamServiceClient,
	questionClient questionv1.QuestionServiceClient,
	sessionClient sessionv1.SessionServiceClient,
	scoringClient scoringv1.ScoringServiceClient,
	config TestConfig,
	userID int,
	metrics map[string]*ResponseTimeMetrics,
	endTime time.Time,
) {
	// Generate UUID yang unik untuk student ID
	studentID := uuid.New().String()

	// Step 1: Mulai sesi ujian
	sessionID, err := startExamSession(sessionClient, config, studentID, metrics["StartSession"])
	if err != nil || sessionID == "" {
		log.Printf("User %d failed to start session: %v", userID, err)
		return
	}

	// Simulasikan jeda setelah memulai sesi (pengguna membaca instruksi)
	simulateUserThinkingTime(config.ThinkTimeMin, config.ThinkTimeMax)

	// Cek apakah sudah melebihi waktu pengujian
	if time.Now().After(endTime) {
		return
	}

	// Step 2: Ambil pertanyaan ujian
	questions, err := getExamQuestions(questionClient, config, metrics["GetExamQuestions"])
	if err != nil || len(questions) == 0 {
		log.Printf("User %d failed to get questions: %v", userID, err)
		return
	}

	// Simulasikan jeda setelah mendapatkan pertanyaan
	simulateUserThinkingTime(config.ThinkTimeMin/2, config.ThinkTimeMax/2)

	// Step 3: Jawab pertanyaan satu per satu
	answeredCount := 0
	for i, question := range questions {
		// Cek apakah sudah melebihi waktu pengujian
		if time.Now().After(endTime) {
			break
		}

		// Simulasikan pengguna membaca dan menjawab soal
		simulateUserThinkingTime(config.AnswerTimeMin, config.AnswerTimeMax)

		// Pilih jawaban acak
		selectedChoice := randomAnswer(question)

		// Kirim jawaban
		err := submitAnswer(sessionClient, config, sessionID, question.Id, selectedChoice, metrics["SubmitAnswer"])
		if err == nil {
			answeredCount++
		}

		// Simulasikan kemungkinan pengguna tidak menjawab semua soal
		// Semakin banyak soal yang dijawab, semakin besar kemungkinan untuk berhenti
		if i > len(questions)/3 && rand.Float32() < (float32(i)/float32(len(questions)))*0.2 {
			break
		}
	}

	// Simulasikan jeda sebelum menyelesaikan ujian
	simulateUserThinkingTime(config.ThinkTimeMin, config.ThinkTimeMax)

	// Cek apakah sudah melebihi waktu pengujian
	if time.Now().After(endTime) {
		return
	}

	// Step 4: Selesaikan sesi ujian
	err = finishExamSession(sessionClient, config, sessionID, metrics["FinishSession"])
	if err != nil {
		log.Printf("User %d failed to finish session: %v", userID, err)
		// Tetap lanjutkan ke perhitungan skor meskipun ada error
	}

	// Simulasikan jeda sebelum melihat skor
	simulateUserThinkingTime(config.ThinkTimeMin/2, config.ThinkTimeMax/2)

	// Cek apakah sudah melebihi waktu pengujian
	if time.Now().After(endTime) {
		return
	}

	// Step 5: Hitung skor ujian
	err = calculateExamScore(scoringClient, config, sessionID, metrics["CalculateScore"])
	if err != nil {
		log.Printf("User %d failed to calculate score: %v", userID, err)
	}
}

// Fungsi untuk memulai sesi ujian
func startExamSession(
	client sessionv1.SessionServiceClient,
	config TestConfig,
	studentID string,
	metrics *ResponseTimeMetrics,
) (string, error) {
	// Simulasikan kemungkinan kegagalan jaringan
	if simulateNetworkFailure(config.NetworkFailureRate) {
		recordMetric(metrics, 0, status.Error(codes.Unavailable, "simulated network failure"))
		return "", fmt.Errorf("network failure")
	}

	// Simulasikan kemungkinan respons lambat
	if simulateSlowResponse(config.SlowResponseRate) {
		time.Sleep(2 * time.Second)
	}

	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), config.RequestTimeout)
	defer cancel()

	resp, err := client.StartSession(ctx, &sessionv1.StartSessionRequest{
		ExamId:    config.ExistingExamID,
		StudentId: studentID,
	})

	duration := time.Since(start)
	recordMetric(metrics, duration, err)

	if err != nil {
		return "", err
	}

	return resp.Id, nil
}

// Fungsi untuk mendapatkan pertanyaan ujian
func getExamQuestions(
	client questionv1.QuestionServiceClient,
	config TestConfig,
	metrics *ResponseTimeMetrics,
) ([]*questionv1.Question, error) {
	// Simulasikan kemungkinan kegagalan jaringan
	if simulateNetworkFailure(config.NetworkFailureRate) {
		recordMetric(metrics, 0, status.Error(codes.Unavailable, "simulated network failure"))
		return nil, fmt.Errorf("network failure")
	}

	// Simulasikan kemungkinan respons lambat
	if simulateSlowResponse(config.SlowResponseRate) {
		time.Sleep(3 * time.Second)
	}

	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), config.RequestTimeout)
	defer cancel()

	resp, err := client.GetExamQuestions(ctx, &questionv1.GetExamQuestionsRequest{
		ExamId:    config.ExistingExamID,
		Randomize: true,
		Limit:     20,
	})

	duration := time.Since(start)
	recordMetric(metrics, duration, err)

	if err != nil {
		return nil, err
	}

	return resp.Questions, nil
}

// Fungsi untuk mengirim jawaban
func submitAnswer(
	client sessionv1.SessionServiceClient,
	config TestConfig,
	sessionID string,
	questionID string,
	selectedChoice string,
	metrics *ResponseTimeMetrics,
) error {
	// Simulasikan kemungkinan kegagalan jaringan
	if simulateNetworkFailure(config.NetworkFailureRate) {
		recordMetric(metrics, 0, status.Error(codes.Unavailable, "simulated network failure"))
		return fmt.Errorf("network failure")
	}

	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), config.RequestTimeout)
	defer cancel()

	_, err := client.SubmitAnswer(ctx, &sessionv1.SubmitAnswerRequest{
		SessionId:      sessionID,
		QuestionId:     questionID,
		SelectedChoice: selectedChoice,
	})

	duration := time.Since(start)
	recordMetric(metrics, duration, err)

	return err
}

// Fungsi untuk menyelesaikan sesi ujian
func finishExamSession(
	client sessionv1.SessionServiceClient,
	config TestConfig,
	sessionID string,
	metrics *ResponseTimeMetrics,
) error {
	// Simulasikan kemungkinan kegagalan jaringan
	if simulateNetworkFailure(config.NetworkFailureRate) {
		recordMetric(metrics, 0, status.Error(codes.Unavailable, "simulated network failure"))
		return fmt.Errorf("network failure")
	}

	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), config.RequestTimeout)
	defer cancel()

	_, err := client.FinishSession(ctx, &sessionv1.FinishSessionRequest{
		Id: sessionID,
	})

	duration := time.Since(start)
	recordMetric(metrics, duration, err)

	return err
}

// Fungsi untuk menghitung skor ujian
func calculateExamScore(
	client scoringv1.ScoringServiceClient,
	config TestConfig,
	sessionID string,
	metrics *ResponseTimeMetrics,
) error {
	// Simulasikan kemungkinan kegagalan jaringan
	if simulateNetworkFailure(config.NetworkFailureRate) {
		recordMetric(metrics, 0, status.Error(codes.Unavailable, "simulated network failure"))
		return fmt.Errorf("network failure")
	}

	// Simulasikan kemungkinan respons lambat
	if simulateSlowResponse(config.SlowResponseRate) {
		time.Sleep(2 * time.Second)
	}

	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), config.RequestTimeout)
	defer cancel()

	_, err := client.CalculateScore(ctx, &scoringv1.CalculateScoreRequest{
		SessionId: sessionID,
	})

	duration := time.Since(start)
	recordMetric(metrics, duration, err)

	return err
}

// Fungsi untuk memverifikasi dan mengaktifkan ujian jika perlu
func verifyAndActivateExam(client examv1.ExamServiceClient, examID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exam, err := client.GetExam(ctx, &examv1.GetExamRequest{
		Id: examID,
	})
	if err != nil {
		return false, err
	}

	if exam.Status.State == examv1.ExamState_EXAM_STATE_ACTIVE {
		return true, nil
	}

	// Jika ujian belum aktif, coba aktifkan
	if exam.Status.State == examv1.ExamState_EXAM_STATE_CREATED {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := client.ActivateExam(ctx, &examv1.ActivateExamRequest{
			Id:       examID,
			ClassIds: exam.ClassIds,
		})
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

// Fungsi untuk mencatat metrik secara thread-safe
func recordMetric(metrics *ResponseTimeMetrics, duration time.Duration, err error) {
	metrics.mutex.Lock()
	defer metrics.mutex.Unlock()

	if err != nil {
		metrics.ErrorCount++
		return
	}

	metrics.Count++
	metrics.Total += duration
	metrics.Samples = append(metrics.Samples, duration)

	// Update min dan max (hanya untuk sample yang valid)
	if metrics.Count == 1 || duration < metrics.Min {
		metrics.Min = duration
	}
	if metrics.Count == 1 || duration > metrics.Max {
		metrics.Max = duration
	}

	// Update average on-the-fly
	metrics.Average = metrics.Total / time.Duration(metrics.Count)
}

// Fungsi untuk mensimulasikan waktu berpikir pengguna
func simulateUserThinkingTime(minSeconds, maxSeconds int) {
	if minSeconds <= 0 || maxSeconds <= 0 || minSeconds > maxSeconds {
		return
	}

	thinkingTime := rand.Intn(maxSeconds-minSeconds+1) + minSeconds
	time.Sleep(time.Duration(thinkingTime) * time.Second)
}

// Fungsi untuk memilih jawaban acak
func randomAnswer(question *questionv1.Question) string {
	if len(question.Choices) == 0 {
		return "A" // Default jika tidak ada pilihan
	}

	choiceIdx := rand.Intn(len(question.Choices))
	return string(rune('A' + choiceIdx))
}

// Fungsi untuk mensimulasikan kegagalan jaringan
func simulateNetworkFailure(failureRate float32) bool {
	return rand.Float32() < failureRate
}

// Fungsi untuk mensimulasikan respons lambat
func simulateSlowResponse(slowRate float32) bool {
	return rand.Float32() < slowRate
}

// Menghitung persentil dari sampel response time
func calculatePercentiles(metrics *ResponseTimeMetrics) {
	metrics.mutex.Lock()
	defer metrics.mutex.Unlock()

	if len(metrics.Samples) == 0 {
		return
	}

	// Sort samples
	sort.Slice(metrics.Samples, func(i, j int) bool {
		return metrics.Samples[i] < metrics.Samples[j]
	})

	// Find min and max (meskipun sudah diupdate on-the-fly, pastikan lagi)
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
		"StartSession":     200 * time.Millisecond,
		"GetExamQuestions": 250 * time.Millisecond,
		"SubmitAnswer":     100 * time.Millisecond,
		"FinishSession":    200 * time.Millisecond,
		"CalculateScore":   300 * time.Millisecond,
		"TotalUserFlow":    10 * time.Minute, // Alur total per pengguna
	}

	// Validasi P95 response time
	slaLimit, exists := slaLimits[metrics.Name]
	if exists {
		assert.LessOrEqual(t, metrics.Percentiles[95], slaLimit,
			"P95 response time for %s exceeds SLA limit of %v", metrics.Name, slaLimit)
	}
}

// Fungsi untuk memantau penggunaan resource
func monitorResources(metrics *[]ResourceMetrics, done chan struct{}, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			// Get CPU usage
			cpuPercent, err := cpu.Percent(0, false)
			cpuUsage := 0.0
			if err == nil && len(cpuPercent) > 0 {
				cpuUsage = cpuPercent[0]
			}

			// Get memory stats
			var memUsed, memFree uint64
			vmStat, err := mem.VirtualMemory()
			if err == nil {
				memUsed = vmStat.Used
				memFree = vmStat.Free
			}

			// Get number of goroutines
			numGoroutines := runtime.NumGoroutine()

			// Record metrics
			*metrics = append(*metrics, ResourceMetrics{
				Timestamp:  time.Now(),
				CPUUsage:   cpuUsage,
				MemoryUsed: memUsed,
				MemoryFree: memFree,
				GoRoutines: numGoroutines,
			})
		}
	}
}

// Fungsi untuk mengekspor hasil ke CSV
func exportResultsToCSV(metrics map[string]*ResponseTimeMetrics, outputDir string) {
	// Ekspor ringkasan metrik
	summaryFile, err := os.Create(fmt.Sprintf("%s/summary_metrics.csv", outputDir))
	if err != nil {
		log.Printf("Failed to create summary file: %v", err)
		return
	}
	defer summaryFile.Close()

	summaryWriter := csv.NewWriter(summaryFile)
	defer summaryWriter.Flush()

	// Tulis header
	summaryWriter.Write([]string{
		"API", "Count", "ErrorCount", "ErrorRate", "Min (ms)", "Max (ms)",
		"Avg (ms)", "P50 (ms)", "P90 (ms)", "P95 (ms)", "P99 (ms)",
	})

	// Tulis data untuk setiap API
	for _, m := range metrics {
		if m.Count == 0 {
			continue
		}

		errorRate := float64(0)
		if m.Count > 0 {
			errorRate = float64(m.ErrorCount) / float64(m.Count+m.ErrorCount) * 100
		}

		summaryWriter.Write([]string{
			m.Name,
			fmt.Sprintf("%d", m.Count),
			fmt.Sprintf("%d", m.ErrorCount),
			fmt.Sprintf("%.2f%%", errorRate),
			fmt.Sprintf("%.2f", float64(m.Min)/float64(time.Millisecond)),
			fmt.Sprintf("%.2f", float64(m.Max)/float64(time.Millisecond)),
			fmt.Sprintf("%.2f", float64(m.Average)/float64(time.Millisecond)),
			fmt.Sprintf("%.2f", float64(m.Percentiles[50])/float64(time.Millisecond)),
			fmt.Sprintf("%.2f", float64(m.Percentiles[90])/float64(time.Millisecond)),
			fmt.Sprintf("%.2f", float64(m.Percentiles[95])/float64(time.Millisecond)),
			fmt.Sprintf("%.2f", float64(m.Percentiles[99])/float64(time.Millisecond)),
		})
	}

	// Ekspor data sampel detail untuk setiap API
	for _, m := range metrics {
		if len(m.Samples) == 0 {
			continue
		}

		detailFile, err := os.Create(fmt.Sprintf("%s/detail_%s.csv", outputDir, m.Name))
		if err != nil {
			log.Printf("Failed to create detail file for %s: %v", m.Name, err)
			continue
		}

		detailWriter := csv.NewWriter(detailFile)
		detailWriter.Write([]string{"Sample", "Duration (ms)"})

		for i, sample := range m.Samples {
			detailWriter.Write([]string{
				fmt.Sprintf("%d", i+1),
				fmt.Sprintf("%.2f", float64(sample)/float64(time.Millisecond)),
			})
		}

		detailWriter.Flush()
		detailFile.Close()
	}
}

// Fungsi untuk mengekspor data resource ke CSV
func exportResourceMetricsToCSV(metrics []ResourceMetrics, outputDir string) {
	resourceFile, err := os.Create(fmt.Sprintf("%s/resource_metrics.csv", outputDir))
	if err != nil {
		log.Printf("Failed to create resource metrics file: %v", err)
		return
	}
	defer resourceFile.Close()

	resourceWriter := csv.NewWriter(resourceFile)
	defer resourceWriter.Flush()

	// Tulis header
	resourceWriter.Write([]string{
		"Timestamp", "CPU Usage (%)", "Memory Used (MB)", "Memory Free (MB)", "Goroutines",
	})

	// Tulis data
	for _, rm := range metrics {
		resourceWriter.Write([]string{
			rm.Timestamp.Format(time.RFC3339),
			fmt.Sprintf("%.2f", rm.CPUUsage),
			fmt.Sprintf("%.2f", float64(rm.MemoryUsed)/1024/1024),
			fmt.Sprintf("%.2f", float64(rm.MemoryFree)/1024/1024),
			fmt.Sprintf("%d", rm.GoRoutines),
		})
	}
}

// createServiceClients membuat koneksi ke semua service
func createServiceClients(examAddr, questionAddr, sessionAddr, scoringAddr string) (
	examv1.ExamServiceClient,
	questionv1.QuestionServiceClient,
	sessionv1.SessionServiceClient,
	scoringv1.ScoringServiceClient,
	[]*grpc.ClientConn,
	error) {

	// Opsi koneksi gRPC
	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	}

	// Buat koneksi ke ExamService
	examConn, err := grpc.Dial(examAddr, dialOptions...)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to connect to exam service: %v", err)
	}

	// Buat koneksi ke QuestionService
	questionConn, err := grpc.Dial(questionAddr, dialOptions...)
	if err != nil {
		examConn.Close()
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to connect to question service: %v", err)
	}

	// Buat koneksi ke SessionService
	sessionConn, err := grpc.Dial(sessionAddr, dialOptions...)
	if err != nil {
		examConn.Close()
		questionConn.Close()
		return nil, nil, nil, nil, nil, fmt.Errorf("failed to connect to session service: %v", err)
	}

	// Buat koneksi ke ScoringService
	scoringConn, err := grpc.Dial(scoringAddr, dialOptions...)
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
