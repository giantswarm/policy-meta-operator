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

type Policy struct {
	ID   edgedb.UUID `edgedb:"id"`
	Name string      `edgedb:"name"`
	Mode string      `edgedb:"mode"`
}

type Exception struct {
	ID       edgedb.UUID `edgedb:"id"`
	Targets  []Target    `edgedb:"targets"`
	Policies []Policy    `edgedb:"policies"`
}

type Target struct {
	ID         edgedb.UUID `edgedb:"id"`
	Names      []string    `edgedb:"names"`
	Namespaces []string    `edgedb:"namespaces"`
	Kind       string      `edgedb:"kind"`
}

//go:embed queries/setupTypes.edgeql
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
