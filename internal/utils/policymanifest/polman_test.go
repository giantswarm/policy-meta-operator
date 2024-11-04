package polman_test

import (
	"github.com/edgedb/edgedb-go"
	policyAPI "github.com/giantswarm/policy-api/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	edgedbutils "github.com/giantswarm/policy-meta-operator/internal/utils/edgedb"
	polman "github.com/giantswarm/policy-meta-operator/internal/utils/policymanifest"
)

var (
	policyName = "disallow-capabilities"
	//	policyState       = "warming"
	workloadName      = "hello-world"
	workloadNamespace = "default"
	workloadKind      = "Deployment"
)

var _ = Describe("PolicyManifest Util Package", func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	var EdgedbExceptions []edgedbutils.Exception

	BeforeEach(func() {
		EdgedbExceptions = []edgedbutils.Exception{
			{
				Targets: []edgedbutils.Target{
					{
						Kind:       workloadKind,
						Names:      []string{workloadName},
						Namespaces: []string{workloadNamespace},
					},
				},
				Policies: []edgedbutils.Policy{
					{
						Name:               policyName,
						DefaultPolicyState: edgedb.OptionalStr{},
						ID:                 edgedb.UUID{},
					},
				},
				ID: edgedb.UUID{},
			},
		}

	})
	Context("When reconciling an edgedb Exception", func() {
		It("should extract the correct PolicyAPI targets", func() {
			testTranslateEdgedbExceptions := polman.TranslateEdgedbExceptions

			resultTargets := testTranslateEdgedbExceptions(EdgedbExceptions)

			expectedTargets := []policyAPI.Target{
				{
					Kind:       workloadKind,
					Names:      []string{workloadName},
					Namespaces: []string{workloadNamespace},
				},
			}

			Expect(resultTargets).To(Equal(expectedTargets))
		})
	})

})
