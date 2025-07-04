from(bucket: "ai_data") |> range(start: 2025-05-01T00:00:00Z, stop: 2025-05-30T00:00:00Z) |> filter(fn: (r) => r.topic =~ /dev/)
