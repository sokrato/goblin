package signal

import (
	"github.com/dlutxx/goblin/utils"
	"testing"
)

func TestSend(t *testing.T) {
	salary := New("amount")
	total := 0
	salary.Connect(func(args utils.Dict) {
		amount, _ := args.Int("amount")
		total += amount
	})
	salary.Send(utils.Dict{"amount": 100})
	if total != 100 {
		t.FailNow()
	}
}
