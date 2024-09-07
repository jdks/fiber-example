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

type EventAction int

func (a *EventAction) UnmarshalJSON(b []byte) error {
	if a == nil {
		return fmt.Errorf("cannot unmarshal into nil EventAction")
	}

	s := strings.Trim(string(b), `"`)
	switch strings.ToLower(s) {
	case logoutLabel:
		*a = LogoutAction
	case loginLabel:
		*a = LoginAction
	case viewLabel:
		*a = ViewAction
	case depositLabel:
		*a = DepositAction
	case withdrawLabel:
		*a = WithdrawAction
	case tradeLabel:
		*a = TradeAction
	case sendLabel:
		*a = SendAction
	default:
		return fmt.Errorf("unknown EventAction: %s", s)
	}
	return nil
}

func (a EventAction) MarshalJSON() ([]byte, error) {
	var label string
	switch a {
	case LogoutAction:
		label = logoutLabel
	case LoginAction:
		label = loginLabel
	case ViewAction:
		label = viewLabel
	case DepositAction:
		label = depositLabel
	case WithdrawAction:
		label = withdrawLabel
	case TradeAction:
		label = tradeLabel
	case SendAction:
		label = sendLabel
	default:
		return nil, fmt.Errorf("unknown EventAction: %d", a)
	}
	return json.Marshal(label)
}
