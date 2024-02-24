// Package kv provides a simple key-value store backed by a Raptor SQLite database.
package kv

import (
	"context"
	"time"

	"github.com/maddiesch/go-raptor"
	"github.com/maddiesch/go-raptor/statement"
	"github.com/maddiesch/go-raptor/statement/conditional"
)

const (
	KVTableName = "com.maddiesch.kv-store"

	keyName   = "Key"
	valueName = "Value"
)

func Prepare(ctx context.Context, db raptor.DB) error {
	stmt := statement.CreateTable(KVTableName).IfNotExists().PrimaryKey(keyName, statement.ColumnTypeText).Column(
		statement.Column(valueName, statement.ColumnTypeBlob).NotNull(),
		statement.Column("CreatedAt", statement.ColumnTypeInteger).NotNull().Default(`CURRENT_TIMESTAMP`),
		statement.Column("UpdatedAt", statement.ColumnTypeInteger).NotNull(),
	)

	return raptor.Transact(ctx, db, func(tx raptor.DB) error {
		if _, err := raptor.ExecStatement(ctx, tx, stmt); err != nil {
			return err
		}

		return nil
	})
}

func Set(ctx context.Context, db raptor.Executor, key string, value []byte) error {
	stmt := statement.Insert().OrReplace().Into(KVTableName).Value(valueName, value).Value(keyName, key).Value("UpdatedAt", time.Now().Unix())
	_, err := raptor.ExecStatement(ctx, db, stmt)
	return err
}

func Get(ctx context.Context, db raptor.Querier, key string) (b []byte, err error) {
	stmt := statement.Select(valueName).From(KVTableName).Where(conditional.Equal(keyName, key)).Limit(1)
	err = raptor.QueryRowStatement(ctx, db, stmt).Scan(&b)
	return
}

func Delete(ctx context.Context, db raptor.Executor, key string) error {
	stmt := statement.Delete().From(KVTableName).Where(conditional.Equal(keyName, key))
	_, err := raptor.ExecStatement(ctx, db, stmt)
	return err
}

func Exists(ctx context.Context, db raptor.Querier, key string) bool {
	if e, err := ExistsErr(ctx, db, key); err != nil {
		return false
	} else {
		return e
	}
}

func ExistsErr(ctx context.Context, db raptor.Querier, key string) (bool, error) {
	stmt := statement.Select("1").From(KVTableName).Where(conditional.Equal(keyName, key)).Limit(1)
	var exists bool
	err := raptor.QueryRowStatement(ctx, db, statement.Exists(stmt)).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
