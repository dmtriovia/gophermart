package getorder

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dmitrovia/gophermart/internal/functions/validatef"
	"github.com/dmitrovia/gophermart/internal/logger"
	"github.com/dmitrovia/gophermart/internal/models/apimodels"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/ordermodel"
	"github.com/dmitrovia/gophermart/internal/models/handlerattr/getorderattr"
	"github.com/dmitrovia/gophermart/internal/service"
	"github.com/gorilla/mux"
)

type GetOrder struct {
	orderService service.OrderService
	attr         *getorderattr.GetOrderAttr
}

func NewGetOrderHandler(
	s service.OrderService,
	inAttr *getorderattr.GetOrderAttr,
) *GetOrder {
	return &GetOrder{orderService: s, attr: inAttr}
}

func (h *GetOrder) GetOrderHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	reqAttr := &apimodels.InGetOrder{}

	getReqData(req, reqAttr)

	isValid := validate(reqAttr, h.attr)
	if !isValid {
		writer.WriteHeader(http.StatusUnprocessableEntity)

		return
	}

	isExist, order, err := orderIsExist(h, reqAttr)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		logger.DoInfoLogFromErr("GetOrderHandler->orderIsExist",
			err, h.attr.GetLogger())

		return
	}

	if !isExist {
		writer.WriteHeader(http.StatusNoContent)

		return
	}

	marshal, err := formResponeBody(order)
	if err != nil {
		logger.DoInfoLogFromErr(
			"GetOrderHandler->formResponeBody",
			err, h.attr.GetLogger())
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = writer.Write(*marshal)
	if err != nil {
		logger.DoInfoLogFromErr(
			"GetOrderHandler->writer.Write",
			err, h.attr.GetLogger())
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
}

func orderIsExist(
	handler *GetOrder,
	reqAttr *apimodels.InGetOrder,
) (bool, *ordermodel.Order, error) {
	isExist, order, err := handler.orderService.OrderIsExist(
		reqAttr.Identifier)
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
	reqAttr *apimodels.InGetOrder,
) {
	reqAttr.Identifier = mux.Vars(req)["number"]
}

func validate(reqAttr *apimodels.InGetOrder,
	attr *getorderattr.GetOrderAttr,
) bool {
	res, _ := validatef.IsMatchesTemplate(
		reqAttr.Identifier, attr.GetValidIdentOrderPattern())

	return res
}

func formResponeBody(
	order *ordermodel.Order,
) (*[]byte, error) {
	outOrder := &apimodels.OutGetOrder{}
	outOrder.SetOutGetOrder(
		order.GetIdentifier(),
		order.GetStatus(),
		order.GetAccrual())

	orderMarshall, err := json.Marshal(outOrder)
	if err != nil {
		return nil,
			fmt.Errorf("formResponeBody->Marshal: %w",
				err)
	}

	return &orderMarshall, nil
}
