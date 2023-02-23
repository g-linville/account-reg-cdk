package main

import (
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
)

// example tests. To run these tests, uncomment this file along with the
// example resource in account-reg-cdk_test.go
func TestAccountRegCdkStack(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewAccountRegCdkStack(app, "MyStack", nil)

	// THEN
	template := assertions.Template_FromStack(stack)

	template.HasResourceProperties(jsii.String("AWS::SQS::Queue"), map[string]interface{}{
		"VisibilityTimeout": 300,
	})
}
