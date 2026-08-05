import logging
from typing import Dict, List

import pykube
import pytest
from pytest_helm_charts.clusters import Cluster

logger = logging.getLogger(__name__)

app_name = "policy-meta-operator"

selector = {"app.kubernetes.io/name": app_name}


@pytest.mark.smoke
def test_api_working(kube_cluster: Cluster) -> None:
    assert kube_cluster.kube_client is not None
    assert len(pykube.Node.objects(kube_cluster.kube_client)) >= 1


@pytest.mark.smoke
def test_cluster_info(
    kube_cluster: Cluster, cluster_type: str, test_extra_info: Dict[str, str]
) -> None:
    logger.info(f"Running on cluster type {cluster_type}")
    key = "external_cluster_type"
    if key in test_extra_info:
        logger.info(f"{key} is {test_extra_info[key]}")
    assert kube_cluster.kube_client is not None
    assert cluster_type != ""


@pytest.mark.smoke
def test_deployment_installed(kube_cluster: Cluster) -> None:
    deployments: List[pykube.Deployment] = list(
        pykube.Deployment.objects(
            kube_cluster.kube_client, namespace=pykube.all
        ).filter(selector=selector)
    )
    assert len(deployments) == 1

    deployment = deployments[0].obj
    assert deployment["spec"]["replicas"] == 1

    containers = deployment["spec"]["template"]["spec"]["containers"]
    assert len(containers) == 1
    assert f"giantswarm/{app_name}" in containers[0]["image"]


@pytest.mark.smoke
def test_rbac_installed(kube_cluster: Cluster) -> None:
    service_accounts = list(
        pykube.ServiceAccount.objects(
            kube_cluster.kube_client, namespace=pykube.all
        ).filter(selector=selector)
    )
    assert len(service_accounts) == 1

    cluster_roles = list(
        pykube.ClusterRole.objects(kube_cluster.kube_client).filter(selector=selector)
    )
    assert len(cluster_roles) == 1

    cluster_role_bindings = list(
        pykube.ClusterRoleBinding.objects(kube_cluster.kube_client).filter(
            selector=selector
        )
    )
    assert len(cluster_role_bindings) == 1

    binding = cluster_role_bindings[0].obj
    assert binding["roleRef"]["name"] == cluster_roles[0].name
    assert binding["subjects"][0]["name"] == service_accounts[0].name
    assert binding["subjects"][0]["namespace"] == service_accounts[0].namespace
