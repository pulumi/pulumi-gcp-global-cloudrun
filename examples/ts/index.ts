import * as pulumi from "@pulumi/pulumi";
import * as cloudrun from "@pulumi/gcp-global-cloudrun";

const conf = new pulumi.Config()
const project = conf.require("project")

const deployment = new cloudrun.Deployment("my-sample-deployment", {
    projectId: project,

    imageName: "gcr.io/ahmetb-public/zoneprinter",
    serviceName: "demo-service-ts"
});

export const ip = deployment.ipAddress;
