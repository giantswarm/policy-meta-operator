package utils

import (
	"context"
	_ "embed"

	"github.com/edgedb/edgedb-go"
	edgedbutils "github.com/giantswarm/policy-meta-operator/internal/utils/edgedb"
)

//go:embed insertAutomatedException.edgeql
var insertAutomatedExceptionQuery string

//go:embed insertPolicyException.edgeql
var insertPolicyExceptionQuery string

func InsertAutomatedException(ctx context.Context, client *edgedb.Client, args ...interface{}) (edgedbutils.Exception, error) {
	var result edgedbutils.Exception

	err := client.QuerySingle(
		ctx,
		insertAutomatedExceptionQuery,
		&result,
		args...,
	)

	return result, err
}

func InsertPolicyException(ctx context.Context, client *edgedb.Client, args ...interface{}) (edgedbutils.Exception, error) {
	var result edgedbutils.Exception

	err := client.QuerySingle(
		ctx,
		insertPolicyExceptionQuery,
		&result,
		args...,
	)

	return result, err
}

// func GetPolicyExceptionsFromEdgeDB(ctx context.Context, client *edgedb.Client) []PolicyException {
// 	// Select users.
// 	var output []PolicyException
// 	query := "SELECT PolicyException {name, counter, last_reconciliation}"
// 	err := client.Query(ctx, query, &output)
// 	if err != nil {
// 		log.Log.Error(err, "Error querying for PolicyException")
// 	}

// 	return output
// }
