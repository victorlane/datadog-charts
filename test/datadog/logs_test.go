package datadog

import (
	"testing"

	"github.com/DataDog/helm-charts/test/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
)

func Test_logsSpecificDDURLOverride(t *testing.T) {
	manifest, err := common.RenderChart(t, common.HelmCommand{
		ReleaseName: "datadog",
		ChartPath:   "../../charts/datadog",
		ShowOnly:    []string{"templates/daemonset.yaml"},
		Values:      []string{"../../charts/datadog/values.yaml"},
		Overrides: map[string]string{
			"datadog.apiKeyExistingSecret":                   "datadog-secret",
			"datadog.site":                                   "datadoghq.eu",
			"datadog.logs.enabled":                           "true",
			"datadog.logs.dd_url":                            "https://logs.example.com",
			"datadog.operator.enabled":                       "false",
			"datadog.kubeStateMetricsEnabled":                "false",
			"datadog.csi.enabled":                            "false",
			"datadog.autoscaling.workload.enabled":           "false",
			"clusterAgent.metricsProvider.useDatadogMetrics": "false",
		},
	})
	require.NoError(t, err)

	var daemonset appsv1.DaemonSet
	common.Unmarshal(t, manifest, &daemonset)

	agentContainer, found := getContainer(t, daemonset.Spec.Template.Spec.Containers, "agent")
	require.True(t, found)

	envs := getEnvVarMap(agentContainer.Env)
	assert.Equal(t, "datadoghq.eu", envs["DD_SITE"])
	assert.Equal(t, "https://logs.example.com", envs["DD_LOGS_CONFIG_LOGS_DD_URL"])
	assert.NotContains(t, envs, "DD_DD_URL")
}

func Test_logsDDSSLEnableOverride(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "enabled",
			value: "true",
		},
		{
			name:  "disabled",
			value: "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifest, err := common.RenderChart(t, common.HelmCommand{
				ReleaseName: "datadog",
				ChartPath:   "../../charts/datadog",
				ShowOnly:    []string{"templates/daemonset.yaml"},
				Values:      []string{"../../charts/datadog/values.yaml"},
				Overrides: map[string]string{
					"datadog.apiKeyExistingSecret":                   "datadog-secret",
					"datadog.logs.enabled":                           "true",
					"datadog.logs.dd_ssl_enable":                     tt.value,
					"datadog.operator.enabled":                       "false",
					"datadog.kubeStateMetricsEnabled":                "false",
					"datadog.csi.enabled":                            "false",
					"datadog.autoscaling.workload.enabled":           "false",
					"clusterAgent.metricsProvider.useDatadogMetrics": "false",
				},
			})
			require.NoError(t, err)

			var daemonset appsv1.DaemonSet
			common.Unmarshal(t, manifest, &daemonset)

			agentContainer, found := getContainer(t, daemonset.Spec.Template.Spec.Containers, "agent")
			require.True(t, found)

			envs := getEnvVarMap(agentContainer.Env)
			assert.Equal(t, tt.value, envs["DD_LOGS_CONFIG_DD_SSL_ENABLE"])
		})
	}
}
