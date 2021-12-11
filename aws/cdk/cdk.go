package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2alpha/v2"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2integrationsalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/joho/godotenv"
)

type CdkStackProps struct {
	awscdk.StackProps
}

func NewCdkStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	os.Chdir("../lambda")
	mainCmd := exec.Command("go", "build", "-o", "main", "main.go")
	mainCmd.Env = append(os.Environ(), "GOOS=linux", "CGO_ENABLED=0", "GOARCH=amd64")
	_, err := mainCmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	env := envForLambdaExecution()
	functionName := jsii.String("admob-reporter-function")
	function := awslambda.NewFunction(stack, functionName, &awslambda.FunctionProps{
		FunctionName: functionName,
		Runtime:      awslambda.Runtime_GO_1_X(),
		Code:         awslambda.Code_FromAsset(jsii.String("../lambda"), &awss3assets.AssetOptions{}),
		Architecture: awslambda.Architecture_X86_64(),
		Handler:      jsii.String("main"),
		Environment:  env,
	})

	target := awseventstargets.NewLambdaFunction(function, &awseventstargets.LambdaFunctionProps{})
	targets := []awsevents.IRuleTarget{target}
	awsevents.NewRule(stack, jsii.String("admob-reporter-rule"), &awsevents.RuleProps{
		Schedule: awsevents.Schedule_Expression((*env)["CRON"]),
		Targets:  &targets,
	})

	integration := awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("admob-reporter-integration"), function, &awscdkapigatewayv2integrationsalpha.HttpLambdaIntegrationProps{})
	httpApi := awscdkapigatewayv2alpha.NewHttpApi(stack, jsii.String("admob-reporter-api"), &awscdkapigatewayv2alpha.HttpApiProps{
		ApiName: jsii.String("admob-reporter-api"),
	})
	httpApi.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Integration: integration,
		Path:        jsii.String("/post"),
		Methods: &[]awscdkapigatewayv2alpha.HttpMethod{
			awscdkapigatewayv2alpha.HttpMethod_POST,
		},
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewCdkStack(app, "AdmobReporterCdkStack", &CdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}

func envForLambdaExecution() *map[string]*string {
	// Required environment variables for lambda execution
	keys := []string{
		"ADMOB_PUBLISHER_ID",
		"ADMOB_OAUTH2_REFRESH_TOKEN",
		"SLACK_WEBHOOK_URL",
		"CRON",
	}

	// Check whether all of the required environment variables are set
	someMissingKey := false
	for _, key := range keys {
		if os.Getenv(key) == "" {
			someMissingKey = true
			break
		}
	}

	// Try to load from .env file if some keys are missing
	if someMissingKey {
		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}
	}

	// Build env map
	env := make(map[string]*string)
	for _, key := range keys {
		env[key] = jsii.String(os.Getenv(key))
	}

	timezone := os.Getenv("TZ")
	if timezone != "" {
		env["TZ"] = jsii.String(timezone)
	}

	return &env
}
