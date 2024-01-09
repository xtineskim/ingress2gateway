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
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

func Test_ToGateway(t *testing.T) {
	iPrefix := networkingv1.PathTypePrefix
	// gPathPrefix := gatewayv1beta1.PathMatchPathPrefix
	// isPathType := networkingv1.PathTypeImplementationSpecific

	testCases := []struct {
		name                     string
		ingresses                []networkingv1.Ingress
		expectedGatewayResources i2gw.GatewayResources
		expectedErrors           field.ErrorList
	}{
		{
			name: "test",
			ingresses: []networkingv1.Ingress{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "cilium-ingress-basic",
						Namespace: "default",
					},
					Spec: networkingv1.IngressSpec{
						IngressClassName: ptrTo("ingress-cilium"),
						Rules: []networkingv1.IngressRule{{
							Host: "test.mydomain.com",
							IngressRuleValue: networkingv1.IngressRuleValue{
								HTTP: &networkingv1.HTTPIngressRuleValue{
									Paths: []networkingv1.HTTPIngressPath{{
										Path:     "/",
										PathType: &iPrefix,
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
			expectedGatewayResources: i2gw.GatewayResources{
				Gateways: map[types.NamespacedName]gatewayv1beta1.Gateway{
					{
						Namespace: "default",
						Name:      "cilium-ingress",
					}: {
						ObjectMeta: metav1.ObjectMeta{Name: "cilium-ingress", Namespace: "default"},
						Spec: gatewayv1beta1.GatewaySpec{
							GatewayClassName: "cilium",
							Listeners: []gatewayv1beta1.Listener{{
								Name:     "test-mydomain-com-http",
								Port:     80,
								Protocol: gatewayv1beta1.HTTPProtocolType,
								Hostname: ptrTo(gatewayv1beta1.Hostname("test.mydomain.com")),
							}},
						},
					},
				},
				HTTPRoutes: map[types.NamespacedName]gatewayv1beta1.HTTPRoute{
					{Namespace: "default", Name: "cilium-ingress-basic"}: {
						ObjectMeta: metav1.ObjectMeta{Name: "cilium-ingress-basic", Namespace: "default"},
						Spec: gatewayv1beta1.HTTPRouteSpec{
							CommonRouteSpec: gatewayv1beta1.CommonRouteSpec{
								ParentRefs: []gatewayv1beta1.ParentReference{{
									Name: "cilium",
								}},
							},
							Hostnames: []gatewayv1beta1.Hostname{"test.mydomain.com"},
						},
					},
				},
			},
			expectedErrors: field.ErrorList{},
		},
	}
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
