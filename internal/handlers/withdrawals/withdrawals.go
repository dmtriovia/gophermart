package withdrawals

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dmitrovia/gophermart/internal/logger"
	"github.com/dmitrovia/gophermart/internal/models/apimodels"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/accounthistorymodel"
	"github.com/dmitrovia/gophermart/internal/models/handlerattr/withdrawalsattr"
	"github.com/dmitrovia/gophermart/internal/service"
)

type Withdrawals struct {
	accService service.AccountService
	attr       *withdrawalsattr.WithdrawalsAttr
}

func NewWithdrawalsHandler(
	accs service.AccountService,
	inAttr *withdrawalsattr.WithdrawalsAttr,
) *Withdrawals {
	return &Withdrawals{
		accService: accs,
		attr:       inAttr,
	}
}

func (h *Withdrawals) WithdrawalsHandler(
	writer http.ResponseWriter,
	_ *http.Request,
) {
	accHists, err := GetAccountHistoryByClient(h)
	if err != nil {
		logger.DoInfoLogFromErr(
			"WithdrawalsHandler->GetAccountHistoryByClient",
			err, h.attr.GetLogger())
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	if len(*accHists) == 0 {
		writer.WriteHeader(http.StatusNoContent)
	}

	marshal, err := formResponeBody(accHists)
	if err != nil {
		logger.DoInfoLogFromErr(
			"WithdrawalsHandler->formResponeBody",
			err, h.attr.GetLogger())
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = writer.Write(*marshal)
	if err != nil {
		logger.DoInfoLogFromErr(
			"WithdrawalsHandler->writer.Write",
			err, h.attr.GetLogger())
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
}

func GetAccountHistoryByClient(
	handler *Withdrawals,
) (*[]accounthistorymodel.AccountHistory, error) {
	accHists,
		scanErrors,
		err := handler.accService.GetAccountHistoryByClient(
		handler.attr.GetSessionUser().GetID())
	if err != nil {
		return nil, fmt.Errorf(
			"getAccHistoryByClient->GetAccountHistoryByClient: %w",
			err)
	}

	if len(*scanErrors) != 0 {
		logger.DoInfoLogFromErr(
			"getAccHistoryByClient->GetAccountHistoryByClient",
			err, handler.attr.GetLogger())
	}

	return accHists, nil
}

func formResponeBody(
	accHists *[]accounthistorymodel.AccountHistory,
) (*[]byte, error) {
	marshal := make([]apimodels.OutWithdrawals,
		0, len(*accHists))

	for _, hist := range *accHists {
		tmp := apimodels.OutWithdrawals{}
		tmp.OrderIdentifier = hist.GetOrder().GetIdentifier()
		tmp.PointsWriteOff = hist.GetpointsWriteOff()
		tmp.Createddate = hist.GetCreateddate()

		marshal = append(marshal, tmp)
	}

	accHistsMarshall, err := json.Marshal(marshal)
	if err != nil {
		return nil,
			fmt.Errorf("formResponeBody->Marshal: %w",
				err)
	}

	return &accHistsMarshall, nil
}
