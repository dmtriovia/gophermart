package getstatusfromcalcsystemattr

import (
	"net/http"

	"go.uber.org/zap"
)

type GetStatusFromCalcSystemAttr struct {
	client      *http.Client
	method      string
	contentType string
	urlForReq   string
	defURL      string
	zapLogger   *zap.Logger
}

func (r *GetStatusFromCalcSystemAttr) Init(
	logger *zap.Logger,
	accrualSystemAddress string,
) {
	r.client = &http.Client{}
	r.method = http.MethodGet
	r.defURL = "http://" + accrualSystemAddress +
		"/api/orders/"
	r.contentType = "text/plain"
	r.zapLogger = logger
}

func (
	r *GetStatusFromCalcSystemAttr) SetURLForReq(
	url string,
) {
	r.urlForReq = url
}

func (
	r *GetStatusFromCalcSystemAttr,
) GetLogger() *zap.Logger {
	return r.zapLogger
}

func (
	r *GetStatusFromCalcSystemAttr,
) GetClient() *http.Client {
	return r.client
}

func (
	r *GetStatusFromCalcSystemAttr) GetMethod() string {
	return r.method
}

func (
	r *GetStatusFromCalcSystemAttr) GetContentType() string {
	return r.contentType
}

func (
	r *GetStatusFromCalcSystemAttr) GetDefURL() string {
	return r.defURL
}

func (
	r *GetStatusFromCalcSystemAttr) GetURLForReq() string {
	return r.urlForReq
}
