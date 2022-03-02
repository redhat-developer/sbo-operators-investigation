# AWS SDK S3 App
## Info
Example S3 app built off of [AWS SDK for Go](https://github.com/aws/aws-sdk-go). This app also utilizes [minikube](https://minikube.sigs.k8s.io/docs/start/) and [Docker](https://docs.docker.com/get-docker/).

## About
This application provides test functionality of Amazon Web Services using [AWS SDK for Go](https://github.com/aws/aws-sdk-go) inside a Kubernetes container. The purpose of the app is to test AWS binding capabilities inside of Kubernetes both with and without the [Service Binding Operator](https://github.com/redhat-developer/service-binding-operator). This particular app binds to AWS S3 and offers the ability to upload, delete, retrieve, and describe bucket objects directly from a Kubernetes container.
## Setup
### IAM Credentials
[IAM credentials](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html) are needed to access AWS services along with the region .  Both the **Access Key ID** and **Secret Access Key** are needed. If using credentials from awscli, both keys can be found at `~/.aws/credentials`.  Both these keys need to be converted to base64 and inserted in the `secret.yaml` under `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` respectively.  


### Bucket
An [AWS Bucket](https://docs.aws.amazon.com/AmazonS3/latest/userguide/GetStartedWithS3.html) will need to be setup beforehand. Please refer to the AWS S3 Documentation on how to [Create An S3 Bucket](https://docs.aws.amazon.com/AmazonS3/latest/userguide/creating-bucket.html). If using a pre-existing bucket, ensure that your AWS account has [proper access and roles to the bucket](https://docs.aws.amazon.com/AmazonS3/latest/userguide/s3-access-control.html).

From the bucket you will need the **Bucket name** along with an **Object Key**. The [Bucket's Region](https://docs.aws.amazon.com/AmazonS3/latest/API/API_control_Region.html) is also needed. In `secret.yaml`, add the S3 Bucket's region to the `AWS_REGION` data.
## Build
1. Start minikube;

   `minikube start`

2. Navigate to:

   `~/sbo-services-investigations/aws/s3/app`

3. Run the command:
   
   `eval $(minikube docker-env)`

4. Build the Docker image using the command:

   `docker build -t aws-sdk-test .`

### App Without SBO

 Run the commands:
   
   `kubectl apply -f secret.yaml `
   
   `kubectl apply -f app.yaml`

### App With SBO
   
 Run the commands:
   
   `kubectl apply -f secret.yaml `
   
   `kubectl apply -f app-sb.yaml`
   
   `kubectl apply -f service-binding.yaml`
## Run App

Once the app has been built, run the command with the namespace the app is installed:

`kubectl get pods -n <namespace>`

Get the name of the pod and run the command, replacing `aws-sdk-pod` with the name of the pod.

`kubectl exec --stdin -t <aws-sdk-pod> -n <namespace>  -- sh`

The app can be run to upload, delete, get, and describe objects within the S3 bucket.

### Upload

Run the command with the name of the bucket being used and the key you want for the test file:

`go run aws-sdk-test.go -u -b <bucket> -k <key> -d 10m`

If the application runs successfully, a new test file will be uploaded to the defined bucket and key.

### Delete

Run the command with the name of the bucket being used and the key you want to delete:

`go run aws-sdk-test.go -r -b <bucket> -k <key> -d 10m`

If the application runs successfully, the selected key will be deleted from the bucket.

### Get

Run the command with the name of the bucket being used and the key of the item to be downloaded:

`go run aws-sdk-test.go -g -b <bucket> -k <key> -d 10m`

If the application runs successfully, a new file will be created in the working directory with the contents of the file in the S3 bucket.

### Info

Run the command with the name of the bucket being used and the key you want to retrieve info on:

`go run aws-sdk-test.go -i -b <bucket> -k <key> -d 10m`

If the application runs successfully, the info of the S3 object will be displayed.
