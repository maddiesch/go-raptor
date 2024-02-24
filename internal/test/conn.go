package test

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/maddiesch/go-raptor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "embed"
)

//go:embed test_data.sql
var testDataSql string

func Setup(t TestingT) (*raptor.Conn, context.Context) {
	nh := md5.Sum([]byte(t.Name()))
	ctx := context.Background()
	conn, err := raptor.New(fmt.Sprintf("file:db-%s?mode=memory&cache=shared", hex.EncodeToString(nh[:])))
	require.NoError(t, err, "failed to open test database")

	t.Cleanup(func() {
		assert.NoError(t, conn.Close(), "failed to close test database")
	})

	_, err = conn.Exec(ctx, testDataSql)
	require.NoError(t, err, "failed to prepare the test database")

	conn.SetLogger(&TestQueryLogger{TestingT: t})

	return conn, ctx
}

type Person struct {
	ID        int64
	FirstName string
	LastName  string
}

type Pet struct {
	ID       int64
	PersonID int64  `db:"ParentID"`
	Kind     string `db:"Type"`
	Name     string
	Age      *int
	Metadata string `db:"-"`
}
