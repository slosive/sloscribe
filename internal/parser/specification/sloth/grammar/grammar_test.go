package grammar

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGrammar(t *testing.T) {
	t.Parallel()

	t.Run("Successfully parse sloth definitions for slo description in multiline comment", func(t *testing.T) {
		spec, err := Eval(`@sloth.slo name common
@sloth.slo description Common SLO
based on availability
for Kubernetes apiserver
HTTP request responses`)
		require.NoError(t, err)
		require.NotEmpty(t, spec.SLOs)
		require.Len(t, spec.SLOs, 1)

		assert.EqualValues(t, `Common SLO
based on availability
for Kubernetes apiserver
HTTP request responses`, spec.SLOs[0].Description)
	})

	t.Run("Fail to parse sloth definitions for slo name with empty value", func(t *testing.T) {
		_, err := Eval(`@sloth.slo name`)
		require.Error(t, err)
	})

	t.Run("Successfully parse sloth definitions for slo name with too many whitespaces", func(t *testing.T) {
		spec, err := Eval(`@sloth.slo 		name		 requests_availability`)
		require.NoError(t, err)
		require.NotEmpty(t, spec.SLOs)
		require.Len(t, spec.SLOs, 1)

		assert.EqualValues(t, "requests_availability", spec.SLOs[0].Name)
	})

	t.Run("Successfully parse sloth definitions from source string", func(t *testing.T) {
		spec, err := Eval(`@sloth.slo name requests_availability
@sloth.slo objective 99.9
@sloth.slo description Common SLO based on availability for Kubernetes apiserver HTTP request responses
@sloth.slo labels test value
@sloth.slo labels test1 value @sloth.alerting name K8sApiserverAvailabilityAlert
@sloth.alerting.ticket disable true
@sloth.sli error_query sum(rate(apiserver_request_total{code=~"(5..|429)"}[{{.window}}]))
`)
		require.NoError(t, err)
		require.NotEmpty(t, spec.SLOs)
		require.Len(t, spec.SLOs, 1)

		assert.EqualValues(t, "requests_availability", spec.SLOs[0].Name)
		assert.EqualValues(t, 99.9, spec.SLOs[0].Objective)
		assert.EqualValues(t, map[string]string{"test": "value", "test1": "value"}, spec.SLOs[0].Labels)
		assert.EqualValues(t, "Common SLO based on availability for Kubernetes apiserver HTTP request responses", spec.SLOs[0].Description)
		assert.EqualValues(t, "K8sApiserverAvailabilityAlert", spec.SLOs[0].Alerting.Name)
		assert.True(t, spec.SLOs[0].Alerting.TicketAlert.Disable)
		assert.EqualValues(t, "sum(rate(apiserver_request_total{code=~\"(5..|429)\"}[{{.window}}]))", spec.SLOs[0].SLI.Events.ErrorQuery)
	})
}
