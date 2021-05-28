package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	cloudrun "github.com/stack72/pulumi-globalgcpcloudrun/sdk/go/globalgcpcloudrun"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		c := config.New(ctx, "")
		project := c.Require("project")

		deployment, err := cloudrun.NewDeployment(ctx, "demo-deployment-go", &cloudrun.DeploymentArgs{
			ImageName:   pulumi.String("gcr.io/ahmetb-public/zoneprinter"),
			ServiceName: "demo-service-ts",
			ProjectId:   project,
		})
		if err != nil {
			return err
		}

		ctx.Export("ip", deployment.IpAddress)

		return nil
	})
}
