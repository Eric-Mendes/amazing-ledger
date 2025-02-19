package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test the creation of a new synthetic report object
func TestSyntheticReport(t *testing.T) {
	accountLiquidacao, _ := NewAnalyticAccount("assets.bacen.conta_liquidacao")

	type wants struct {
		results     []AccountResult
		totalCredit int64
		totalDebit  int64
		err         error
	}

	tests := []struct {
		name    string
		account string
		wants   wants
	}{
		{
			name: "Successfully creates a synthetic report",
			wants: wants{
				results: []AccountResult{
					{
						Account: accountLiquidacao,
						Credit:  200,
						Debit:   300,
					},
				},
				totalCredit: 200,
				totalDebit:  300,
				err:         nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSyntheticReport(tt.wants.totalCredit, tt.wants.totalDebit, tt.wants.results)

			assert.Nil(t, err)
			assert.Equal(t, len(tt.wants.results), len(got.Results))
			assert.Equal(t, tt.wants.totalCredit, got.TotalCredit)
			assert.Equal(t, tt.wants.totalDebit, got.TotalDebit)
		})
	}

}
