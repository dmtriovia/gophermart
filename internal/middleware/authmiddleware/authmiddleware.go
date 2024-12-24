package authmiddleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dmitrovia/gophermart/internal/logger"
	"github.com/dmitrovia/gophermart/internal/models/middlewareattr"
	"github.com/golang-jwt/jwt/v4"
)

const tokenLen = 2

var errUnexpectedMethod = errors.New("data is empty")

var errUserNotExist = errors.New("user is not exist")

func AuthMiddleware(next http.Handler,
	attr *middlewareattr.AuthMiddlewareAttr,
) http.Handler {
	return http.HandlerFunc(
		func(writer http.ResponseWriter, reader *http.Request) {
			authHeader := reader.Header.Get("Authorization")

			if authHeader == "" {
				setErrStr(writer, attr, "header Authorization is empty")

				return
			}

			authToken := strings.Split(authHeader, " ")
			isBearer := authToken[0] == "Bearer"
			isLenValid := len(authToken) == tokenLen

			if !isLenValid || !isBearer {
				setErrStr(writer, attr, "Invalid token format")

				return
			}

			token, err := parseToken(authToken[1], attr)
			if err != nil {
				setErr(writer, attr, err)

				return
			}

			isValid, err := isValidToken(token, attr)
			if err != nil {
				setErr(writer, attr, err)

				return
			}

			if !isValid {
				setErrStr(writer, attr, "token is invalid")

				return
			}

			next.ServeHTTP(writer, reader)
		})
}

func isValidToken(token *jwt.Token,
	attr *middlewareattr.AuthMiddlewareAttr,
) (bool, error) {
	if !token.Valid {
		return false, nil
	}

	claims, oka := token.Claims.(jwt.MapClaims)
	if !oka {
		return false, nil
	}

	timeNow := float64(time.Now().Unix())
	claimsExp, oka := claims["exp"].(float64)

	if !oka {
		return false, nil
	}

	if timeNow > claimsExp {
		return false, nil
	}

	login, ok := claims["login"].(string)
	if !ok {
		return false, nil
	}

	exist, _, err := attr.GetAuthService().UserIsExist(login)
	if err != nil {
		return false, fmt.Errorf(
			"isValidToken->UserIsExist %w", err)
	}

	if exist {
		return false, errUserNotExist
	}

	return true, nil
}

func parseToken(inToken string,
	attr *middlewareattr.AuthMiddlewareAttr,
) (*jwt.Token, error) {
	token, err := jwt.Parse(inToken,
		func(token *jwt.Token) (interface{}, error) {
			_, isHMAC := token.Method.(*jwt.SigningMethodHMAC)

			if !isHMAC {
				headerAlg, oka := token.Header["alg"].(string)

				if !oka {
					return nil, errUnexpectedMethod
				}

				logger.DoInfoLogFromStr("AuthMiddleware",
					"Unexpected signing method "+headerAlg,
					attr.GetLogger())

				return nil, errUnexpectedMethod
			}

			return []byte(attr.GetSecret()), nil
		})

	return token, fmt.Errorf(
		"AuthMiddleware>parseToken %w", err)
}

func setErrStr(writer http.ResponseWriter,
	attr *middlewareattr.AuthMiddlewareAttr,
	txt string,
) {
	writer.WriteHeader(http.StatusUnauthorized)
	logger.DoInfoLogFromStr("AuthMiddleware",
		txt, attr.GetLogger())
}

func setErr(writer http.ResponseWriter,
	attr *middlewareattr.AuthMiddlewareAttr,
	err error,
) {
	writer.WriteHeader(http.StatusUnauthorized)
	logger.DoInfoLogFromErr("AuthMiddleware",
		err, attr.GetLogger())
}
