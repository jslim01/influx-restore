import csv
from datetime import datetime, timezone
from tqdm import tqdm

INPUT_CSV = "tmp_restore/exported.csv"
OUTPUT_LP = "tmp_restore/out_shifted.lp"
ORIGINAL_END_TIME_STR = "2025-05-30T00:00:00Z"

# 기준 시점 계산
ORIGINAL_END_TIME = datetime.fromisoformat(ORIGINAL_END_TIME_STR.replace("Z", "+00:00"))
NOW = datetime.utcnow().replace(tzinfo=timezone.utc)
TIME_DIFF = NOW - ORIGINAL_END_TIME

def escape(value: str) -> str:
    return value.replace(" ", "\\ ").replace(",", "\\,").replace("=", "\\=")

def is_numeric(value: str) -> bool:
    try:
        float(value)
        return True
    except ValueError:
        return False

def is_header_row(row: list[str]) -> bool:
    return {"_time", "_value", "_field"}.issubset(set(row))

record_count = 0
skipped_count = 0
current_headers = []

with open(INPUT_CSV, newline='', encoding='utf-8') as infile, open(OUTPUT_LP, "w", encoding='utf-8') as outfile:
    reader = csv.reader(infile)
    progress = tqdm(reader, desc="변환 중")

    for row in progress:
        if not row or row[0].startswith("#"):
            continue

        # 새로운 헤더 감지
        if is_header_row(row):
            current_headers = row
            continue

        if not current_headers:
            continue  # 헤더 없으면 건너뜀

        if len(row) != len(current_headers):
            skipped_count += 1
            continue

        record = dict(zip(current_headers, row))

        try:
            measurement = record.get("_measurement", "")
            field_name = record.get("_field", "")
            field_value = record.get("_value", "")
            time_str = record.get("_time", "")

            if not measurement or not field_name or not time_str:
                skipped_count += 1
                continue

            # 시간 이동
            timestamp = datetime.fromisoformat(time_str.replace("Z", "+00:00"))
            new_timestamp = timestamp + TIME_DIFF
            timestamp_ns = int(new_timestamp.timestamp() * 1e9)

            # 태그 구성
            tags = []
            for key, val in record.items():
                if key.startswith("_") or key in {"result", "table"}:
                    continue
                if val and val != "None":
                    tags.append(f"{escape(key)}={escape(val)}")
            tag_str = ",".join(tags)

            # 필드 값
            if not is_numeric(field_value):
                field_value = f'"{field_value}"'

            # 라인 작성
            line = (
                f"{measurement},{tag_str} {field_name}={field_value} {timestamp_ns}"
                if tag_str else
                f"{measurement} {field_name}={field_value} {timestamp_ns}"
            )

            outfile.write(line + "\n")
            record_count += 1

        except Exception as e:
            skipped_count += 1
            continue

print(f"🎯 변환 완료: {record_count}행 처리됨, {skipped_count}행 건너뜀")
