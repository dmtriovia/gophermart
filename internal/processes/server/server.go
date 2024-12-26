package server

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/dmitrovia/gophermart/internal/functions/validatef"
	"github.com/dmitrovia/gophermart/internal/logger"
	"github.com/dmitrovia/gophermart/internal/migrator"
	"github.com/dmitrovia/gophermart/internal/models/serverattr"
)

var errParseFlags = errors.New("addr is not valid")

//go:embed db/migrations/*.sql
var MigrationsFS embed.FS

func RunProcess() {
	attr := &serverattr.ServerAttr{}

	err := attr.Init()
	if err != nil {
		fmt.Println("RunProcess->attr.Init: %w", err)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(attr.GetWaitSecRespDB()))

	defer cancel()

	err = UseMigrations(attr)
	if err != nil {
		logger.DoInfoLogFromErr(
			"UseMigrations", err, attr.GetLogger())

		return
	}

	err = attr.SetPgxConn(ctx)
	if err != nil {
		logger.DoInfoLogFromErr(
			"SetPgxConn", err, attr.GetLogger())

		return
	}

	initiateFlags(attr)

	err = initSystemAttrs(attr)
	if err != nil {
		logger.DoInfoLogFromErr(
			"initSystemAttrs", err, attr.GetLogger())

		return
	}

	err = runServer(attr)
	if err != nil {
		logger.DoInfoLogFromErr(
			"runServer", err, attr.GetLogger())

		return
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
	if err != nil {
		return fmt.Errorf(
			"runServer->GetServer().ListenAndServe() %w", err)
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
