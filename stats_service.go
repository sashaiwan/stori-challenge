package main

type DebitCreditStats struct {
	Average float64
	Count   int
}

type MonthlyStats struct {
	Credit DebitCreditStats
	Debit  DebitCreditStats
}

type TransactionStats struct {
	StatsByMonth map[string]MonthlyStats
	TotalBalance float64
}

// Receives a transactions array and returns a TransactionStats struct
//
//	[month]: {
//	  credit: {
//	    average: float64,
//	    count: int,
//	  },
//	  debit: {
//	    average: float64,
//	    count: int,
//	  }
//	},
//	totalBalance: float64
func getTransactionStats(transactions []Transaction) TransactionStats {
	stats := TransactionStats{
		StatsByMonth: make(map[string]MonthlyStats),
		TotalBalance: 0,
	}

	debitByMonth := make(map[string][]float64)
	creditByMonth := make(map[string][]float64)

	for _, transaction := range transactions {
		// Date is in YYYY-MM-DD format
		month := transaction.Date[5:7]

		if _, exists := stats.StatsByMonth[month]; !exists {
			stats.StatsByMonth[month] = MonthlyStats{
				Credit: DebitCreditStats{Average: 0, Count: 0},
				Debit:  DebitCreditStats{Average: 0, Count: 0},
			}
		}

		monthStats := stats.StatsByMonth[month]
		if transaction.Type == Credit {
			stats.TotalBalance += transaction.Amount
			monthStats.Credit.Count += 1
			creditByMonth[month] = append(creditByMonth[month], transaction.Amount)
		} else {
			// Debit transactions
			stats.TotalBalance -= transaction.Amount
			monthStats.Debit.Count += 1
			debitByMonth[month] = append(debitByMonth[month], transaction.Amount)
		}

		stats.StatsByMonth[month] = monthStats
	}

	for month := range stats.StatsByMonth {
		monthStats := stats.StatsByMonth[month]

		if debits := debitByMonth[month]; len(debits) > 0 {
			var sum float64
			for _, amount := range debits {
				sum += amount
			}
			monthStats.Debit.Average = sum / float64(len(debits))
		}

		if credits := creditByMonth[month]; len(credits) > 0 {
			var sum float64
			for _, amount := range credits {
				sum += amount
			}
			monthStats.Credit.Average = sum / float64(len(credits))
		}

		stats.StatsByMonth[month] = monthStats
	}
	// 1. total balance in the account
	// 2. The number of transactions grouped by month
	// 3. The average credit and average debit amounts grouped by month
	return stats
}
