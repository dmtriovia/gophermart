package balance

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dmitrovia/gophermart/internal/logger"
	"github.com/dmitrovia/gophermart/internal/models/apimodels"
	"github.com/dmitrovia/gophermart/internal/models/handlerattr/balanceattr"
	"github.com/dmitrovia/gophermart/internal/service"
)

type Balance struct {
	serv service.AccountService
	attr *balanceattr.BalanceAttr
}

func NewBalanceHandler(
	s service.AccountService,
	inAttr *balanceattr.BalanceAttr,
) *Balance {
	return &Balance{serv: s, attr: inAttr}
}

func (h *Balance) BalanceHandler(
	writer http.ResponseWriter,
	_ *http.Request,
) {
	balance, err := getBalanceByClient(h)
	if err != nil {
		logger.DoInfoLogFromErr(
			"BalanceHandler->GetBalanceByClient",
			err, h.attr.GetLogger())
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	marshal, err := formResponeBody(balance)
	if err != nil {
		logger.DoInfoLogFromErr(
			"BalanceHandler->formResponeBody",
			err, h.attr.GetLogger())
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = writer.Write(*marshal)
	if err != nil {
		logger.DoInfoLogFromErr(
			"BalanceHandler->writer.Write",
			err, h.attr.GetLogger())
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
}

func getBalanceByClient(
	handler *Balance,
) (*apimodels.OutBalance, error) {
	balance := &apimodels.OutBalance{}

	acc, err := handler.serv.GetAccountByClient(
		handler.attr.GetSessionUser().GetID())
	if err != nil {
		return nil, fmt.Errorf(
			"GetBalanceByClient->GetAccountByClient: %w",
			err)
	}

	balance.SetOutBalance(acc.GetPoints(), acc.GetWithdrawn())

	return balance, nil
}

func formResponeBody(
	balance *apimodels.OutBalance,
) (*[]byte, error) {
	balanceMarshall, err := json.Marshal(balance)
	if err != nil {
		return nil,
			fmt.Errorf("formResponeBody->Marshal: %w",
				err)
	}

	return &balanceMarshall, nil
}
