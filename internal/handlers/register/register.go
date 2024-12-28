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
	statusBR  = http.StatusBadRequest
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
	regUser := &apimodels.InRegisterUser{}

	err := getReqData(req, regUser)
	if err != nil {
		setErr(writer, h.attr, err, statusBR, "getReqData")

		return
	}

	isValid := validate(regUser)
	if !isValid {
		setErr(writer, h.attr, err, statusBR, "validate")

		return
	}

	exist, _, err := h.authService.UserIsExist(regUser.Login)
	if err != nil {
		setErr(writer, h.attr, err, statusISE, "UserIsExist")

		return
	}

	if exist {
		writer.WriteHeader(http.StatusConflict)

		return
	}

	err = CreateUser(h, regUser)
	if err != nil {
		setErr(writer, h.attr, err, statusISE, "CreateUser")

		return
	}

	token, err := generateToken(regUser.Login, h.attr)
	if err != nil {
		setErr(writer, h.attr, err, statusISE, "generateToken")

		return
	}

	writer.Header().Set("Authorization", token)
	writer.WriteHeader(http.StatusOK)
}

func CreateUser(handler *Register,
	regUser *apimodels.InRegisterUser,
) error {
	passwHash, err := cryptPass(regUser.Password)
	if err != nil {
		return fmt.Errorf(
			"CreateUser->cryptPass %w", err)
	}

	user := &usermodel.User{}
	user.SetLogin(regUser.Login)
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
	status int,
	method string,
) {
	writer.WriteHeader(status)
	logger.DoInfoLogFromErr("register->"+method,
		err, inAttr.GetLogger())
}

func validate(user *apimodels.InRegisterUser) bool {
	if user.Login == "" || user.Password == "" {
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
	user *apimodels.InRegisterUser,
) error {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("getReqData->io.ReadAll %w", err)
	}

	if len(body) == 0 {
		return fmt.Errorf("getReqData: %w", errEmptyData)
	}

	err = json.Unmarshal(body, user)
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
