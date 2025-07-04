CONTAINER_NAME="influxdb"
CONTAINER_LP_PATH="/restore/out_shifted.lp"
ORG_NAME="prod"
BUCKET_NAME="new_ai_data"
TOKEN="admintoken123"


echo "📤 Writing data to InfluxDB..."

docker exec "$CONTAINER_NAME" influx write \
  --org "$ORG_NAME" \
  --bucket "$BUCKET_NAME" \
  --token "$TOKEN" \
  --file "$CONTAINER_LP_PATH"

echo "✅ Done!"