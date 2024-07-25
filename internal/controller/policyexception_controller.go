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
	"time"

	"github.com/edgedb/edgedb-go"

	policyAPI "github.com/giantswarm/policy-api/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/giantswarm/policy-meta-operator/internal/utils"
)

// PolicyExceptionReconciler reconciles a AutomatedException object
type PolicyExceptionReconciler struct {
	client.Client
	EdgeDBClient *edgedb.Client
	Scheme       *runtime.Scheme
}

func (r *PolicyExceptionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	log.Log.Info("Reconciling PolicyException")

	// Create exception with counter
	_, err := utils.InsertPolicyException(ctx, r.EdgeDBClient, req.Name, int64(1), time.Now())
	if err != nil {
		log.Log.Error(err, "Error inserting policy exception in database")
	}

	// Select PolicyExceptions
	policyExceptions := utils.GetPolicyExceptionsFromEdgeDB(ctx, r.EdgeDBClient)

	for _, policyException := range policyExceptions {
		log.Log.Info(fmt.Sprintf("PolicyException: %s, Counter: %d", policyException.Name, policyException.Counter))
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PolicyExceptionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&policyAPI.PolicyException{}).
		Complete(r)
}
