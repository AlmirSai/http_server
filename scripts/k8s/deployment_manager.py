#!/usr/bin/env python3

from datetime import datetime
from typing import Dict, List, Optional
import logging

from kubernetes import client
from .base_manager import BaseK8sManager


class DeploymentManager(BaseK8sManager):
    def list_deployments(self, namespace: str = 'default') -> List[Dict[str, str]]:
        try:
            deployments = self.apps_v1.list_namespaced_deployment(namespace)
            return [{
                'name': dep.metadata.name,
                'namespace': dep.metadata.namespace,
                'replicas': str(dep.spec.replicas),
                'available': str(dep.status.available_replicas or 0),
                'strategy': dep.spec.strategy.type,
                'image': self._get_deployment_images(dep),
                'age': self._get_age(dep.metadata.creation_timestamp)
            } for dep in deployments.items]
        except Exception as e:
            logging.error(f"Error listing deployments: {e}")
            return []

    def _get_deployment_images(self, deployment: client.V1Deployment) -> str:
        if not deployment.spec.template.spec.containers:
            return 'None'
        return ', '.join(
            container.image
            for container in deployment.spec.template.spec.containers
        )

    def scale_deployment(self, name: str, namespace: str, replicas: int) -> bool:
        try:
            self.apps_v1.patch_namespaced_deployment_scale(
                name=name,
                namespace=namespace,
                body={'spec': {'replicas': replicas}}
            )
            logging.info(f"Scaled deployment {name} to {replicas} replicas")
            return True
        except Exception as e:
            logging.error(f"Error scaling deployment: {e}")
            return False

    def rollback_deployment(
        self,
        name: str,
        namespace: str,
        revision: Optional[int] = None) -> bool:
        try:
            if revision:
                deployment_history = self.apps_v1.read_namespaced_deployment_history(
                    name=name, namespace=namespace)
                for revision_history in deployment_history.history:
                    if revision_history.revision == revision:
                        self.apps_v1.patch_namespaced_deployment(
                            name=name,
                            namespace=namespace,
                            body=revision_history.template
                        )
                        break
            else:
                # Rollback to the previous revision
                self.apps_v1.patch_namespaced_deployment(
                    name=name,
                    namespace=namespace,
                    body={'spec': {'template': {'metadata': {'annotations': {
                        'kubectl.kubernetes.io/restartedAt': datetime.now().isoformat()}}}}}
                )
            logging.info(f"Rolled back deployment {name} in namespace {namespace}")
            return True
        except Exception as e:
            logging.error(f"Error rolling back deployment: {e}")
            return False

    def update_deployment_image(
        self,
        name: str,
        namespace: str,
        container_name: str,
        new_image: str) -> bool:
        try:
            deployment = self.apps_v1.read_namespaced_deployment(name, namespace)
            for container in deployment.spec.template.spec.containers:
                if container.name == container_name:
                    container.image = new_image
                    break
            else:
                logging.error(
                    f"Container {container_name} not found in deployment {name}")
                return False

            self.apps_v1.patch_namespaced_deployment(
                name=name,
                namespace=namespace,
                body=deployment
            )
            logging.info(f"Updated image for container {
                         container_name} in deployment {name}")
            return True
        except Exception as e:
            logging.error(f"Error updating deployment image: {e}")
            return False
