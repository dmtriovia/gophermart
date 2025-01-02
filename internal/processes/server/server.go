package server

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dmitrovia/gophermart/internal/functions/validatef"
	"github.com/dmitrovia/gophermart/internal/logger"
	"github.com/dmitrovia/gophermart/internal/migrator"
	"github.com/dmitrovia/gophermart/internal/models/processesattr/serverattr"
)

var errParseFlags = errors.New("addr is not valid")

//go:embed db/migrations/*.sql
var MigrationsFS embed.FS

func RunProcess(waitG *sync.WaitGroup) {
	fmt.Println("Run server process")

	attr := &serverattr.ServerAttr{}

	err := attr.PreInit()
	if err != nil {
		logger.DoInfoLogFromErr(
			"RP->attr.preInit", err, attr.GetLogger())
	}

	ctxDB, cancel := context.WithTimeout(
		context.Background(),
		attr.GetWaitSecRespDB())

	defer cancel()

	initiateFlags(attr)

	err = initSystemAttrs(attr)
	if err != nil {
		logger.DoInfoLogFromErr(
			"RP->initSystemAttrs", err, attr.GetLogger())

		return
	}

	err = attr.SetPgxConn(ctxDB)
	if err != nil {
		logger.DoInfoLogFromErr(
			"RP->SetPgxConn", err, attr.GetLogger())

		return
	}

	err = attr.Init()
	if err != nil {
		fmt.Println("RP->attr.Init: %w", err)
	}

	err = UseMigrations(attr)
	if err != nil {
		logger.DoInfoLogFromErr(
			"RP->UseMigrations", err, attr.GetLogger())

		return
	}

	go waitClose(attr, waitG)

	err = runServer(attr)
	if err != nil {
		logger.DoInfoLogFromErr(
			"RP->runServer", err, attr.GetLogger())

		return
	}

	fmt.Println("End server process")
}

func waitClose(
	attr *serverattr.ServerAttr,
	waitG *sync.WaitGroup,
) {
	channelCancel := make(chan os.Signal, 1)
	signal.Notify(channelCancel,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT)

	for {
		_, ok := <-channelCancel
		if ok {
			err := attr.GetServer().Shutdown(context.TODO())
			if err != nil {
				fmt.Println("RP->Shutdown: %w", err)
			}

			waitG.Done()

			return
		}
	}
}

func initiateFlags(attr *serverattr.ServerAttr) {
	flag.StringVar(attr.GetDatabaseURL(),
		"d", attr.GetDefDatabaseURL(),
		"database connection address.")
	flag.StringVar(attr.GetAccrualSystemAddress(),
		"r", attr.GetDefAccSysAddr(),
		"address of the accrual calculation system")
	flag.StringVar(attr.GetRunAddress(),
		"a", attr.GetDefPort(), "Port to listen on.")
}

func initSystemAttrs(attr *serverattr.ServerAttr) error {
	RunAddress := os.Getenv("ADDRESS")
	DatabaseURL := os.Getenv("DATABASE_URI")
	AccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")

	if RunAddress != "" {
		res, err := validatef.IsMatchesTemplate(
			RunAddress, attr.GetValidAddrPattern())
		if err != nil {
			return fmt.Errorf(
				"initSystemAttrs->IsMatchesTemplate: %w",
				err)
		}

		if !res {
			return errParseFlags
		}

		attr.SetRunAddress(RunAddress)
	}

	if DatabaseURL != "" {
		attr.SetDatabaseURL(DatabaseURL)
	}

	if AccrualSystemAddress != "" {
		attr.SetAccrualSystemAddress(AccrualSystemAddress)
	}

	return nil
}

func runServer(attr *serverattr.ServerAttr) error {
	err := attr.GetServer().ListenAndServe()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf(
			"runServer->GetServer.ListenAndServe %w", err)
	}

	return nil
}

func UseMigrations(attr *serverattr.ServerAttr) error {
	migrator, err := migrator.MustGetNewMigrator(
		MigrationsFS, attr.GetmigrationsDir())
	if err != nil {
		return fmt.Errorf("useMigrations->MustGetNewMigrator %w",
			err)
	}

	conn, err := sql.Open("postgres", *attr.GetDatabaseURL())
	if err != nil {
		return fmt.Errorf("useMigrations->sql.Open %w", err)
	}

	defer conn.Close()

	err = migrator.ApplyMigrations(conn)
	if err != nil {
		return fmt.Errorf("useMigrations->ApplyMigrations %w",
			err)
	}

	return nil
}
