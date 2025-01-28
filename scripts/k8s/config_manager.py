#!/usr/bin/env python3

from base64 import b64encode
from typing import Dict, List
import logging

from kubernetes import client
from .base_manager import BaseK8sManager


class ConfigManager(BaseK8sManager):
    def list_configmaps(
        self,
        namespace: str = 'default') -> List[Dict[str, List[str] | str]]:
        try:
            configmaps = self.v1.list_namespaced_config_map(namespace)
            return [{
                'name': cm.metadata.name,
                'namespace': cm.metadata.namespace,
                'data_keys': list(cm.data.keys()) if cm.data else [],
                'age': self._get_age(cm.metadata.creation_timestamp)
            } for cm in configmaps.items]
        except Exception as e:
            logging.error(f"Error listing configmaps: {e}")
            return []

    def list_secrets(
        self,
        namespace: str = 'default') -> List[Dict[str, List[str] | str]]:
        try:
            secrets = self.v1.list_namespaced_secret(namespace)
            return [{
                'name': secret.metadata.name,
                'namespace': secret.metadata.namespace,
                'type': secret.type,
                'data_keys': list(secret.data.keys()) if secret.data else [],
                'age': self._get_age(secret.metadata.creation_timestamp)
            } for secret in secrets.items]
        except Exception as e:
            logging.error(f"Error listing secrets: {e}")
            return []

    def create_configmap(self, name: str, namespace: str, data: Dict[str, str]) -> bool:
        try:
            body = client.V1ConfigMap(
                metadata=client.V1ObjectMeta(name=name),
                data=data
            )
            self.v1.create_namespaced_config_map(
                namespace=namespace, body=body)
            logging.info(f"Created ConfigMap {name} in namespace {namespace}")
            return True
        except Exception as e:
            logging.error(f"Error creating ConfigMap: {e}")
            return False

    def create_secret(
        self,
        name: str,
        namespace: str,
        data: Dict[str, str], secret_type: str = 'Opaque') -> bool:
        try:
            encoded_data = {k: b64encode(v.encode()).decode()
                            for k, v in data.items()}
            body = client.V1Secret(
                metadata=client.V1ObjectMeta(name=name),
                type=secret_type,
                data=encoded_data
            )
            self.v1.create_namespaced_secret(namespace=namespace, body=body)
            logging.info(f"Created Secret {name} in namespace {namespace}")
            return True
        except Exception as e:
            logging.error(f"Error creating Secret: {e}")
            return False

    def update_configmap(self, name: str, namespace: str, data: Dict[str, str]) -> bool:
        try:
            configmap = self.v1.read_namespaced_config_map(name, namespace)
            configmap.data = data
            self.v1.patch_namespaced_config_map(
                name=name,
                namespace=namespace,
                body=configmap
            )
            logging.info(f"Updated ConfigMap {name} in namespace {namespace}")
            return True
        except Exception as e:
            logging.error(f"Error updating ConfigMap: {e}")
            return False

    def update_secret(self, name: str, namespace: str, data: Dict[str, str]) -> bool:
        try:
            secret = self.v1.read_namespaced_secret(name, namespace)
            encoded_data = {k: b64encode(v.encode()).decode()
                            for k, v in data.items()}
            secret.data = encoded_data
            self.v1.patch_namespaced_secret(
                name=name,
                namespace=namespace,
                body=secret
            )
            logging.info(f"Updated Secret {name} in namespace {namespace}")
            return True
        except Exception as e:
            logging.error(f"Error updating Secret: {e}")
            return False

    def delete_configmap(self, name: str, namespace: str) -> bool:
        try:
            self.v1.delete_namespaced_config_map(
                name=name,
                namespace=namespace
            )
            logging.info(f"Deleted ConfigMap {name} in namespace {namespace}")
            return True
        except Exception as e:
            logging.error(f"Error deleting ConfigMap: {e}")
            return False

    def delete_secret(self, name: str, namespace: str) -> bool:
        try:
            self.v1.delete_namespaced_secret(
                name=name,
                namespace=namespace
            )
            logging.info(f"Deleted Secret {name} in namespace {namespace}")
            return True
        except Exception as e:
            logging.error(f"Error deleting Secret: {e}")
            return False
