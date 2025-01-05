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
	"github.com/dmitrovia/gophermart/internal/handlers/setorder"
	"github.com/dmitrovia/gophermart/internal/handlers/withdraw"
	"github.com/dmitrovia/gophermart/internal/handlers/withdrawals"
	"github.com/dmitrovia/gophermart/internal/middleware/authmiddleware"
	"github.com/dmitrovia/gophermart/internal/middleware/loggermiddleware"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/usermodel"
	"github.com/dmitrovia/gophermart/internal/models/handlerattr/balanceattr"
	"github.com/dmitrovia/gophermart/internal/models/handlerattr/getorderattr"
	"github.com/dmitrovia/gophermart/internal/models/handlerattr/getordersattr"
	"github.com/dmitrovia/gophermart/internal/models/handlerattr/loginattr"
	"github.com/dmitrovia/gophermart/internal/models/handlerattr/registerattr"
	"github.com/dmitrovia/gophermart/internal/models/handlerattr/setorderattr"
	"github.com/dmitrovia/gophermart/internal/models/handlerattr/withdrawalsattr"
	"github.com/dmitrovia/gophermart/internal/models/handlerattr/withdrawattr"
	"github.com/dmitrovia/gophermart/internal/models/middlewareattr/authmiddlewareattr"
	"github.com/dmitrovia/gophermart/internal/service/accountservice"
	"github.com/dmitrovia/gophermart/internal/service/authservice"
	"github.com/dmitrovia/gophermart/internal/service/calculateservice"
	"github.com/dmitrovia/gophermart/internal/service/orderservice"
	"github.com/dmitrovia/gophermart/internal/storage/accountstorage"
	"github.com/dmitrovia/gophermart/internal/storage/orderstorage"
	"github.com/dmitrovia/gophermart/internal/storage/userstorage"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

const initReadTimeout = 15

const initWriteTimeout = 15

const initIdleTimeout = 60

type ServerAttr struct {
	runAddress       string
	databaseURL      string
	validAddrPattern string
	apiURL           string
	migrationsDir    string
	server           *http.Server
	zapLogger        *zap.Logger
	orderStorage     *orderstorage.OrderStorage
	userStorage      *userstorage.UserStorage
	accountStorage   *accountstorage.AccountStorage
	accountService   *accountservice.AccountService
	authService      *authservice.AuthService
	orderService     *orderservice.OrderService
	calculateService *calculateservice.CalculateService
	pgxConn          *pgx.Conn
	waitSecRespDB    time.Duration
	defReadTimeout   time.Duration
	defWriteTimeout  time.Duration
	defIdleTimeout   time.Duration
	withdrawalsAttr  *withdrawalsattr.WithdrawalsAttr
	loginAttr        *loginattr.LoginAttr
	rigsterAttr      *registerattr.RegisterAttr
	setOrderAttr     *setorderattr.SetOrderAttr
	getOrdersAttr    *getordersattr.GetOrdersAttr
	getOrderAttr     *getorderattr.GetOrderAttr
	balanceAttr      *balanceattr.BalanceAttr
	withdrawAttr     *withdrawattr.WithdrawAttr
	authMidAttr      *authmiddlewareattr.AuthMiddlewareAttr
	sessionUser      *usermodel.User
}

func (p *ServerAttr) Init() error {
	p.sessionUser = &usermodel.User{}

	p.defReadTimeout,
		p.defWriteTimeout = initReadTimeout*time.Second,
		initWriteTimeout*time.Second
	p.defIdleTimeout = initIdleTimeout * time.Second
	p.apiURL = "/api/user/"
	p.migrationsDir = "db/migrations"
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

	initHandlersAttr(p)
	p.authMidAttr = &authmiddlewareattr.AuthMiddlewareAttr{}
	p.authMidAttr.Init(p.zapLogger,
		p.authService, p.sessionUser)

	mux := mux.NewRouter()
	initAPIMethods(mux, p)

	p.server = &http.Server{
		Addr:         p.runAddress,
		Handler:      mux,
		ErrorLog:     nil,
		ReadTimeout:  p.defReadTimeout,
		WriteTimeout: p.defWriteTimeout,
		IdleTimeout:  p.defIdleTimeout,
	}

	return nil
}

