CREATE TABLE IF NOT EXISTS users (
    id             TEXT        PRIMARY KEY,
    is_active      BOOLEAN     DEFAULT True,
    last_commit    TIMESTAMP,
    last_status    TEXT,
    is_enabled     BOOLEAN     DEFAULT True
);

CREATE TABLE IF NOT EXISTS raw_data (
    id          TEXT            NOT NULL,
    "timestamp" BIGINT          NOT NULL CHECK("timestamp" > 0),
    "type"      VARCHAR(15),

    accelX      REAL,
    accelY      REAL,
    accelZ      REAL,

    gyroX       REAL,
    gyroY       REAL,
    gyroZ       REAL,

    qX          NUMERIC(4, 3),
    qY          NUMERIC(4, 3),
    qZ          NUMERIC(4, 3),
    qW          NUMERIC(4, 3),

    yaw         NUMERIC(5, 2),
    pitch       NUMERIC(5, 2),
    roll        NUMERIC(5, 2)
);

CREATE TABLE IF NOT EXISTS dataset (
    id          TEXT            NOT NULL,
    "timestamp" BIGINT          NOT NULL CHECK("timestamp" > 0),
    "type"      VARCHAR(15),

    accelX      REAL,
    accelY      REAL,
    accelZ      REAL,

    gyroX       REAL,
    gyroY       REAL,
    gyroZ       REAL,

    qX          NUMERIC(4, 3),
    qY          NUMERIC(4, 3),
    qZ          NUMERIC(4, 3),
    qW          NUMERIC(4, 3),

    yaw         NUMERIC(5, 2),
    pitch       NUMERIC(5, 2),
    roll        NUMERIC(5, 2)
);