package eks

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func (ts *Tester) createK8sClientSet() (err error) {
	ts.lg.Info("loading *restclient.Config")
	cfg := ts.createClientConfig()
	if cfg == nil {
		ts.lg.Warn("*restclient.Config is nil; reading kubeconfig")
		cfg, err = clientcmd.BuildConfigFromFlags("", ts.cfg.KubeConfigPath)
		if err != nil {
			ts.lg.Warn("failed to read kubeconfig", zap.Error(err))
			return err
		}
		ts.lg.Info("loaded *restclient.Config from kubeconfig")
	} else {
		ts.lg.Info("loaded *restclient.Config from eksconfig")
	}

	ts.lg.Info("creating k8s client set")
	ts.k8sClientSet, err = clientset.NewForConfig(cfg)
	if err == nil {
		ts.lg.Info("created k8s client set")
		return nil
	}

	ts.lg.Warn("failed to create k8s client set", zap.Error(err))
	return err
}

const authProviderName = "eks"

func (ts *Tester) createClientConfig() *restclient.Config {
	if ts.cfg.Name == "" {
		return nil
	}
	if ts.cfg.Region == "" {
		return nil
	}
	if ts.cfg.Status.ClusterAPIServerEndpoint == "" {
		return nil
	}
	if ts.cfg.Status.ClusterCADecoded == "" {
		return nil
	}
	return &restclient.Config{
		Host: ts.cfg.Status.ClusterAPIServerEndpoint,
		TLSClientConfig: restclient.TLSClientConfig{
			CAData: []byte(ts.cfg.Status.ClusterCADecoded),
		},
		AuthProvider: &clientcmdapi.AuthProviderConfig{
			Name: authProviderName,
			Config: map[string]string{
				"region":       ts.cfg.Region,
				"cluster-name": ts.cfg.Name,
			},
		},
	}
}

func init() {
	restclient.RegisterAuthProviderPlugin(authProviderName, newAuthProvider)
}

func newAuthProvider(_ string, config map[string]string, _ restclient.AuthProviderConfigPersister) (restclient.AuthProvider, error) {
	awsRegion, ok := config["region"]
	if !ok {
		return nil, fmt.Errorf("'clientcmdapi.AuthProviderConfig' does not include 'region' key %+v", config)
	}
	clusterName, ok := config["cluster-name"]
	if !ok {
		return nil, fmt.Errorf("'clientcmdapi.AuthProviderConfig' does not include 'cluster-name' key %+v", config)
	}

	sess := session.Must(session.NewSession(aws.NewConfig().WithRegion(awsRegion)))
	return &eksAuthProvider{ts: newTokenSource(sess, clusterName)}, nil
}

type eksAuthProvider struct {
	ts oauth2.TokenSource
}

func (p *eksAuthProvider) WrapTransport(rt http.RoundTripper) http.RoundTripper {
	return &oauth2.Transport{
		Source: p.ts,
		Base:   rt,
	}
}

func (p *eksAuthProvider) Login() error {
	return nil
}

func newTokenSource(sess *session.Session, clusterName string) oauth2.TokenSource {
	return &eksTokenSource{sess: sess, clusterName: clusterName}
}

type eksTokenSource struct {
	sess        *session.Session
	clusterName string
}

// Reference
// https://github.com/kubernetes-sigs/aws-iam-authenticator/blob/master/README.md#api-authorization-from-outside-a-cluster
// https://github.com/kubernetes-sigs/aws-iam-authenticator/blob/master/pkg/token/token.go
const (
	v1Prefix        = "k8s-aws-v1."
	clusterIDHeader = "x-k8s-aws-id"
)

func (s *eksTokenSource) Token() (*oauth2.Token, error) {
	stsAPI := sts.New(s.sess)
	request, _ := stsAPI.GetCallerIdentityRequest(&sts.GetCallerIdentityInput{})
	request.HTTPRequest.Header.Add(clusterIDHeader, s.clusterName)

	payload, err := request.Presign(60)
	if err != nil {
		return nil, err
	}
	token := v1Prefix + base64.RawURLEncoding.EncodeToString([]byte(payload))
	tokenExpiration := time.Now().Local().Add(14 * time.Minute)
	return &oauth2.Token{
		AccessToken: token,
		TokenType:   "Bearer",
		Expiry:      tokenExpiration,
	}, nil
}

func (ts *Tester) getPods(ns string) (*v1.PodList, error) {
	return ts.k8sClientSet.CoreV1().Pods(ns).List(metav1.ListOptions{})
}

func (ts *Tester) getAllNodes() (*v1.NodeList, error) {
	return ts.k8sClientSet.CoreV1().Nodes().List(metav1.ListOptions{})
}
