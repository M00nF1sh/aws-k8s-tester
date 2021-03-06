package eks

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	awscfn "github.com/aws/aws-k8s-tester/pkg/aws/cloudformation"
	awsiam "github.com/aws/aws-k8s-tester/pkg/aws/iam"
	"github.com/aws/aws-k8s-tester/version"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"go.uber.org/zap"
)

// TemplateClusterRole is the CloudFormation template for EKS cluster role.
//
// e.g.
//   Error creating load balancer (will retry): failed to ensure load balancer for service eks-*/hello-world-service: Error creating load balancer: "AccessDenied: User: arn:aws:sts::404174646922:assumed-role/eks-*-cluster-role/* is not authorized to perform: ec2:DescribeAccountAttributes\n\tstatus code: 403"
//
// TODO: scope down (e.g. ec2:DescribeAccountAttributes, ec2:DescribeInternetGateways)
// mng, fargate, etc. may use other roles
const TemplateClusterRole = `
---
AWSTemplateFormatVersion: '2010-09-09'
Description: 'Amazon EKS Cluster Role Basic'

Parameters:

  RoleName:
    Description: EKS Role name
    Type: String
    Default: aws-k8s-tester-eks-role

  RoleServicePrincipals:
    Description: EKS Role Service Principals
    Type: CommaDelimitedList
    Default: 'ec2.amazonaws.com,eks.amazonaws.com,eks-fargate-pods.amazonaws.com'

  RoleManagedPolicyARNs:
    Description: EKS Role managed policy ARNs
    Type: CommaDelimitedList
    Default: 'arn:aws:iam::aws:policy/AmazonEKSServicePolicy,arn:aws:iam::aws:policy/AmazonEKSClusterPolicy,arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy,arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy,arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly,arn:aws:iam::aws:policy/ElasticLoadBalancingFullAccess'

Resources:

  Role:
    Type: AWS::IAM::Role
    Properties:
      RoleName: !Ref RoleName
      AssumeRolePolicyDocument:
        Version: '2012-10-17'
        Statement:
        - Effect: Allow
          Principal:
            Service: !Ref RoleServicePrincipals
          Action:
          - sts:AssumeRole
      ManagedPolicyArns: !Ref RoleManagedPolicyARNs
      Path: /
      Policies:
      # https://github.com/kubernetes-sigs/aws-alb-ingress-controller/blob/master/docs/examples/iam-policy.json
      - PolicyName: !Join ['-', [!Ref RoleName, 'alb-policy']]
        PolicyDocument:
          Version: '2012-10-17'
          Statement:
          - Effect: Allow
            Action:
            - acm:DescribeCertificate
            - acm:ListCertificates
            - acm:GetCertificate
            Resource: "*"
          - Effect: Allow
            Action:
            - ec2:AuthorizeSecurityGroupIngress
            - ec2:CreateSecurityGroup
            - ec2:CreateTags
            - ec2:DeleteTags
            - ec2:DeleteSecurityGroup
            - ec2:DescribeAccountAttributes
            - ec2:DescribeAddresses
            - ec2:DescribeInstances
            - ec2:DescribeInstanceStatus
            - ec2:DescribeInternetGateways
            - ec2:DescribeNetworkInterfaces
            - ec2:DescribeSecurityGroups
            - ec2:DescribeSubnets
            - ec2:DescribeTags
            - ec2:DescribeVpcs
            - ec2:ModifyInstanceAttribute
            - ec2:ModifyNetworkInterfaceAttribute
            - ec2:RevokeSecurityGroupIngress
            Resource: "*"
          - Effect: Allow
            Action:
            - elasticloadbalancing:AddListenerCertificates
            - elasticloadbalancing:AddTags
            - elasticloadbalancing:CreateListener
            - elasticloadbalancing:CreateLoadBalancer
            - elasticloadbalancing:CreateRule
            - elasticloadbalancing:CreateTargetGroup
            - elasticloadbalancing:DeleteListener
            - elasticloadbalancing:DeleteLoadBalancer
            - elasticloadbalancing:DeleteRule
            - elasticloadbalancing:DeleteTargetGroup
            - elasticloadbalancing:DeregisterTargets
            - elasticloadbalancing:DescribeListenerCertificates
            - elasticloadbalancing:DescribeListeners
            - elasticloadbalancing:DescribeLoadBalancers
            - elasticloadbalancing:DescribeLoadBalancerAttributes
            - elasticloadbalancing:DescribeRules
            - elasticloadbalancing:DescribeSSLPolicies
            - elasticloadbalancing:DescribeTags
            - elasticloadbalancing:DescribeTargetGroups
            - elasticloadbalancing:DescribeTargetGroupAttributes
            - elasticloadbalancing:DescribeTargetHealth
            - elasticloadbalancing:ModifyListener
            - elasticloadbalancing:ModifyLoadBalancerAttributes
            - elasticloadbalancing:ModifyRule
            - elasticloadbalancing:ModifyTargetGroup
            - elasticloadbalancing:ModifyTargetGroupAttributes
            - elasticloadbalancing:RegisterTargets
            - elasticloadbalancing:RemoveListenerCertificates
            - elasticloadbalancing:RemoveTags
            - elasticloadbalancing:SetIpAddressType
            - elasticloadbalancing:SetSecurityGroups
            - elasticloadbalancing:SetSubnets
            - elasticloadbalancing:SetWebACL
            Resource: "*"
          - Effect: Allow
            Action:
            - iam:CreateServiceLinkedRole
            - iam:GetServerCertificate
            - iam:ListServerCertificates
            Resource: "*"
          - Effect: Allow
            Action:
            - cognito-idp:DescribeUserPoolClient
            Resource: "*"
          - Effect: Allow
            Action:
            - waf-regional:GetWebACLForResource
            - waf-regional:GetWebACL
            - waf-regional:AssociateWebACL
            - waf-regional:DisassociateWebACL
            Resource: "*"
          - Effect: Allow
            Action:
            - tag:GetResources
            - tag:TagResources
            Resource: "*"
          - Effect: Allow
            Action:
            - waf:GetWebACL
            Resource: "*"
          - Effect: Allow
            Action:
            - shield:DescribeProtection
            - shield:GetSubscriptionState
            - shield:DeleteProtection
            - shield:CreateProtection
            - shield:DescribeSubscription
            - shield:ListProtections
            Resource: "*"

Outputs:

  RoleARN:
    Description: Role ARN that EKS uses to create AWS resources for Kubernetes
    Value: !GetAtt Role.Arn

`

