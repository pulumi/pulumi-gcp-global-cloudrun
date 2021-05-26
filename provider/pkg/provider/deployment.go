// Copyright 2016-2021, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"fmt"
	"github.com/pulumi/pulumi-gcp/sdk/v5/go/gcp/cloudrun"
	"github.com/pulumi/pulumi-gcp/sdk/v5/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type DeploymentArgs struct {
	ImageName   pulumi.StringInput `pulumi:"imageName"`
	ProjectID   string             `pulumi:"projectId"`
	ServiceName string             `pulumi:"serviceName"`
}

// The Deployment component resource.
type Deployment struct {
	pulumi.ResourceState

	IPAddress pulumi.StringOutput `pulumi:"ipAddress"`
}

func NewDeployment(ctx *pulumi.Context,
	name string, args *DeploymentArgs, opts ...pulumi.ResourceOption) (*Deployment, error) {
	if args == nil {
		args = &DeploymentArgs{}
	}

	component := &Deployment{}
	err := ctx.RegisterComponentResource(GCPCloudRunGlobalToken, name, component, opts...)
	if err != nil {
		return nil, err
	}

	// get a list of the current regions
	locations, err := cloudrun.GetLocations(ctx, &cloudrun.GetLocationsArgs{
		Project: &args.ProjectID,
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("error getting CloudRun locations: %v", err)
	}

	var negs []*compute.RegionNetworkEndpointGroup
	for _, location := range locations.Locations {
		service, err := cloudrun.NewService(ctx, fmt.Sprintf("%s-%s-service", name, location), &cloudrun.ServiceArgs{
			Location: pulumi.String(location),
			Project:  pulumi.String(args.ProjectID),
			Template: &cloudrun.ServiceTemplateArgs{
				Spec: &cloudrun.ServiceTemplateSpecArgs{
					Containers: cloudrun.ServiceTemplateSpecContainerArray{
						&cloudrun.ServiceTemplateSpecContainerArgs{
							Image: args.ImageName,
						},
					},
				},
			},
		}, pulumi.Parent(component))
		if err != nil {
			return nil, fmt.Errorf("error creating CloudRun Service in %s: %v", location, err)
		}

		_, err = cloudrun.NewIamMember(ctx, fmt.Sprintf("%s-%s-iam-member", name, location), &cloudrun.IamMemberArgs{
			Location: service.Location,
			Project:  service.Project,
			Service:  service.Name,
			Role:     pulumi.String("roles/run.invoker"),
			Member:   pulumi.String("allUsers"),
		}, pulumi.Parent(service))
		if err != nil {
			return nil, fmt.Errorf("error creating CloudRun IAM Member in %s: %v", location, err)
		}

		neg, err := compute.NewRegionNetworkEndpointGroup(ctx, fmt.Sprintf("%s-%s-neg", name, location), &compute.RegionNetworkEndpointGroupArgs{
			NetworkEndpointType: pulumi.String("SERVERLESS"),
			Region:              service.Location,
			Project:             service.Project,
			CloudRun: &compute.RegionNetworkEndpointGroupCloudRunArgs{
				Service: service.Name,
			},
		}, pulumi.Parent(service))
		if err != nil {
			return nil, fmt.Errorf("error creating Compute Region Network Endpoint Group in %s: %v", location, err)
		}

		negs = append(negs, neg)
	}

	ipv4Address, err := compute.NewGlobalAddress(ctx, fmt.Sprintf("%s-ipv4-global-address", name), &compute.GlobalAddressArgs{
		Project:   pulumi.String(args.ProjectID),
		IpVersion: pulumi.String("IPV4"),
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("error creating IPV4 Global Address: %v", err)
	}

	backendArgs := &compute.BackendServiceArgs{
		Project:        pulumi.String(args.ProjectID),
		EnableCdn:      pulumi.BoolPtr(false),
		SecurityPolicy: pulumi.String(""),
		LogConfig: &compute.BackendServiceLogConfigArgs{
			Enable:     pulumi.BoolPtr(true),
			SampleRate: pulumi.Float64Ptr(1.0),
		},
	}

	var x compute.BackendServiceBackendArray
	for _, neg := range negs {
		x = append(x, &compute.BackendServiceBackendArgs{
			Group: neg.ID(),
		})
	}
	backendArgs.Backends = x
	backendServices, err := compute.NewBackendService(ctx, fmt.Sprintf("%s-backend-service", name), backendArgs)
	if err != nil {
		return nil, fmt.Errorf("error creating Backend Service: %v", err)
	}

	computeUrlMap, err := compute.NewURLMap(ctx, fmt.Sprintf("%s-url-map", name), &compute.URLMapArgs{
		Project:        pulumi.String(args.ProjectID),
		DefaultService: backendServices.ID(),
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("error creating URL Map: %v", err)
	}

	targetHttpProxy, err := compute.NewTargetHttpProxy(ctx, fmt.Sprintf("%s-target-http-proxy", name), &compute.TargetHttpProxyArgs{
		Project: pulumi.String(args.ProjectID),
		UrlMap:  computeUrlMap.ID(),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating Target HTTP Proxy: %v", err)
	}

	_, err = compute.NewGlobalForwardingRule(ctx, fmt.Sprintf("%s-global-forwarding-rule", name), &compute.GlobalForwardingRuleArgs{
		Project:   pulumi.String(args.ProjectID),
		Target:    targetHttpProxy.ID(),
		IpAddress: ipv4Address.Address,
		PortRange: pulumi.String("80"),
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("error creating Global Forwarding Rule: %v", err)
	}

	component.IPAddress = ipv4Address.Address

	if err := ctx.RegisterResourceOutputs(component, pulumi.Map{
		"ipAddress": ipv4Address.Address,
	}); err != nil {
		return nil, err
	}

	return component, nil
}
