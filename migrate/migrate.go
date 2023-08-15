package migrate

import (
	"context"

	"github.com/maddiesch/go-raptor"
	"github.com/maddiesch/go-raptor/statement"
	"github.com/maddiesch/go-raptor/statement/conditional"
)

const (
	MigrationTableName = "raptor_migrations"
)

type Migration struct {
	Name string
	Up   []string
	Down []string
}

func Up(ctx context.Context, db raptor.DB, m ...Migration) error {
	createTableStatement := statement.CreateTable(MigrationTableName).IfNotExists().PrimaryKey("name", statement.ColumnTypeText)
	_, err := raptor.ExecStatement(ctx, db, createTableStatement)
	if err != nil {
		return err
	}

	for _, mig := range m {
		queryExists := statement.Exists(
			statement.Select("name").From(MigrationTableName).Where(conditional.Equal("name", mig.Name)),
		)
		var exists bool
		if err := raptor.QueryRowStatement(ctx, db, queryExists).Scan(&exists); err != nil {
			return err
		}
		if exists {
			continue
		}

		err := db.Transact(ctx, func(d raptor.DB) error {
			for _, q := range mig.Up {
				if _, err := d.Exec(ctx, q); err != nil {
					return err
				}
			}

			_, err := raptor.ExecStatement(ctx, d, statement.Insert().Into(MigrationTableName).Value("name", mig.Name))

			return err
		})
		if err != nil {
			return err
		}
	}

	return nil
}
