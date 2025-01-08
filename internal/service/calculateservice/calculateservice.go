package calculateservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dmitrovia/gophermart/internal/logger"
	"github.com/dmitrovia/gophermart/internal/models/apimodels"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/accounthistorymodel"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/accountmodel"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/ordermodel"
	"github.com/dmitrovia/gophermart/internal/models/endpointsattr/getstatusfromcalcsystemattr"
	"github.com/dmitrovia/gophermart/internal/storage"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

var errEmptyData = errors.New("data is empty")

var errStatusInternalServerError = errors.New(
	"error when accessing the service")

var errStatusNoContent = errors.New(
	"the order is not registered in the payment system")

type CalculateService struct {
	accRepo               storage.AccountStorage
	accOrder              storage.OrderStorage
	ctxDurationDB         time.Duration
	ctxDurationOutService time.Duration
	pgxConn               *pgx.Conn
}

func NewCalculateService(
	stor storage.AccountStorage,
	acco storage.OrderStorage,
	ctxDurationDB time.Duration,
	ctxDurationOutService time.Duration,
	pgxC *pgx.Conn,
) *CalculateService {
	return &CalculateService{
		pgxConn:               pgxC,
		accRepo:               stor,
		ctxDurationDB:         ctxDurationDB,
		ctxDurationOutService: ctxDurationOutService,
		accOrder:              acco,
	}
}

func (s *CalculateService) CalculatePoints(
	acc *accountmodel.Account,
	order *ordermodel.Order,
	points float32,
) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		s.ctxDurationDB)
	defer cancel()

	tranz, err := s.pgxConn.Begin(ctx)
	if err != nil {
		return fmt.Errorf(
			"CalculatePoints->s.pgxConn.Begin %w", err)
	}

	_, err = s.accRepo.ChangePointsByID(
		&ctx, tranz, acc.GetID(), points, "-")
	if err != nil {
		return rollbackWithErr(&ctx, tranz, err,
			"CalculatePoints->ChangePointsByID")
	}

	_, err = s.accRepo.ChangeWithdrawnByID(
		&ctx, tranz, acc.GetID(), points, "+")
	if err != nil {
		return rollbackWithErr(&ctx, tranz, err,
			"CalculatePoints->ChangeWithdrawnByID")
	}

	_, err = s.accOrder.ChangePointsWriteOffByID(
		&ctx, tranz, acc.GetID(), points, "+")
	if err != nil {
		return rollbackWithErr(&ctx, tranz, err,
			"CalculatePoints->ChangePointsWriteOffByID")
	}

	err = createAccountHistory(&ctx, tranz, s, order, points)
	if err != nil {
		return rollbackWithErr(&ctx, tranz, err,
			"CalculatePoints->CreateAccountHistory")
	}

	err = tranz.Commit(ctx)
	if err != nil {
		return fmt.Errorf("CalculatePoints->tranz.Commit %w", err)
	}

	return nil
}

func rollbackWithErr(ctx *context.Context,
	tranz pgx.Tx,
	err error,
	nameF string,
) error {
	errR := tranz.Rollback(*ctx)
	if errR != nil {
		return fmt.Errorf(
			"rollbackWithErr->Rollback: %w",
			errR)
	}

	return fmt.Errorf(
		"%s: %w",
		nameF, err)
}

func createAccountHistory(
	ctx *context.Context,
	tranz pgx.Tx,
	service *CalculateService,
	order *ordermodel.Order,
	points float32,
) error {
	accHist := accounthistorymodel.AccountHistory{}
	accHist.SetOrder(order)
	accHist.SetpointsWriteOff(&points)

	err := service.accRepo.CreateAccountHistory(ctx,
		tranz, &accHist)
	if err != nil {
		return fmt.Errorf(
			"CreateAccountHistory->accRepo.CreateAccountHistory %w",
			err)
	}

	return nil
}

func (
	s *CalculateService,
) UpdateStatusOrdersAndCalculatePoints(
	attr *getstatusfromcalcsystemattr.
		GetStatusFromCalcSystemAttr,
) error {
	ctxDB, cancel := context.WithTimeout(
		context.Background(), s.ctxDurationDB)
	defer cancel()

	statuses := "'NEW','REGISTERED','PROCESSING'"
	funcName := "UpdateStatusOrdersAndCalculatePoints"

	orders, scanErrors,
		err := s.accOrder.GetOrdersByStatuses(&ctxDB, statuses)
	if err != nil {
		return fmt.Errorf("UpdateStatusOrdersAndCalculatePoints"+
			"->GetOrdersByStatuses: %w", err)
	}

	if len(*scanErrors) != 0 {
		for _, err := range *scanErrors {
			doLog(funcName+"->Scan", err, attr.GetLogger())
		}
	}

	for _, order := range *orders {
		ctxReq, cancel1 := context.WithTimeout(
			context.Background(), s.ctxDurationOutService)
		defer cancel1()

		tmp := apimodels.InGetStatusFromCalcSystem{}
		tmp.Set(order.GetIdentifier())
		attr.SetURLForReq(attr.GetDefURL() + *tmp.Identifier)

		response, err := getStatusFromCalcSystem(&ctxReq, attr)
		if err != nil {
			doLog(funcName+"->getStatusFromCalcSystem",
				err, attr.GetLogger())

			break
		}

		err = processResponse(ctxDB, s, response, &order)
		if err != nil {
			doLog(funcName+"->processResponse",
				err, attr.GetLogger())

			continue
		}

		err = response.Body.Close()
		if err != nil {
			doLog(funcName+"->esponse.Body.Close",
				err, attr.GetLogger())

			continue
		}
	}

	return nil
}

