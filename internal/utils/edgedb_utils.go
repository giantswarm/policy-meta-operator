package utils

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/edgedb/edgedb-go"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	EDGEDB_DSN = "edgedb://?password_file=/etc/edgedb/password&tls_ca_file_env=EDGEDB_TLS_CA_FILE&user_env=EDGEDB_USER&host_env=EDGEDB_HOST&port_env=EDGEDB_PORT"
)

type AutomatedException struct {
	ID                 edgedb.UUID             `edgedb:"id"`
	Name               string                  `edgedb:"name"`
	LastReconciliation edgedb.OptionalDateTime `edgedb:"last_reconciliation"`
	Counter            edgedb.OptionalInt64    `edgedb:"counter"`
}

type PolicyException struct {
	ID                 edgedb.UUID             `edgedb:"id"`
	Name               string                  `edgedb:"name"`
	LastReconciliation edgedb.OptionalDateTime `edgedb:"last_reconciliation"`
	Counter            edgedb.OptionalInt64    `edgedb:"counter"`
}

//go:embed setupAutomatedExceptionType.edgeql
var setupAutomatedExceptionTypeQuery string

//go:embed setupAutomatedExceptionType.edgeql
var setupPolicyExceptionTypeQuery string

//go:embed insertAutomatedException.edgeql
var insertAutomatedExceptionQuery string

//go:embed insertPolicyException.edgeql
var insertPolicyExceptionQuery string

func GetEDGEDBClient(ctx context.Context, opts edgedb.Options) *edgedb.Client {
	_ = log.FromContext(ctx)
	client, err := edgedb.CreateClientDSN(ctx, EDGEDB_DSN, opts)
	if err != nil {
		log.Log.Error(err, "Error creating edgedb client")
	}
	return client
}

func CloseClient(client *edgedb.Client) {
	err := client.Close()
	if err != nil {
		log.Log.Error(err, "Error closing edgedb client")
	}
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

func SetupPolicyExceptionType(ctx context.Context, client *edgedb.Client) (edgedb.Optional, error) {
	var result edgedb.Optional

	err := client.QuerySingle(
		ctx,
		setupPolicyExceptionTypeQuery,
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

func InsertPolicyException(ctx context.Context, client *edgedb.Client, args ...interface{}) (PolicyException, error) {
	var result PolicyException

	err := client.QuerySingle(
		ctx,
		insertPolicyExceptionQuery,
		&result,
		args...,
	)

	return result, err
}

func GetPolicyExceptionsFromEdgeDB(ctx context.Context, client *edgedb.Client) {
	// Select users.
	var output []PolicyException
	args := map[string]interface{}{"name": "my-deployment-exception"}
	query := "SELECT PolicyException {name, counter, last_reconciliation} FILTER .name = <str>$name"
	err := client.Query(ctx, query, &output, args)

	fmt.Println(output)
	if err != nil {
		log.Log.Error(err, "Error querying for PolicyException")
	}
}
