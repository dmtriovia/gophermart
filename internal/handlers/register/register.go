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
	"github.com/dmitrovia/gophermart/internal/models/bizmodels"
	"github.com/dmitrovia/gophermart/internal/models/handlerattr"
	"github.com/dmitrovia/gophermart/internal/service"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Register struct {
	serv service.AuthService
	attr *handlerattr.RegisterAttr
}

var errEmptyData = errors.New("data is empty")

func NewRegisterHandler(
	s service.AuthService,
	inAttr *handlerattr.RegisterAttr,
) *Register {
	return &Register{serv: s, attr: inAttr}
}

func (h *Register) RegisterHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	regUser := &apimodels.RegisterUser{}

	err := getReqData(req, regUser)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		logger.DoInfoLog("register->getReqData",
			err, h.attr.GetLogger())

		return
	}

	isValid := validate(regUser)
	if !isValid {
		writer.WriteHeader(http.StatusBadRequest)
		logger.DoInfoLog("register->validate",
			err, h.attr.GetLogger())

		return
	}

	exist, _, err := h.serv.UserIsExist(regUser.Login)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		logger.DoInfoLog("register->UserIsExist",
			err, h.attr.GetLogger())

		return
	}

	if exist {
		writer.WriteHeader(http.StatusConflict)

		return
	}

	passwHash, err := cryptPass(regUser.Password)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		logger.DoInfoLog("register->cryptPass",
			err, h.attr.GetLogger())

		return
	}

	user := &bizmodels.User{
		Login:    regUser.Login,
		Password: passwHash,
	}

	err = h.serv.CreateUser(user)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		logger.DoInfoLog("register->CreateUser",
			err, h.attr.GetLogger())

		return
	}

	token, err := generateToken(regUser.Login, h.attr)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		logger.DoInfoLog("register->generateToken",
			err, h.attr.GetLogger())

		return
	}

	writer.Header().Set("Authorization", token)
	writer.WriteHeader(http.StatusOK)
}

func validate(user *apimodels.RegisterUser) bool {
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
	user *apimodels.RegisterUser,
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
	attr *handlerattr.RegisterAttr,
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
