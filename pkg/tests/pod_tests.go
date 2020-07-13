package tests

import (
	"crypto/tls"
	"encoding/json"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/client"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/controller"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/monitors/metadata"
	"net/http"
)

var _ =  KubeDescribe( "[Pods]" , func() {
	ginkgo.It("Check `Codeready Operator` integrity", func() {
		codeready, err := GetPodByLabel(CodeReadyOperatorLabel)
		if err != nil {
			panic(err)
		}

		Expect(codeready).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred(), "failed to get information from pod %v\n", CodeReadyOperatorLabel)
	})

	ginkgo.It("Check `Postgres DB` integrity", func() {
		Logger.Info("Getting information and metrics from Postgres DB pod")
		postgres, err := GetPodByLabel(PostgresLabel)

		Expect(postgres).NotTo(BeNil())
		if err != nil {
			Logger.Panic("Error on getting information about postgres pod.")
		}

		if err := DescribePod(postgres); err != nil {
			Logger.Fatal("Failed to set metadata about postgres pod.")
		}

		Expect(err).NotTo(HaveOccurred(), "failed to get information from pod %v\n", PostgresLabel)
	})

	ginkgo.It("Check `Keycloak Server` integrity", func() {
		Logger.Info("Getting information and metrics from Keycloak Server pod")
		keycloak, err := GetPodByLabel(KeycloakLabel)

		Expect(keycloak).NotTo(BeNil())
		if err != nil {
			Logger.Panic("Error on getting information about keycloak pod.")
		}

		if err := DescribePod(keycloak); err != nil {
			Logger.Fatal("Failed to set metadata about keycloak pod.")
		}

		Expect(err).NotTo(HaveOccurred(), "failed to get information from pod %v\n", KeycloakLabel)
	})

	ginkgo.It("Check `Devfile Registry` integrity", func() {
		Logger.Info("Getting information and metrics from Devfile Registry pod")
		devFile, err := GetPodByLabel(DevFileLabel)

		Expect(devFile).NotTo(BeNil())

		if err != nil {
			Logger.Panic("Error on getting information about devFile pod.")
		}

		if err := DescribePod(devFile); err != nil {
			Logger.Fatal("Failed to set metadata about devFile pod.")
		}

		Expect(err).NotTo(HaveOccurred(), "failed to get information from pod %v\n", DevFileLabel)
	})

	ginkgo.It("Check `Plugin Registry` integrity", func() {
		Logger.Info("Getting information and metrics from Plugin Registry pod")
		pluginRegistry, err := GetPodByLabel(PluginRegistryLabel)

		Expect(pluginRegistry).NotTo(BeNil())
		if err != nil {
			Logger.Panic("Error on getting information about pluginRegistry pod.")
		}

		if err := DescribePod(pluginRegistry); err != nil {
			Logger.Fatal("Failed to set metadata about pluginRegistry pod.")
		}

		Expect(err).NotTo(HaveOccurred(), "failed to get information from pod %v\n", PluginRegistryLabel)
	})

	ginkgo.It("Check `Codeready server` integrity", func() {
		Logger.Info("Getting information and metrics from Codeready server pod")
		codeready, err := GetPodByLabel(CodReadyServerLabel)
		Expect(codeready).NotTo(BeNil())

		if err != nil {
			Logger.Panic("Error on getting information about codeready pod.")
		}

		if err := DescribePod(codeready); err != nil {
			Logger.Fatal("Failed to set metadata about codeready pod.")
		}

		Expect(err).NotTo(HaveOccurred(), "failed to get information from pod %v\n", CodReadyServerLabel)
	})

	ginkgo.It("Check if Codeready Server is already up", func() {
		var t CHE
		Logger.Info("Checking if Code Ready API server is up")
		k8sClient, err := client.NewK8sClient()

		if err != nil {
			panic(err)
		}

		ctrl:= controller.NewTestHarnessController(k8sClient)

		che, err := ctrl.GetCustomResource()
		Expect(che).NotTo(BeNil())
		Expect(che).NotTo(BeNil())
		client := &http.Client{Transport:  &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		}

		cheUrl := che.Status.CheURL
		Expect(cheUrl).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred())

		resp, err := client.Get(cheUrl +"/api/system/state")

		if err != nil {
			logrus.Error(err)
		}

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&t)
		if err != nil {
			metadata.Instance.CodeReadyServerIsUp = false
		} else {
			metadata.Instance.CodeReadyServerIsUp = true
		}

		Expect(err).NotTo(HaveOccurred())
	})
})