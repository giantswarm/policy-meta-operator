package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/edgedb/edgedb-go"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	utils "github.com/giantswarm/policy-meta-operator/internal/utils"
	polman "github.com/giantswarm/policy-meta-operator/internal/utils/policymanifest"
)

type PolicyManifestReconciler struct {
	client.Client
	EdgeDBClient *edgedb.Client
	Scheme       *runtime.Scheme
}

func (r *PolicyManifestReconciler) Reconcile(ctx context.Context) error {
	_ = log.FromContext(ctx)

	for {
		// List all Policies from edgedb
		policies := []string{"disallow-capabilities-strict", "require-run-as-nonroot"}

		for _, policy := range policies {
			policyManifest, err := polman.CreatePolicyManifest(ctx, r.EdgeDBClient, policy)
			if err != nil {
				log.Log.Error(err, "Error creating PolicyManifest from edgedb data")
			}

			c := utils.ClientHelper{Client: r.Client}
			if _, err := c.CreateOrUpdate(ctx, &policyManifest); err != nil {
				// Error creating or updating PolicyManifest
				log.Log.Error(err, "unable to create or update PolicyManifest")
				return err
			} else {
				log.Log.Info(fmt.Sprintf("PolicyManifest reconciled: %s", policyManifest.Name))
			}
		}
		time.Sleep(30 * time.Second)
	}
}
