package economy

func makeMockAccounts() *mockAccounts {
	return &mockAccounts{
		accounts: make(map[int64]int64),
		rejects:  make(map[int64]bool),
	}
}

type mockAccounts struct {
	accounts map[int64]int64
	rejects  map[int64]bool
}

func (ma *mockAccounts) Credit(accountID, funds int64) {
	cur := ma.accounts[accountID]
	ma.accounts[accountID] = cur + funds
}

func (ma *mockAccounts) DebitIfPossible(accountID, funds int64) bool {
	if ma.rejects[accountID] {
		return false
	}
	cur := ma.accounts[accountID]
	ma.accounts[accountID] = cur - funds
	return true
}
