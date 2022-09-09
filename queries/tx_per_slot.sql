SELECT
    slot,
    max(assumeNotNull(stats.num_successful_transactions)) + max(assumeNotNull(stats.num_failed_transactions)) AS maxTxs
FROM slot_status
WHERE type = 'frozen'
GROUP BY slot
ORDER BY slot
