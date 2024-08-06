package polman

import (
	"context"
	_ "embed"

	"github.com/edgedb/edgedb-go"
	policyAPI "github.com/giantswarm/policy-api/api/v1alpha1"

	edgedbutils "github.com/giantswarm/policy-meta-operator/internal/utils/edgedb"
)

const (
	ComponentName = "policy-meta-operator"
	ManagedBy     = "app.kubernetes.io/managed-by"
)

//go:embed getPolicyManifestQuery.edgeql
var getPolicyManifestQuery string

func getPolicyManifest(ctx context.Context, client *edgedb.Client, args ...interface{}) (edgedbutils.PolicyManifest, error) {
	var result edgedbutils.PolicyManifest

	err := client.QuerySingle(
		ctx,
		getPolicyManifestQuery,
		&result,
		args...,
	)

	return result, err

}

// translateEdgedbExceptions translates edgedb Exceptions to policyAPI.Target type
func translateEdgedbExceptions(edgedbExceptions []edgedbutils.Exception) []policyAPI.Target {
	var newExceptions []policyAPI.Target
	for _, exception := range edgedbExceptions {
		for _, target := range exception.Targets {
			newExceptions = append(newExceptions, policyAPI.Target{
				Kind:       target.Kind,
				Names:      target.Names,
				Namespaces: target.Namespaces,
			})
		}
	}

	return newExceptions
}

func CreatePolicyManifest(ctx context.Context, client *edgedb.Client, args string) (policyAPI.PolicyManifest, error) {
	var policyManifest policyAPI.PolicyManifest

	edgedbPolman, err := getPolicyManifest(ctx, client, args)
	if err != nil {
		return policyManifest, err
	}

	// Set GroupVersionKind
	policyManifest.SetGroupVersionKind(policyAPI.GroupVersion.WithKind("PolicyManifest"))

	// Set Arguments
	// TODO: REPLACE THIS WHEN WE HAVE POLICYCONFIGS
	policyManifest.Spec.Args = []string{"test args"}

	// Set Name
	policyManifest.Name = edgedbPolman.Name

	// Set Labels
	policyManifest.Labels = make(map[string]string)
	policyManifest.Labels[ManagedBy] = ComponentName

	// Set Exceptions
	policyManifest.Spec.Exceptions = translateEdgedbExceptions(edgedbPolman.PolicyExceptions)
	policyManifest.Spec.AutomatedExceptions = translateEdgedbExceptions(edgedbPolman.AutomatedExceptions)

	// Set Mode
	policyManifest.Spec.Mode = edgedbPolman.Mode

	return policyManifest, err
}
