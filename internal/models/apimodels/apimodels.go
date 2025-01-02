package apimodels

import "time"

type InRegisterUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type InLoginUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type InSetOrder struct {
	Identifier string `json:"identifier"`
}

type InWithdraw struct {
	PointsWriteOff  float32 `json:"sum"`
	OrderIdentifier string  `json:"order"`
}

type InGetOrder struct {
	Identifier string
}

type OutGetOrders struct {
	Identifier  *string    `json:"number"`
	Createddate *time.Time `json:"uploaded_at"`
	Status      *string    `json:"status"`
	Accrual     *int32     `json:"accrual,omitempty"`
}

func (o *OutGetOrders) SetOutGetOrders(
	identifier *string,
	createddate *time.Time,
	status *string,
	accrual *int32,
) {
	o.Identifier = identifier
	o.Createddate = createddate
	o.Status = status
	o.Accrual = accrual
}

type OutBalance struct {
	Points    *float32 `json:"current"`
	Withdrawn *float32 `json:"withdrawn"`
}

func (o *OutBalance) SetOutBalance(
	points *float32,
	withdrawn *float32,
) {
	o.Points = points
	o.Withdrawn = withdrawn
}

type OutWithdrawals struct {
	PointsWriteOff  *float32   `json:"sum"`
	OrderIdentifier *string    `json:"order"`
	Createddate     *time.Time `json:"processed_at"`
}

func (o *OutWithdrawals) SetOutWithdrawals(
	pointsWriteOff *float32,
	orderIdentifier *string,
	createddate *time.Time,
) {
	o.PointsWriteOff = pointsWriteOff
	o.OrderIdentifier = orderIdentifier
	o.Createddate = createddate
}

type OutGetOrder struct {
	Identifier *string `json:"number"`
	Status     *string `json:"status"`
	Accrual    *int32  `json:"accrual,omitempty"`
}

func (o *OutGetOrder) SetOutGetOrder(
	identifier *string,
	status *string,
	accrual *int32,
) {
	o.Identifier = identifier
	o.Status = status
	o.Accrual = accrual
}
