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
	"strconv"
	"time"

	"github.com/edgedb/edgedb-go"

	policyAPI "github.com/giantswarm/policy-api/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/giantswarm/policy-meta-operator/internal/utils"
)

// AutomatedExceptionReconciler reconciles an AutomatedException object
type AutomatedExceptionReconciler struct {
	client.Client
	Scheme       *runtime.Scheme
	EdgeDBClient *edgedb.Client
}

func (r *AutomatedExceptionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	log.Log.Info("Reconciling AutomatedException")

	// Create exception with counter
	_, err := utils.InsertAutomatedException(ctx, r.EdgeDBClient, req.Name, int64(1), time.Now())
	if err != nil {
		log.Log.Error(err, "Error inserting exception in database")
	}

	// Select users.
	var automatedExceptions []utils.AutomatedException
	query := "select AutomatedException{name, last_reconciliation, counter}"
	err = r.EdgeDBClient.Query(ctx, query, &automatedExceptions)
	if err != nil {
		log.Log.Error(err, "error making query")
	}

	for _, exception := range automatedExceptions {
		time, _ := exception.LastReconciliation.Get()
		counter, _ := exception.Counter.Get()
		log.Log.Info(fmt.Sprintf("Exception: %s\nLast reconciliation: %s\nCounter: %s", exception.Name, time, strconv.FormatInt(counter, 10)))
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AutomatedExceptionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&policyAPI.AutomatedException{}).
		Complete(r)
}
