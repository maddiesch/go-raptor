package raptor

import "context"

func Transact(ctx context.Context, db TxBroker, fn func(DB) error) error {
	return db.Transact(ctx, fn)
}

func TransactV[V any](ctx context.Context, db TxBroker, fn func(DB) (V, error)) (V, error) {
	var v V

	err := db.Transact(ctx, func(d DB) (err error) {
		v, err = fn(d)
		return
	})

	return v, err
}

func TransactV2[V1, V2 any](ctx context.Context, db TxBroker, fn func(DB) (V1, V2, error)) (V1, V2, error) {
	var v1 V1
	var v2 V2

	err := db.Transact(ctx, func(d DB) (err error) {
		v1, v2, err = fn(d)
		return
	})

	return v1, v2, err
}
