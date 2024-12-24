package serverattr

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dmitrovia/gophermart/internal/handlers/balance"
	"github.com/dmitrovia/gophermart/internal/handlers/getorders"
	"github.com/dmitrovia/gophermart/internal/handlers/login"
	"github.com/dmitrovia/gophermart/internal/handlers/notallowed"
	"github.com/dmitrovia/gophermart/internal/handlers/register"
	"github.com/dmitrovia/gophermart/internal/handlers/setorders"
	"github.com/dmitrovia/gophermart/internal/handlers/withdraw"
	"github.com/dmitrovia/gophermart/internal/handlers/withdrawals"
	"github.com/dmitrovia/gophermart/internal/logger"
	"github.com/dmitrovia/gophermart/internal/middleware/loggermiddleware"
	"github.com/dmitrovia/gophermart/internal/models/handlerattr"
	"github.com/dmitrovia/gophermart/internal/models/middlewareattr"
	"github.com/dmitrovia/gophermart/internal/service/accountservice"
	"github.com/dmitrovia/gophermart/internal/service/authservice"
	"github.com/dmitrovia/gophermart/internal/service/orderservice"
	"github.com/dmitrovia/gophermart/internal/storage/postgrestorage"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type ServerAttr struct {
	runAddress           string
	databaseURL          string
	accrualSystemAddress string
	defPORT              string
	defAccSysAddr        string
	defDatabaseURL       string
	validAddrPattern     string
	server               *http.Server
	zapLogger            *zap.Logger
	zapLogLevel          string
	postgreStorage       *postgrestorage.PostgreStorage
	accountService       *accountservice.AccountService
	authService          *authservice.AuthService
	orderSerice          *orderservice.OrderService
	pgxConn              *pgx.Conn
	waitSecRespDB        int
	defReadTimeout       int
	defWriteTimeout      int
	defIdleTimeout       int
	apiURL               string
	migrationsDir        string
	loginAttr            *handlerattr.LoginAttr
	rigsterAttr          *handlerattr.RegisterAttr
	authMidAttr          *middlewareattr.AuthMiddlewareAttr
}

func (p *ServerAttr) Init() error {
	p.defPORT = "localhost:8080"
	p.defAccSysAddr = ""
	p.defDatabaseURL = ""
	p.validAddrPattern = "^[a-zA-Z/ ]{1,100}:[0-9]{1,10}$"
	p.waitSecRespDB = 10
	p.defReadTimeout = 15
	p.defWriteTimeout = 15
	p.defIdleTimeout = 60
	p.apiURL = "/api/user/"
	p.migrationsDir = "db/migrations"
	p.postgreStorage = &postgrestorage.PostgreStorage{}
	p.postgreStorage.Initiate(p.pgxConn)
	p.accountService = accountservice.NewAccountService(
		p.postgreStorage, p.waitSecRespDB)
	p.authService = authservice.NewAuthService(
		p.postgreStorage, p.waitSecRespDB)
	p.orderSerice = orderservice.NewOrderService(
		p.postgreStorage, p.waitSecRespDB)
	p.zapLogLevel = "info"

	logger, err := logger.Initialize(p.zapLogLevel)
	if err != nil {
		return fmt.Errorf(
			"ServerAttr->Init->logger.Initialize %w",
			err)
	}

	p.zapLogger = logger

	p.loginAttr = &handlerattr.LoginAttr{}
	p.rigsterAttr = &handlerattr.RegisterAttr{}
	p.authMidAttr = &middlewareattr.AuthMiddlewareAttr{}

	p.loginAttr.Init(logger)
	p.rigsterAttr.Init(logger)
	p.authMidAttr.Init(logger, p.authService)

	mux := mux.NewRouter()
	initAPIMethods(mux, p)
	p.server = &http.Server{
		Addr:         p.runAddress,
		Handler:      mux,
		ErrorLog:     nil,
		ReadTimeout:  time.Duration(p.defReadTimeout),
		WriteTimeout: time.Duration(p.defWriteTimeout),
		IdleTimeout:  time.Duration(p.defIdleTimeout),
	}

	return nil
}

func initAPIMethods(
	mux *mux.Router,
	attr *ServerAttr,
) {
	get := http.MethodGet
	post := http.MethodPost

	getOrder := getorders.NewGetOrderHandler(
		attr.orderSerice).GetOrderHandler
	balance := balance.NewGetOrderHandler(
		attr.accountService).BalanceHandler
	withdrawals := withdrawals.NewWithdrawalsHandler(
		attr.accountService).WithdrawalsHandler
	hNotAllowed := notallowed.NotAllowed{}
	register := register.NewRegisterHandler(
		attr.authService, attr.rigsterAttr).RegisterHandler
	login := login.NewLoginHandler(
		attr.authService, attr.loginAttr).LoginHandler
	setOrder := setorders.NewSetOrderHandler(
		attr.orderSerice).SetOrderHandler
	withdraw := withdraw.NewWithdrawHandler(
		attr.accountService).WithdrawHandler

	setMethod(get, "orders", mux, attr, getOrder)
	setMethod(get, "balance", mux, attr, balance)
	setMethod(get, "withdrawals", mux, attr, withdrawals)
	setMethod(post, "register", mux, attr, register)
	setMethod(post, "login", mux, attr, login)
	setMethod(post, "orders", mux, attr, setOrder)
	setMethod(post, "withdraw", mux, attr, withdraw)

	mux.MethodNotAllowedHandler = hNotAllowed
}

func setMethod(
	method string,
	url string,
	mux *mux.Router,
	attr *ServerAttr,
	handler func(http.ResponseWriter, *http.Request),
) {
	tmp := mux.Methods(method).Subrouter()
	tmp.HandleFunc(attr.apiURL+url,
		handler)
	tmp.Use(
		loggermiddleware.RequestLogger(attr.zapLogger))
}

func (p *ServerAttr) GetLogger() *zap.Logger {
	return p.zapLogger
}

func (p *ServerAttr) GetmigrationsDir() string {
	return p.migrationsDir
}

func (p *ServerAttr) GetValidAddrPattern() string {
	return p.validAddrPattern
}

func (p *ServerAttr) GetServer() *http.Server {
	return p.server
}

func (p *ServerAttr) GetDefPort() string {
	return p.defPORT
}

func (p *ServerAttr) GetRunAddress() *string {
	return &p.runAddress
}

func (p *ServerAttr) GetDefAccSysAddr() string {
	return p.defAccSysAddr
}

func (p *ServerAttr) GetAccrualSystemAddress() *string {
	return &p.accrualSystemAddress
}

func (p *ServerAttr) GetDefDatabaseURL() string {
	return p.defDatabaseURL
}

func (p *ServerAttr) GetDatabaseURL() *string {
	return &p.databaseURL
}

func (p *ServerAttr) GetWaitSecRespDB() int {
	return p.waitSecRespDB
}

func (p *ServerAttr) SetRunAddress(addr string) {
	p.runAddress = addr
}

func (p *ServerAttr) SetDatabaseURL(databaseURL string) {
	p.databaseURL = databaseURL
}

func (p *ServerAttr) SetAccrualSystemAddress(
	accSysAddr string,
) {
	p.accrualSystemAddress = accSysAddr
}

func (p *ServerAttr) SetPgxConn(
	ctx context.Context,
) error {
	dbConn, err := pgx.Connect(ctx, p.databaseURL)
	if err != nil {
		return fmt.Errorf("initiateDBconn->pgx.Connect %w", err)
	}

	p.pgxConn = dbConn

	return nil
}
