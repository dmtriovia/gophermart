package register_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dmitrovia/gophermart/internal/handlers/register"
	"github.com/dmitrovia/gophermart/internal/logger"
	"github.com/dmitrovia/gophermart/internal/models/apimodels"
	"github.com/dmitrovia/gophermart/internal/models/handlerattr/registerattr"
	"github.com/dmitrovia/gophermart/internal/service/accountservice"
	"github.com/dmitrovia/gophermart/internal/service/authservice"
	"github.com/dmitrovia/gophermart/internal/storage/accountstorage"
	"github.com/dmitrovia/gophermart/internal/storage/userstorage"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const url string = "http://localhost:8080"

const stok int = http.StatusOK

const post string = "POST"

const databaseURL = "postgres://postgres:postgres@" +
	"localhost:5432/praktikum?sslmode=disable"

const waitSecRespDB = 10 * time.Second

type testData struct {
	tn     string
	login  string
	pass   string
	expcod int
	exbody string
	meth   string
}

func getTestData() *[]testData {
	return &[]testData{
		{
			meth: post, tn: "1", login: "dmitrovia2",
			pass: "temppass", expcod: stok, exbody: "",
		},
	}
}

func TestRegisterHandler(t *testing.T) {
	t.Helper()
	t.Parallel()

	testCases := getTestData()

	handler, err := Init()
	if err != nil {
		return
	}

	for _, test := range *testCases {
		t.Run(http.MethodPost, func(t *testing.T) {
			t.Parallel()

			reqData, err := formReqBody(&test)
			if err != nil {
				fmt.Println(err)

				return
			}

			req, err := http.NewRequestWithContext(
				context.Background(),
				test.meth,
				url+"/api/user/register", reqData)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			newr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/api/user/register",
				handler.RegisterHandler)
			router.ServeHTTP(newr, req)
			status := newr.Code
			body, _ := io.ReadAll(newr.Body)

			assert.Equal(t,
				test.expcod,
				status, test.tn+": Response code didn't match expected")

			if test.exbody != "" {
				assert.JSONEq(t, test.exbody, string(body))
			}
		})
	}
}

func Init() (*register.Register,
	error,
) {
	accStorage := &accountstorage.AccountStorage{}
	userStorage := &userstorage.UserStorage{}
	attr := &registerattr.RegisterAttr{}

	ctxDB, cancel := context.WithTimeout(
		context.Background(), waitSecRespDB)

	defer cancel()

	pgxConn, err := NewPgxConn(ctxDB)
	if err != nil {
		fmt.Println(err)

		return nil, err
	}

	logger, err := newLogger()
	if err != nil {
		fmt.Println(err)

		return nil, err
	}

	accStorage.Initiate(pgxConn)
	userStorage.Initiate(pgxConn)

	auths := authservice.NewAuthService(
		userStorage, waitSecRespDB)
	accs := accountservice.NewAccountService(
		accStorage, waitSecRespDB)

	attr.Init(logger)

	hand := register.NewRegisterHandler(
		auths, accs, attr)

	return hand, nil
}

func NewPgxConn(
	ctxDB context.Context,
) (*pgx.Conn, error) {
	dbConn, err := pgx.Connect(ctxDB, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("NewPgxConn->pgx.Connect %w", err)
	}

	return dbConn, nil
}

func newLogger() (*zap.Logger, error) {
	logger, err := logger.Initialize("info")
	if err != nil {
		return nil, fmt.Errorf(
			"newLogger->logger.Initialize %w",
			err)
	}

	return logger, nil
}

func formReqBody(
	data *testData,
) (*bytes.Reader, error) {
	register := &apimodels.InRegisterUser{}
	register.Login = data.login
	register.Password = data.pass

	marshall, err := json.Marshal(register)
	if err != nil {
		return nil,
			fmt.Errorf("formReqBody->Marshal: %w",
				err)
	}

	return bytes.NewReader(marshall), nil
}
