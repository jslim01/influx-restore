# influx-time-switch
5월 1일 ~ 30일 까지의 데이터를 추출하고 5월 30일을 기준으로 한달치 타임시리즈를 현재 시간으로 당기는 스크립트입니다.

### 사전 조건
1. backup 폴더를 프로젝트 루트에 생성해야합니다.
2. backup 폴더에 2025-06-04.tar.gz를 위치시킵니다. (nas => 데이터 폴더에 저장되어 있습니다.)
3. uv 설치 및 venv 가상환경이 구성되어 있어야합니다.

```
uv venv
source .venv/bin/activate
uv pip install tqdm
```

### 실행
스크립트 실행 (한달치 시간 스위치: m1 기준 약 18~20분 소요)

```
chmod +x *.sh
./total.sh
```