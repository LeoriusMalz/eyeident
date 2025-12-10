CREATE TABLE IF NOT EXISTS users (
    id             TEXT        PRIMARY KEY,
    is_active      BOOLEAN     DEFAULT True,
    last_commit    TIMESTAMP,
    last_status    TEXT,
    is_enabled     BOOLEAN     DEFAULT True
);

CREATE TABLE IF NOT EXISTS raw_data (
    id          TEXT            NOT NULL,
    datetime    BIGINT          NOT NULL CHECK(datetime > 0),
    "type"      VARCHAR(15),

    accelX      REAL,
    accelY      REAL,
    accelZ      REAL,

    gyroX       REAL,
    gyroY       REAL,
    gyroZ       REAL,

    azimutX     NUMERIC(5, 2),
    azimutY     NUMERIC(5, 2),
    azimutZ     NUMERIC(5, 2),
    cos         NUMERIC(4, 3),

    qX          NUMERIC(4, 3),
    qY          NUMERIC(4, 3),
    qZ          NUMERIC(4, 3),
    qW          NUMERIC(4, 3),

    PRIMARY KEY (id, datetime)
);

CREATE TABLE IF NOT EXISTS dataset (
    id          TEXT            NOT NULL,
    datetime    BIGINT          NOT NULL CHECK(datetime > 0),
    "type"      VARCHAR(15),

    accelX      REAL,
    accelY      REAL,
    accelZ      REAL,

    gyroX       REAL,
    gyroY       REAL,
    gyroZ       REAL,

    azimutX     NUMERIC(5, 2),
    azimutY     NUMERIC(5, 2),
    azimutZ     NUMERIC(5, 2),
    cos         NUMERIC(4, 3),

    qX          NUMERIC(4, 3),
    qY          NUMERIC(4, 3),
    qZ          NUMERIC(4, 3),
    qW          NUMERIC(4, 3),

    PRIMARY KEY (id, datetime)
);