INSERT INTO dataset (id, "timestamp", "type",
                     accelX, accelY, accelZ,
                     gyroX, gyroY, gyroZ,
                     qX, qY, qZ, qW,
                     yaw, pitch, roll)
SELECT
    id, "timestamp", "type",
    accelX, accelY, accelZ,
    gyroX, gyroY, gyroZ,
    qX, qY, qZ, qW,
    yaw, pitch, roll
FROM raw_data
ORDER BY "timestamp", id;