func initHandlersAttr(attr *ServerAttr) {
	attr.loginAttr = &loginattr.LoginAttr{}
	attr.rigsterAttr = &registerattr.RegisterAttr{}
	attr.setOrderAttr = &setorderattr.SetOrderAttr{}
	attr.getOrdersAttr = &getordersattr.GetOrdersAttr{}
	attr.getOrderAttr = &getorderattr.GetOrderAttr{}
	attr.balanceAttr = &balanceattr.BalanceAttr{}
	attr.withdrawAttr = &withdrawattr.WithdrawAttr{}
	attr.withdrawalsAttr = &withdrawalsattr.WithdrawalsAttr{}

	attr.loginAttr.Init(attr.zapLogger)
	attr.rigsterAttr.Init(attr.zapLogger)

	attr.setOrderAttr.Init(attr.zapLogger, attr.sessionUser)
	attr.getOrdersAttr.Init(attr.zapLogger, attr.sessionUser)
	attr.getOrderAttr.Init(attr.GetLogger(), attr.sessionUser)
	attr.balanceAttr.Init(attr.zapLogger, attr.sessionUser)
	attr.withdrawAttr.Init(attr.zapLogger, attr.sessionUser)
	attr.withdrawalsAttr.Init(attr.zapLogger, attr.sessionUser)
}

func initAPIMethods(
	mux *mux.Router,
	attr *ServerAttr,
) {
	get := http.MethodGet
	post := http.MethodPost

	getOrders := getorders.NewGetOrdersHandler(
		attr.orderService, attr.getOrdersAttr).GetOrdersHandler
	balance := balance.NewBalanceHandler(
		attr.accountService, attr.balanceAttr).BalanceHandler
	withdrawals := withdrawals.NewWithdrawalsHandler(
		attr.accountService,
		attr.withdrawalsAttr).WithdrawalsHandler
	hNotAllowed := notallowed.NotAllowed{}
	register := register.NewRegisterHandler(
		attr.authService, attr.accountService,
		attr.rigsterAttr).RegisterHandler
	login := login.NewLoginHandler(
		attr.authService, attr.loginAttr).LoginHandler
	setOrder := setorder.NewSetOrderHandler(
		attr.orderService, attr.setOrderAttr).SetOrderHandler
	withdraw := withdraw.NewWithdrawHandler(
		attr.accountService, attr.orderService,
		attr.calculateService, attr.withdrawAttr).WithdrawHandler

	setMethod(get, "orders", mux, attr, getOrders, true)
	setMethod(get, "balance", mux, attr, balance, true)
	setMethod(get, "withdrawals", mux, attr, withdrawals, true)
	setMethod(post, "register", mux, attr, register, false)
	setMethod(post, "login", mux, attr, login, false)
	setMethod(post, "orders", mux, attr, setOrder, true)
	setMethod(post, "withdraw", mux, attr, withdraw, true)

	mux.MethodNotAllowedHandler = hNotAllowed
}

func setMethod(
	method string,
	url string,
	mux *mux.Router,
	attr *ServerAttr,
	handler func(http.ResponseWriter, *http.Request),
	onlyAuth bool,
) {
	subRouter := mux.Methods(method).Subrouter()
	subRouter.HandleFunc(attr.apiURL+url,
		handler)
	subRouter.Use(
		loggermiddleware.RequestLogger(attr.zapLogger))

	if onlyAuth {
		subRouter.Use(
			authmiddleware.AuthMiddleware(attr.authMidAttr))
	}
}

func (p *ServerAttr) GetLogger() *zap.Logger {
	return p.zapLogger
}

func (p *ServerAttr) SetLogger(logger *zap.Logger) {
	p.zapLogger = logger
}

func (p *ServerAttr) GetmigrationsDir() string {
	return p.migrationsDir
}

func (p *ServerAttr) GetValidAddrPattern() string {
	return p.validAddrPattern
}

func (p *ServerAttr) SetValidAddrPattern(pattern string) {
	p.validAddrPattern = pattern
}

func (p *ServerAttr) GetServer() *http.Server {
	return p.server
}

func (p *ServerAttr) GetRunAddress() *string {
	return &p.runAddress
}

func (p *ServerAttr) GetDatabaseURL() *string {
	return &p.databaseURL
}

func (p *ServerAttr) GetWaitSecRespDB() time.Duration {
	return p.waitSecRespDB
}

func (p *ServerAttr) SetWaitSecRespDB(dur time.Duration) {
	p.waitSecRespDB = dur
}

func (p *ServerAttr) SetRunAddress(addr string) {
	p.runAddress = addr
}

func (p *ServerAttr) SetDatabaseURL(databaseURL string) {
	p.databaseURL = databaseURL
}

func (p *ServerAttr) SetPgxConn(
	ctxDB context.Context,
) error {
	dbConn, err := pgx.Connect(ctxDB, p.databaseURL)
	if err != nil {
		return fmt.Errorf("initiateDBconn->pgx.Connect %w", err)
	}

	p.pgxConn = dbConn

	return nil
}
