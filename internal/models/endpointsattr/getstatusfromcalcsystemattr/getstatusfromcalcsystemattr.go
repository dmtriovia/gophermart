package getstatusfromcalcsystemattr

import (
	"bytes"
	"net/http"
)

type GetStatusFromCalcSystemAttr struct {
	reqData     *bytes.Reader
	client      *http.Client
	method      string
	contentType string
	url         string
}

func (r *GetStatusFromCalcSystemAttr) Init(
	reqData *bytes.Reader,
	client *http.Client,
	method string,
	contentType string,
	url string,
) {
	r.reqData = reqData
	r.client = client
	r.method = method
	r.contentType = contentType
	r.url = url
}

func (
	r *GetStatusFromCalcSystemAttr,
) GetReqData() *bytes.Reader {
	return r.reqData
}

func (
	r *GetStatusFromCalcSystemAttr) GetClient() *http.Client {
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
	r *GetStatusFromCalcSystemAttr) GetURL() string {
	return r.url
}
