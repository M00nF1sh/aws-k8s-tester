package eks

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	pkgaws "github.com/aws/aws-k8s-tester/pkg/aws"
	"github.com/aws/aws-k8s-tester/pkg/logutil"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

var (
	deleteDry    bool
	deleteFailed bool
	deletePrefix string
	deleteAgo    time.Duration
)

func newList() *cobra.Command {
	ac := &cobra.Command{
		Use:   "list <subcommand>",
		Short: "List EKS resources",
	}
	ac.PersistentFlags().StringVar(&region, "region", "us-west-2", "EKS region")
	ac.PersistentFlags().StringVar(&resolverURL, "resolver-url", "", "EKS resolver endpoint URL")
	ac.PersistentFlags().StringVar(&signingName, "signing-name", "", "EKS signing name")

	ac.PersistentFlags().BoolVar(
		&deleteDry,
		"delete-dry",
		true,
		"'true' to delete clusters in dry mode",
	)
	ac.PersistentFlags().BoolVar(
		&deleteFailed,
		"delete-failed",
		false,
		"'true' to clean up failed clusters",
	)
	ac.PersistentFlags().StringVar(
		&deletePrefix,
		"delete-prefix",
		"",
		"Cluster name prefix to match and delete",
	)
	ac.PersistentFlags().DurationVar(
		&deleteAgo,
		"delete-ago",
		0,
		"Duration to delete clusters created x-duration ago",
	)

	ac.AddCommand(
		newListClusters(),
	)
	return ac
}

func newListClusters() *cobra.Command {
	return &cobra.Command{
		Use:   "clusters",
		Short: "List EKS clusters",
		Run:   listClustersFunc,
	}
}

func listClustersFunc(cmd *cobra.Command, args []string) {
	lg, _ := logutil.GetDefaultZapLogger()
	awsCfgEKS := &pkgaws.Config{
		Logger:      lg,
		Region:      region,
		ResolverURL: resolverURL,
		SigningName: signingName,
	}
	ssEKS, _, _, err := pkgaws.New(awsCfgEKS)
	if err != nil {
		panic(err)
	}
	svc := awseks.New(ssEKS)
	clusterNames := make([]string, 0)
	if err = svc.ListClustersPages(&awseks.ListClustersInput{},
		func(output *awseks.ListClustersOutput, lastPage bool) bool {
			for _, name := range output.Clusters {
				clusterNames = append(clusterNames, aws.StringValue(name))
			}
			return true
		}); err != nil {
		panic(err)
	}
	if len(clusterNames) == 0 {
		fmt.Println("'aws-k8s-tester eks list clusters' returned 0 cluster")
		return
	}

	fmt.Printf("Listing %d clusters\n", len(clusterNames))

	for i, name := range clusterNames {
		out, err := svc.DescribeCluster(&awseks.DescribeClusterInput{
			Name: aws.String(name),
		})
		if err != nil {
			fmt.Printf("[%03d] %q failed to describe (%v, retriable %v, throttled %v, error type %v)\n",
				i,
				name,
				err,
				request.IsErrorRetryable(err),
				request.IsErrorThrottle(err),
				reflect.TypeOf(err),
			)

			awsErr, ok := err.(awserr.Error)
			if ok && awsErr.Code() == "ResourceNotFoundException" &&
				strings.HasPrefix(awsErr.Message(), "No cluster found for") {
				fmt.Printf("deleting %q (reason: %v)\n", name, err)
				if !deleteDry {
					_, derr := svc.DeleteCluster(&awseks.DeleteClusterInput{Name: aws.String(name)})
					fmt.Println("deleted", name, derr)
				}
			}

			time.Sleep(3 * time.Second)
			println()
			continue
		}
		if out.Cluster == nil {
			panic(fmt.Errorf("[%03d] %q empty cluster", i, name))
		}

		clus := out.Cluster

		createdAtUTC := aws.TimeValue(clus.CreatedAt).UTC()
		nowUTC := time.Now().UTC()

		fmt.Printf("[%03d] %q [created %v (%q), version %q, status %q, IAM Role %q, VPC %q]\n",
			i,
			name,
			createdAtUTC,
			humanize.RelTime(createdAtUTC, nowUTC, "ago", "from now"),
			aws.StringValue(clus.Version),
			aws.StringValue(clus.Status),
			aws.StringValue(clus.RoleArn),
			aws.StringValue(clus.ResourcesVpcConfig.VpcId),
		)

		if deleteFailed && aws.StringValue(clus.Status) == "FAILED" {
			fmt.Printf("deleting %q (reason: %v)\n", name, aws.StringValue(clus.Status))
			if !deleteDry {
				_, derr := svc.DeleteCluster(&awseks.DeleteClusterInput{Name: aws.String(name)})
				fmt.Println("deleted", name, derr)
			}

			time.Sleep(3 * time.Second)
			println()
			continue
		}

		if len(deletePrefix) > 0 && strings.HasPrefix(name, deletePrefix) {
			fmt.Printf("deleting %q (reason: %q)\n", name, deletePrefix)
			if !deleteDry {
				_, derr := svc.DeleteCluster(&awseks.DeleteClusterInput{Name: aws.String(name)})
				fmt.Println("deleted", name, derr)
			}

			time.Sleep(3 * time.Second)
			println()
			continue
		}

		createDur := nowUTC.Sub(createdAtUTC)
		if deleteAgo > 0 && createDur > deleteAgo {
			fmt.Printf("deleting %q (reason: create at %v, delete ago %v, create duration %v)\n",
				name,
				createdAtUTC,
				deleteAgo,
				createDur,
			)
			if !deleteDry {
				_, derr := svc.DeleteCluster(&awseks.DeleteClusterInput{Name: aws.String(name)})
				fmt.Println("deleted", name, derr)
			}

			time.Sleep(3 * time.Second)
			println()
			continue
		}
	}

	fmt.Println("Successfully listed", len(clusterNames), "clusters")
}
