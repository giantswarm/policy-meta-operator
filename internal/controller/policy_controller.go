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

	"github.com/edgedb/edgedb-go"
	"github.com/pingcap/errors"

	policyAPI "github.com/giantswarm/policy-api/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	edgedbutils "github.com/giantswarm/policy-meta-operator/internal/utils/edgedb"
)

// PolicyReconciler reconciles a AutomatedException object
type PolicyReconciler struct {
	client.Client
	EdgeDBClient *edgedb.Client
	Scheme       *runtime.Scheme
}

func (r *PolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var gsPolicy policyAPI.Policy

	if err := r.Get(ctx, req.NamespacedName, &gsPolicy); err != nil {
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

	_, err := edgedbutils.InsertPolicy(ctx, r.EdgeDBClient, gsPolicy)
	if err != nil {
		log.Log.Error(err, "Error inserting Policy in database")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&policyAPI.Policy{}).
		Complete(r)
}
