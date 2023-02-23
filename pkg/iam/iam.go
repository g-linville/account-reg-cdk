package iam

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type Statement struct {
	Action   []string
	Resource []string
}

func AssumeRole(scope constructs.Construct, name, principalARN string, policy awsiam.ManagedPolicy, externalID any) awsiam.IRole {
	ps := awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions: jsii.Strings("sts:AssumeRole"),
		Conditions: &map[string]any{
			"StringEquals": map[string]any{
				"sts:ExternalId": "something",
			},
		},
		Effect: awsiam.Effect_ALLOW,
	})

	principal := awsiam.NewArnPrincipal(jsii.String(principalARN))
	principal.AddToAssumeRolePolicy(awsiam.NewPolicyDocument(&awsiam.PolicyDocumentProps{
		Statements: &[]awsiam.PolicyStatement{ps},
	}))

	role := awsiam.NewRole(scope, jsii.String(name), &awsiam.RoleProps{
		AssumedBy:       principal,
		ManagedPolicies: &[]awsiam.IManagedPolicy{policy},
	})

	return role
}

func ManagedPolicy(scope constructs.Construct, name string, stmts ...Statement) awsiam.ManagedPolicy {
	var policyStatements []awsiam.PolicyStatement
	for _, stmt := range stmts {
		policyStatements = append(policyStatements, awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Actions:   jsii.Strings(stmt.Action...),
			Effect:    awsiam.Effect_ALLOW,
			Resources: jsii.Strings(stmt.Resource...),
		}))
	}
	return awsiam.NewManagedPolicy(scope, &name, &awsiam.ManagedPolicyProps{
		Statements: &policyStatements,
	})
}
