package eksconfig_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/aws/aws-k8s-tester/eksconfig"
	"github.com/aws/aws-sdk-go/service/eks"
)

func TestEnv(t *testing.T) {
	cfg := eksconfig.NewDefault()
	defer func() {
		os.RemoveAll(cfg.ConfigPath)
		os.RemoveAll(cfg.KubectlCommandsOutputPath)
		os.RemoveAll(cfg.RemoteAccessCommandsOutputPath)
	}()

	kubectlDownloadURL := "https://amazon-eks.s3-us-west-2.amazonaws.com/1.11.5/2018-12-06/bin/linux/amd64/kubectl"
	if runtime.GOOS == "darwin" {
		kubectlDownloadURL = strings.Replace(kubectlDownloadURL, "linux", runtime.GOOS, -1)
	}

	os.Setenv("AWS_K8S_TESTER_EKS_KUBECTL_COMMANDS_OUTPUT_PATH", "hello-kubectl")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_KUBECTL_COMMANDS_OUTPUT_PATH")
	os.Setenv("AWS_K8S_TESTER_EKS_REMOTE_ACCESS_COMMANDS_OUTPUT_PATH", "hello-ssh")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_REMOTE_ACCESS_COMMANDS_OUTPUT_PATH")
	os.Setenv("AWS_K8S_TESTER_EKS_REGION", "us-east-1")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_REGION")
	os.Setenv("AWS_K8S_TESTER_EKS_LOG_LEVEL", "debug")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_LOG_LEVEL")
	os.Setenv("AWS_K8S_TESTER_EKS_KUBECTL_DOWNLOAD_URL", kubectlDownloadURL)
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_KUBECTL_DOWNLOAD_URL")
	os.Setenv("AWS_K8S_TESTER_EKS_KUBECTL_PATH", "/tmp/aws-k8s-tester-test/kubectl")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_KUBECTL_PATH")
	os.Setenv("AWS_K8S_TESTER_EKS_KUBECONFIG_PATH", "/tmp/aws-k8s-tester/kubeconfig2")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_KUBECONFIG_PATH")
	os.Setenv("AWS_K8S_TESTER_EKS_ON_FAILURE_DELETE", "false")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ON_FAILURE_DELETE")
	os.Setenv("AWS_K8S_TESTER_EKS_ON_FAILURE_DELETE_WAIT_SECONDS", "780")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ON_FAILURE_DELETE_WAIT_SECONDS")
	os.Setenv("AWS_K8S_TESTER_EKS_COMMAND_AFTER_CREATE_CLUSTER", "echo hello1")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_COMMAND_AFTER_CREATE_CLUSTER")
	os.Setenv("AWS_K8S_TESTER_EKS_COMMAND_AFTER_CREATE_ADD_ONS", "echo hello2")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_COMMAND_AFTER_CREATE_ADD_ONS")

	os.Setenv("AWS_K8S_TESTER_EKS_PARAMETERS_VPC_CREATE", "false")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_PARAMETERS_VPC_CREATE")
	os.Setenv("AWS_K8S_TESTER_EKS_PARAMETERS_VPC_ID", "vpc-id")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_PARAMETERS_VPC_ID")
	os.Setenv("AWS_K8S_TESTER_EKS_PARAMETERS_VPC_CIDR", "my-cidr")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_PARAMETERS_VPC_CIDR")
	os.Setenv("AWS_K8S_TESTER_EKS_PARAMETERS_PUBLIC_SUBNET_CIDR_1", "public-cidr1")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_PARAMETERS_PUBLIC_SUBNET_CIDR_1")
	os.Setenv("AWS_K8S_TESTER_EKS_PARAMETERS_PUBLIC_SUBNET_CIDR_2", "public-cidr2")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_PARAMETERS_PUBLIC_SUBNET_CIDR_2")
	os.Setenv("AWS_K8S_TESTER_EKS_PARAMETERS_PUBLIC_SUBNET_CIDR_3", "public-cidr3")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_PARAMETERS_PUBLIC_SUBNET_CIDR_3")
	os.Setenv("AWS_K8S_TESTER_EKS_PARAMETERS_PRIVATE_SUBNET_CIDR_1", "private-cidr1")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_PARAMETERS_PRIVATE_SUBNET_CIDR_1")
	os.Setenv("AWS_K8S_TESTER_EKS_PARAMETERS_PRIVATE_SUBNET_CIDR_2", "private-cidr2")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_PARAMETERS_PRIVATE_SUBNET_CIDR_2")
	os.Setenv("AWS_K8S_TESTER_EKS_PARAMETERS_TAGS", "to-delete=2019;hello-world=test")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_PARAMETERS_TAGS")
	os.Setenv("AWS_K8S_TESTER_EKS_PARAMETERS_REQUEST_HEADER_KEY", "eks-options")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_PARAMETERS_REQUEST_HEADER_KEY")
	os.Setenv("AWS_K8S_TESTER_EKS_PARAMETERS_REQUEST_HEADER_VALUE", "kubernetesVersion=1.11")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_PARAMETERS_REQUEST_HEADER_VALUE")
	os.Setenv("AWS_K8S_TESTER_EKS_PARAMETERS_RESOLVER_URL", "amazon.com")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_PARAMETERS_RESOLVER_URL")
	os.Setenv("AWS_K8S_TESTER_EKS_PARAMETERS_SIGNING_NAME", "a")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_PARAMETERS_SIGNING_NAME")
	os.Setenv("AWS_K8S_TESTER_EKS_PARAMETERS_ROLE_CREATE", "false")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_PARAMETERS_ROLE_CREATE")
	os.Setenv("AWS_K8S_TESTER_EKS_PARAMETERS_ROLE_ARN", "cluster-role-arn")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_PARAMETERS_ROLE_ARN")
	os.Setenv("AWS_K8S_TESTER_EKS_PARAMETERS_ROLE_SERVICE_PRINCIPALS", "eks.amazonaws.com,eks-beta-pdx.aws.internal,eks-dev.aws.internal")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_PARAMETERS_ROLE_SERVICE_PRINCIPALS")
	os.Setenv("AWS_K8S_TESTER_EKS_PARAMETERS_ROLE_MANAGED_POLICY_ARNS", "arn:aws:iam::aws:policy/AmazonEKSServicePolicy,arn:aws:iam::aws:policy/AmazonEKSClusterPolicy,arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforSSM")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_PARAMETERS_ROLE_MANAGED_POLICY_ARNS")
	os.Setenv("AWS_K8S_TESTER_EKS_PARAMETERS_VERSION", "1.15")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_PARAMETERS_VERSION")
	os.Setenv("AWS_K8S_TESTER_EKS_PARAMETERS_ENCRYPTION_CMK_CREATE", "false")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_PARAMETERS_ENCRYPTION_CMK_CREATE")
	os.Setenv("AWS_K8S_TESTER_EKS_PARAMETERS_ENCRYPTION_CMK_ARN", "key-arn")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_PARAMETERS_ENCRYPTION_CMK_ARN")

	os.Setenv("AWS_K8S_TESTER_EKS_REMOTE_ACCESS_KEY_CREATE", "false")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_REMOTE_ACCESS_KEY_CREATE")
	os.Setenv("AWS_K8S_TESTER_EKS_REMOTE_ACCESS_KEY_NAME", "a")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_REMOTE_ACCESS_KEY_NAME")
	os.Setenv("AWS_K8S_TESTER_EKS_REMOTE_ACCESS_PRIVATE_KEY_PATH", "a")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_REMOTE_ACCESS_PRIVATE_KEY_PATH")

	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_CREATED", "true")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_CREATED")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_ENABLE", "true")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_ENABLE")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_FETCH_LOGS", "false")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_FETCH_LOGS")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_ROLE_CREATE", "false")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_ROLE_CREATE")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_ROLE_NAME", "mng-role-name")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_ROLE_NAME")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_ROLE_ARN", "mng-role-arn")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_ROLE_ARN")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_ROLE_SERVICE_PRINCIPALS", "ec2.amazonaws.com,eks.amazonaws.com,hello.amazonaws.com")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_ROLE_SERVICE_PRINCIPALS")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_ROLE_MANAGED_POLICY_ARNS", "a,b,c")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_ROLE_MANAGED_POLICY_ARNS")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_REQUEST_HEADER_KEY", "a")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_REQUEST_HEADER_KEY")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_REQUEST_HEADER_VALUE", "b")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_REQUEST_HEADER_VALUE")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_RESOLVER_URL", "a")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_RESOLVER_URL")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_SIGNING_NAME", "a")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_SIGNING_NAME")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_REMOTE_ACCESS_USER_NAME", "a")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_REMOTE_ACCESS_USER_NAME")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_MNGS", `{"mng-test-name-cpu":{"name":"mng-test-name-cpu","tags":{"cpu":"hello-world"},"release-version":"test-ver-cpu","ami-type":"AL2_x86_64","asg-min-size":17,"asg-max-size":99,"asg-desired-capacity":77,"instance-types":["type-cpu-1","type-cpu-2"],"volume-size":40},"mng-test-name-gpu":{"name":"mng-test-name-gpu","tags":{"gpu":"hello-world"},"release-version":"test-ver-gpu","ami-type":"AL2_x86_64_GPU","asg-min-size":30,"asg-max-size":35,"asg-desired-capacity":34,"instance-types":["type-gpu-1","type-gpu-2"],"volume-size":500}}`)
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_MNGS")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_LOGS_DIR", "a")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_MANAGED_NODE_GROUPS_LOGS_DIR")

	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_NLB_HELLO_WORLD_ENABLE", "true")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_NLB_HELLO_WORLD_ENABLE")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_NLB_HELLO_WORLD_CREATED", "true")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_NLB_HELLO_WORLD_CREATED")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_NLB_HELLO_WORLD_NLB_ARN", "invalid")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_NLB_HELLO_WORLD_NLB_ARN")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_NLB_HELLO_WORLD_NLB_NAME", "invalid")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_NLB_HELLO_WORLD_NLB_NAME")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_NLB_HELLO_WORLD_URL", "invalid")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_NLB_HELLO_WORLD_URL")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_NLB_HELLO_WORLD_DEPLOYMENT_REPLICAS", "333")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_NLB_HELLO_WORLD_DEPLOYMENT_REPLICAS")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_NLB_HELLO_WORLD_NAMESPACE", "test-namespace")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_NLB_HELLO_WORLD_NAMESPACE")

	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_ALB_2048_ENABLE", "true")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_ALB_2048_ENABLE")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_ALB_2048_CREATED", "true")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_ALB_2048_CREATED")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_ALB_2048_ALB_ARN", "invalid")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_ALB_2048_ALB_ARN")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_ALB_2048_ALB_NAME", "invalid")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_ALB_2048_ALB_NAME")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_ALB_2048_URL", "invalid")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_ALB_2048_URL")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_ALB_2048_DEPLOYMENT_REPLICAS_ALB", "333")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_ALB_2048_DEPLOYMENT_REPLICAS_ALB")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_ALB_2048_DEPLOYMENT_REPLICAS_2048", "555")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_ALB_2048_DEPLOYMENT_REPLICAS_2048")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_ALB_2048_NAMESPACE", "test-namespace")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_ALB_2048_NAMESPACE")

	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_JOB_PI_ENABLE", "true")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_JOB_PI_ENABLE")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_JOB_PI_CREATED", "true")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_JOB_PI_CREATED")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_JOB_PI_NAMESPACE", "hello1")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_JOB_PI_NAMESPACE")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_JOB_PI_COMPLETES", "100")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_JOB_PI_COMPLETES")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_JOB_PI_PARALLELS", "10")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_JOB_PI_PARALLELS")

	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_JOB_ECHO_ENABLE", "true")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_JOB_ECHO_ENABLE")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_JOB_ECHO_CREATED", "true")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_JOB_ECHO_CREATED")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_JOB_ECHO_NAMESPACE", "hello2")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_JOB_ECHO_NAMESPACE")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_JOB_ECHO_COMPLETES", "1000")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_JOB_ECHO_COMPLETES")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_JOB_ECHO_PARALLELS", "100")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_JOB_ECHO_PARALLELS")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_JOB_ECHO_ECHO_SIZE", "10000")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_JOB_ECHO_ECHO_SIZE")

	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_CRON_JOB_ENABLE", "true")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_CRON_JOB_ENABLE")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_CRON_JOB_CREATED", "true")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_CRON_JOB_CREATED")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_CRON_JOB_NAMESPACE", "hello3")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_CRON_JOB_NAMESPACE")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_CRON_JOB_SCHEDULE", "*/1 * * * *")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_CRON_JOB_SCHEDULE")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_CRON_JOB_COMPLETES", "100")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_CRON_JOB_COMPLETES")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_CRON_JOB_PARALLELS", "10")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_CRON_JOB_PARALLELS")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_CRON_JOB_SUCCESSFUL_JOBS_HISTORY_LIMIT", "100")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_CRON_JOB_SUCCESSFUL_JOBS_HISTORY_LIMIT")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_CRON_JOB_FAILED_JOBS_HISTORY_LIMIT", "1000")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_CRON_JOB_FAILED_JOBS_HISTORY_LIMIT")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_CRON_JOB_ECHO_SIZE", "10000")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_CRON_JOB_ECHO_SIZE")

	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_SECRETS_ENABLE", "true")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_SECRETS_ENABLE")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_SECRETS_CREATED", "true")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_SECRETS_CREATED")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_SECRETS_NAMESPACE", "hello")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_SECRETS_NAMESPACE")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_SECRETS_OBJECTS", "5")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_SECRETS_OBJECTS")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_SECRETS_SIZE", "10")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_SECRETS_SIZE")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_SECRETS_SECRET_QPS", "10")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_SECRETS_SECRET_QPS")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_SECRETS_SECRET_BURST", "10")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_SECRETS_SECRET_BURST")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_SECRETS_POD_QPS", "10")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_SECRETS_POD_QPS")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_SECRETS_POD_BURST", "10")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_SECRETS_POD_BURST")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_SECRETS_WRITES_RESULT_PATH", "writes.csv")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_SECRETS_WRITES_RESULT_PATH")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_SECRETS_READS_RESULT_PATH", "reads.csv")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_SECRETS_READS_RESULT_PATH")

	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_ENABLE", "true")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_ENABLE")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_CREATED", "true")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_CREATED")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_NAMESPACE", "hello")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_NAMESPACE")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_ROLE_NAME", "hello")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_ROLE_NAME")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_ROLE_MANAGED_POLICY_ARNS", "a,b,c")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_ROLE_MANAGED_POLICY_ARNS")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_SERVICE_ACCOUNT_NAME", "hello")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_SERVICE_ACCOUNT_NAME")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_CONFIG_MAP_NAME", "hello")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_CONFIG_MAP_NAME")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_CONFIG_MAP_SCRIPT_FILE_NAME", "hello.sh")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_CONFIG_MAP_SCRIPT_FILE_NAME")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_S3_BUCKET_NAME", "hello")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_S3_BUCKET_NAME")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_S3_KEY", "hello")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_S3_KEY")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_DEPLOYMENT_NAME", "hello-deployment")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_DEPLOYMENT_NAME")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_DEPLOYMENT_RESULT_PATH", "hello-deployment.log")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_IRSA_DEPLOYMENT_RESULT_PATH")

	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_FARGATE_ENABLE", "true")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_FARGATE_ENABLE")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_FARGATE_CREATED", "true")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_FARGATE_CREATED")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_FARGATE_NAMESPACE", "hello")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_FARGATE_NAMESPACE")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_FARGATE_ROLE_NAME", "hello")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_FARGATE_ROLE_NAME")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_FARGATE_ROLE_SERVICE_PRINCIPALS", "a,b,c")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_FARGATE_ROLE_SERVICE_PRINCIPALS")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_FARGATE_ROLE_MANAGED_POLICY_ARNS", "a,b,c")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_FARGATE_ROLE_MANAGED_POLICY_ARNS")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_FARGATE_PROFILE_NAME", "hello")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_FARGATE_PROFILE_NAME")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_FARGATE_SECRET_NAME", "HELLO-SECRET")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_FARGATE_SECRET_NAME")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_FARGATE_POD_NAME", "fargate-pod")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_FARGATE_POD_NAME")
	os.Setenv("AWS_K8S_TESTER_EKS_ADD_ON_FARGATE_CONTAINER_NAME", "fargate-container")
	defer os.Unsetenv("AWS_K8S_TESTER_EKS_ADD_ON_FARGATE_CONTAINER_NAME")

	if err := cfg.UpdateFromEnvs(); err != nil {
		t.Fatal(err)
	}

	if cfg.KubectlCommandsOutputPath != "hello-kubectl" {
		t.Fatalf("unexpected %q", cfg.KubectlCommandsOutputPath)
	}
	if cfg.RemoteAccessCommandsOutputPath != "hello-ssh" {
		t.Fatalf("unexpected %q", cfg.RemoteAccessCommandsOutputPath)
	}
	if cfg.Region != "us-east-1" {
		t.Fatalf("unexpected %q", cfg.Region)
	}
	if cfg.LogLevel != "debug" {
		t.Fatalf("unexpected %q", cfg.LogLevel)
	}
	if cfg.KubectlDownloadURL != kubectlDownloadURL {
		t.Fatalf("unexpected KubectlDownloadURL %q", cfg.KubectlDownloadURL)
	}
	if cfg.KubectlPath != "/tmp/aws-k8s-tester-test/kubectl" {
		t.Fatalf("unexpected KubectlPath %q", cfg.KubectlPath)
	}
	if cfg.KubeConfigPath != "/tmp/aws-k8s-tester/kubeconfig2" {
		t.Fatalf("unexpected KubeConfigPath %q", cfg.KubeConfigPath)
	}
	if cfg.OnFailureDelete {
		t.Fatalf("unexpected OnFailureDelete %v", cfg.OnFailureDelete)
	}
	if cfg.OnFailureDeleteWaitSeconds != 780 {
		t.Fatalf("unexpected OnFailureDeleteWaitSeconds %d", cfg.OnFailureDeleteWaitSeconds)
	}
	if cfg.CommandAfterCreateCluster != "echo hello1" {
		t.Fatalf("unexpected CommandAfterCreateCluster %q", cfg.CommandAfterCreateCluster)
	}
	if cfg.CommandAfterCreateAddOns != "echo hello2" {
		t.Fatalf("unexpected CommandAfterCreateAddOns %q", cfg.CommandAfterCreateAddOns)
	}

	if cfg.Parameters.VPCCreate {
		t.Fatalf("unexpected Parameters.VPCCreate %v", cfg.Parameters.VPCCreate)
	}
	if cfg.Parameters.VPCID != "vpc-id" {
		t.Fatalf("unexpected Parameters.VPCID %q", cfg.Parameters.VPCID)
	}
	if cfg.Parameters.VPCCIDR != "my-cidr" {
		t.Fatalf("unexpected Parameters.VPCCIDR %q", cfg.Parameters.VPCCIDR)
	}
	if cfg.Parameters.PublicSubnetCIDR1 != "public-cidr1" {
		t.Fatalf("unexpected Parameters.PublicSubnetCIDR1 %q", cfg.Parameters.PublicSubnetCIDR1)
	}
	if cfg.Parameters.PublicSubnetCIDR2 != "public-cidr2" {
		t.Fatalf("unexpected Parameters.PublicSubnetCIDR2 %q", cfg.Parameters.PublicSubnetCIDR2)
	}
	if cfg.Parameters.PublicSubnetCIDR3 != "public-cidr3" {
		t.Fatalf("unexpected Parameters.PublicSubnetCIDR3 %q", cfg.Parameters.PublicSubnetCIDR3)
	}
	if cfg.Parameters.PrivateSubnetCIDR1 != "private-cidr1" {
		t.Fatalf("unexpected Parameters.PrivateSubnetCIDR1 %q", cfg.Parameters.PrivateSubnetCIDR1)
	}
	if cfg.Parameters.PrivateSubnetCIDR2 != "private-cidr2" {
		t.Fatalf("unexpected Parameters.PrivateSubnetCIDR2 %q", cfg.Parameters.PrivateSubnetCIDR2)
	}
	expectedTags := map[string]string{"to-delete": "2019", "hello-world": "test"}
	if !reflect.DeepEqual(cfg.Parameters.Tags, expectedTags) {
		t.Fatalf("Tags expected %v, got %v", expectedTags, cfg.Parameters.Tags)
	}
	if cfg.Parameters.RequestHeaderKey != "eks-options" {
		t.Fatalf("unexpected Parameters.RequestHeaderKey %q", cfg.Parameters.RequestHeaderKey)
	}
	if cfg.Parameters.RequestHeaderValue != "kubernetesVersion=1.11" {
		t.Fatalf("unexpected Parameters.RequestHeaderValue %q", cfg.Parameters.RequestHeaderValue)
	}
	if cfg.Parameters.ResolverURL != "amazon.com" {
		t.Fatalf("unexpected Parameters.ResolverURL %q", cfg.Parameters.ResolverURL)
	}
	if cfg.Parameters.SigningName != "a" {
		t.Fatalf("unexpected Parameters.SigningName %q", cfg.Parameters.SigningName)
	}
	if cfg.Parameters.RoleCreate {
		t.Fatalf("unexpected Parameters.RoleCreate %v", cfg.Parameters.RoleCreate)
	}
	if cfg.Parameters.RoleARN != "cluster-role-arn" {
		t.Fatalf("unexpected Parameters.RoleARN %q", cfg.Parameters.RoleARN)
	}
	expectedRoleServicePrincipals := []string{
		"eks.amazonaws.com",
		"eks-beta-pdx.aws.internal",
		"eks-dev.aws.internal",
	}
	if !reflect.DeepEqual(expectedRoleServicePrincipals, cfg.Parameters.RoleServicePrincipals) {
		t.Fatalf("unexpected Parameters.RoleServicePrincipals %+v", cfg.Parameters.RoleServicePrincipals)
	}
	expectedRoleManagedPolicyARNs := []string{
		"arn:aws:iam::aws:policy/AmazonEKSServicePolicy",
		"arn:aws:iam::aws:policy/AmazonEKSClusterPolicy",
		"arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforSSM",
	}
	if !reflect.DeepEqual(expectedRoleManagedPolicyARNs, cfg.Parameters.RoleManagedPolicyARNs) {
		t.Fatalf("unexpected Parameters.RoleManagedPolicyARNs %+v", cfg.Parameters.RoleManagedPolicyARNs)
	}
	if cfg.Parameters.Version != "1.15" {
		t.Fatalf("unexpected Parameters.Version %q", cfg.Parameters.Version)
	}
	if cfg.Parameters.EncryptionCMKCreate {
		t.Fatalf("unexpected Parameters.EncryptionCMKCreate %v", cfg.Parameters.EncryptionCMKCreate)
	}
	if cfg.Parameters.EncryptionCMKARN != "key-arn" {
		t.Fatalf("unexpected Parameters.EncryptionCMKARN %q", cfg.Parameters.EncryptionCMKARN)
	}

	if cfg.RemoteAccessKeyCreate {
		t.Fatalf("unexpected cfg.RemoteAccessKeyCreate %v", cfg.RemoteAccessKeyCreate)
	}
	if cfg.RemoteAccessKeyName != "a" {
		t.Fatalf("unexpected cfg.RemoteAccessKeyName %q", cfg.RemoteAccessKeyName)
	}
	if cfg.RemoteAccessPrivateKeyPath != "a" {
		t.Fatalf("unexpected cfg.RemoteAccessPrivateKeyPath %q", cfg.RemoteAccessPrivateKeyPath)
	}

	if cfg.AddOnManagedNodeGroups.Created { // read-only must be ignored
		t.Fatalf("unexpected cfg.AddOnManagedNodeGroups.Created %v", cfg.AddOnManagedNodeGroups.Created)
	}
	if !cfg.AddOnManagedNodeGroups.Enable {
		t.Fatalf("unexpected cfg.AddOnManagedNodeGroups.Enable %v", cfg.AddOnManagedNodeGroups.Enable)
	}
	if cfg.AddOnManagedNodeGroups.FetchLogs {
		t.Fatalf("unexpected cfg.AddOnManagedNodeGroups.FetchLogs %v", cfg.AddOnManagedNodeGroups.FetchLogs)
	}
	if cfg.AddOnManagedNodeGroups.RoleCreate {
		t.Fatalf("unexpected AddOnManagedNodeGroups.RoleCreate %v", cfg.AddOnManagedNodeGroups.RoleCreate)
	}
	if cfg.AddOnManagedNodeGroups.RoleName != "mng-role-name" {
		t.Fatalf("unexpected cfg.AddOnManagedNodeGroups.RoleName %q", cfg.AddOnManagedNodeGroups.RoleName)
	}
	if cfg.AddOnManagedNodeGroups.RoleARN != "mng-role-arn" {
		t.Fatalf("unexpected cfg.AddOnManagedNodeGroups.RoleARN %q", cfg.AddOnManagedNodeGroups.RoleARN)
	}
	expectedManagedNodeGroupRoleServicePrincipals := []string{
		"ec2.amazonaws.com",
		"eks.amazonaws.com",
		"hello.amazonaws.com",
	}
	if !reflect.DeepEqual(expectedManagedNodeGroupRoleServicePrincipals, cfg.AddOnManagedNodeGroups.RoleServicePrincipals) {
		t.Fatalf("unexpected cfg.AddOnManagedNodeGroups.RoleServicePrincipals %+v", cfg.AddOnManagedNodeGroups.RoleServicePrincipals)
	}
	expectedManagedNodeGroupRoleManagedPolicyARNs := []string{
		"a",
		"b",
		"c",
	}
	if !reflect.DeepEqual(expectedManagedNodeGroupRoleManagedPolicyARNs, cfg.AddOnManagedNodeGroups.RoleManagedPolicyARNs) {
		t.Fatalf("unexpected cfg.AddOnManagedNodeGroups.RoleManagedPolicyARNs %+v", cfg.AddOnManagedNodeGroups.RoleManagedPolicyARNs)
	}
	if cfg.AddOnManagedNodeGroups.RequestHeaderKey != "a" {
		t.Fatalf("unexpected cfg.AddOnManagedNodeGroups.RequestHeaderKey %q", cfg.AddOnManagedNodeGroups.RequestHeaderKey)
	}
	if cfg.AddOnManagedNodeGroups.RequestHeaderValue != "b" {
		t.Fatalf("unexpected cfg.AddOnManagedNodeGroups.RequestHeaderValue %q", cfg.AddOnManagedNodeGroups.RequestHeaderValue)
	}
	if cfg.AddOnManagedNodeGroups.ResolverURL != "a" {
		t.Fatalf("unexpected cfg.AddOnManagedNodeGroups.ResolverURL %q", cfg.AddOnManagedNodeGroups.ResolverURL)
	}
	if cfg.AddOnManagedNodeGroups.SigningName != "a" {
		t.Fatalf("unexpected cfg.AddOnManagedNodeGroups.SigningName %q", cfg.AddOnManagedNodeGroups.SigningName)
	}
	if cfg.AddOnManagedNodeGroups.RemoteAccessUserName != "a" {
		t.Fatalf("unexpected cfg.AddOnManagedNodeGroups.RemoteAccessUserName %q", cfg.AddOnManagedNodeGroups.RemoteAccessUserName)
	}
	cpuName, gpuName := "mng-test-name-cpu", "mng-test-name-gpu"
	expectedMNGs := map[string]eksconfig.MNG{
		cpuName: {
			Name:               cpuName,
			Tags:               map[string]string{"cpu": "hello-world"},
			ReleaseVersion:     "test-ver-cpu",
			AMIType:            "AL2_x86_64",
			ASGMinSize:         17,
			ASGMaxSize:         99,
			ASGDesiredCapacity: 77,
			InstanceTypes:      []string{"type-cpu-1", "type-cpu-2"},
			VolumeSize:         40,
		},
		gpuName: {
			Name:               gpuName,
			Tags:               map[string]string{"gpu": "hello-world"},
			ReleaseVersion:     "test-ver-gpu",
			AMIType:            eks.AMITypesAl2X8664Gpu,
			ASGMinSize:         30,
			ASGMaxSize:         35,
			ASGDesiredCapacity: 34,
			InstanceTypes:      []string{"type-gpu-1", "type-gpu-2"},
			VolumeSize:         500,
		},
	}
	if !reflect.DeepEqual(cfg.AddOnManagedNodeGroups.MNGs, expectedMNGs) {
		t.Fatalf("expected cfg.AddOnManagedNodeGroups.MNGs %+v, got %+v", expectedMNGs, cfg.AddOnManagedNodeGroups.MNGs)
	}
	if cfg.AddOnManagedNodeGroups.LogsDir != "a" {
		t.Fatalf("unexpected cfg.AddOnManagedNodeGroups.LogsDir %q", cfg.AddOnManagedNodeGroups.LogsDir)
	}

	if cfg.AddOnNLBHelloWorld.Created { // read-only must be ignored
		t.Fatalf("unexpected cfg.AddOnNLBHelloWorld.Created %v", cfg.AddOnNLBHelloWorld.Created)
	}
	if !cfg.AddOnNLBHelloWorld.Enable {
		t.Fatalf("unexpected cfg.AddOnNLBHelloWorld.Enable %v", cfg.AddOnNLBHelloWorld.Enable)
	}
	if cfg.AddOnNLBHelloWorld.NLBARN != "" { // env should be ignored for read-only
		t.Fatalf("unexpected cfg.AddOnNLBHelloWorld.NLBARN %q", cfg.AddOnNLBHelloWorld.NLBARN)
	}
	if cfg.AddOnNLBHelloWorld.NLBName != "" { // env should be ignored for read-only
		t.Fatalf("unexpected cfg.AddOnNLBHelloWorld.NLBName %q", cfg.AddOnNLBHelloWorld.NLBName)
	}
	if cfg.AddOnNLBHelloWorld.URL != "" { // env should be ignored for read-only
		t.Fatalf("unexpected cfg.AddOnNLBHelloWorld.URL %q", cfg.AddOnNLBHelloWorld.URL)
	}
	if cfg.AddOnNLBHelloWorld.DeploymentReplicas != 333 {
		t.Fatalf("unexpected cfg.AddOnNLBHelloWorld.DeploymentReplicas %d", cfg.AddOnNLBHelloWorld.DeploymentReplicas)
	}
	if cfg.AddOnNLBHelloWorld.Namespace != "test-namespace" {
		t.Fatalf("unexpected cfg.AddOnNLBHelloWorld.Namespace %q", cfg.AddOnNLBHelloWorld.Namespace)
	}

	if cfg.AddOnALB2048.Created { // read-only must be ignored
		t.Fatalf("unexpected cfg.AddOnALB2048.Created %v", cfg.AddOnALB2048.Created)
	}
	if !cfg.AddOnALB2048.Enable {
		t.Fatalf("unexpected cfg.AddOnALB2048.Enable %v", cfg.AddOnALB2048.Enable)
	}
	if cfg.AddOnALB2048.ALBARN != "" { // env should be ignored for read-only
		t.Fatalf("unexpected cfg.AddOnALB2048.ALBARN %q", cfg.AddOnALB2048.ALBARN)
	}
	if cfg.AddOnALB2048.ALBName != "" { // env should be ignored for read-only
		t.Fatalf("unexpected cfg.AddOnALB2048.ALBName %q", cfg.AddOnALB2048.ALBName)
	}
	if cfg.AddOnALB2048.URL != "" { // env should be ignored for read-only
		t.Fatalf("unexpected cfg.AddOnALB2048.URL %q", cfg.AddOnALB2048.URL)
	}
	if cfg.AddOnALB2048.DeploymentReplicasALB != 333 {
		t.Fatalf("unexpected cfg.AddOnALB2048.DeploymentReplicasALB %d", cfg.AddOnALB2048.DeploymentReplicasALB)
	}
	if cfg.AddOnALB2048.DeploymentReplicas2048 != 555 {
		t.Fatalf("unexpected cfg.AddOnALB2048.DeploymentReplicas2048 %d", cfg.AddOnALB2048.DeploymentReplicas2048)
	}
	if cfg.AddOnALB2048.Namespace != "test-namespace" {
		t.Fatalf("unexpected cfg.AddOnALB2048.Namespace %q", cfg.AddOnALB2048.Namespace)
	}

	if cfg.AddOnJobPi.Created { // read-only must be ignored
		t.Fatalf("unexpected cfg.AddOnJobPi.Created %v", cfg.AddOnJobPi.Created)
	}
	if !cfg.AddOnJobPi.Enable {
		t.Fatalf("unexpected cfg.AddOnJobPi.Enable %v", cfg.AddOnJobPi.Enable)
	}
	if cfg.AddOnJobPi.Namespace != "hello1" {
		t.Fatalf("unexpected cfg.AddOnJobPi.Namespace %q", cfg.AddOnJobPi.Namespace)
	}
	if cfg.AddOnJobPi.Completes != 100 {
		t.Fatalf("unexpected cfg.AddOnJobPi.Completes %v", cfg.AddOnJobPi.Completes)
	}
	if cfg.AddOnJobPi.Parallels != 10 {
		t.Fatalf("unexpected cfg.AddOnJobPi.Parallels %v", cfg.AddOnJobPi.Parallels)
	}

	if cfg.AddOnJobEcho.Created { // read-only must be ignored
		t.Fatalf("unexpected cfg.AddOnJobEcho.Created %v", cfg.AddOnJobEcho.Created)
	}
	if !cfg.AddOnJobEcho.Enable {
		t.Fatalf("unexpected cfg.AddOnJobEcho.Enable %v", cfg.AddOnJobEcho.Enable)
	}
	if cfg.AddOnJobEcho.Namespace != "hello2" {
		t.Fatalf("unexpected cfg.AddOnJobEcho.Namespace %q", cfg.AddOnJobEcho.Namespace)
	}
	if cfg.AddOnJobEcho.Completes != 1000 {
		t.Fatalf("unexpected cfg.AddOnJobEcho.Completes %v", cfg.AddOnJobEcho.Completes)
	}
	if cfg.AddOnJobEcho.Parallels != 100 {
		t.Fatalf("unexpected cfg.AddOnJobEcho.Parallels %v", cfg.AddOnJobEcho.Parallels)
	}
	if cfg.AddOnJobEcho.EchoSize != 10000 {
		t.Fatalf("unexpected cfg.AddOnJobEcho.EchoSize %v", cfg.AddOnJobEcho.EchoSize)
	}

	if cfg.AddOnCronJob.Created { // read-only must be ignored
		t.Fatalf("unexpected cfg.AddOnCronJob.Created %v", cfg.AddOnCronJob.Created)
	}
	if !cfg.AddOnCronJob.Enable {
		t.Fatalf("unexpected cfg.AddOnCronJob.Enable %v", cfg.AddOnCronJob.Enable)
	}
	if cfg.AddOnCronJob.Namespace != "hello3" {
		t.Fatalf("unexpected cfg.AddOnCronJob.Namespace %q", cfg.AddOnCronJob.Namespace)
	}
	if cfg.AddOnCronJob.Schedule != "*/1 * * * *" {
		t.Fatalf("unexpected cfg.AddOnCronJob.Schedule %q", cfg.AddOnCronJob.Schedule)
	}
	if cfg.AddOnCronJob.Completes != 100 {
		t.Fatalf("unexpected cfg.AddOnCronJob.Completes %v", cfg.AddOnCronJob.Completes)
	}
	if cfg.AddOnCronJob.Parallels != 10 {
		t.Fatalf("unexpected cfg.AddOnCronJob.Parallels %v", cfg.AddOnCronJob.Parallels)
	}
	if cfg.AddOnCronJob.SuccessfulJobsHistoryLimit != 100 {
		t.Fatalf("unexpected cfg.AddOnCronJob.SuccessfulJobsHistoryLimit %v", cfg.AddOnCronJob.SuccessfulJobsHistoryLimit)
	}
	if cfg.AddOnCronJob.FailedJobsHistoryLimit != 1000 {
		t.Fatalf("unexpected cfg.AddOnCronJob.FailedJobsHistoryLimit %v", cfg.AddOnCronJob.FailedJobsHistoryLimit)
	}
	if cfg.AddOnCronJob.EchoSize != 10000 {
		t.Fatalf("unexpected cfg.AddOnCronJob.EchoSize %v", cfg.AddOnCronJob.EchoSize)
	}

	if cfg.AddOnSecrets.Created { // read-only must be ignored
		t.Fatalf("unexpected cfg.AddOnSecrets.Created %v", cfg.AddOnSecrets.Created)
	}
	if !cfg.AddOnSecrets.Enable {
		t.Fatalf("unexpected cfg.AddOnSecrets.Enable %v", cfg.AddOnSecrets.Enable)
	}
	if cfg.AddOnSecrets.Namespace != "hello" {
		t.Fatalf("unexpected cfg.AddOnSecrets.Namespace %q", cfg.AddOnSecrets.Namespace)
	}
	if cfg.AddOnSecrets.Objects != 5 {
		t.Fatalf("unexpected cfg.AddOnSecrets.Objects %v", cfg.AddOnSecrets.Objects)
	}
	if cfg.AddOnSecrets.Size != 10 {
		t.Fatalf("unexpected cfg.AddOnSecrets.Size %v", cfg.AddOnSecrets.Size)
	}
	if cfg.AddOnSecrets.SecretQPS != 10 {
		t.Fatalf("unexpected cfg.AddOnSecrets.SecretQPS %v", cfg.AddOnSecrets.SecretQPS)
	}
	if cfg.AddOnSecrets.SecretBurst != 10 {
		t.Fatalf("unexpected cfg.AddOnSecrets.SecretBurst %v", cfg.AddOnSecrets.SecretBurst)
	}
	if cfg.AddOnSecrets.PodQPS != 10 {
		t.Fatalf("unexpected cfg.AddOnSecrets.PodQPS %v", cfg.AddOnSecrets.PodQPS)
	}
	if cfg.AddOnSecrets.PodBurst != 10 {
		t.Fatalf("unexpected cfg.AddOnSecrets.PodBurst %v", cfg.AddOnSecrets.PodBurst)
	}
	if cfg.AddOnSecrets.WritesResultPath != "writes.csv" {
		t.Fatalf("unexpected cfg.AddOnSecrets.WritesResultPath %q", cfg.AddOnSecrets.WritesResultPath)
	}
	if cfg.AddOnSecrets.ReadsResultPath != "reads.csv" {
		t.Fatalf("unexpected cfg.AddOnSecrets.ReadsResultPath %q", cfg.AddOnSecrets.ReadsResultPath)
	}

	if cfg.AddOnIRSA.Created { // read-only must be ignored
		t.Fatalf("unexpected cfg.AddOnIRSA.Created %v", cfg.AddOnIRSA.Created)
	}
	if !cfg.AddOnIRSA.Enable {
		t.Fatalf("unexpected cfg.AddOnIRSA.Enable %v", cfg.AddOnIRSA.Enable)
	}
	if cfg.AddOnIRSA.Namespace != "hello" {
		t.Fatalf("unexpected cfg.AddOnIRSA.Namespace %q", cfg.AddOnIRSA.Namespace)
	}
	if cfg.AddOnIRSA.RoleName != "hello" {
		t.Fatalf("unexpected cfg.AddOnIRSA.RoleName %q", cfg.AddOnIRSA.RoleName)
	}
	expectedAddOnIRSARoleManagedPolicyARNs := []string{"a", "b", "c"}
	if !reflect.DeepEqual(cfg.AddOnIRSA.RoleManagedPolicyARNs, expectedAddOnIRSARoleManagedPolicyARNs) {
		t.Fatalf("unexpected cfg.AddOnIRSA.RoleManagedPolicyARNs %q", cfg.AddOnIRSA.RoleManagedPolicyARNs)
	}
	if cfg.AddOnIRSA.ServiceAccountName != "hello" {
		t.Fatalf("unexpected cfg.AddOnIRSA.ServiceAccountName %q", cfg.AddOnIRSA.ServiceAccountName)
	}
	if cfg.AddOnIRSA.ConfigMapName != "hello" {
		t.Fatalf("unexpected cfg.AddOnIRSA.ConfigMapName %q", cfg.AddOnIRSA.ConfigMapName)
	}
	if cfg.AddOnIRSA.ConfigMapScriptFileName != "hello.sh" {
		t.Fatalf("unexpected cfg.AddOnIRSA.ConfigMapScriptFileName %q", cfg.AddOnIRSA.ConfigMapScriptFileName)
	}
	if cfg.AddOnIRSA.S3BucketName != "hello" {
		t.Fatalf("unexpected cfg.AddOnIRSA.S3BucketName %q", cfg.AddOnIRSA.S3BucketName)
	}
	if cfg.AddOnIRSA.S3Key != "hello" {
		t.Fatalf("unexpected cfg.AddOnIRSA.S3Key %q", cfg.AddOnIRSA.S3Key)
	}
	if cfg.AddOnIRSA.DeploymentName != "hello-deployment" {
		t.Fatalf("unexpected cfg.AddOnIRSA.DeploymentName %q", cfg.AddOnIRSA.DeploymentName)
	}
	if cfg.AddOnIRSA.DeploymentResultPath != "hello-deployment.log" {
		t.Fatalf("unexpected cfg.AddOnIRSA.DeploymentResultPath %q", cfg.AddOnIRSA.DeploymentResultPath)
	}

	if cfg.AddOnFargate.Created { // read-only must be ignored
		t.Fatalf("unexpected cfg.AddOnFargate.Created %v", cfg.AddOnFargate.Created)
	}
	if !cfg.AddOnFargate.Enable {
		t.Fatalf("unexpected cfg.AddOnFargate.Enable %v", cfg.AddOnFargate.Enable)
	}
	if cfg.AddOnFargate.Namespace != "hello" {
		t.Fatalf("unexpected cfg.AddOnFargate.Namespace %q", cfg.AddOnFargate.Namespace)
	}
	if cfg.AddOnFargate.RoleName != "hello" {
		t.Fatalf("unexpected cfg.AddOnFargate.RoleName %q", cfg.AddOnFargate.RoleName)
	}
	expectedAddOnFargateRoleServicePrincipals := []string{"a", "b", "c"}
	if !reflect.DeepEqual(cfg.AddOnFargate.RoleServicePrincipals, expectedAddOnFargateRoleServicePrincipals) {
		t.Fatalf("unexpected cfg.AddOnFargate.RoleServicePrincipals %q", cfg.AddOnFargate.RoleServicePrincipals)
	}
	expectedAddOnFargateRoleManagedPolicyARNs := []string{"a", "b", "c"}
	if !reflect.DeepEqual(cfg.AddOnFargate.RoleManagedPolicyARNs, expectedAddOnFargateRoleManagedPolicyARNs) {
		t.Fatalf("unexpected cfg.AddOnFargate.RoleManagedPolicyARNs %q", cfg.AddOnFargate.RoleManagedPolicyARNs)
	}
	if cfg.AddOnFargate.ProfileName != "hello" {
		t.Fatalf("unexpected cfg.AddOnFargate.ProfileName %q", cfg.AddOnFargate.ProfileName)
	}
	if cfg.AddOnFargate.SecretName != "HELLO-SECRET" {
		t.Fatalf("unexpected cfg.AddOnFargate.SecretName %q", cfg.AddOnFargate.SecretName)
	}
	if cfg.AddOnFargate.PodName != "fargate-pod" {
		t.Fatalf("unexpected cfg.AddOnFargate.PodName %q", cfg.AddOnFargate.PodName)
	}
	if cfg.AddOnFargate.ContainerName != "fargate-container" {
		t.Fatalf("unexpected cfg.AddOnFargate.ContainerName %q", cfg.AddOnFargate.ContainerName)
	}

	cfg.Parameters.RoleManagedPolicyARNs = nil
	cfg.Parameters.RoleServicePrincipals = nil
	cfg.AddOnManagedNodeGroups.RoleName = ""
	cfg.AddOnManagedNodeGroups.RoleManagedPolicyARNs = nil
	cfg.AddOnManagedNodeGroups.RoleServicePrincipals = nil
	if err := cfg.ValidateAndSetDefaults(); err != nil {
		t.Fatal(err)
	}
	cfg.AddOnNLBHelloWorld.Enable = false
	cfg.AddOnALB2048.Enable = false
	cfg.AddOnJobEcho.Enable = false
	cfg.AddOnJobPi.Enable = false
	if err := cfg.ValidateAndSetDefaults(); err != nil {
		t.Fatal(err)
	}

	if cfg.AddOnFargate.SecretName != "hellosecret" {
		t.Fatalf("unexpected cfg.AddOnFargate.SecretName %q", cfg.AddOnFargate.SecretName)
	}

	d, err := ioutil.ReadFile(cfg.ConfigPath)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(d))
}
