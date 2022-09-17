package blockstore

import (
	"testing"

	"github.com/certusone/radiance/fixtures"
	"github.com/certusone/radiance/pkg/shred"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataShredsToEntries_Mainnet_Genesis(t *testing.T) {
	rawShreds := fixtures.DataShreds(t, "mainnet", 0)
	shreds := parseShreds(t, rawShreds)
	meta := &SlotMeta{
		Consumed:           3,
		Received:           3,
		LastIndex:          2,
		NumEntryEndIndexes: 1,
		EntryEndIndexes:    []uint32{2},
	}
	entries, err := DataShredsToEntries(meta, shreds)
	require.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, shreds, entries[0].Shreds)
	assert.Equal(t, 3080, len(entries[0].Raw))
	assert.Equal(t, []shred.Entry{
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("BEgmsE3RjuKda9PpnF8r5Ch4HJMHttEhTt38jsyuKuaV"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("7GSbH4g5eTo5jehuAcQmKuoFnTNaYdFRWYJraMCYJ4pY"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("cXQPc1rqMHwxtDfLEJ9CdFWguD9LMnuKaAadwgnoAU9"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("GsdubVEKzBBxdwpRrYPsCoF6TxkwLkJC7vaoAcWvQ6zj"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("GVPeccPmHPAD7sN5za67jvuvRbyJDrZMdFPJgc5k4ApF"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("HCKdcqbc9xQjg4dAqRtJSte4teNRNGuRdQaXbcFqwmsi"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("8JC6mNJbS8JFEm7bhEgkRzk9n6aueX49wXm1hKbMRXtD"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("2weUJw23EwjYkDdH5TWP4iBwgXPW88YDgtLqFUnq5iW2"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("HDUsEkxvxps4KZHHcg1g1Kvkx6xHHQEjb6diEWEbykLs"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("EW44ssuR2hgVYr22WYvRhyWpSwvqpj9FqR268EmyWrp1"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("EZ8KLZWb4fJTqty1PYjmiZjBa7LG1vZiVP385Tf1wGSv"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("4rFoMK6rVXTwKcyeG7GKg7CGTTJpUaLV6nfk4aApHncA"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("2XfkNnbeLgiWFJG88cDvie63Ay4DhF9yVhpaLwVmZvYM"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("Dv8SnjbyB3FHhoct34FNUmsxzgkVTG82nxDaMKLXnHZv"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("8eMT19ctrhoh7buqug6FL5yaDybEhpdGspEfW4Vmr8B2"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("9N8vjyY7mV6Ln77LaG4MZJWdGDApi6gw2Mf3EZ7h74AM"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("AVMHML49UUPxhcVY2QEpULbNXHNxfEHKj1J2cAbEs9Vp"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("DhfgF92wssf9sussfcf2A9gZ9rTJs9k4CuDKf5K9aoMm"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("Hzg7Ptw7fH8gqfaN5FjADcB4kw1XyfZgeCTJhcz9Dt9Y"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("7UDgQWeCGsdjgKDkLDgQ8TA5uuWhUnk6Yf3UwjKRAvS4"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("2ZDGg77CJd9PrGakiApavtY8nTwM6a35ENbuLe1sLXqJ"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("GUDUf6Dj2PZ9qsySXE85qWRj2Ax7peVnTh8AQa4wEHq2"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("FBWpsKo4XZ5jxHh2L3yr17jpe9JdZd7wTDzskA6zYUQW"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("8o6zDK3dC1rJ6UvaZazcYRACTX1zjAvUsGjua2YMWJtX"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("93krGb8CjfAzFWXbcXMUgaiAkLuWQ9wPiwpgG4L153Ji"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("14ngqzcDAwQSkGyYXbV8KHtMtQSup7FFSa7qHS74qS65"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("EdokMAJe1FEkXJeoL4CB1FvKMan63PJ3AgjZzLT5ZHF4"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("2gJtJPAmHY4VCQkF7EmpsBGNRi3122uFebcPJZAVvozf"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("ALffShUMdNB2548HAKU5dwdd1K8JNiApytn8RBKFA5MY"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("3YEHfjSheJX87bME3wNvXN6kEfuHHwsxn7EjaMTRFE9Z"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("2t1F4YsuJSYHNwougVtqHgJHZ7yBswcKpxJnUUYigE1Q"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("8vXxmsQaMQqE3zMmS1FF3UHDho9rVY6sotp5BSNEbS9X"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("CZuyDW5rZFcDHChQRS6FAAcTQHjt5Xh1CL6sKQoMPFND"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("5gws32wDVTpuFrNhCqZwxxu5B8CKvSW9yf4VrqbNWjH7"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("5PJC6VVWQnvqaioJ2wevRFhmoQNCBX9s5h1Kbj6JE5no"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("6WAAcuMZHv1rd1b6v92fqny4dsGSY7Hy6EJpGMmAAznh"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("5joYwpeat27k6XXs3NzJfV7hkw4juX7CCUZvb98dTWJt"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("CNqKzPgR3HUEdf7j6bqj63z86ov5oPTjEitJX8vMucsP"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("HTCaWhSSNviX6YDAyMteKpcgQkAY4GzaMvmKXdkqsvx5"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("4yEgDCGrWidGYjWXL5B3Hcfg7GhwLUdpBYbausUxT2Md"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("HmVzH2rs2oMw2jka3faNT4StoUQhKEVNRsucSbFJgUzQ"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("8UiBLqWrzUhEEtEbf4UhtikFD4qWqzE8aGoimbubT33r"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("uw3ierG673dDLNHTLBAYk6vGp2XzP6MnabjzasqkRLg"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("HkVC9Dm1xm2A4tiFtqodjRBawX38Y4WhPHn5dxDLHowT"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("FMfLjob5dfTmXPuY97TuuWBRZk4qEzkYao6TdZQXSNU9"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("6brZ9hfpcpqq9AUNSa1WYrDWwUazR7uUYDutNn7AxrYq"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("GddzudgJRbP7pZYw7s62VwBzeyjQTcKC6mkN8Ke7Vc6p"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("Gwk1E6CeD1a4Az9Y6Dw9gm1cXWnYFAGioj3UrmttPFGs"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("497LZxPFN8UMparCsDRM2Xwctz9Yr6ggRRxrRWYiYUnr"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("GgV4yq8mU1H8DPqegPqEoyU66jc8dvrrvy5P1kURpUmB"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("mgHymgwdvc5fNukape9sMCAd5HTKoEQQTwxNNtb8Rat"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("2LPLYYGtvfzVBXS7nYYYhBezo36PpvpPJR3TrNfAj3tB"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("FsHZMmXjqNd11HHuCFUTifCC1PzDa5Qa73Et1Pq5CgEN"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("8uq7LJ8GNCoNcsTBgqWLNnxAmR9G2954fdJLw7d22ze8"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("3LaPqPMbS57jESrHL6j8i4CDSCFqioHYM4nw8pMHkyH7"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("5H8mKSpdwXfa7B7NcS7NSRA5qFvmYsnSmfFj1URwF9VV"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("C1dX3PxkzPUeRX4WhsLPgULnisGHWtx1iEw7gTUqbEQe"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("CPosYdr3P8mndtQXF4zkr11KcVuXucMjqrvL7sxJjz4L"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("GrLSqSMZPyeQDySx9RGmFwhgLa67jrhEXsVpZuSwQ51c"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("CBmrfRRb2zghqbHyeiPdfXWGrC5naMT3USR35aeLVYRH"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("4EzStQd1KbsBrmHhfR5r8a6NhaYFAN5bi7t5tVmEMh5U"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("ExLNu3PUdfsyq4DCWYb7VaafoBguxj4bRgyxbfuACQaq"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("Ap9zjPARQp8PRFGFirio2WP7zghbALResHdS9u4FwXZr"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
		{
			NumHashes: 12500,
			Hash:      solana.MustHashFromBase58("4sGjMW1sUnHzSxGspuhpqLDx6wiyjNtZAMdL4VZHirAn"),
			NumTxns:   0,
			Txns:      []solana.Transaction{},
		},
	}, entries[0].Entries)
}

func parseShreds(t testing.TB, raw [][]byte) (shreds []shred.Shred) {
	shreds = make([]shred.Shred, len(raw))
	for i, buf := range raw {
		shreds[i] = shred.NewShredFromSerialized(buf)
		require.NotNil(t, shreds[i], "invalid shred %d", i)
	}
	return shreds
}
