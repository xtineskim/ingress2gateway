/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cilium

import (
	"testing"

	"github.com/kubernetes-sigs/ingress2gateway/pkg/i2gw"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

func Test_ToGateway(t *testing.T) {
	// iPrefix := networkingv1.PathTypePrefix
	gPathPrefix := gatewayv1beta1.PathMatchPathPrefix
	isPathType := networkingv1.PathTypeImplementationSpecific

	testCases := []struct {
		name                     string
		ingresses                map[types.NamespacedName]*networkingv1.Ingress
		expectedGatewayResources i2gw.GatewayResources
		expectedErrors           field.ErrorList
	}{{
		name: "ImplementationSpecific HTTPRouteMatching",
		ingresses: map[types.NamespacedName]*networkingv1.Ingress{
			{Namespace: "default", Name: "implementation-specific-regex"}: {
				ObjectMeta: metav1.ObjectMeta{
					Name:      "implementation-specific-regex",
					Namespace: "default",
				},
				Spec: networkingv1.IngressSpec{
					IngressClassName: ptrTo("ingress-nginx"),
					Rules: []networkingv1.IngressRule{{
						Host: "test.mydomain.com",
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{{
									Path:     "/~/echo/**/test",
									PathType: &isPathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: "test",
											Port: networkingv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								}},
							},
						},
					}},
				},
			},
		},
		expectedGatewayResources: i2gw.GatewayResources{},
		expectedErrors: field.ErrorList{
			{
				Type:     field.ErrorTypeInvalid,
				Field:    "spec.rules[0].http.paths[0].pathType",
				BadValue: ptr.To("ImplementationSpecific"),
				Detail:   "implementationSpecific path type is not supported in generic translation, and your provider does not provide custom support to translate it",
			},
		}}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			provider := NewProvider(&i2gw.ProviderConf{})
			resources := i2gw.InputResources{
				Ingresses:       tc.ingresses,
				CustomResources: nil,
			}

			gatewayResources, errs := provider.ToGatewayAPI(resources)

			if len(gatewayResources.HTTPRoutes) != len(tc.expectedGatewayResources.HTTPRoutes) {
				t.Errorf("Expected %d HTTPRoutes, got %d: %+v", len(tc.expectedGatewayResources.HTTPRoutes), len(gatewayResources.HTTPRoutes), gatewayResources.HTTPRoutes)
			} else {

			}

		})

	}
}

func ptrTo[T any](a T) *T {
	return &a
}
