package edgedbutils

import (
	"context"
	_ "embed"

	"github.com/edgedb/edgedb-go"
	policyAPI "github.com/giantswarm/policy-api/api/v1alpha1"
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

func translateTargetsToEdgedbTypes(targets []policyAPI.Target) []Target {
	var edgedbTarget []Target

	for _, target := range targets {
		edgedbTarget = append(edgedbTarget, Target{
			Names:      target.Names,
			Namespaces: target.Namespaces,
			Kind:       target.Kind,
		})
	}

	return edgedbTarget
}

//go:embed insertPolicyException.edgeql
var insertPolicyExceptionQuery string

func InsertPolicyException(ctx context.Context, client *edgedb.Client, policyException policyAPI.PolicyException) (Exception, error) {
	var edgedbException Exception

	// Temporary hard code fields
	policyName := policyException.Spec.Policies
	targetNames := translateTargetsToEdgedbTypes(policyException.Spec.Targets)[0].Names
	targetKind := translateTargetsToEdgedbTypes(policyException.Spec.Targets)[0].Kind
	targetNamespaces := translateTargetsToEdgedbTypes(policyException.Spec.Targets)[0].Namespaces
	policyExceptionName := policyException.Name

	params := []interface{}{
		policyName,
		targetNames,
		targetNamespaces,
		targetKind,
		policyExceptionName,
	}

	err := client.QuerySingle(
		ctx,
		insertPolicyExceptionQuery,
		&edgedbException,
		params...,
	)

	return edgedbException, err
}
