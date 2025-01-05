package calculateservice

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dmitrovia/gophermart/internal/models/bizmodels/accounthistorymodel"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/accountmodel"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/ordermodel"
	gsfcs "github.com/dmitrovia/gophermart/internal/models/endpointsattr/getstatusfromcalcsystemattr"
	"github.com/dmitrovia/gophermart/internal/storage"
	"github.com/jackc/pgx/v5"
)

type CalculateService struct {
	accRepo     storage.AccountStorage
	accOrder    storage.OrderStorage
	ctxDuration time.Duration
	pgxConn     *pgx.Conn
}

func NewCalculateService(
	stor storage.AccountStorage,
	acco storage.OrderStorage,
	ctxDur time.Duration,
	pgxC *pgx.Conn,
) *CalculateService {
	return &CalculateService{
		pgxConn:     pgxC,
		accRepo:     stor,
		ctxDuration: ctxDur,
		accOrder:    acco,
	}
}

func (s *CalculateService) CalculatePoints(
	acc *accountmodel.Account,
	order *ordermodel.Order,
	points float32,
) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		s.ctxDuration)
	defer cancel()

	tranz, err := s.pgxConn.Begin(ctx)
	if err != nil {
		return fmt.Errorf(
			"CalculatePoints->s.pgxConn.Begin %w", err)
	}

	_, err = s.accRepo.MinusPointsByID(
		&ctx, acc.GetID(), points)
	if err != nil {
		return fmt.Errorf(
			"CalculatePoints->MinusPointsByID: %w",
			err)
	}

	_, err = s.accRepo.PlusWithdrawnByID(
		&ctx, acc.GetID(), points)
	if err != nil {
		return fmt.Errorf(
			"CalculatePoints->PlusWithdrawnByID: %w", err)
	}

	_, err = s.accOrder.PlusPointsWriteOffByID(
		&ctx, acc.GetID(), points)
	if err != nil {
		return fmt.Errorf(
			"CalculatePoints->PlusPointsWriteOffByID: %w",
			err)
	}

	err = createAccountHistory(s, &ctx, order, points)
	if err != nil {
		return fmt.Errorf(
			"CalculatePoints->CreateAccountHistory: %w",
			err)
	}

	err = tranz.Commit(ctx)
	if err != nil {
		return fmt.Errorf("CalculatePoints->tranz.Commit %w", err)
	}

	return nil
}

func createAccountHistory(
	service *CalculateService,
	ctx *context.Context,
	order *ordermodel.Order,
	points float32,
) error {
	accHist := accounthistorymodel.AccountHistory{}
	accHist.SetOrder(order)
	accHist.SetpointsWriteOff(&points)

	err := service.accRepo.CreateAccountHistory(ctx, &accHist)
	if err != nil {
		return fmt.Errorf(
			"CreateAccountHistory->accRepo.CreateAccountHistory %w",
			err)
	}

	return nil
}

func (
	s *CalculateService,
) UpdateStatusOrdersAndCalculatePoints() error {
	return nil
}

func getStatusFromCalcSystem(
	ctx *context.Context,
	attr *gsfcs.GetStatusFromCalcSystemAttr,
) (
	*http.Response, error,
) {
	req, err := http.NewRequestWithContext(
		*ctx,
		attr.GetMethod(),
		attr.GetURL(),
		attr.GetReqData())
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
