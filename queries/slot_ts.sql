WITH CAST(timestamp, 'Nullable(DateTime)') AS ts
SELECT
    slot,
    minIf(ts, type = 'firstShredReceived') AS first,
    minIf(ts, type = 'completed') AS completed,
    minIf(ts, type = 'optimisticConfirmation') AS confirmed
FROM slot_status
GROUP BY slot
ORDER BY slot ASC
