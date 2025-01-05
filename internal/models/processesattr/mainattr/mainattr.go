package mainattr

import (
	"fmt"
	"time"

	"github.com/dmitrovia/gophermart/internal/logger"
	"github.com/dmitrovia/gophermart/internal/models/processesattr/calcsysattr"
	"github.com/dmitrovia/gophermart/internal/models/processesattr/serverattr"
	"go.uber.org/zap"
)

const initwaitSecRespDB = 10

const initWaitSecRespCalcService = 10

const initInetervalCallCalcService = 2

type MainAttr struct {
	calcProcAttr         *calcsysattr.CalcSysAttr
	serverProcAttr       *serverattr.ServerAttr
	zapLogger            *zap.Logger
	defDatabaseURL       string
	defPORT              string
	defAccSysAddr        string
	runAddress           string
	databaseURL          string
	accrualSystemAddress string
}

func (p *MainAttr) GetRunAddress() *string {
	return &p.runAddress
}

func (p *MainAttr) SetRunAddress(addr string) {
	p.runAddress = addr
}

func (p *MainAttr) GetAccrualSystemAddress() *string {
	return &p.accrualSystemAddress
}

func (p *MainAttr) SetAccrualSystemAddress(addr string) {
	p.accrualSystemAddress = addr
}

func (p *MainAttr) GetDefDatabaseURL() string {
	return p.defDatabaseURL
}

func (p *MainAttr) GetDefAccrualSystemAddress() string {
	return p.defAccSysAddr
}

func (p *MainAttr) GetDefRunAddress() string {
	return p.defPORT
}

func (p *MainAttr) GetDatabaseURL() *string {
	return &p.databaseURL
}

func (
	p *MainAttr,
) GetCalcProcAttr() *calcsysattr.CalcSysAttr {
	return p.calcProcAttr
}

func (
	p *MainAttr,
) GetServerProcAttr() *serverattr.ServerAttr {
	return p.serverProcAttr
}

func (p *MainAttr) GetLogger() *zap.Logger {
	return p.zapLogger
}

func (p *MainAttr) PreInitProcesses() {
	p.serverProcAttr.SetDatabaseURL(p.databaseURL)
	p.calcProcAttr.SetDatabaseURL(p.databaseURL)

	p.serverProcAttr.SetRunAddress(p.runAddress)
	p.calcProcAttr.SetAccrualSystemAddress(
		p.accrualSystemAddress)

	p.calcProcAttr.SetLogger(p.zapLogger)
	p.serverProcAttr.SetLogger(p.zapLogger)
}

func (p *MainAttr) Init() error {
	p.calcProcAttr = &calcsysattr.CalcSysAttr{}
	p.serverProcAttr = &serverattr.ServerAttr{}

	zapLogLevel := "info"
	pattern := "^[a-zA-Z/ ]{1,100}:[0-9]{1,10}$"

	interval := initwaitSecRespDB * time.Second
	p.calcProcAttr.SetWaitSecRespDB(interval)
	p.serverProcAttr.SetWaitSecRespDB(interval)
	p.GetCalcProcAttr().SetInetervalCallCalcService(
		initInetervalCallCalcService * time.Second)
	p.GetCalcProcAttr().SetWaitSecRespCalcService(
		initWaitSecRespCalcService * time.Second)

	p.defAccSysAddr = "localhost:8090"
	p.defDatabaseURL = "postgres://postgres:postgres@" +
		"localhost:5432/praktikum?sslmode=disable"
	p.GetServerProcAttr().SetValidAddrPattern(pattern)
	p.defPORT = "localhost:8080"

	logger, err := logger.Initialize(zapLogLevel)
	if err != nil {
		return fmt.Errorf(
			"PreInit->logger.Initialize %w",
			err)
	}

	p.zapLogger = logger

	return nil
}
