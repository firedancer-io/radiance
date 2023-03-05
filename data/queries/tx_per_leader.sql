select leader, avg(stats.num_successful_transactions + assumeNotNull(stats.num_failed_transactions)) as avgTxs
from slot_status
group by leader
order by avgTxs desc
