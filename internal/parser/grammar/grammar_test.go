package grammar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInfoGrammar(t *testing.T) {
	t.Run("Successfully parse application version,name,url from source string", func(t *testing.T) {
		slos, err := Eval(`@sloth name requests_availability
@sloth objective 99.9
@sloth description Common SLO based on availability for Kubernetes apiserver HTTP request responses
@sloth labels test value
@sloth labels test1 value @sloth.alerting name K8sApiserverAvailabilityAlert
@sloth.alerting.ticket disable true
@sloth.sli error_query sum(rate(apiserver_request_total{code=~"(5..|429)"}[{{.window}}]))
`)
		require.NoError(t, err)

		slo := slos["requests_availability"]
		assert.EqualValues(t, "requests_availability", slo.Name)
		assert.EqualValues(t, 99.9, slo.Objective)
		assert.EqualValues(t, map[string]string{"test": "value", "test1": "value"}, slo.Labels)
		assert.EqualValues(t, "Common SLO based on availability for Kubernetes apiserver HTTP request responses", slo.Description)
		assert.EqualValues(t, "K8sApiserverAvailabilityAlert", slo.Alerting.Name)
		assert.True(t, slo.Alerting.TicketAlert.Disable)
		assert.EqualValues(t, "sum(rate(apiserver_request_total{code=~\"(5..|429)\"}[{{.window}}]))", slo.SLI.Events.ErrorQuery)
	})
}

