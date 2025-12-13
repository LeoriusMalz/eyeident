SELECT id, "timestamp", "type",
       MAX(accelX), MAX(accelY), MAX(accelZ),
       MAX(gyroX), MAX(gyroY), MAX(gyroZ),
       MAX(qX), MAX(qY), MAX(qZ), MAX(qW),
       MAX(yaw), MAX(pitch), MAX(roll) FROM dataset
WHERE id = ANY($1)
    AND "type" = ANY($2)
    AND "timestamp" BETWEEN $3 AND $4
GROUP BY id, "timestamp", "type"
ORDER BY id, "timestamp"
LIMIT 200;
