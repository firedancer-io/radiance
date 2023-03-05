CREATE TABLE IF NOT EXISTS slot_status_queue
(
    slot                                UInt64,
    timestamp                           UInt64,
    delay                               UInt64,
    type                                Enum8(
        'unspecified' = 0,
        'firstShredReceived' = 1,
        'completed' = 2,
        'createdBank' = 3,
        'frozen' = 4,
        'dead' = 5,
        'optimisticConfirmation' = 6,
        'root' = 7
        ),

    source                              String,
    leader                              String,

    parent                              UInt64,

    "stats.num_transaction_entries"     Nullable(UInt64),
    "stats.num_successful_transactions" Nullable(UInt64),
    "stats.num_failed_transactions"     Nullable(UInt64),
    "stats.max_transactions_per_entry"  Nullable(UInt64),

    err                                 String
) ENGINE = Kafka()
      SETTINGS
          kafka_broker_list = '<snip>:30036',
          kafka_topic_list = 'certus.radiance.slot_status',
          kafka_group_name = 'heimdall-chdev1',
          kafka_format = 'Protobuf',
          kafka_max_block_size = 32,
          format_schema = 'network.proto:SlotStatus';

CREATE TABLE IF NOT EXISTS slot_status
(
    slot                                UInt64,
    timestamp                           DateTime64(3),
    delay                               UInt64,
    type                                Enum8(
        'unspecified' = 0,
        'firstShredReceived' = 1,
        'completed' = 2,
        'createdBank' = 3,
        'frozen' = 4,
        'dead' = 5,
        'optimisticConfirmation' = 6,
        'root' = 7
        ),


    source                              String,
    leader                              String,

    parent                              UInt64,

    "stats.num_transaction_entries"     Nullable(UInt64),
    "stats.num_successful_transactions" Nullable(UInt64),
    "stats.num_failed_transactions"     Nullable(UInt64),
    "stats.max_transactions_per_entry"  Nullable(UInt64),

    err                                 String
) ENGINE = MergeTree()
      PARTITION BY toDate(timestamp)
      ORDER BY (slot, type, leader);

CREATE MATERIALIZED VIEW IF NOT EXISTS slot_status_view TO slot_status
AS
SELECT fromUnixTimestamp64Milli(timestamp) AS timestamp,
       * EXCEPT ( timestamp )
FROM slot_status_queue;
