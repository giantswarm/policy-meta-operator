package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/edgedb/edgedb-go"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/giantswarm/policy-meta-operator/internal/utils"
	edgedbutils "github.com/giantswarm/policy-meta-operator/internal/utils/edgedb"
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
		policies, err := edgedbutils.ListPoliciesNames(ctx, r.EdgeDBClient)
		if err != nil {
			log.Log.Error(err, "Error fetching Policies from edgedb")
		} else if len(policies) == 0 {
			log.Log.Info("No Policies found in edgedb")
		}
		for _, policy := range policies {
			policyManifest, err := polman.CreatePolicyManifest(ctx, r.EdgeDBClient, policy.Name)
			if err != nil {
				log.Log.Error(err, "Error creating PolicyManifest from edgedb data")
			}

			c := utils.ClientHelper{Client: r.Client}
			if op, err := c.CreateOrUpdate(ctx, &policyManifest); err != nil {
				// Error creating or updating PolicyManifest
				log.Log.Error(err, "unable to create or update PolicyManifest")
				return err
			} else {
				switch {
				case op == "created":
					log.Log.Info(fmt.Sprintf("Created PolicyManifest %s", policyManifest.Name))
				case op == "updated":
					log.Log.Info(fmt.Sprintf("Updated PolicyManifest %s", policyManifest.Name))
				}
			}
		}
		time.Sleep(30 * time.Second)
	}
}
