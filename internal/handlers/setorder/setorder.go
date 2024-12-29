package setorder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

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
	reqOrder := &apimodels.InSetOrder{}

	err := getReqData(req, reqOrder)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		logger.DoInfoLogFromErr("SetOrderHandler->getReqData",
			err, h.attr.GetLogger())

		return
	}

	isValid := validate(reqOrder, h.attr)
	if !isValid {
		writer.WriteHeader(http.StatusUnprocessableEntity)

		return
	}

	exist, order, err := h.serv.OrderIsExist(
		reqOrder.Identifier)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		logger.DoInfoLogFromErr("SetOrderHandler->OrderIsExist",
			err, h.attr.GetLogger())

		return
	}

	if exist {
		orderClient := order.GetClient().GetLogin()
		sessionClient := h.attr.GetSessionUser().GetLogin()

		if orderClient == sessionClient {
			writer.WriteHeader(http.StatusOK)

			return
		}

		writer.WriteHeader(http.StatusConflict)

		return
	}

	err = createOrder(reqOrder, h)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		logger.DoInfoLogFromErr("SetOrderHandler->createOrder",
			err, h.attr.GetLogger())

		return
	}

	writer.WriteHeader(http.StatusAccepted)
}

func createOrder(reqOrder *apimodels.InSetOrder,
	hand *SetOrders,
) error {
	order := &ordermodel.Order{}

	order.SetIdentifier(reqOrder.Identifier)
	order.SetClient(hand.attr.GetSessionUser())
	order.SetStatus(ordermodel.OrderStatusNew)

	err := hand.serv.CreateOrder(order)
	if err != nil {
		return fmt.Errorf(
			"createOrder->h.serv.CreateOrder %w", err)
	}

	return nil
}

func validate(order *apimodels.InSetOrder,
	attr *setorderattr.SetOrderAttr,
) bool {
	res, _ := validatef.IsMatchesTemplate(
		order.Identifier, attr.GetValidIdentOrderPattern())

	return res
}

func getReqData(
	req *http.Request,
	order *apimodels.InSetOrder,
) error {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("getReqData->io.ReadAll %w", err)
	}

	if len(body) == 0 {
		return fmt.Errorf("getReqData: %w", errEmptyData)
	}

	err = json.Unmarshal(body, order)
	if err != nil {
		return fmt.Errorf("getReqData->json.Unmarshal %w", err)
	}

	err = req.Body.Close()
	if err != nil {
		return fmt.Errorf("getReqData->req.Body.Close() %w", err)
	}

	return nil
}