func (ts *Tester) createClusterRole() error {
	if !ts.cfg.Parameters.RoleCreate {
		ts.lg.Info("Parameters.RoleCreate false; skipping creation")
		return awsiam.Validate(
			ts.lg,
			ts.iamAPI,
			ts.cfg.Parameters.RoleName,
			[]string{"eks.amazonaws.com"},
			[]string{
				"arn:aws:iam::aws:policy/AmazonEKSServicePolicy",
				"arn:aws:iam::aws:policy/AmazonEKSClusterPolicy",
			},
		)
	}
	if ts.cfg.Parameters.RoleCFNStackID != "" &&
		ts.cfg.Parameters.RoleARN != "" {
		ts.lg.Info("role already created; no need to create a new one")
		return nil
	}
	if ts.cfg.Parameters.RoleName == "" {
		return errors.New("cannot create a cluster role with an empty Parameters.RoleName")
	}

	tmpl := TemplateClusterRole

	// role ARN is empty, create a default role
	// otherwise, use the existing one
	ts.lg.Info("creating a new role", zap.String("cluster-role-name", ts.cfg.Parameters.RoleName))
	stackInput := &cloudformation.CreateStackInput{
		StackName:    aws.String(ts.cfg.Parameters.RoleName),
		Capabilities: aws.StringSlice([]string{"CAPABILITY_NAMED_IAM"}),
		OnFailure:    aws.String(cloudformation.OnFailureDelete),
		TemplateBody: aws.String(tmpl),
		Tags: awscfn.NewTags(map[string]string{
			"Kind":                   "aws-k8s-tester",
			"Name":                   ts.cfg.Name,
			"aws-k8s-tester-version": version.ReleaseVersion,
		}),
		Parameters: []*cloudformation.Parameter{
			{
				ParameterKey:   aws.String("RoleName"),
				ParameterValue: aws.String(ts.cfg.Parameters.RoleName),
			},
		},
	}
	if len(ts.cfg.Parameters.RoleServicePrincipals) > 0 {
		ts.lg.Info("creating a new cluster role with custom service principals",
			zap.Strings("service-principals", ts.cfg.Parameters.RoleServicePrincipals),
		)
		stackInput.Parameters = append(stackInput.Parameters, &cloudformation.Parameter{
			ParameterKey:   aws.String("RoleServicePrincipals"),
			ParameterValue: aws.String(strings.Join(ts.cfg.Parameters.RoleServicePrincipals, ",")),
		})
	}
	if len(ts.cfg.Parameters.RoleManagedPolicyARNs) > 0 {
		ts.lg.Info("creating a new cluster role with custom managed role policies",
			zap.Strings("policy-arns", ts.cfg.Parameters.RoleManagedPolicyARNs),
		)
		stackInput.Parameters = append(stackInput.Parameters, &cloudformation.Parameter{
			ParameterKey:   aws.String("RoleManagedPolicyARNs"),
			ParameterValue: aws.String(strings.Join(ts.cfg.Parameters.RoleManagedPolicyARNs, ",")),
		})
	}
	stackOutput, err := ts.cfnAPI.CreateStack(stackInput)
	if err != nil {
		return err
	}
	ts.cfg.Parameters.RoleCFNStackID = aws.StringValue(stackOutput.StackId)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	ch := awscfn.Poll(
		ctx,
		ts.stopCreationCh,
		ts.interruptSig,
		ts.lg,
		ts.cfnAPI,
		ts.cfg.Parameters.RoleCFNStackID,
		cloudformation.ResourceStatusCreateComplete,
		25*time.Second,
		10*time.Second,
	)
	var st awscfn.StackStatus
	for st = range ch {
		if st.Error != nil {
			cancel()
			ts.cfg.RecordStatus(fmt.Sprintf("failed to create role (%v)", st.Error))
			ts.lg.Warn("polling errror", zap.Error(st.Error))
		}
	}
	cancel()
	if st.Error != nil {
		return st.Error
	}
	// update status after creating a new IAM role
	for _, o := range st.Stack.Outputs {
		switch k := aws.StringValue(o.OutputKey); k {
		case "RoleARN":
			ts.cfg.Parameters.RoleARN = aws.StringValue(o.OutputValue)
		default:
			return fmt.Errorf("unexpected OutputKey %q from %q", k, ts.cfg.Parameters.RoleCFNStackID)
		}
	}

	ts.lg.Info("created a new role",
		zap.String("cluster-role-cfn-stack-id", ts.cfg.Parameters.RoleCFNStackID),
		zap.String("cluster-role-arn", ts.cfg.Parameters.RoleARN),
	)
	return ts.cfg.Sync()
}

