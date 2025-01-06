package setorder

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/dmitrovia/gophermart/internal/functions/validatef"
	"github.com/dmitrovia/gophermart/internal/logger"
	"github.com/dmitrovia/gophermart/internal/models/apimodels"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/ordermodel"
	"github.com/dmitrovia/gophermart/internal/models/handlerattr/setorderattr"
	"github.com/dmitrovia/gophermart/internal/service"
)

var errEmptyData = errors.New("data is empty")

type SetOrders struct {
	serv service.OrderService
	attr *setorderattr.SetOrderAttr
}

func NewSetOrderHandler(
	s service.OrderService,
	inAttr *setorderattr.SetOrderAttr,
) *SetOrders {
	return &SetOrders{serv: s, attr: inAttr}
}

func (h *SetOrders) SetOrderHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	reqAttr := &apimodels.InSetOrder{}

	err := getReqData(req, reqAttr)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		logger.DoInfoLogFromErr("SetOrderHandler->getReqData",
			err, h.attr.GetLogger())

		return
	}

	isValid := validate(reqAttr, h.attr)
	if !isValid {
		writer.WriteHeader(http.StatusUnprocessableEntity)

		return
	}

	exist, order, err := h.serv.OrderIsExist(
		reqAttr.Identifier)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		logger.DoInfoLogFromErr("SetOrderHandler->OrderIsExist",
			err, h.attr.GetLogger())

		return
	}

	if exist {
		orderClient := order.GetClient().GetLogin()
		sessionClient := h.attr.GetSessionUser().GetLogin()

		if *orderClient == *sessionClient {
			writer.WriteHeader(http.StatusOK)

			return
		}

		writer.WriteHeader(http.StatusConflict)

		return
	}

	err = createOrder(reqAttr, h)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		logger.DoInfoLogFromErr("SetOrderHandler->createOrder",
			err, h.attr.GetLogger())

		return
	}

	writer.WriteHeader(http.StatusAccepted)
}

func createOrder(reqAttr *apimodels.InSetOrder,
	hand *SetOrders,
) error {
	order := &ordermodel.Order{}

	order.SetIdentifier(&reqAttr.Identifier)
	order.SetClient(hand.attr.GetSessionUser())

	status := ordermodel.OrderStatusNew
	order.SetStatus(&status)

	err := hand.serv.CreateOrder(order)
	if err != nil {
		return fmt.Errorf(
			"createOrder->h.serv.CreateOrder %w", err)
	}

	return nil
}

func validate(reqAttr *apimodels.InSetOrder,
	attr *setorderattr.SetOrderAttr,
) bool {
	res, _ := validatef.IsMatchesTemplate(
		reqAttr.Identifier, attr.GetValidIdentOrderPattern())

	if !res {
		return res
	}

	value, err := strconv.Atoi(reqAttr.Identifier)
	if err != nil {
		return false
	}

	resLuna := validatef.IsValidLuna(value)

	return resLuna
}

func getReqData(
	req *http.Request,
	reqAttr *apimodels.InSetOrder,
) error {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("getReqData->io.ReadAll %w", err)
	}

	if len(body) == 0 {
		return fmt.Errorf("getReqData: %w", errEmptyData)
	}

	reqAttr.Identifier = string(body)

	return nil
}
