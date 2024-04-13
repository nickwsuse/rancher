package secrets

import (
	"github.com/rancher/shepherd/pkg/namegenerator"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	namespace      = "default"
	labelKey       = "label1"
	labelVal       = "autoLabel"
	descKey        = "field.cattle.io/description"
	descVal        = "automated secret description"
	name           = "steve-secret"
	annoKey        = "anno1"
	annoVal        = "automated annotation"
	dataKey        = "foo"
	dataVal        = "bar"
	updatedAnnoKey = "newAnno"
	updatedAnnoVal = "updated annotation"
)

var (
	secretName = namegenerator.AppendRandomString(name)

	secret = coreV1.Secret{
		Type: coreV1.SecretTypeOpaque,
		ObjectMeta: metav1.ObjectMeta{
			Name:        secretName,
			Namespace:   namespace,
			Annotations: map[string]string{annoKey: annoVal},
			Labels:      map[string]string{labelKey: labelVal},
		},
		Data: map[string][]byte{dataKey: []byte(dataVal)},
	}
)

func getSecretLabelsAndAnnotations(actualResources map[string]string) map[string]string {
	expectedResources := map[string]string{}

	for resource := range actualResources {
		if _, found := actualResources[resource]; found {
			expectedResources[resource] = actualResources[resource]
		}
	}

	return expectedResources
}
