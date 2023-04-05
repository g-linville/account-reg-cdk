# Lambda function and Cloud Formation template for AWS Account Registration
This contains the code for the AWS Lambda function for AWS Account Registration flow for Hub. Additionally, it houses the cloud formation template used by users to create the role ARN.

## IMPORTANT
For now, every time the template is updated, the version of the template should be bumped. For example, if the file is `template/template-v5.json` and a change is being made to the file, be sure to bump the version to `template/template-v6.json`. 

There may be times when an update for dev/staging would not require this template version to be bumped, but if a version is being used in prod, the VERSION MUST BE BUMPED.

## Workflows
There are 4 workflows:
1. `build.yaml` - Ensures the app builds on PR
2. `deploy_lambda_dev.yaml` - Deploys the code to AWS Lambda on each merge/push.
3. `deploy_template_dev.yaml` - Pushes a new version of the Cloud Formation template to the S3 Bucket on each merge/push.
4. `publish_lambda_dev.yaml` - Creates a new published version of the Lambda. This is run manually, and only needs to be run to lock in a version of the lambda for continued use. In order to use the new published version, the URL using this Lambda should be updated.
