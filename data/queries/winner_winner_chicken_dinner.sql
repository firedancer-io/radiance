SELECT
    winner,
    count(),
    bar(count(), 1, 10000, 20)
FROM
    (
        SELECT
            slot,
            any(source) AS winner
        FROM
            (
                SELECT
                    source,
                    slot,
                    minIf(timestamp, type = 'completed') AS firstCompleted
                FROM slot_status
                WHERE toDate(timestamp) = today()
                GROUP BY
                    source,
                    slot
                ORDER BY
                    slot ASC,
                    firstCompleted ASC
                )
        GROUP BY slot
        ORDER BY slot ASC
        )
GROUP BY winner
ORDER BY count() DESC
