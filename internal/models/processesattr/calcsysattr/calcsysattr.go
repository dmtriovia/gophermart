package calcsysattr

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitrovia/gophermart/internal/service/accountservice"
	"github.com/dmitrovia/gophermart/internal/service/authservice"
	"github.com/dmitrovia/gophermart/internal/service/calculateservice"
	"github.com/dmitrovia/gophermart/internal/service/orderservice"
	"github.com/dmitrovia/gophermart/internal/storage/accountstorage"
	"github.com/dmitrovia/gophermart/internal/storage/orderstorage"
	"github.com/dmitrovia/gophermart/internal/storage/userstorage"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type CalcSysAttr struct {
	databaseURL              string
	accrualSystemAddress     string
	zapLogger                *zap.Logger
	orderStorage             *orderstorage.OrderStorage
	userStorage              *userstorage.UserStorage
	accountStorage           *accountstorage.AccountStorage
	accountService           *accountservice.AccountService
	authService              *authservice.AuthService
	orderService             *orderservice.OrderService
	calculateService         *calculateservice.CalculateService
	pgxConn                  *pgx.Conn
	waitSecRespDB            time.Duration
	waitSecRespCalcService   time.Duration
	inetervalCallCalcService time.Duration
}

func (p *CalcSysAttr) Init() {
	p.accountStorage = &accountstorage.AccountStorage{}
	p.orderStorage = &orderstorage.OrderStorage{}
	p.userStorage = &userstorage.UserStorage{}
	p.accountStorage.Initiate(p.pgxConn)
	p.orderStorage.Initiate(p.pgxConn)
	p.userStorage.Initiate(p.pgxConn)
	p.accountService = accountservice.NewAccountService(
		p.accountStorage, p.waitSecRespDB)
	p.authService = authservice.NewAuthService(
		p.userStorage, p.waitSecRespDB)
	p.orderService = orderservice.NewOrderService(
		p.orderStorage, p.waitSecRespDB)
	p.calculateService = calculateservice.NewCalculateService(
		p.accountStorage,
		p.orderStorage, p.waitSecRespDB, p.pgxConn)
}

func (
	p *CalcSysAttr,
) GetCalculateService() *calculateservice.CalculateService {
	return p.calculateService
}

func (p *CalcSysAttr) GetLogger() *zap.Logger {
	return p.zapLogger
}

func (p *CalcSysAttr) SetLogger(logger *zap.Logger) {
	p.zapLogger = logger
}

func (p *CalcSysAttr) GetAccrualSystemAddress() *string {
	return &p.accrualSystemAddress
}

func (p *CalcSysAttr) GetDatabaseURL() *string {
	return &p.databaseURL
}

func (p *CalcSysAttr) GetWaitSecRespDB() time.Duration {
	return p.waitSecRespDB
}

func (p *CalcSysAttr) SetWaitSecRespDB(dur time.Duration) {
	p.waitSecRespDB = dur
}

func (
	p *CalcSysAttr,
) GetInetervalCallCalcService() time.Duration {
	return p.inetervalCallCalcService
}

func (
	p *CalcSysAttr,
) SetInetervalCallCalcService(dur time.Duration) {
	p.inetervalCallCalcService = dur
}

func (
	p *CalcSysAttr,
) GetWaitSecRespCalcService() time.Duration {
	return p.waitSecRespCalcService
}

func (
	p *CalcSysAttr,
) SetWaitSecRespCalcService(dur time.Duration) {
	p.waitSecRespCalcService = dur
}

func (p *CalcSysAttr) SetDatabaseURL(databaseURL string) {
	p.databaseURL = databaseURL
}

func (p *CalcSysAttr) SetAccrualSystemAddress(
	accSysAddr string,
) {
	p.accrualSystemAddress = accSysAddr
}

func (p *CalcSysAttr) SetPgxConn(
	ctxDB context.Context,
) error {
	dbConn, err := pgx.Connect(ctxDB, p.databaseURL)
	if err != nil {
		return fmt.Errorf("initiateDBconn->pgx.Connect %w", err)
	}

	p.pgxConn = dbConn

	return nil
}
