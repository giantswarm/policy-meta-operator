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

	edgedbutils "github.com/giantswarm/policy-meta-operator/internal/utils/edgedb"
)

// PolicyReconciler reconciles a AutomatedException object
type ClusterPolicyReconciler struct {
	client.Client
	EdgeDBClient *edgedb.Client
	Scheme       *runtime.Scheme
}

func (r *ClusterPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var clusterPolicy kyvernoV1.ClusterPolicy

	if err := r.Get(ctx, req.NamespacedName, &clusterPolicy); err != nil {
		if !errors.IsNotFound(err) {
			// Error fetching the report
			log.Log.Error(err, "unable to fetch Policy")
		} else {
			// Policy not found, make sure we don't have it in edgedb either
			err := edgedbutils.DeletePolicy(ctx, r.EdgeDBClient, req.Name)
			if err != nil {
				log.Log.Error(err, "Error deleting Policy from database")
			}
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if clusterPolicy.Spec.ValidationFailureAction.Enforce() {
		if clusterPolicy.HasValidate() {
			log.Log.Info("Kyverno ClusterPolicy is in enforce mode and has validating rules, adding it to edgedb")

			policyRuleNames := extractRuleNames(clusterPolicy)
			policyTargetKinds := extractTargetKinds(clusterPolicy)

			if len(policyRuleNames) == 0 || len(policyTargetKinds) == 0 {
				log.Log.Error(errors.New("Error extracting rule names or target kinds"), fmt.Sprintf("Error extracting rule names or target kinds from Kyverno ClusterPolicy %s", clusterPolicy.Name))
				return ctrl.Result{}, nil
			}

			_, err := edgedbutils.InsertKyvernoClusterPolicy(ctx, r.EdgeDBClient, clusterPolicy.Name, policyRuleNames, policyTargetKinds)
			if err != nil {
				log.Log.Error(err, "Error inserting Kyverno ClusterPolicy in database")
			}

		}
	}

	return ctrl.Result{}, nil
}

func extractRuleNames(kyvernoPolicy kyvernoV1.ClusterPolicy) []string {
	var ruleNames []string

	for _, rule := range kyvernoPolicy.Spec.Rules {
		ruleNames = append(ruleNames, rule.Name)
	}

	return ruleNames
}

func extractTargetKinds(kyvernoPolicy kyvernoV1.ClusterPolicy) []string {
	var targetKinds []string

	for _, rule := range kyvernoPolicy.Spec.Rules {
		for _, match := range rule.MatchResources.Any {
			targetKinds = append(targetKinds, match.ResourceDescription.Kinds...)
		}
		for _, match := range rule.MatchResources.All {
			targetKinds = append(targetKinds, match.ResourceDescription.Kinds...)
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
