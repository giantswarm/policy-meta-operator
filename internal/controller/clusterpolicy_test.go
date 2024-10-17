package controller_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/giantswarm/policy-meta-operator/internal/controller"

	kyvernoV1 "github.com/kyverno/kyverno/api/kyverno/v1"
)

var _ = Describe("Kyverno ClusterPolicy Controller", func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	var clusterPolicy *kyvernoV1.ClusterPolicy

	BeforeEach(func() {
		// Load the Cluster Policy YAML file
		filePath := filepath.Join("..", "..", "tests", "manifests", "disallow-capabilities-strict.yaml")
		yamlFile, err := os.ReadFile(filePath)
		Expect(err).NotTo(HaveOccurred(), "Failed to read YAML file")

		// Parse the YAML into the clusterPolicy object
		clusterPolicy = &kyvernoV1.ClusterPolicy{}
		err = yaml.Unmarshal(yamlFile, clusterPolicy)
		Expect(err).NotTo(HaveOccurred(), "Failed to unmarshal YAML")

	})
	Context("When reconciling a Kyverno Cluster Policy", func() {
		It("should extract rule names", func() {
			testExtractRuleNamesFunc := controller.ExtractRuleNames
			expectedRules := []string{"require-drop-all", "adding-capabilities-strict", "autogen-require-drop-all", "autogen-cronjob-require-drop-all", "autogen-adding-capabilities-strict", "autogen-cronjob-adding-capabilities-strict"}
			ruleNames := testExtractRuleNamesFunc(*clusterPolicy)

			Expect(ruleNames).To(Equal(expectedRules))
		})

		It("should extract target kinds", func() {
			testExtractTargetKindsFunc := controller.ExtractTargetKinds
			expectedKinds := []string{"Pod", "DaemonSet", "Deployment", "Job", "ReplicaSet", "ReplicationController", "StatefulSet", "CronJob"}
			targetKinds := testExtractTargetKindsFunc(*clusterPolicy)

			Expect(targetKinds).To(Equal(expectedKinds))
		})

		It("should know that GS workloads don't need to be exempted", func() {
			testShouldExcludeGiantSwarmResources := controller.ShouldExcludeGiantSwarmResources

			shouldExcludeGSResources := testShouldExcludeGiantSwarmResources(*clusterPolicy)

			// We have the team label set, result should be False
			Expect(shouldExcludeGSResources).To(BeFalse())
		})
	})

})
