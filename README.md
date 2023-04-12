# Lambda function and Cloud Formation template for AWS Account Registration
This contains the code for the AWS Lambda function for AWS Account Registration flow for Hub. Additionally, it houses the cloud formation template used by users to create the role ARN.

## IMPORTANT
For now, every time the template is updated, the version of the template should be bumped. For example, if the file is `template/template-v5.json` and a change is being made to the file, be sure to bump the version to `template/template-v6.json`. 

There may be times when an update for dev/staging would not require this template version to be bumped, but if a version is being used in prod, the VERSION MUST BE BUMPED.

## Workflows
There are 3 workflow types:
1. `build.yaml` - Ensures the app builds on PR
2. `deploy_lambda_*.yaml` - Deploys the code for the Lambda function to the S3 bucket.
3. `deploy_template_*.yaml` - Pushes a new version of the Cloud Formation template to the S3 bucket.

The `dev` versions of the deploy workflows are run on each merge/push.
The `stg` versions of the deploy workflows run on a schedule each night.
The `prod` versions of the deploy workflows should be run manually when needed.
Each of them can also be run manually when needed.