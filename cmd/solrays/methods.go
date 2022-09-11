package main

var (
	// https://docs.solana.com/developing/clients/jsonrpc-api
	methodWhitelist = []string{
		// Stable
		"getAccountInfo",
		"getBalance",
		"getBlock",
		"getBlockHeight",
		"getBlockProduction",
		"getBlockCommitment",
		"getBlocks",
		"getBlocksWithLimit",
		"getBlockTime",
		"getClusterNodes",
		"getEpochInfo",
		"getEpochSchedule",
		"getFeeForMessage",
		"getFirstAvailableBlock",
		"getGenesisHash",
		"getHealth",
		"getHighestSnapshotSlot",
		"getIdentity",
		"getInflationGovernor",
		"getInflationRate",
		"getInflationReward",
		"getLargestAccounts",
		"getLatestBlockhash",
		"getLeaderSchedule",
		"getMaxRetransmitSlot",
		"getMaxShredInsertSlot",
		"getMinimumBalanceForRentExemption",
		"getMultipleAccounts",
		"getProgramAccounts",
		"getRecentPerformanceSamples",
		"getSignaturesForAddress",
		"getSignatureStatuses",
		"getSlot",
		"getSlotLeader",
		"getSlotLeaders",
		"getStakeActivation",
		"getStakeMinimumDelegation",
		"getSupply",
		"getTokenAccountBalance",
		"getTokenAccountsByDelegate",
		"getTokenAccountsByOwner",
		"getTokenLargestAccounts",
		"getTokenSupply",
		"getTransaction",
		"getTransactionCount",
		"getVersion",
		"getVoteAccounts",
		"isBlockhashValid",
		"minimumLedgerSlot",
		"requestAirdrop",
		"sendTransaction",
		"simulateTransaction",

		// Unstable
		"blockSubscribe",
		"blockUnsubscribe",
		"slotsUpdatesSubscribe",
		"slotsUpdatesUnsubscribe",
		"voteSubscribe",
		"voteUnsubscribe",

		// Deprecated
		"getConfirmedBlock",
		"getConfirmedBlocks",
		"getConfirmedBlocksWithLimit",
		"getConfirmedSignaturesForAddress2",
		"getConfirmedTransaction",
		"getFeeCalculatorForBlockhash",
		"getFeeRateGovernor",
		"getFees",
		"getRecentBlockhash",
		"getSnapshotSlot",
	}

	methodWhitelistMap = make(map[string]bool)
)

func init() {
	for _, v := range methodWhitelist {
		methodWhitelistMap[v] = true
	}
}
