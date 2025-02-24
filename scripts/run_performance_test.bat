@REM scripts/run_performance_test.bat
@echo off

REM Start services
echo Starting services...
start /B go run scripts/start_services.go

REM Wait for services to start
echo Waiting for services to start...
timeout /t 5 /nobreak

REM Run performance test
echo Running performance test...
go test -v ./tests/unit -run TestExamSessionPerformance

REM Cleanup
echo Cleaning up...
taskkill /F /IM "main.exe" > nul 2>&1