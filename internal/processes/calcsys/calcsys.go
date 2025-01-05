package calcsys

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/dmitrovia/gophermart/internal/logger"
	"github.com/dmitrovia/gophermart/internal/models/processesattr/calcsysattr"
)

func RunProcess(waitG *sync.WaitGroup,
	attr *calcsysattr.CalcSysAttr,
) {
	fmt.Println("Run calcsys")

	ctxDB, cancel := context.WithTimeout(
		context.Background(), attr.GetWaitSecRespDB())

	defer cancel()

	err := attr.SetPgxConn(ctxDB)
	if err != nil {
		doErr(waitG, err, "RP->SetPgxConn", attr)

		return
	}

	attr.Init()

	updateStatusOrdersAndCalculatePoints(attr, waitG)

	fmt.Println("End calcsys")
}

func doErr(waitG *sync.WaitGroup,
	err error,
	errMsg string,
	attr *calcsysattr.CalcSysAttr,
) {
	logger.DoInfoLogFromErr(
		errMsg, err, attr.GetLogger())
	waitG.Done()
}

func updateStatusOrdersAndCalculatePoints(
	attr *calcsysattr.CalcSysAttr,
	waitG *sync.WaitGroup,
) {
	channelCancel := make(chan os.Signal, 1)
	signal.Notify(channelCancel,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT)

	for {
		select {
		case <-channelCancel:
			waitG.Done()

			return
		case <-time.After(
			attr.GetInetervalCallCalcService()):
			err := attr.GetCalculateService().
				UpdateStatusOrdersAndCalculatePoints()
			if err != nil {
				logger.DoInfoLogFromErr(
					"updateStatusOrdersAndCalculatePoints",
					err,
					attr.GetLogger())
			}
		}
	}
}
