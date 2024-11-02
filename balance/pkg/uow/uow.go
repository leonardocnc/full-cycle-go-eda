package uow

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type RepositoryFactory func(tx *sql.Tx) interface{}

type UowInterface interface {
	Register(name string, fc RepositoryFactory)
	GetRepository(ctx context.Context, name string) (interface{}, error)
	Do(ctx context.Context, fn func(uow *Uow) error) error
	CommitOrRollback() error
	Rollback() error
	UnRegister(name string)
}

type Uow struct {
	Db           *sql.DB
	Tx           *sql.Tx
	Repositories map[string]RepositoryFactory
}

func NewUow(ctx context.Context, db *sql.DB) *Uow {
	return &Uow{
		Db:           db,
		Repositories: make(map[string]RepositoryFactory),
	}
}

func (u *Uow) Register(name string, fc RepositoryFactory) {
	u.Repositories[name] = fc
}

func (u *Uow) UnRegister(name string) {
	delete(u.Repositories, name)
}

func (u *Uow) GetRepository(ctx context.Context, name string) (interface{}, error) {
	if u.Tx == nil {
		if err := u.beginTransaction(ctx); err != nil {
			return nil, err
		}
	}
	repo, exists := u.Repositories[name]
	if !exists {
		return nil, fmt.Errorf("repository not found: %s", name)
	}
	return repo(u.Tx), nil
}

func (u *Uow) Do(ctx context.Context, fn func(uow *Uow) error) error {
	if u.Tx != nil {
		return fmt.Errorf("transaction already started")
	}
	if err := u.beginTransaction(ctx); err != nil {
		return err
	}
	if err := fn(u); err != nil {
		return u.rollbackWithError(err)
	}
	return u.CommitOrRollback()
}

func (u *Uow) beginTransaction(ctx context.Context) error {
	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	u.Tx = tx
	return nil
}

func (u *Uow) rollbackWithError(originalErr error) error {
	if err := u.Rollback(); err != nil {
		return fmt.Errorf("original error: %s, rollback error: %s", originalErr.Error(), err.Error())
	}
	return originalErr
}

func (u *Uow) Rollback() error {
	if u.Tx == nil {
		return errors.New("no transaction to rollback")
	}
	if err := u.Tx.Rollback(); err != nil {
		return err
	}
	u.Tx = nil
	return nil
}

func (u *Uow) CommitOrRollback() error {
	if err := u.Tx.Commit(); err != nil {
		return u.rollbackWithError(err)
	}
	u.Tx = nil
	return nil
}
