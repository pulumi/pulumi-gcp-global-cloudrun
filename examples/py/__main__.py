"""A Python Pulumi program"""

import pulumi
import pulumi_gcp_global_cloudrun as cloudrun

config = pulumi.Config()
project = config.require("project")

deployment = cloudrun.Deployment("my-sample-deployment",
                                 project_id=project,
                                 image_name="gcr.io/ahmetb-public/zoneprinter",
                                 service_name="demo-service-py")

pulumi.export('ip', deployment.ip_address)
