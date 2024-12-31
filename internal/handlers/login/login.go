package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dmitrovia/gophermart/internal/logger"
	"github.com/dmitrovia/gophermart/internal/models/apimodels"
	"github.com/dmitrovia/gophermart/internal/models/handlerattr/loginattr"
	"github.com/dmitrovia/gophermart/internal/service"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var errEmptyData = errors.New("data is empty")

type Login struct {
	serv service.AuthService
	attr *loginattr.LoginAttr
}

func NewLoginHandler(
	s service.AuthService,
	inAttr *loginattr.LoginAttr,
) *Login {
	return &Login{serv: s, attr: inAttr}
}

func (h *Login) LoginHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	reqAttr := &apimodels.InLoginUser{}

	err := getReqData(req, reqAttr)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		logger.DoInfoLogFromErr("login->getReqData",
			err, h.attr.GetLogger())

		return
	}

	isValid := validate(reqAttr)
	if !isValid {
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	exist, user, err := h.serv.UserIsExist(reqAttr.Login)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		logger.DoInfoLogFromErr("login->UserIsExist",
			err, h.attr.GetLogger())

		return
	}

	if !exist {
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = checkPass(user.GetPassword(), reqAttr.Password)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		logger.DoInfoLogFromErr("login->checkPass",
			err, h.attr.GetLogger())

		return
	}

	token, err := generateToken(reqAttr.Login, h.attr)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		logger.DoInfoLogFromErr("login->generateToken",
			err, h.attr.GetLogger())

		return
	}

	writer.Header().Set("Authorization", token)
	writer.WriteHeader(http.StatusOK)
}

func generateToken(
	id string,
	attr *loginattr.LoginAttr,
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

func checkPass(hash string, pass string) error {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash), []byte(pass))
	if err != nil {
		return fmt.Errorf(
			"checkPass->bcrypt.CompareHashAndPassword %w", err)
	}

	return nil
}

func validate(reqAttr *apimodels.InLoginUser) bool {
	if reqAttr.Login == "" || reqAttr.Password == "" {
		return false
	}

	return true
}

func getReqData(
	req *http.Request,
	reqAttr *apimodels.InLoginUser,
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
