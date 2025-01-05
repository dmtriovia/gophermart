package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/dmitrovia/gophermart/internal/functions/validatef"
	"github.com/dmitrovia/gophermart/internal/logger"
	"github.com/dmitrovia/gophermart/internal/models/processesattr/mainattr"
	"github.com/dmitrovia/gophermart/internal/processes/calcsys"
	"github.com/dmitrovia/gophermart/internal/processes/server"
)

const procCnt = 2

var errParseFlags = errors.New("addr is not valid")

func main() {
	fmt.Println("Run main")

	mainAttr := &mainattr.MainAttr{}

	err := mainAttr.Init()
	if err != nil {
		logger.DoInfoLogFromErr(
			"main->Init", err, mainAttr.GetLogger())

		return
	}

	initiateFlags(mainAttr)

	err = initSystemAttrs(mainAttr)
	if err != nil {
		logger.DoInfoLogFromErr(
			"main->initSystemAttrs", err, mainAttr.GetLogger())

		return
	}

	mainAttr.PreInitProcesses()

	waitGroup := new(sync.WaitGroup)
	go server.RunProcess(waitGroup,
		mainAttr.GetServerProcAttr())
	go calcsys.RunProcess(waitGroup,
		mainAttr.GetCalcProcAttr())

	waitGroup.Add(procCnt)
	waitGroup.Wait()
	fmt.Println("End main")
}

func initiateFlags(matr *mainattr.MainAttr) {
	flag.StringVar(
		matr.GetDatabaseURL(),
		"d", matr.GetDefDatabaseURL(),
		"database connection address.",
	)
	flag.StringVar(
		matr.GetRunAddress(),
		"a", matr.GetDefRunAddress(),
		"Port to listen on.",
	)
	flag.StringVar(
		matr.GetAccrualSystemAddress(),
		"r", matr.GetDefAccrualSystemAddress(),
		"address of the accrual calculation system",
	)
	flag.Parse()
}

func initSystemAttrs(matr *mainattr.MainAttr) error {
	RunAddress := os.Getenv("RUN_ADDRESS")
	DatabaseURL := os.Getenv("DATABASE_URI")
	AccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")

	if RunAddress != "" {
		res, err := validatef.IsMatchesTemplate(
			RunAddress,
			matr.GetServerProcAttr().GetValidAddrPattern())
		if err != nil {
			return fmt.Errorf(
				"initSystemAttrs->IsMatchesTemplate: %w",
				err)
		}

		if !res {
			return errParseFlags
		}

		matr.GetServerProcAttr().SetRunAddress(RunAddress)
	}

	if DatabaseURL != "" {
		matr.GetServerProcAttr().SetDatabaseURL(DatabaseURL)
		matr.GetCalcProcAttr().SetDatabaseURL(DatabaseURL)
	}

	if AccrualSystemAddress != "" {
		matr.GetCalcProcAttr().SetAccrualSystemAddress(
			AccrualSystemAddress)
	}

	return nil
}
