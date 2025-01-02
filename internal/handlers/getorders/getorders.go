package getorders

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dmitrovia/gophermart/internal/logger"
	"github.com/dmitrovia/gophermart/internal/models/apimodels"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/ordermodel"
	"github.com/dmitrovia/gophermart/internal/models/handlerattr/getordersattr"
	"github.com/dmitrovia/gophermart/internal/service"
)

type GetOrders struct {
	serv service.OrderService
	attr *getordersattr.GetOrdersAttr
}

func NewGetOrdersHandler(
	s service.OrderService,
	inAttr *getordersattr.GetOrdersAttr,
) *GetOrders {
	return &GetOrders{serv: s, attr: inAttr}
}

func (h *GetOrders) GetOrderHandler(
	writer http.ResponseWriter,
	_ *http.Request,
) {
	orders, err := getOrdersByClient(h)
	if err != nil {
		logger.DoInfoLogFromErr(
			"GetOrdersHandler->GetOrdersByClient",
			err, h.attr.GetLogger())
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	if len(*orders) == 0 {
		writer.WriteHeader(http.StatusNoContent)

		return
	}

	marshal, err := formResponeBody(orders)
	if err != nil {
		logger.DoInfoLogFromErr(
			"GetOrdersHandler->formResponeBody",
			err, h.attr.GetLogger())
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = writer.Write(*marshal)
	if err != nil {
		logger.DoInfoLogFromErr(
			"GetOrdersHandler->writer.Write",
			err, h.attr.GetLogger())
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
}

func getOrdersByClient(
	handler *GetOrders,
) (*[]ordermodel.Order, error) {
	orders, scanErrors, err := handler.serv.GetOrdersByClient(
		handler.attr.GetSessionUser().GetID())
	if err != nil {
		return nil, fmt.Errorf(
			"GetOrdersByClient->GetOrdersByClient: %w",
			err)
	}

	if len(*scanErrors) != 0 {
		for _, err := range *scanErrors {
			logger.DoInfoLogFromErr(
				"GetOrdersByClient->GetOrdersByClient",
				err, handler.attr.GetLogger())
		}
	}

	return orders, nil
}

func formResponeBody(
	orders *[]ordermodel.Order,
) (*[]byte, error) {
	marshal := make([]apimodels.OutGetOrders, 0, len(*orders))

	for _, order := range *orders {
		tmp := apimodels.OutGetOrders{}
		tmp.SetOutGetOrders(order.GetIdentifier(),
			order.GetCreateddate(),
			order.GetStatus(),
			order.GetAccrual())

		marshal = append(marshal, tmp)
	}

	ordersMarshall, err := json.Marshal(marshal)
	if err != nil {
		return nil,
			fmt.Errorf("formResponeBody->Marshal: %w",
				err)
	}

	return &ordersMarshall, nil
}
