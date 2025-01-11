// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

package apis

import (
	"context"
	"database/sql"
	"time"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/model"
)

func (r ApiRequest) log(format string, a ...any) {
	configuration.Log(r.p.DbName, r.p.UserName, format, a...)
}

func (r ApiRequest) trace(format string, a ...any) {
	configuration.Trace(r.p.DbName, r.p.UserName, format, a...)
}

func (r ApiRequest) startTiming() time.Time {
	return configuration.StartTiming()
}

func (r ApiRequest) stopTiming(funcName string, start time.Time) {
	configuration.StopTiming(r.p.DbName, r.p.UserName, funcName, start)
}

func (r ApiRequest) logAndReturnError(format string, a ...any) error {
	return configuration.LogAndReturnError(r.p.DbName, r.p.UserName, format, a...)
}

func (r ApiRequest) execContext(query string, args ...any) error {
	_, err := r.db.ExecContext(r.ctx, query, args...)
	return err
}

func (r ApiRequest) queryContext(query string, args ...any) (*sql.Rows, error) {
	return r.db.QueryContext(r.ctx, query, args...)
}

func (r ApiRequest) beginTransaction() error {
	r.trace("beginTransaction()")
	if err := r.execContext(getBeginTransactionQuery()); err != nil {
		return r.logAndReturnError("Begin transaction error: %v.", err)
	}
	r.log("Begin transaction.")
	return nil
}

// function is available for mocking
var getBeginTransactionQuery = func() string {
	return "BEGIN"
}

func (r ApiRequest) commitTransaction() error {
	r.trace("commitTransaction()")
	if err := r.execContext(getCommitTransactionQuery()); err != nil {
		return r.logAndReturnError("Commit transaction error: %v.", err)
	}
	r.log("Commit transaction.")
	return nil
}

// function is available for mocking
var getCommitTransactionQuery = func() string {
	return "COMMIT"
}

func (r ApiRequest) rollbackTransaction() error {
	r.trace("rollbackTransaction()")
	if err := r.execContext(getRollbackTransactionQuery()); err != nil {
		return r.logAndReturnError("Rollback transaction error: %v.", err)
	}
	r.log("ROLLBACK transaction.")
	return nil
}

// function is available for mocking
var getRollbackTransactionQuery = func() string {
	return "ROLLBACK"
}

func createContextAndApiRequest(ct context.Context, p ApiParameters, uri string) (request ApiRequest, cancelFunc context.CancelFunc, error error) {
	ctx, cancel := configuration.GetContextWithTimeOut(ct, p.DbName)
	db, err := database.GetDbByName(p.DbName)
	r := ApiRequest{
		ctx:         ctx,
		p:           p,
		db:          db,
		ctxChan:     make(chan GridResponse, 1),
		transaction: model.GetNewRowWithUuid(),
	}
	return r, cancel, err
}
