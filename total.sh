#!/bin/bash

echo "🖥️ 5월 1일 ~ 5월 30일 데이터를 복구하고 5월 30일을 현재시간으로 당깁니다.. m1 기준 약 18~20분 소요"

set -e  # 에러 발생 시 스크립트 중단

START_TIME=$(date +%s)

echo "🚀 Docker Compose 시작..."
docker-compose up -d
echo "⏳ InfluxDB 준비 대기 중..."

# InfluxDB 헬스체크 루프 (서비스 이름과 포트는 필요에 따라 수정)
SERVICE_NAME="influxdb"
MAX_WAIT=60
WAIT_TIME=0

# InfluxDB 컨테이너에서 헬스 API 응답 확인 (8086 포트 기준)
until docker-compose exec "$SERVICE_NAME" curl -s http://localhost:8086/health | grep -q '"status":"pass"'; do
    sleep 2
    WAIT_TIME=$((WAIT_TIME + 2))
    echo "⌛ 기다리는 중... (${WAIT_TIME}s / ${MAX_WAIT}s)"
    if [ "$WAIT_TIME" -ge "$MAX_WAIT" ]; then
        echo "❌ InfluxDB 준비 시간 초과"
        exit 1
    fi
done

echo "✅ InfluxDB 준비 완료 - 안정화를 위해 5초 추가 대기..."
sleep 5

echo "🚀 Step 1: 백업 데이터 복구 실행..."
./task_1_restore_script.sh
echo "✅ 복구 완료"

echo "🚀 Step 2: 특정 구간 데이터 추출 스크립트 실행..."
./task_2_export_csv.sh
echo "✅ 추출 완료"

echo "🚀 Step 3: 시간 조정 스크립트 실행..."
./task_3_time_shift.sh
echo "✅ 시간 조정 완료"

echo "🚀 Step 4: 데이터 쓰기 스크립트 실행..."
./task_4_write_csv.sh
echo "✅ 데이터 쓰기 완료"

END_TIME=$(date +%s)
ELAPSED=$((END_TIME - START_TIME))

# 실행 시간 출력 (시:분:초 형식)
HOURS=$((ELAPSED / 3600))
MINUTES=$(((ELAPSED % 3600) / 60))
SECONDS=$((ELAPSED % 60))

printf "🎉 전체 작업 완료! 총 소요 시간: %02d:%02d:%02d\n" "$HOURS" "$MINUTES" "$SECONDS"
