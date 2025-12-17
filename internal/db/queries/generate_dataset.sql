COPY (
    SELECT id, "timestamp", "type",
           MAX(accelX) AS accelX, MAX(accelY) AS accelY, MAX(accelZ) AS accelZ,
           MAX(gyroX) AS gyroX, MAX(gyroY) AS gyroY, MAX(gyroZ) AS gyroZ,
           MAX(qX) AS qX, MAX(qY) AS qY, MAX(qZ) AS qZ, MAX(qW) AS qW,
           MAX(yaw) AS yaw, MAX(pitch) AS pitch, MAX(roll) AS roll
    FROM dataset
    WHERE id IN (%s)
      AND "type" IN (%s)
      AND "timestamp" BETWEEN %d AND %d
    GROUP BY id, "timestamp", "type"
    ORDER BY id, "timestamp" DESC
    LIMIT %d
) TO STDOUT WITH CSV HEADER