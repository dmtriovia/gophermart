package withdraw

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dmitrovia/gophermart/internal/functions/validatef"
	"github.com/dmitrovia/gophermart/internal/logger"
	"github.com/dmitrovia/gophermart/internal/models/apimodels"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/accountmodel"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/ordermodel"
	"github.com/dmitrovia/gophermart/internal/models/handlerattr/withdrawattr"
	"github.com/dmitrovia/gophermart/internal/service"
)

var errEmptyData = errors.New("data is empty")

const (
	statusISE = http.StatusInternalServerError
)

type Withdraw struct {
	orderService   service.OrderService
	accountService service.AccountService
	attr           *withdrawattr.WithdrawAttr
}

func NewWithdrawHandler(
	accs service.AccountService,
	ords service.OrderService,
	inAttr *withdrawattr.WithdrawAttr,
) *Withdraw {
	return &Withdraw{
		accountService: accs,
		orderService:   ords,
		attr:           inAttr,
	}
}

func (h *Withdraw) WithdrawHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	reqWithdraw := &apimodels.InWithdraw{}

	err := getReqData(req, reqWithdraw)
	if err != nil {
		setErr(writer, h.attr, err, "getReqData")

		return
	}

	isValid := validate(reqWithdraw, h.attr)
	if !isValid {
		writer.WriteHeader(http.StatusUnprocessableEntity)

		return
	}

	isExist, order, err := orderIsExist(h, reqWithdraw)
	if err != nil {
		setErr(writer, h.attr, err, "orderIsExist")

		return
	}

	belongsToSessionUser := order.
		GetClient().GetID() == h.attr.GetSessionUser().GetID()
	if !isExist || !belongsToSessionUser {
		writer.WriteHeader(http.StatusUnprocessableEntity)

		return
	}

	acc, err := getAccountByClient(h)
	if err != nil {
		setErr(writer, h.attr, err, "getAccByClient")

		return
	}

	enough := checkEnoughFunds(acc, reqWithdraw)
	if !enough {
		writer.WriteHeader(http.StatusPaymentRequired)

		return
	}

	err = writeOffPoints(acc, order)
	if err != nil {
		setErr(writer, h.attr, err, "writeOffPoints")

		return
	}

	writer.WriteHeader(http.StatusOK)
}

func setErr(writer http.ResponseWriter,
	inAttr *withdrawattr.WithdrawAttr,
	err error,
	method string,
) {
	writer.WriteHeader(statusISE)
	logger.DoInfoLogFromErr("WithdrawHandler->"+method,
		err, inAttr.GetLogger())
}

func writeOffPoints(*accountmodel.Account,
	*ordermodel.Order,
) error {
	return nil
}

func getAccountByClient(
	handler *Withdraw,
) (*accountmodel.Account, error) {
	acc, err := handler.accountService.GetAccountByClient(
		handler.attr.GetSessionUser().GetID())
	if err != nil {
		return nil, fmt.Errorf(
			"getAccountByClient->getAccountByClient: %w",
			err)
	}

	return acc, nil
}

func checkEnoughFunds(acc *accountmodel.Account,
	withdraw *apimodels.InWithdraw,
) bool {
	clientPoints := acc.GetPoints()
	withdrawPoints := withdraw.PointsWriteOff

	return clientPoints >= withdrawPoints
}

func orderIsExist(
	handler *Withdraw,
	withdraw *apimodels.InWithdraw,
) (bool, *ordermodel.Order, error) {
	isExist, order, err := handler.orderService.OrderIsExist(
		withdraw.OrderIdentifier)
	if err != nil {
		return false, nil, fmt.Errorf(
			"OrderIsExist->orderService.OrderIsExist: %w",
			err)
	}

	if isExist {
		return true, order, nil
	}

	return false, nil, nil
}

func getReqData(
	req *http.Request,
	withdraw *apimodels.InWithdraw,
) error {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("getReqData->io.ReadAll %w", err)
	}

	if len(body) == 0 {
		return fmt.Errorf("getReqData: %w", errEmptyData)
	}

	err = json.Unmarshal(body, withdraw)
	if err != nil {
		return fmt.Errorf("getReqData->json.Unmarshal %w", err)
	}

	err = req.Body.Close()
	if err != nil {
		return fmt.Errorf("getReqData->req.Body.Close() %w", err)
	}

	return nil
}

func validate(withdraw *apimodels.InWithdraw,
	attr *withdrawattr.WithdrawAttr,
) bool {
	res, _ := validatef.IsMatchesTemplate(
		withdraw.OrderIdentifier,
		attr.GetValidIdentOrderPattern())

	return res
}
