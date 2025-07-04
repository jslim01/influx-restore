CONTAINER_NAME="influxdb"
TOKEN="admintoken123"
FLUX_FILE="/scripts/query.flux"

echo "📤 Exporting data to CSV... (12GB 저장 공간 필요)"

docker exec "$CONTAINER_NAME" sh -c "influx query --file $FLUX_FILE --token $TOKEN --raw > /restore/exported.csv"

echo "✅ Done!"