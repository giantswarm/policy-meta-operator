package edgedbutils

import (
	"context"

	_ "embed"

	"github.com/edgedb/edgedb-go"

	policyAPI "github.com/giantswarm/policy-api/api/v1alpha1"
)

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

//go:embed queries/insertKyvernoClusterPolicy.edgeql
var insertKyvernoClusterPolicyQuery string

func InsertKyvernoClusterPolicy(ctx context.Context, client *edgedb.Client, name string, ruleNames []string, targetKinds []string, category string) (KyvernoClusterPolicy, error) {
	var kyvernoClusterPolicy KyvernoClusterPolicy

	err := client.QuerySingle(
		ctx,
		insertKyvernoClusterPolicyQuery,
		&kyvernoClusterPolicy,
		name,
		ruleNames,
		targetKinds,
		category,
	)

	return kyvernoClusterPolicy, err
}

func DeleteKyvernoClusterPolicy(ctx context.Context, client *edgedb.Client, policyName string) error {
	var kyvernoClusterPolicy KyvernoClusterPolicy

	err := client.QuerySingle(
		ctx,
		"DELETE KyvernoClusterPolicy FILTER .name = <str>$0 LIMIT 1",
		&kyvernoClusterPolicy,
		policyName,
	)

	return err
}

//go:embed queries/insertPolicy.edgeql
var insertPolicyQuery string

func InsertPolicy(ctx context.Context, client *edgedb.Client, policy policyAPI.Policy) (Policy, error) {
	var edgedbPolicy Policy

	err := client.QuerySingle(
		ctx,
		insertPolicyQuery,
		&edgedbPolicy,
		policy.Name,
		policy.Spec.DefaultPolicyState,
	)

	return edgedbPolicy, err
}

//go:embed queries/insertPolicyConfig.edgeql
var insertPolicyConfigQuery string

func InsertPolicyConfig(ctx context.Context, client *edgedb.Client, policyConfig policyAPI.PolicyConfig) (PolicyConfig, error) {
	var edgedbPolicyConfig PolicyConfig

	err := client.QuerySingle(
		ctx,
		insertPolicyConfigQuery,
		&edgedbPolicyConfig,
		policyConfig.Name,
		policyConfig.Spec.PolicyState,
		policyConfig.Spec.PolicyName,
	)

	return edgedbPolicyConfig, err
}

//go:embed queries/insertAutomatedException.edgeql
var insertAutomatedExceptionQuery string

func InsertAutomatedException(ctx context.Context, client *edgedb.Client, automatedException policyAPI.AutomatedException) (Exception, error) {
	var edgedbException Exception
	var targetIDs []edgedb.UUID

	// Create Policies in edgedb if they don't exist
	err := createPoliciesIfNonExistent(ctx, client, automatedException.Spec.Policies)
	if err != nil {
		return edgedbException, err
	}

	// Create Targets in edgedb if they don't exist, otherwise select them so we can append them to the policy exception
	targetIDs, err = createTargetsIfNonExistent(ctx, client, automatedException.Spec.Targets)
	if err != nil {
		return edgedbException, err
	}

	params := []interface{}{
		automatedException.Spec.Policies,
		targetIDs,
		automatedException,
	}

	err = client.QuerySingle(
		ctx,
		insertAutomatedExceptionQuery,
		&edgedbException,
		params...,
	)

	return edgedbException, err
}

//go:embed queries/insertPolicyException.edgeql
var insertPolicyExceptionQuery string

func InsertPolicyException(ctx context.Context, client *edgedb.Client, policyException policyAPI.PolicyException) (Exception, error) {
	var edgedbException Exception
	var targetIDs []edgedb.UUID

	// Create Policies in edgedb if they don't exist
	err := createPoliciesIfNonExistent(ctx, client, policyException.Spec.Policies)
	if err != nil {
		return edgedbException, err
	}

	// Create Targets in edgedb if they don't exist, otherwise select them so we can append them to the policy exception
	targetIDs, err = createTargetsIfNonExistent(ctx, client, policyException.Spec.Targets)
	if err != nil {
		return edgedbException, err
	}

	params := []interface{}{
		policyException.Spec.Policies,
		targetIDs,
		policyException.Name,
	}

	err = client.QuerySingle(
		ctx,
		insertPolicyExceptionQuery,
		&edgedbException,
		params...,
	)

	return edgedbException, err
}

//go:embed queries/insertTarget.edgeql
var insertTargetQuery string

func createTargetsIfNonExistent(ctx context.Context, client *edgedb.Client, targets []policyAPI.Target) ([]edgedb.UUID, error) {
	newTargets := translateTargetsToEdgedbTypes(targets)
	var outputIDs []edgedb.UUID

	for _, target := range newTargets {
		var result Target

		params := []interface{}{
			target.Names,
			target.Namespaces,
			target.Kind,
		}

		err := client.QuerySingle(
			ctx,
			insertTargetQuery,
			&result,
			params...,
		)
		if err != nil {
			return outputIDs, err
		}

		outputIDs = append(outputIDs, result.ID)
	}

	return outputIDs, nil
}

func createPoliciesIfNonExistent(ctx context.Context, client *edgedb.Client, policies []string) error {

	for _, policy := range policies {
		var result []Policy
		err := client.Query(
			ctx,
			"SELECT Policy {name} FILTER .name = <str>$0",
			&result,
			policy,
		)
		if err != nil {
			return err
		}

		// Check if result is empty
		var newPolicy Policy
		if len(result) == 0 {
			err := client.QuerySingle(
				ctx,
				"INSERT Policy {name := <str>$0}",
				&newPolicy,
				policy,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func DeletePolicyException(ctx context.Context, client *edgedb.Client, policyExceptionName string) error {
	var edgedbException Exception

	err := client.QuerySingle(
		ctx,
		"DELETE PolicyException FILTER .name = <str>$0 LIMIT 1",
		&edgedbException,
		policyExceptionName,
	)

	return err
}

func DeletePolicy(ctx context.Context, client *edgedb.Client, policyName string) error {
	var edgedbPolicy Policy

	err := client.QuerySingle(
		ctx,
		"DELETE Policy FILTER .name = <str>$0 LIMIT 1",
		&edgedbPolicy,
		policyName,
	)

	return err
}

func DeletePolicyConfig(ctx context.Context, client *edgedb.Client, policyConfigName string) error {
	var edgedbPolicy PolicyConfig

	err := client.QuerySingle(
		ctx,
		"DELETE PolicyConfig FILTER .name = <str>$0 LIMIT 1",
		&edgedbPolicy,
		policyConfigName,
	)

	return err
}

func ListPoliciesNames(ctx context.Context, client *edgedb.Client) ([]Policy, error) {
	var policies []Policy

	err := client.Query(
		ctx,
		"SELECT Policy{name}",
		&policies,
	)

	return policies, err
}
