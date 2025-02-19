package entities

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
)

func TestNewTransaction(t *testing.T) {
	id := uuid.New()
	event := uint32(1)
	company := "abc"
	competenceDate := time.Now()
	metadata := json.RawMessage(`{}`)

	e11, _ := NewEntry(uuid.New(), vos.DebitOperation, "liability.clients.available.111", vos.NextAccountVersion, 123, metadata)
	e12, _ := NewEntry(uuid.New(), vos.CreditOperation, "liability.clients.available.222", vos.NextAccountVersion, 123, metadata)
	validTwoEntries := []Entry{e11, e12}

	e21, _ := NewEntry(uuid.New(), vos.DebitOperation, "liability.clients.available.333", vos.NextAccountVersion, 400, metadata)
	e22, _ := NewEntry(uuid.New(), vos.CreditOperation, "liability.clients.available.444", vos.NextAccountVersion, 300, metadata)
	e23, _ := NewEntry(uuid.New(), vos.CreditOperation, "liability.clients.available.555", vos.NextAccountVersion, 100, metadata)
	validThreeEntries := []Entry{e21, e22, e23}

	e31, _ := NewEntry(uuid.New(), vos.DebitOperation, "liability.clients.available.333", vos.NextAccountVersion, 100, metadata)
	e32, _ := NewEntry(uuid.New(), vos.DebitOperation, "liability.clients.available.333", vos.Version(2), 200, metadata)
	e33, _ := NewEntry(uuid.New(), vos.DebitOperation, "liability.clients.available.333", vos.Version(3), 300, metadata)
	e34, _ := NewEntry(uuid.New(), vos.CreditOperation, "liability.clients.available.444", vos.Version(4), 600, metadata)
	validFourEntries := []Entry{e31, e32, e33, e34}

	testCases := []struct {
		name                string
		id                  uuid.UUID
		entries             []Entry
		expectedTransaction Transaction
		expectedErr         error
	}{
		{
			name:                "Invalid entries number when the transaction has no entries",
			id:                  id,
			entries:             []Entry{},
			expectedTransaction: Transaction{},
			expectedErr:         app.ErrInvalidEntriesNumber,
		},
		{
			name:                "Invalid entries number when the transaction has 1 entry",
			id:                  id,
			entries:             []Entry{e11},
			expectedTransaction: Transaction{},
			expectedErr:         app.ErrInvalidEntriesNumber,
		},
		{
			name:    "Valid transaction with 2 entries",
			id:      id,
			entries: validTwoEntries,
			expectedTransaction: Transaction{
				ID:             id,
				Entries:        validTwoEntries,
				Event:          event,
				Company:        company,
				CompetenceDate: competenceDate,
			},
			expectedErr: nil,
		},
		{
			name:    "Valid transaction with 2 entries - unordered entries",
			id:      id,
			entries: []Entry{e12, e11},
			expectedTransaction: Transaction{
				ID:             id,
				Entries:        validTwoEntries,
				Event:          event,
				Company:        company,
				CompetenceDate: competenceDate,
			},
			expectedErr: nil,
		},
		{
			name:    "Valid transaction with 3 entries",
			id:      id,
			entries: validThreeEntries,
			expectedTransaction: Transaction{
				ID:             id,
				Entries:        validThreeEntries,
				Event:          event,
				Company:        company,
				CompetenceDate: competenceDate,
			},
			expectedErr: nil,
		},
		{
			name:    "Valid transaction with 3 entries - unordered entries",
			id:      id,
			entries: []Entry{e23, e22, e21},
			expectedTransaction: Transaction{
				ID:             id,
				Entries:        validThreeEntries,
				Event:          event,
				Company:        company,
				CompetenceDate: competenceDate,
			},
			expectedErr: nil,
		},
		{
			name:    "Valid transaction with 4 entries - unordered entries by version",
			id:      id,
			entries: []Entry{e33, e34, e32, e31},
			expectedTransaction: Transaction{
				ID:             id,
				Entries:        validFourEntries,
				Event:          event,
				Company:        company,
				CompetenceDate: competenceDate,
			},
			expectedErr: nil,
		},
		{
			name:                "Invalid transaction with 2 entries and balance != 0",
			id:                  id,
			entries:             []Entry{e11, e22},
			expectedTransaction: Transaction{},
			expectedErr:         app.ErrInvalidBalance,
		},
		{
			name:                "Invalid transaction with 3 entries and balance != 0",
			id:                  id,
			entries:             []Entry{e11, e12, e21},
			expectedTransaction: Transaction{},
			expectedErr:         app.ErrInvalidBalance,
		},
		{
			name:                "Invalid transaction with empty ID",
			id:                  uuid.Nil,
			entries:             validTwoEntries,
			expectedTransaction: Transaction{},
			expectedErr:         app.ErrInvalidTransactionID,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTransaction(tt.id, event, company, competenceDate, tt.entries...)

			assert.ErrorIs(t, err, tt.expectedErr)
			assert.Equal(t, tt.expectedTransaction, got)
		})
	}
}
