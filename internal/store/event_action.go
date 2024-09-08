package store

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	// Low risk
	LogoutAction EventAction = iota + 1
	LoginAction
	ViewAction

	// Medium risk
	DepositAction

	// High risk
	WithdrawAction
	TradeAction
	SendAction
)

const (
	logoutLabel   = "logout"
	loginLabel    = "login"
	viewLabel     = "view"
	depositLabel  = "deposit"
	withdrawLabel = "withdraw"
	tradeLabel    = "trade"
	sendLabel     = "send"
)

var EventActionMap = map[string]EventAction{
	logoutLabel:   LogoutAction,
	loginLabel:    LoginAction,
	viewLabel:     ViewAction,
	depositLabel:  DepositAction,
	withdrawLabel: WithdrawAction,
	tradeLabel:    TradeAction,
	sendLabel:     SendAction,
}
var EventActionLookup = map[EventAction]string{
	LogoutAction:   logoutLabel,
	LoginAction:    loginLabel,
	ViewAction:     viewLabel,
	DepositAction:  depositLabel,
	WithdrawAction: withdrawLabel,
	TradeAction:    tradeLabel,
	SendAction:     sendLabel,
}

var EventActionLabels = []string{
	logoutLabel,
	loginLabel,
	viewLabel,
	depositLabel,
	withdrawLabel,
	tradeLabel,
	sendLabel,
}

type EventAction int

func (a *EventAction) UnmarshalJSON(b []byte) error {
	if a == nil {
		return fmt.Errorf("cannot unmarshal into nil EventAction")
	}

	s := strings.Trim(string(b), `"`)
	if _, ok := EventActionMap[s]; !ok {
		return fmt.Errorf("unknown EventAction: %s", s)
	}
	*a = EventActionMap[s]
	return nil
}

func (a EventAction) MarshalJSON() ([]byte, error) {
	if _, ok := EventActionLookup[a]; !ok {
		return nil, fmt.Errorf("unknown EventAction: %d", a)
	}
	return json.Marshal(EventActionLookup[a])
}
