SELECT slot,
       timestamp,
       runningDifference(toUnixTimestamp64Milli(timestamp)) AS deltaMs,
       bar(deltaMs, 0, 10000, 10)                           AS deltaMsB,
       type
FROM (
         SELECT *
         FROM slot_status
         WHERE (slot = 138604605)
           AND (source = 'val4.ffm1')
         ORDER BY timestamp ASC,
                  type ASC
         )
