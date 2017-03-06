package integration

import (
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Upgrade", func() {
	BeforeEach(func() {
		dir := setupTestWorkingDirWithVersion("v1.2.2")
		os.Chdir(dir)
	})
	Describe("Upgrading a cluster using offline mode", func() {
		Context("Using a minikube layout", func() {
			Context("Using Ubuntu 16.04", func() {
				ItOnAWS("should be upgraded [slow] [upgrade]", func(aws infrastructureProvisioner) {
					WithMiniInfrastructure(Ubuntu1604LTS, aws, func(node NodeDeets, sshKey string) {
						// Install previous version cluster
						err := installKismaticMini(node, sshKey)
						Expect(err).ToNot(HaveOccurred())

						// Extract current version of kismatic
						pwd, err := os.Getwd()
						Expect(err).ToNot(HaveOccurred())
						err = extractCurrentKismatic(pwd)
						Expect(err).ToNot(HaveOccurred())

						// Perform upgrade
						cmd := exec.Command("./kismatic", "upgrade", "offline", "-f", "kismatic-testing.yaml")
						cmd.Stderr = os.Stderr
						cmd.Stdout = os.Stdout
						err = cmd.Run()
						Expect(err).ToNot(HaveOccurred())
					})
				})
			})

			Context("Using CentOS 7", func() {
				ItOnAWS("should be upgraded [slow] [upgrade]", func(aws infrastructureProvisioner) {
					WithMiniInfrastructure(CentOS7, aws, func(node NodeDeets, sshKey string) {
						// Install previous version cluster
						err := installKismaticMini(node, sshKey)
						Expect(err).ToNot(HaveOccurred())

						// Extract new version of kismatic
						pwd, err := os.Getwd()
						Expect(err).ToNot(HaveOccurred())
						err = extractCurrentKismatic(pwd)
						Expect(err).ToNot(HaveOccurred())

						// Perform upgrade
						cmd := exec.Command("./kismatic", "upgrade", "offline", "-f", "kismatic-testing.yaml")
						cmd.Stderr = os.Stderr
						cmd.Stdout = os.Stdout
						err = cmd.Run()
						Expect(err).ToNot(HaveOccurred())
					})
				})
			})
		})

		// This spec will be used for testing non-destructive kismatic features on
		// an upgraded cluster.
		// This spec is open to modification when new assertions have to be made.
		Context("Using a skunkworks cluster", func() {
			ItOnAWS("should result in an upgraded cluster [slow] [upgrade]", func(aws infrastructureProvisioner) {
				WithInfrastructureAndDNS(NodeCount{Etcd: 3, Master: 2, Worker: 3, Ingress: 2, Storage: 2}, CentOS7, aws, func(nodes provisionedNodes, sshKey string) {
					opts := installOptions{allowPackageInstallation: true}
					err := installKismatic(nodes, opts, sshKey)
					FailIfError(err)

					pwd, err := os.Getwd()
					FailIfError(err)
					err = extractCurrentKismatic(pwd)
					FailIfError(err)

					cmd := exec.Command("./kismatic", "upgrade", "offline", "-f", "kismatic-testing.yaml")
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					err = cmd.Run()
					FailIfError(err)

					assertClusterVersionIsCurrent()

					sub := SubDescribe("Using an upgraded cluster")
					defer sub.Check()

					sub.It("should allow adding a new storage volume", func() error {
						planFile, err := os.Open("kismatic-testing.yaml")
						if err != nil {
							return err
						}
						return createVolume(planFile, "test-vol", 1, 1, "")
					})

					sub.It("should allow adding a worker node", func() error {
						return nil
					})

					sub.It("should have an accessible dashboard", func() error {
						return nil
					})

					sub.It("should be able to deploy a workload with ingress", func() error {
						return nil
					})

					sub.It("should not have kube-apiserver systemd service", func() error {
						return nil
					})
				})
			})
		})
	})
})