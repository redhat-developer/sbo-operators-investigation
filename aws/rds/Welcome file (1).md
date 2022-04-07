# **AWS ACK RDS App**

## Info

Example ACK RDS app built off of [AWS ACK RDS](https://github.com/aws-controllers-k8s/rds-controller). This app also utilizes [minikube](https://minikube.sigs.k8s.io/docs/start/) and [Docker](https://docs.docker.com/get-docker/).

## About

This application tests the API functionality of [AWS RDS](https://aws.amazon.com/rds/?trk=c0fcea17-fb6a-4c27-ad98-192318a276ff&sc_channel=ps&sc_campaign=acquisition&sc_medium=ACQ-P|PS-GO|Brand|Desktop|SU|Database|Solution|US|EN|Text&s_kwcid=AL!4422!3!548665196304!e!!g!!amazon%20relational%20db&ef_id=EAIaIQobChMIpJKl5a-C9wIVbfHjBx0QYAOwEAAYASABEgJe7PD_BwE:G:s&s_kwcid=AL!4422!3!548665196304!e!!g!!amazon%20relational%20db) using the RDS controller from [Amazon Controllers for Kubernetes](https://github.com/aws-controllers-k8s/community).  The app walks through the process of setting up a PostgreSQL database instance using ACK along with a sample application to test API calls to the created database.

## Setup
### AWS ACK RDS
[AWS ACK RDS](https://github.com/redhat-developer/sbo-services-investigation/tree/main/aws/s3/app) must be installed on the cluster. Refer to the AWS ACK [installation guide](https://aws-controllers-k8s.github.io/community/docs/user-docs/install/).

### IAM Credentials
A `secret.yaml` name `aws-creds` must live on the cluster AWS ACK RDS is installed on. The secret must contain the user's [IAM credentials](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html) with the keywords `AWS_REGION`, `AWS_ACCESS_KEY_ID`, and `AWS_SECRET_ACCESS_KEY` for the region, access key and secret access key respectfully.

### DB Username & Password
In the `ack-postgres-chart/values.yaml` file, file in the `dbcreds` field with your desired `username` and `password`. Fill in both `name` fields under `dbinstance` & `dbcreds` for the database name and name of the secrets respectfully.

### Postgres Secret
In `app/rds-chart-app/values.yaml`, change the fields underneath `postgresql` with the appropriate DB username, password, and type (in this case, `postgres`). `host` and `port` will be filled in during the installation.

## Installation
1. Using helm in the `~/rds` directory of the project, run the command:

	`helm install ack-postgres-chart -n [namespace] --generate-name`

	where `[namespace]` is the namespace where ACK RDS is installed on.

2. Run the command:

	`kubectl get dbinstances -n [namespace]`

3. Wait for the dbinstance to be created. Grab the dbisntance that was created and run:

	`kubectl describe dbinstance [dbinstance-name] -n [namespace]`

4. In the dbinstance description, grab the `.Endpoint.Address` & `.Endpoint.Port` and place them in the `/rds/app/rds-chart-app/values.yaml` file under `postgresql.sed.host` & `postgresql.sed.port` respectfully.

5. In the `~/rds/app` folder, run the command:

	`helm install rds-chart-app -n [namespace] --generate-name`

6. Inside the `~/rds/app` folder, run the following commands:

	`kubectl apply -f service-binding.yaml -n [namespace]`

	`kubectl expose deployment aws-rds-sbo -n aws-rds --type=NodePort --name=rds-app-svc`

7. Use minikube to start the service:

	`minikube service -n [namespace] rds-app-svc`

