package edgedbutils

import (
	"context"
	_ "embed"

	"github.com/edgedb/edgedb-go"
)

type PolicyManifest struct {
	ID         edgedb.UUID `edgedb:"id"`
	Exceptions []Exception `edgedb:"exceptions"`
	Name       string      `edgedb:"name"`
}

type Exception struct {
	ID      edgedb.UUID `edgedb:"id"`
	Name    string      `edgedb:"name"`
	Targets []Target    `edgedb:"targets"`
}

type Target struct {
	ID         edgedb.UUID `edgedb:"id"`
	Names      []string    `edgedb:"names"`
	Namespaces []string    `edgedb:"namespaces"`
	Kind       string      `edgedb:"kind"`
}

//go:embed setupTypes.edgeql
var setupTypesQuery string

func SetupTypes(ctx context.Context, client *edgedb.Client) (edgedb.Optional, error) {
	var result edgedb.Optional

	err := client.QuerySingle(
		ctx,
		setupTypesQuery,
		&result,
	)

	return result, err
}
