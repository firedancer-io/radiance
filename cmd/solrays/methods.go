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
		"getFeeCalculatorForBlockhash",
		"getFeeRateGovernor",
		"getFees",
		"getFirstAvailableBlock",
		"getGenesisHash",
		"getHealth",
		"getIdentity",
		"getInflationGovernor",
		"getInflationRate",
		"getInflationReward",
		"getLargestAccounts",
		"getLeaderSchedule",
		"getMaxRetransmitSlot",
		"getMaxShredInsertSlot",
		"getMinimumBalanceForRentExemption",
		"getMultipleAccounts",
		"getProgramAccounts",
		"getRecentBlockhash",
		"getRecentPerformanceSamples",
		"getSignaturesForAddress",
		"getSignatureStatuses",
		"getSlot",
		"getSlotLeader",
		"getSlotLeaders",
		"getStakeActivation",
		"getSnapshotSlot",
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
		"minimumLedgerSlot",
		"requestAirdrop",
		"sendTransaction",
		"simulateTransaction",

		// Deprecated
		"getConfirmedBlock",
		"getConfirmedBlocks",
		"getConfirmedBlocksWithLimit",
		"getConfirmedSignaturesForAddress2",
		"getConfirmedTransaction",
	}

	methodWhitelistMap = make(map[string]bool)
)

func init() {
	for _, v := range methodWhitelist {
		methodWhitelistMap[v] = true
	}
}
