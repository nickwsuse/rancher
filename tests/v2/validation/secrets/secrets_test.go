package secrets

import (
	"testing"

	"github.com/rancher/machine/libmachine/log"
	"github.com/rancher/shepherd/clients/rancher"
	management "github.com/rancher/shepherd/clients/rancher/generated/management/v3"
	"github.com/rancher/shepherd/extensions/clusters"
	"github.com/rancher/shepherd/extensions/secrets"
	"github.com/rancher/shepherd/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SecretTestSuite struct {
	suite.Suite
	client  *rancher.Client
	session *session.Session
	cluster *management.Cluster
}

func (s *SecretTestSuite) TearDownSuite() {
	s.session.Cleanup()
}

func (s *SecretTestSuite) SetupSuite() {
	s.session = session.NewSession()

	client, err := rancher.NewClient("", s.session)
	require.NoError(s.T(), err)
	s.client = client

	clusterName := client.RancherConfig.ClusterName
	require.NotEmptyf(s.T(), clusterName, "Cluster name to install should be set")
	clusterID, err := clusters.GetClusterIDByName(s.client, clusterName)
	require.NoError(s.T(), err, "Error getting cluster ID")
	s.cluster, err = s.client.Management.Cluster.ByID(clusterID)
	require.NoError(s.T(), err)
}

func (s *SecretTestSuite) TestSecretCreate() {
	subSession := s.session.NewSession()
	defer subSession.Cleanup()

	log.Info("Creating a secret")
	secret.Labels = map[string]string{labelKey: labelVal}
	secretObj, err := s.client.Steve.SteveType(secrets.SecretSteveType).Create(secret)
	require.NoError(s.T(), err)

	log.Info("Validating secret was created with correct resource values")
	labels := getSecretLabelsAndAnnotations(secretObj.ObjectMeta.Labels)
	assert.Contains(s.T(), secretName, secretObj.Name)
	assert.Equal(s.T(), labels, secretObj.Labels)
}

func (s *SecretTestSuite) TestSecretUpdate() {
	subSession := s.session.NewSession()
	defer subSession.Cleanup()

	secretClient := s.client.Steve.SteveType(secrets.SecretSteveType)
	newAnnotations := map[string]string{updatedAnnoKey: updatedAnnoVal, descKey: descVal}

	log.Info("Creating a secret")
	secret.Labels = map[string]string{labelKey: labelVal}
	secretObj, err := s.client.Steve.SteveType(secrets.SecretSteveType).Create(secret)
	require.NoError(s.T(), err)

	log.Info("Updating the secret")
	newSecret := secretObj
	newSecret.ObjectMeta.Annotations = newAnnotations
	updatedSecretObj, err := secretClient.Update(secretObj, newSecret)
	require.NoError(s.T(), err)

	log.Info("Validating secret was properly updated")
	expectedAnnotations := getSecretLabelsAndAnnotations(updatedSecretObj.ObjectMeta.Annotations)
	assert.Equal(s.T(), expectedAnnotations, newAnnotations)
}

func (s *SecretTestSuite) TestSecretDelete() {
	subSession := s.session.NewSession()
	defer subSession.Cleanup()

	secretClient := s.client.Steve.SteveType(secrets.SecretSteveType)

	log.Info("Creating a secret")
	secret.Labels = map[string]string{labelKey: labelVal}
	secretObj, err := s.client.Steve.SteveType(secrets.SecretSteveType).Create(secret)
	require.NoError(s.T(), err)

	log.Info("Deleting the secret")
	err = secretClient.Delete(secretObj)
	require.NoError(s.T(), err)

	log.Info("Validating secret was deleted")
	secretByID, err := secretClient.ByID(secretObj.ID)
	require.Error(s.T(), err)
	assert.Nil(s.T(), secretByID)
	assert.ErrorContains(s.T(), err, "not found")
}

func TestSecretSuite(t *testing.T) {
	suite.Run(t, new(SecretTestSuite))
}