func processResponse(
	ctx context.Context,
	service *CalculateService,
	response *http.Response,
	order *ordermodel.Order,
) error {
	code := response.StatusCode

	switch code {
	case http.StatusOK:
		err := processStatusOK(ctx, service, response, order)
		if err != nil {
			return fmt.Errorf("processResponse->processStatusOK %w",
				err)
		}
	case http.StatusNoContent:
		return fmt.Errorf("StatusNoContent %w",
			errStatusNoContent)
	case http.StatusInternalServerError:
		return fmt.Errorf("StatusInternalServerError %w",
			errStatusInternalServerError)
	}

	return nil
}

func processStatusOK(
	ctx context.Context, service *CalculateService,
	response *http.Response, order *ordermodel.Order,
) error {
	respData := &apimodels.OutGetStatusFromCalcSystem{}

	err := getRespData(response, respData)
	if err != nil {
		return fmt.Errorf("processStatusOK->getRespData %w", err)
	}

	registered := ordermodel.OrderStatusRegistered
	processing := ordermodel.OrderStatusProcessing
	invalid := ordermodel.OrderStatusInvalid
	processed := ordermodel.OrderStatusProcessed

	if *order.GetStatus() == *respData.Status {
		return nil
	}

	if *respData.Status == registered ||
		*respData.Status == processing ||
		*respData.Status == invalid {
		_, err := service.accOrder.UpdateStatusByID(&ctx,
			order.GetID(), *respData.Status)
		if err != nil {
			return fmt.Errorf(
				"processStatusOK->UpdateStatusByID %w", err)
		}
	} else if *respData.Status == processed {
		err := setProcessed(ctx, service, order, respData)
		if err != nil {
			return fmt.Errorf(
				"processStatusOK->setProcessed %w", err)
		}
	}

	return nil
}

func setProcessed(ctx context.Context,
	service *CalculateService,
	order *ordermodel.Order,
	respData *apimodels.OutGetStatusFromCalcSystem,
) error {
	processed := ordermodel.OrderStatusProcessed

	if respData.Accrual == nil {
		*respData.Accrual = 0
	}

	tranz, err := service.pgxConn.Begin(ctx)
	if err != nil {
		return fmt.Errorf(
			"setProcessed->s.pgxConn.Begin %w", err)
	}

	_, err = service.accOrder.UpdateStatusAccrualByID(&ctx,
		tranz, order.GetID(), *respData.Accrual, processed)
	if err != nil {
		return rollbackWithErr(&ctx, tranz, err,
			"setProcessed->UpdateStatusAccrualByID")
	}

	acc, err := service.accRepo.GetAccountByClient(&ctx,
		order.GetClient().GetID())
	if err != nil {
		return rollbackWithErr(&ctx, tranz, err,
			"setProcessed->GetAccountByClient")
	}

	_, err = service.accRepo.ChangePointsByID(&ctx, tranz,
		acc.GetID(), *respData.Accrual, "+")
	if err != nil {
		return rollbackWithErr(&ctx, tranz, err,
			"setProcessed->ChangePointsByID")
	}

	err = tranz.Commit(ctx)
	if err != nil {
		return fmt.Errorf("setProcessed->tranz.Commit %w", err)
	}

	return nil
}

func getRespData(
	response *http.Response,
	respData *apimodels.OutGetStatusFromCalcSystem,
) error {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("getRespData->io.ReadAll %w", err)
	}

	if len(body) == 0 {
		return fmt.Errorf("getRespData: %w", errEmptyData)
	}

	err = json.Unmarshal(body, respData)
	if err != nil {
		return fmt.Errorf("getRespData->json.Unmarshal %w", err)
	}

	return nil
}

func doLog(msgText string,
	err error,
	zLogger *zap.Logger,
) {
	logger.DoInfoLogFromErr(
		msgText,
		err, zLogger)
}

func getStatusFromCalcSystem(
	ctx *context.Context,
	attr *getstatusfromcalcsystemattr.
		GetStatusFromCalcSystemAttr,
) (
	*http.Response, error,
) {
	req, err := http.NewRequestWithContext(
		*ctx,
		attr.GetMethod(),
		attr.GetURLForReq(),
		strings.NewReader(""))
	if err != nil {
		return nil,
			fmt.Errorf(
				"GetStatusFromCalcSystem->NewRequestWithContext: %w",
				err)
	}

	req.Header.Set("Content-Type", attr.GetContentType())

	resp, err := attr.GetClient().Do(req)
	if err != nil {
		return nil,
			fmt.Errorf("GetStatusFromCalcSystem->Do: %w",
				err)
	}

	return resp, nil
}
