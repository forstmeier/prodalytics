# todalytics

> Productivity analytics monitor 🧮

## About

`todalytics` is the root application of the data collection system.  

Collect and store **Todoist** task activity in DynamoDB for analysis to improve productivity.  

## Setup

### Prerequisites

Several packages are required for launching and managing the `todalytics` stack.  

- [jq](https://stedolan.github.io/jq/) - version `jq-1.6`  
- [AWS CLI](https://aws.amazon.com/cli/) - version `aws-cli/1.19.53 Python/3.8.10 Linux/5.11.0-36-generic botocore/1.20.53`  

### Installation

Follow the steps below to configure the required CloudFormation resources in your AWS account.  

- Download the most recent `release.zip` file from the [releases](https://github.com/forstmeier/todalytics/releases) page  
- Extract the contents below into your desired folder
	- `cft.yaml`: the full CloudFormation template definition for the required AWS resources  
	- `events.zip`: an AWS Lambda binary pre-compiled and zipped  
	- `config.json`: configuration file with user-provided or generated information  
	- `start_app`: a Bash script file used to launch the CloudFormation stack  
- All Bash scripts are used to manage the `todalytics` service and reference the `config.json` which should live in the same directory  
- Create a [Todoist App](https://developer.todoist.com/appconsole.html)  
	- Copy the **Client secret** value and add it to the `todoist.client_secret` field in the `config.json` file  
	- Click "Create test token" but don't use the token for anything  
	- Under **Watched Events** select _item:added_, _item:updated_, _item:completed_, _item:uncompleted_, and _item:deleted_ for the webhook  
- Get an API token from your Todoist account  
	- Click "Settings" -> "Integrations" and under **API token** click "Copy to clipboard"  
	- Add the API token value to the `config.json` file under the `todoist.api_token` field  

## Usage

Follow the steps below to launch, configure, and interact with the `todalytics` application.  

1. Run the `start_app` script in the folder the `release.zip` file was extracted into  
	a. This script will optionally populate the `aws.s3.artifact_bucket` and `aws.dynamodb.table_name` values in `config.json` with pre-existing resource names if the `-b` or `-t` flags receive argument values  
	b. If no flag values are received, an AWS S3 bucket and AWS DynamoDB table will be created and the `config.json` file will be updated accordingly  
2. Copy the `EventsAPIEndpoint` valaue from the CloudFormation stack outputs and add it to the [Todoist App](https://developer.todoist.com/appconsole.html) management console **Webhook callback URL** field  
3. Begin using **Todoist** and the `todalytics` app will populate the table constantly  

## Roadmap

Additional events and potentially some data enrichment features may be added depending on the usefulness of collecting this data.  

## Contribute

There are a few tools required to begin working on the `todalytics` codebase. The indicated versions are what the application was built using - other versions or operating systems have not been tested. See the contributing and code of conduct resources for specifics.  

- [Go](https://golang.org/dl/) - version `go version go1.16 linux/amd64`  
- [Git](https://git-scm.com/downloads) - version `git version 2.25.1`  
- [jq](https://stedolan.github.io/jq/) - version `jq-1.6`  
- [AWS CLI](https://aws.amazon.com/cli/) - version `aws-cli/1.19.53 Python/3.8.10 Linux/5.11.0-36-generic botocore/1.20.53`  

Scripts stored in the `bin/` folder are typically used for working with the `todalytics` stack during development. A `config.json` file needs to be added at `etc/config/config.json` with user-provided pre-existing S3 buckets added to the respective `"REPLACE"` field values.  

```json
{
	"aws": {
		"cloudformation": {
			"stack_name": "todalytics"
		},
		"s3": {
			"artifact_bucket": "REPLACE"
		},
		"dynamodb": {
			"table_name": "REPLACE"
		}
	},
	"todoist": {
		"api_token": "REPLACE",
		"client_secret": "REPLACE"
	}
}
```
