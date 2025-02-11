@echo off
setlocal

REM Generate exam service
protoc --go_out=. --go_opt=paths=source_relative ^
    --go-grpc_out=. --go-grpc_opt=paths=source_relative ^
    api/proto/exam/v1/exam.proto

@REM REM Generate question service
@REM protoc --go_out=. --go_opt=paths=source_relative ^
@REM     --go-grpc_out=. --go-grpc_opt=paths=source_relative ^
@REM     api/proto/question/v1/question.proto
@REM
@REM REM Generate session service
@REM protoc --go_out=. --go_opt=paths=source_relative ^
@REM     --go-grpc_out=. --go-grpc_opt=paths=source_relative ^
@REM     api/proto/session/v1/session.proto
@REM
@REM REM Generate scoring service
@REM protoc --go_out=. --go_opt=paths=source_relative ^
@REM     --go-grpc_out=. --go-grpc_opt=paths=source_relative ^
@REM     api/proto/scoring/v1/scoring.proto

echo Protocol buffer code generation completed.
endlocal