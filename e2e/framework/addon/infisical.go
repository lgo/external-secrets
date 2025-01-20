/*
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

package addon

import (
	"fmt"

	vault "github.com/hashicorp/vault/api"

	// nolint
	"github.com/onsi/ginkgo/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/external-secrets/external-secrets-e2e/framework/util"
)

type Infisical struct {
	chart        *HelmChart
	Namespace    string
	PodName      string
	VaultClient  *vault.Client
	VaultURL     string
	VaultMtlsURL string

	RootToken          string
	VaultServerCA      []byte
	ServerCert         []byte
	ServerKey          []byte
	VaultClientCA      []byte
	ClientCert         []byte
	ClientKey          []byte
	JWTPubkey          []byte
	JWTPrivKey         []byte
	JWTToken           string
	JWTRole            string
	JWTPath            string
	JWTK8sPath         string
	KubernetesAuthPath string
	KubernetesAuthRole string

	AppRoleSecret string
	AppRoleID     string
	AppRolePath   string
}

func NewInfisical(namespace string) *Infisical {
	repo := "infisical-" + namespace
	return &Infisical{
		chart: &HelmChart{
			Namespace:    namespace,
			ReleaseName:  fmt.Sprintf("infisical-%s", namespace), // avoid cluster role collision
			Chart:        fmt.Sprintf("%s/infisical-standalone", repo),
			ChartVersion: "1.3.0",
			Repo: ChartRepo{
				Name: repo,
				URL:  "https://dl.cloudsmith.io/public/infisical/helm-charts/helm/charts/",
			},
			Values: []string{"/k8s/infisical.values.yaml"},
		},
		Namespace: namespace,
	}
}

func (l *Infisical) Install() error {
	ginkgo.By("Installing infisical in " + l.Namespace)
	err := l.chart.Install()
	if err != nil {
		return err
	}

	err = l.initInfisical()
	if err != nil {
		return err
	}

	return nil
}

func (l *Infisical) initInfisical() error {
	// Infisical becomes ready after its dependencies are ready and it has completed migrations.
	err := util.WaitForPodsReady(l.chart.config.KubeClientSet, 1, l.Namespace, metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=infisical-standalone",
	})
	if err != nil {
		return fmt.Errorf("error waiting for infisical to be ready: %w", err)
	}

	return nil
}

func (l *Infisical) Logs() error {
	return l.chart.Logs()
}

func (l *Infisical) Uninstall() error {
	return l.chart.Uninstall()
}

func (l *Infisical) Setup(cfg *Config) error {
	return l.chart.Setup(cfg)
}
