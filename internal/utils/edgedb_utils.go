package utils

import (
	"context"
	_ "embed"

	"github.com/edgedb/edgedb-go"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	EDGEDB_DSN = "edgedb://?password_env=EDGEDB_PASSWORD&tls_ca_file_env=EDGEDB_TLS_CA_FILE&user_env=EDGEDB_USER&host_env=EDGEDB_HOST&port_env=EDGEDB_PORT"
)

type AutomatedException struct {
	ID                 edgedb.UUID             `edgedb:"id"`
	Name               string                  `edgedb:"name"`
	LastReconciliation edgedb.OptionalDateTime `edgedb:"last_reconciliation"`
	Counter            edgedb.OptionalInt64    `edgedb:"counter"`
}

//go:embed setupAutomatedExceptionType.edgeql
var setupAutomatedExceptionTypeQuery string

//go:embed insertAutomatedException.edgeql
var insertAutomatedExceptionQuery string

func GetEDGEDBClient(ctx context.Context, opts edgedb.Options) *edgedb.Client {
	_ = log.FromContext(ctx)
	client, err := edgedb.CreateClientDSN(ctx, EDGEDB_DSN, opts)
	if err != nil {
		log.Log.Error(err, "Error creating edgedb client")
	}
	return client
}

func SetupAutomatedExceptionType(ctx context.Context, client *edgedb.Client) (edgedb.Optional, error) {
	var result edgedb.Optional

	err := client.QuerySingle(
		ctx,
		setupAutomatedExceptionTypeQuery,
		&result,
	)

	return result, err
}

func InsertAutomatedException(ctx context.Context, client *edgedb.Client, args ...interface{}) (AutomatedException, error) {
	var result AutomatedException

	err := client.QuerySingle(
		ctx,
		insertAutomatedExceptionQuery,
		&result,
		args...,
	)

	return result, err
}
