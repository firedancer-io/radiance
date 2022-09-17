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
	shreds := parseShreds(t, rawShreds, 1)
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

func TestDataShredsToEntries_Mainnet_Recent(t *testing.T) {
	rawShreds := fixtures.DataShreds(t, "mainnet", 102815960)
	shreds := parseShreds(t, rawShreds, 2)
	meta := &SlotMeta{
		Consumed:           1427,
		Received:           1427,
		LastIndex:          1426,
		NumEntryEndIndexes: 574,
		EntryEndIndexes: []uint32{
			0, 1, 2, 3, 5, 16, 17, 18, 20, 21, 22, 23, 24, 26, 28, 30, 31, 34, 35, 37, 38, 39, 41, 42, 43, 44, 45, 46,
			47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 59, 60, 62, 63, 69, 74, 75, 77, 81, 82, 83, 86, 87, 88, 89, 90, 91,
			92, 93, 94, 95, 101, 102, 103, 104, 105, 107, 108, 114, 115, 116, 117, 123, 124, 127, 133, 136, 137, 138,
			143, 144, 145, 146, 147, 148, 152, 154, 155, 156, 159, 169, 170, 171, 177, 178, 180, 182, 183, 184, 185,
			186, 187, 188, 189, 190, 191, 196, 200, 201, 202, 203, 204, 205, 206, 207, 210, 211, 212, 213, 214, 215,
			216, 217, 218, 219, 221, 223, 232, 235, 237, 239, 242, 243, 245, 246, 249, 251, 253, 254, 255, 260, 265,
			285, 286, 287, 292, 293, 301, 302, 303, 304, 305, 313, 314, 315, 327, 334, 335, 343, 344, 345, 346, 355,
			364, 370, 371, 372, 373, 379, 380, 381, 382, 389, 398, 405, 406, 417, 418, 425, 426, 427, 428, 429, 430,
			439, 440, 441, 450, 452, 453, 454, 455, 463, 464, 465, 466, 467, 468, 478, 479, 480, 481, 482, 483, 484,
			495, 496, 497, 498, 499, 500, 501, 514, 516, 517, 532, 533, 534, 535, 536, 537, 547, 549, 558, 559, 560,
			572, 573, 574, 575, 586, 587, 588, 589, 590, 591, 602, 603, 613, 616, 625, 626, 627, 628, 629, 630, 636,
			637, 638, 644, 645, 646, 653, 654, 655, 656, 663, 664, 665, 666, 675, 677, 684, 685, 686, 687, 688, 689,
			690, 695, 700, 702, 708, 716, 722, 725, 726, 727, 731, 733, 735, 741, 742, 743, 745, 746, 747, 748, 751,
			752, 754, 759, 766, 770, 776, 777, 778, 779, 780, 781, 786, 791, 799, 802, 810, 813, 820, 821, 822, 823,
			824, 825, 826, 827, 828, 829, 830, 838, 846, 848, 849, 850, 851, 852, 861, 862, 863, 873, 874, 875, 876,
			877, 878, 879, 880, 881, 882, 892, 893, 894, 895, 896, 907, 909, 910, 911, 912, 913, 914, 915, 916, 917,
			918, 919, 920, 921, 933, 934, 935, 936, 937, 938, 939, 940, 954, 955, 956, 957, 958, 959, 975, 976, 977,
			978, 994, 997, 998, 999, 1013, 1014, 1015, 1016, 1017, 1018, 1019, 1020, 1021, 1022, 1039, 1040, 1041, 1056,
			1059, 1060, 1061, 1077, 1078, 1079, 1080, 1081, 1082, 1083, 1084, 1085, 1086, 1087, 1090, 1091, 1092, 1093,
			1094, 1095, 1096, 1097, 1098, 1099, 1100, 1101, 1115, 1117, 1118, 1119, 1122, 1123, 1124, 1125, 1126, 1138,
			1139, 1140, 1141, 1142, 1143, 1144, 1145, 1146, 1147, 1148, 1149, 1150, 1151, 1164, 1165, 1166, 1167, 1168,
			1182, 1184, 1189, 1190, 1191, 1192, 1193, 1194, 1195, 1196, 1197, 1198, 1199, 1213, 1214, 1215, 1220, 1222,
			1223, 1224, 1226, 1227, 1228, 1229, 1230, 1231, 1232, 1233, 1234, 1235, 1236, 1237, 1238, 1239, 1240, 1241,
			1242, 1243, 1244, 1245, 1246, 1247, 1262, 1265, 1266, 1267, 1272, 1273, 1274, 1275, 1276, 1277, 1278, 1279,
			1289, 1292, 1293, 1294, 1295, 1296, 1297, 1298, 1299, 1300, 1301, 1302, 1303, 1304, 1305, 1306, 1307, 1308,
			1309, 1310, 1311, 1322, 1325, 1326, 1327, 1328, 1329, 1330, 1331, 1332, 1333, 1334, 1335, 1336, 1337, 1338,
			1347, 1350, 1351, 1352, 1353, 1354, 1355, 1356, 1357, 1358, 1359, 1360, 1367, 1368, 1369, 1370, 1371, 1372,
			1373, 1374, 1379, 1384, 1385, 1386, 1388, 1389, 1394, 1395, 1396, 1397, 1398, 1403, 1404, 1407, 1416, 1417,
			1422, 1423, 1425, 1426,
		},
	}
	entries, err := DataShredsToEntries(meta, shreds)
	require.NoError(t, err)
	assert.Equal(t, 574, len(entries))
}

func parseShreds(t testing.TB, raw [][]byte, version int) (shreds []shred.Shred) {
	shreds = make([]shred.Shred, len(raw))
	for i, buf := range raw {
		shreds[i] = shred.NewShredFromSerialized(buf, version)
		require.NotNil(t, shreds[i], "invalid shred %d", i)
		// Forgetting this assert cost me half an hour of time
		assert.Equal(t, shreds[i].CommonHeader().Index, uint32(i))
	}
	return shreds
}
