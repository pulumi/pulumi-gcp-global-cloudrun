import * as pulumi from "@pulumi/pulumi";
import * as globalcloudrun from "@stack72/pulumi-globalgcpcloudrun";

const conf = new pulumi.Config()
const project = conf.require("project")

const deployment = new globalcloudrun.Deployment("my-sample-deployment", {
    projectId: project,

    imageName: "gcr.io/ahmetb-public/zoneprinter",
    serviceName: "demo-service-ts"
});

export const ip = deployment.ipAddress;
