package register

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dmitrovia/gophermart/internal/logger"
	"github.com/dmitrovia/gophermart/internal/models/apimodels"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/accountmodel"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/usermodel"
	"github.com/dmitrovia/gophermart/internal/models/handlerattr/registerattr"
	"github.com/dmitrovia/gophermart/internal/service"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Register struct {
	authService    service.AuthService
	accountService service.AccountService
	attr           *registerattr.RegisterAttr
}

var errEmptyData = errors.New("data is empty")

const (
	statusISE = http.StatusInternalServerError
)

func NewRegisterHandler(
	authS service.AuthService,
	accS service.AccountService,
	inAttr *registerattr.RegisterAttr,
) *Register {
	return &Register{
		authService:    authS,
		accountService: accS,
		attr:           inAttr,
	}
}

func (h *Register) RegisterHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	reqAttr := &apimodels.InRegisterUser{}

	err := getReqData(req, reqAttr)
	if err != nil {
		setErr(writer, h.attr, err, "getReqData")

		return
	}

	isValid := validate(reqAttr)
	if !isValid {
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	exist, _, err := h.authService.UserIsExist(reqAttr.Login)
	if err != nil {
		setErr(writer, h.attr, err, "UserIsExist")

		return
	}

	if exist {
		writer.WriteHeader(http.StatusConflict)

		return
	}

	err = createUser(h, reqAttr)
	if err != nil {
		setErr(writer, h.attr, err, "CreateUser")

		return
	}

	token, err := generateToken(reqAttr.Login, h.attr)
	if err != nil {
		setErr(writer, h.attr, err, "generateToken")

		return
	}

	writer.Header().Set("Authorization", token)
	writer.WriteHeader(http.StatusOK)
}

func createUser(handler *Register,
	reqAttr *apimodels.InRegisterUser,
) error {
	passwHash, err := cryptPass(reqAttr.Password)
	if err != nil {
		return fmt.Errorf(
			"CreateUser->cryptPass %w", err)
	}

	user := &usermodel.User{}
	user.SetLogin(reqAttr.Login)
	user.SetPassword(passwHash)

	err = handler.authService.CreateUser(user)
	if err != nil {
		return fmt.Errorf(
			"CreateUser->authService.CreateUser %w", err)
	}

	acc := &accountmodel.Account{}
	acc.SetClient(user)

	err = handler.accountService.CreateAccount(acc)
	if err != nil {
		return fmt.Errorf(
			"CreateUser->accountService.CreateAccount %w", err)
	}

	return nil
}

func setErr(writer http.ResponseWriter,
	inAttr *registerattr.RegisterAttr,
	err error,
	method string,
) {
	writer.WriteHeader(statusISE)
	logger.DoInfoLogFromErr("register->"+method,
		err, inAttr.GetLogger())
}

func validate(reqAttr *apimodels.InRegisterUser) bool {
	if reqAttr.Login == "" || reqAttr.Password == "" {
		return false
	}

	return true
}

func cryptPass(pass string) (string, error) {
	passwHash, err := bcrypt.GenerateFromPassword(
		[]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf(
			"cryptPass->GenerateFromPassword %w", err)
	}

	return string(passwHash), nil
}

func getReqData(
	req *http.Request,
	reqAttr *apimodels.InRegisterUser,
) error {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("getReqData->io.ReadAll %w", err)
	}

	if len(body) == 0 {
		return fmt.Errorf("getReqData: %w", errEmptyData)
	}

	err = json.Unmarshal(body, reqAttr)
	if err != nil {
		return fmt.Errorf("getReqData->json.Unmarshal %w", err)
	}

	err = req.Body.Close()
	if err != nil {
		return fmt.Errorf("getReqData->req.Body.Close() %w", err)
	}

	return nil
}

func generateToken(
	id string,
	attr *registerattr.RegisterAttr,
) (string, error) {
	generateToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.MapClaims{
			"id": id,
			"exp": time.Now().Add(
				time.Hour * time.Duration(
					attr.GetTokenExpHour())).Unix(),
		})

	token, err := generateToken.SignedString(
		[]byte(attr.GetSecret()))
	if err != nil {
		return token, fmt.Errorf(
			"generateToken->generateToken.SignedString: %w",
			errEmptyData)
	}

	return token, nil
}
