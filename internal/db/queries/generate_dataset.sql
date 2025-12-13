COPY (
    SELECT id, "timestamp", "type",
           MAX(accelX), MAX(accelY), MAX(accelZ),
           MAX(gyroX), MAX(gyroY), MAX(gyroZ),
           MAX(qX), MAX(qY), MAX(qZ), MAX(qW),
           MAX(yaw), MAX(pitch), MAX(roll)
    FROM dataset
    WHERE id IN (%s)
      AND "type" IN (%s)
      AND "timestamp" BETWEEN %d AND %d
    GROUP BY id, "timestamp", "type"
    ORDER BY id, "timestamp"
    LIMIT 100
) TO STDOUT WITH CSV HEADER