package server

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dmitrovia/gophermart/internal/logger"
	"github.com/dmitrovia/gophermart/internal/migrator"
	"github.com/dmitrovia/gophermart/internal/models/processesattr/serverattr"
)

//go:embed db/migrations/*.sql
var MigrationsFS embed.FS

func RunProcess(
	waitG *sync.WaitGroup,
	attr *serverattr.ServerAttr,
) {
	fmt.Println("Run server process")

	ctxDB, cancel := context.WithTimeout(
		context.Background(), attr.GetWaitSecRespDB())

	defer cancel()

	err := attr.SetPgxConn(ctxDB)
	if err != nil {
		doErr(waitG, err, "RP->SetPgxConn", attr)

		return
	}

	err = attr.Init()
	if err != nil {
		doErr(waitG, err, "RP->Init", attr)

		return
	}

	err = UseMigrations(attr)
	if err != nil {
		doErr(waitG, err, "RP->UseMigrations", attr)

		return
	}

	go waitClose(attr, waitG)

	err = runServer(attr)
	if err != nil {
		doErr(waitG, err, "RP->runServer", attr)

		return
	}

	fmt.Println("End server process")
}

func doErr(waitG *sync.WaitGroup,
	err error,
	errMsg string,
	attr *serverattr.ServerAttr,
) {
	logger.DoInfoLogFromErr(
		errMsg, err, attr.GetLogger())
	waitG.Done()
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
