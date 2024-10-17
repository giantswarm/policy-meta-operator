/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"

	"github.com/edgedb/edgedb-go"
	"github.com/pingcap/errors"

	kyvernoV1 "github.com/kyverno/kyverno/api/kyverno/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/giantswarm/policy-meta-operator/internal/utils"
	edgedbutils "github.com/giantswarm/policy-meta-operator/internal/utils/edgedb"
)

// PolicyReconciler reconciles a AutomatedException object
type ClusterPolicyReconciler struct {
	client.Client
	EdgeDBClient     *edgedb.Client
	Scheme           *runtime.Scheme
	MaxJitterPercent int
}

var (
	//GiantSwarm team label
	GSTeamLabel = "application.giantswarm.io/team"
	//Policy API Exemption label
	//TODO: Move to Policy API
	PolicyAPIExemptionLabel = "policy.giantswarm.io/giantswarm-exempt"
)

func (r *ClusterPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var clusterPolicy kyvernoV1.ClusterPolicy

	if err := r.Get(ctx, req.NamespacedName, &clusterPolicy); err != nil {
		if !errors.IsNotFound(err) {
			// Error fetching the report
			log.Log.Error(err, "unable to fetch Kyverno ClusterPolicy")
		} else {
			// Policy not found, make sure we don't have it in edgedb either
			err := edgedbutils.DeleteKyvernoClusterPolicy(ctx, r.EdgeDBClient, req.Name)
			if err != nil {
				log.Log.Error(err, "Error deleting ClusterPolicy from database")
			}
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if clusterPolicy.Spec.ValidationFailureAction.Enforce() {
		if clusterPolicy.HasValidate() {
			log.Log.Info("Kyverno ClusterPolicy is in enforce mode and has validating rules, adding it to edgedb")

			policyRuleNames := extractRuleNames(clusterPolicy)
			policyTargetKinds := extractTargetKinds(clusterPolicy)
			exempted := isExempted(clusterPolicy)

			if len(policyRuleNames) == 0 || len(policyTargetKinds) == 0 {
				log.Log.Error(errors.New("Error extracting rule names or target kinds"), fmt.Sprintf("Error extracting rule names or target kinds from Kyverno ClusterPolicy %s", clusterPolicy.Name))
				return ctrl.Result{}, nil
			}

			_, err := edgedbutils.InsertKyvernoClusterPolicy(ctx, r.EdgeDBClient, clusterPolicy.Name, policyRuleNames, policyTargetKinds, exempted)
			if err != nil {
				log.Log.Error(err, "Error inserting Kyverno ClusterPolicy in database")
			}

		}
	}

	return utils.JitterRequeue(utils.DefaultRequeueDuration, r.MaxJitterPercent, log.Log), nil
}

func shouldExcludeGiantSwarmResources(clusterPolicy kyvernoV1.ClusterPolicy) bool {
	exempted := true

	// Check if team label exist
	if _, ok := clusterPolicy.Labels[GSTeamLabel]; ok {
		// If team label exists, this policy comes from Giant Swarm, so our workloads are not exempt; they should satisfy the policy, ship an exception, or be excluded within the policy itself.
		exempted = false
	}
	
	// Check if the policy has a label enabling or disabling GS exemption.
	if  gsExemptLabelValue, ok := clusterPolicy.Labels[PolicyAPIExemptionLabel]; ok {
	    gsExempt, err = strconv.ParseBool(gsExemptLabelValue)
	    if err != nil {
	        // The label value is garbage. Complain and error out, or default the behavior
	    }
	}
	
	return gsExempt
	return exempted
}

func extractRuleNames(kyvernoPolicy kyvernoV1.ClusterPolicy) []string {
	var ruleNames []string

	for _, rule := range kyvernoPolicy.Spec.Rules {
		ruleNames = append(ruleNames, rule.Name)
	}

	for _, autogen := range kyvernoPolicy.Status.Autogen.Rules {
		ruleNames = append(ruleNames, autogen.Name)
	}

	return ruleNames
}

func extractTargetKinds(kyvernoPolicy kyvernoV1.ClusterPolicy) []string {
	var targetMap = make(map[string]bool)
	var targetKinds []string

	for _, rule := range kyvernoPolicy.Spec.Rules {
		for _, match := range rule.MatchResources.Any {
			// Deduplicate before storing in targetKinds
			for _, kind := range match.ResourceDescription.Kinds {
				if _, ok := targetMap[kind]; !ok {
					targetMap[kind] = true
					targetKinds = append(targetKinds, kind)
				}
			}
		}
		for _, match := range rule.MatchResources.All {
			for _, kind := range match.ResourceDescription.Kinds {
				if _, ok := targetMap[kind]; !ok {
					targetMap[kind] = true
					targetKinds = append(targetKinds, kind)
				}
			}
		}
	}

	// Duplicate to get target kinds from autogen rules
	for _, rule := range kyvernoPolicy.Status.Autogen.Rules {
		for _, match := range rule.MatchResources.Any {
			for _, kind := range match.ResourceDescription.Kinds {
				if _, ok := targetMap[kind]; !ok {
					targetMap[kind] = true
					targetKinds = append(targetKinds, kind)
				}
			}
		}
		for _, match := range rule.MatchResources.All {
			for _, kind := range match.ResourceDescription.Kinds {
				if _, ok := targetMap[kind]; !ok {
					targetMap[kind] = true
					targetKinds = append(targetKinds, kind)
				}
			}
		}
	}

	return targetKinds
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kyvernoV1.ClusterPolicy{}).
		Complete(r)
}
