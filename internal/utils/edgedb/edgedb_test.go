package edgedbutils_test

import (
	"github.com/edgedb/edgedb-go"
	policyAPI "github.com/giantswarm/policy-api/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	edgedbutils "github.com/giantswarm/policy-meta-operator/internal/utils/edgedb"
)

var (
	workloadName      = "hello-world"
	workloadNamespace = "default"
	workloadKind      = "Deployment"
)

var _ = Describe("edgedb Util Package", func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	var TestTargets []policyAPI.Target

	BeforeEach(func() {
		TestTargets = []policyAPI.Target{
			{
				Kind:       workloadKind,
				Names:      []string{workloadName},
				Namespaces: []string{workloadNamespace},
			},
		}

	})
	Context("When reconciling a Policy API Target array", func() {
		It("it should extract the correct edgedb Targets", func() {
			testTranslateTargets := edgedbutils.TranslateTargetsToEdgedbTypes

			resultTargets := testTranslateTargets(TestTargets)

			expectedTargets := []edgedbutils.Target{
				{
					ID:         edgedb.UUID{},
					Names:      []string{workloadName},
					Namespaces: []string{workloadNamespace},
					Kind:       workloadKind,
				},
			}

			Expect(resultTargets).To(Equal(expectedTargets))
		})
	})

})
