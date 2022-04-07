# **Cloud Managed Service: AWS Managed Services**
## Goal
Part of ongoing [cloud services investigation](https://docs.google.com/document/d/1bZlBrwl4xBlu458Dpr9ZjzDgwzTShnYXXjNatofQtJo/edit#). To provide an integration of our [Service Binding Operator](https://github.com/redhat-developer/service-binding-operator) to various managed services offered by AWS (S3, SNS, RDS, etc.).

## **Issue/ Cause for Investigation**

APPSVC-547: Bind to External Cloud Services.

Expanded on by APPSVC-1092: Implement SBO in ACK RDS Operator.

## **Investigation Repository**
**[sbo-services-investigation/aws](https://github.com/redhat-developer/sbo-services-investigation/tree/main/aws)**

## **Services Investigated**

[Amazon Web Services (AWS)](https://aws.amazon.com/)

-   [AWS RDS](https://aws.amazon.com/rds/)
    
-   [AWS S3](https://aws.amazon.com/s3/?nc2=type_a)
    
-   [AWS CloudTrail](https://aws.amazon.com/cloudtrail/?nc2=type_a)

### AWS RDS
**Credentials Needed**:  [IAM Credentials](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html) to make for managing RDS, [SED and DB Username & Password](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_ConnectToPostgreSQLInstance.html) to access the DB Instance itself.

### AWS S3
**Credentials Needed**:  [IAM Credentials](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html) to manage buckets and access / edit contents.

### AWS CloudTrail
**Credentials Needed**:  [IAM Credentials](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html) to create and view trails.

## **Binding Solutions Investigated**
### **[Amazon Controllers for Kubernetes](https://github.com/aws-controllers-k8s/community) ([ACK RDS](https://github.com/aws-controllers-k8s/rds-controller))**

-   **Binding Info Needed**: [IAM Credentials](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html)
    
-   **Advantages**: API Calls can be in the form of operator defined CRD’s, secrets are the same for all controllers, (RDS) controller automatically creates SED for DB created
    
-   **Disadvantages**: Each service needs its own controller, (RDS) user must manually grab SED info and create a secret for DB binding
    
-   **Application Repository**: [here](https://github.com/redhat-developer/sbo-services-investigation/tree/main/aws/rds)
    
- **Solution**: Have the controller create a secret with SED info, add .status.binding, create PR for the changes

- **Benefits of SBO**: Developers can manage RDS & databases with secrets created by cluster admins,  SED secrets can be automatically created and bound to an application upon DB creation

### **[AWS SDK](https://aws.amazon.com/tools/) ([SDK for Go](https://aws.amazon.com/sdk-for-go/))**

-   **Binding Info Needed**: [IAM Credentials](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html)
    
-   **Advantages**: SDK can access all AWS services from one application from the same secret, application can access multiple services simultaneously
    
-   **Disadvantages**: Not a fully kubernetes contained solution, application must be tailored for each service/API call the user needs
    
-   **Application Repository**: [here](https://github.com/redhat-developer/sbo-services-investigation/tree/main/aws/s3/app)
    

- **Solution**: Use SBO to bind IAM credentials as environment variables through [Direct Binding](https://redhat-developer.github.io/service-binding-operator/userguide/exposing-binding-data/direct-secret-reference.html)

- **Benefits of SBO**: One secret used to bind to all AWS services supported by SDK, removes manual volume mount / environment variable creation for deployments using AWS SDK

## Results

### AWS ACK
AWS ACK controllers provide a kubernetes-based approach to API calls in which SBO can be leveraged to reduce user input and error. Moving forward, a PR will be created for ACK RDS to automatically create and bind DB Instance SED's.

### AWS SDK
AWS SDK provides an all-in-one solution to connecting to AWS services from a single secret. Applications can be made modular to perform specialized tasks/routines on AWS services, and Docker allows for these applications to be deployed onto kubernetes clusters. SBO can leverage AWS SDK through Direct Binding, decreasing manual environment variable configuration
