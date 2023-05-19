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

//	t.Run("Successfully parse application semver version v1.0.0", func(t *testing.T) {
//		app, err := EvalInfo(`@aloe version v1.0.0
//@aloe name cli
//@aloe url https://tfadeyi.github.io`)
//		require.NoError(t, err)
//		assert.EqualValues(t, "v1.0.0", app.Version)
//		assert.EqualValues(t, "cli", app.Name)
//		assert.EqualValues(t, "https://tfadeyi.github.io", app.BaseUrl)
//	})
//	t.Run("Successfully parse application semver version v1.0.0-alpha1", func(t *testing.T) {
//		app, err := EvalInfo(`@aloe version v1.0.0-alpha1
//@aloe name cli
//@aloe url https://tfadeyi.github.io`)
//		require.NoError(t, err)
//		assert.EqualValues(t, "v1.0.0-alpha1", app.Version)
//		assert.EqualValues(t, "cli", app.Name)
//		assert.EqualValues(t, "https://tfadeyi.github.io", app.BaseUrl)
//	})
//	t.Run("Fails to parse application info if the version is missing", func(t *testing.T) {
//		_, err := EvalInfo(`@aloe name cli
//@aloe url https://tfadeyi.github.io`)
//		require.ErrorIs(t, err, ErrMissingRequiredField)
//	})
//	t.Run("Fails to parse application info if the name is missing", func(t *testing.T) {
//		_, err := EvalInfo(`@aloe version v1.0.0-alpha1
//@aloe url https://tfadeyi.github.io`)
//		require.ErrorIs(t, err, ErrMissingRequiredField)
//	})
//	t.Run("Fails to parse application info if the url is missing", func(t *testing.T) {
//		_, err := EvalInfo(`@aloe version v1.0.0-alpha1
//@aloe name cli`)
//		require.ErrorIs(t, err, ErrMissingRequiredField)
//	})
//	t.Run("Fails to parse invalid source string", func(t *testing.T) {
//		_, err := EvalInfo(`please stop writing bad code`)
//		require.ErrorIs(t, err, ErrParseSource)
//	})
//
//	//t.Run("Successfully fetch long Description from source string", func(t *testing.T) {
//	//	ast, err := parseApplicationInfo("",
//	//		"@title \"Super Microservice\"\n"+
//	//			"@description `brand new flashy silver bullet microservice,\n"+
//	//			"will solve all business needs for the next 30 years.` ")
//	//	require.NoError(t, err)
//	//	desc, ok := ast.getGeneralAttribute("@description")
//	//	assert.True(t, ok)
//	//	assert.EqualValues(t, "brand new flashy silver bullet microservice,\nwill solve all business needs for the next 30 years.", desc)
//	//})
//
//	//	t.Run("Successfully fetch all application info fields from comment", func(t *testing.T) {
//	//		ast, err := parseApplicationInfo("", `
//	//@asyncapi application title "microservice"
//	//@asyncapi application title "microservice"`)
//	//		//@[license]
//	//		//	@name "MIT"
//	//		//	@url "license url"
//	//		//@[server prod]
//	//		//	@protocol "kafka"
//	//		//@[server stage]
//	//		//	@protocol "amqp"
//	//		//	@[variable port]
//	//		//		@enum "[8080,8081]"
//	//		require.NoError(t, err)
//	//
//	//		//title, ok := ast.getGeneralAttribute(titleAttr)
//	//		//assert.True(t, ok)
//	//		//assert.EqualValues(t, "microservice", title)
//	//
//	//		//title := ast.Blocks[0].Info.getTitle()
//	//		//assert.True(t, ok)
//	//		log.Printf("%+v", ast.Info[0].Info)
//	//		assert.EqualValues(t, "microservice", ast)
//	//		//version := ast.Blocks[0].Info.getVersion()
//	//		////assert.True(t, ok)
//	//		//assert.EqualValues(t, "v1.0.0", version)
//	//
//	//		//version, ok := ast.getGeneralAttribute(versionAttr)
//	//		//assert.True(t, ok)
//	//		//assert.EqualValues(t, "v1.0.0", version)
//	//
//	//		//desc, ok := ast.getGeneralAttribute("@description")
//	//		//assert.True(t, ok)
//	//		//assert.EqualValues(t, "short description", desc)
//	//		//
//	//		//tos, ok := ast.getGeneralAttribute(tosAttr)
//	//		//assert.True(t, ok)
//	//		//assert.EqualValues(t, "my tos", tos)
//	//		//
//	//		//cName, ok := ast.getGeneralAttribute(conNameAttr)
//	//		//assert.True(t, ok)
//	//		//assert.EqualValues(t, "mr.robot", cName)
//	//		//
//	//		//cEmail, ok := ast.getGeneralAttribute(conEmailAttr)
//	//		//assert.True(t, ok)
//	//		//assert.EqualValues(t, "robot@gmail.com", cEmail)
//	//		//
//	//		//cURL, ok := ast.getGeneralAttribute(conURLAttr)
//	//		//assert.True(t, ok)
//	//		//assert.EqualValues(t, "www.foo.bar", cURL)
//	//		//
//	//		//lName, ok := ast.getLicenseAttribute("@name")
//	//		//assert.True(t, ok)
//	//		//assert.EqualValues(t, "MIT", lName)
//	//		//
//	//		//lURL, ok := ast.getLicenseAttribute(urlAttr)
//	//		//assert.True(t, ok)
//	//		//assert.EqualValues(t, "license url", lURL)
//	//		//
//	//		//servers := ast.getServers()
//	//		//assert.Len(t, servers, 2)
//	//		//assert.Contains(t, servers, "prod")
//	//		//assert.Contains(t, servers, "stage")
//	//		//
//	//		//protocol, ok := ast.getServerAttribute("prod", protocolAttr)
//	//		//assert.True(t, ok)
//	//		//assert.EqualValues(t, "kafka", protocol)
//	//		//
//	//		//protocol, ok = ast.getServerAttribute("stage", protocolAttr)
//	//		//assert.True(t, ok)
//	//		//assert.EqualValues(t, "amqp", protocol)
//	//	})
//
//}
