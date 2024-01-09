/*
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
	"github.com/kubernetes-sigs/ingress2gateway/pkg/i2gw"
	"github.com/kubernetes-sigs/ingress2gateway/pkg/i2gw/providers/common"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type converter struct {
	conf                          *i2gw.ProviderConf
	featuresParsers               []i2gw.FeatureParser
	implementationSpecificOptions i2gw.ProviderImplementationSpecificOptions
}

func newConverter(conf *i2gw.ProviderConf) *converter {
	return &converter{
		conf:                          conf,
		featuresParsers:               []i2gw.FeatureParser{},
		implementationSpecificOptions: i2gw.ProviderImplementationSpecificOptions{},
	}
}

func (c *converter) ToGatewayAPI(resources i2gw.InputResources) (i2gw.GatewayResources, field.ErrorList) {
	gatewayResources, errs := common.ToGateway(resources.Ingresses, i2gw.ProviderImplementationSpecificOptions{})
	if len(errs) > 0 {
		return i2gw.GatewayResources{}, errs
	}
	return gatewayResources, errs
}