func (ts *Tester) deleteClusterRole() error {
	if !ts.cfg.Parameters.RoleCreate {
		ts.lg.Info("Parameters.RoleCreate false; skipping deletion")
		return nil
	}
	if ts.cfg.Parameters.RoleCFNStackID == "" {
		ts.lg.Info("empty role CFN stack ID; no need to delete role")
		return nil
	}

	ts.lg.Info("deleting role CFN stack", zap.String("cluster-role-cfn-stack-id", ts.cfg.Parameters.RoleCFNStackID))
	_, err := ts.cfnAPI.DeleteStack(&cloudformation.DeleteStackInput{
		StackName: aws.String(ts.cfg.Parameters.RoleCFNStackID),
	})
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	ch := awscfn.Poll(
		ctx,
		make(chan struct{}),  // do not exit on stop
		make(chan os.Signal), // do not exit on stop
		ts.lg,
		ts.cfnAPI,
		ts.cfg.Parameters.RoleCFNStackID,
		cloudformation.ResourceStatusDeleteComplete,
		25*time.Second,
		10*time.Second,
	)
	var st awscfn.StackStatus
	for st = range ch {
		if st.Error != nil {
			cancel()
			ts.cfg.RecordStatus(fmt.Sprintf("failed to delete role (%v)", st.Error))
			ts.lg.Warn("polling errror", zap.Error(st.Error))
		}
	}
	cancel()
	if st.Error != nil {
		return st.Error
	}
	ts.lg.Info("deleted a role",
		zap.String("cluster-role-cfn-stack-id", ts.cfg.Parameters.RoleCFNStackID),
		zap.String("cluster-role-arn", ts.cfg.Parameters.RoleARN),
		zap.String("cluster-role-name", ts.cfg.Parameters.RoleName),
	)
	return ts.cfg.Sync()
}
