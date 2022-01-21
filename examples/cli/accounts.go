package main

func makeAccounts() *accounts {
	return &accounts{
		accounts: make(map[int64]int64),
	}
}

type accounts struct {
	accounts map[int64]int64
}

func (ma *accounts) Credit(accountID, funds int64) {
	cur := ma.accounts[accountID]
	ma.accounts[accountID] = cur + funds
}

func (ma *accounts) DebitIfPossible(accountID, funds int64) bool {
	cur := ma.accounts[accountID]
	if cur < funds {
		return false
	}
	ma.accounts[accountID] = cur - funds
	return true
}
