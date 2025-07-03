# influx-restore
influx restore

1. backup 폴더를 생성해야합니다.
2. backup 폴더에 2025-06-04.tar.gz를 위치시킵니다.

3. influxdb 실행
```bash
docker-compose up
```
4. 복구 스크립트 실행
```
chmod +x restore_script.sh
./restore_script.sh
```