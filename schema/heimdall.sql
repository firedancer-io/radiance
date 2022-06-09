CREATE TABLE heimdall_queue
(
    bankSlot                                       UInt64,
    bankID                                         UInt64,
    bankParentHash                                 String,
    feePayer                                       String,
    signature                                      String,
    program                                        String,
    "timings.serialize_us"                         UInt64,
    "timings.create_vm_us"                         UInt64,
    "timings.execute_us"                           UInt64,
    "timings.deserialize_us"                       UInt64,
    "timings.get_or_create_executor_us"            UInt64,
    "timings.changed_account_count"                UInt64,
    "timings.total_account_count"                  UInt64,
    "timings.total_data_size"                      UInt64,
    "timings.data_size_changed"                    UInt64,
    "timings.create_executor_register_syscalls_us" UInt64,
    "timings.create_executor_load_elf_us"          UInt64,
    "timings.create_executor_verify_code_us"       UInt64,
    "timings.create_executor_jit_compile_us"       UInt64
) ENGINE = Kafka('<...>:30036', 'certus.radiance.heimdall', 'heimdall-chdev1',
           'Protobuf')
SETTINGS format_schema = 'heimdall.proto:Observation';

CREATE TABLE IF NOT EXISTS heimdall
(
    bankSlot                                       UInt64,
    bankID                                         UInt64,
    bankParentHash                                 String,
    feePayer                                       String,
    signature                                      String,
    program                                        String,
    "timings.serialize_us"                         UInt64,
    "timings.create_vm_us"                         UInt64,
    "timings.execute_us"                           UInt64,
    "timings.deserialize_us"                       UInt64,
    "timings.get_or_create_executor_us"            UInt64,
    "timings.changed_account_count"                UInt64,
    "timings.total_account_count"                  UInt64,
    "timings.total_data_size"                      UInt64,
    "timings.data_size_changed"                    UInt64,
    "timings.create_executor_register_syscalls_us" UInt64,
    "timings.create_executor_load_elf_us"          UInt64,
    "timings.create_executor_verify_code_us"       UInt64,
    "timings.create_executor_jit_compile_us"       UInt64
) ENGINE = MergeTree() ORDER BY bankSlot;

CREATE MATERIALIZED VIEW IF NOT EXISTS heimdall_view TO heimdall
AS
SELECT *
FROM heimdall_queue;
