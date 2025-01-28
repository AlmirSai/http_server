#!/usr/bin/env python3

from typing import Dict, List, Optional
import logging

from kubernetes import client
from .base_manager import BaseK8sManager


class ServiceManager(BaseK8sManager):
    def list_services(self, namespace: str = 'default') -> List[Dict[str, str]]:
        try:
            services = self.v1.list_namespaced_service(namespace)
            return [{
                'name': svc.metadata.name,
                'namespace': svc.metadata.namespace,
                'type': svc.spec.type,
                'cluster_ip': svc.spec.cluster_ip,
                'external_ip': self._get_external_ip(svc),
                'ports': self._get_service_ports(svc),
                'age': self._get_age(svc.metadata.creation_timestamp)
            } for svc in services.items]
        except Exception as e:
            logging.error(f"Error listing services: {e}")
            return []

    def _get_external_ip(self, service: client.V1Service) -> str:
        if service.status.load_balancer.ingress:
            return service.status.load_balancer.ingress[0].ip or 'Pending'
        return 'None'

    def _get_service_ports(self, service: client.V1Service) -> str:
        if not service.spec.ports:
            return 'None'
        ports = []
        for port in service.spec.ports:
            port_str = f"{port.port}"
            if port.node_port:
                port_str += f":{port.node_port}"
            if port.name:
                port_str += f"/{port.name}"
            ports.append(port_str)
        return ', '.join(ports)

    def create_service(
        self,
        name: str,
        namespace: str,
        port: int,
        target_port: int,
        selector: Dict[str, str],
        service_type: str = 'ClusterIP'
    ) -> bool:
        try:
            body = client.V1Service(
                metadata=client.V1ObjectMeta(
                    name=name
                ),
                spec=client.V1ServiceSpec(
                    selector=selector,
                    ports=[client.V1ServicePort(
                        port=port,
                        target_port=target_port
                    )],
                    type=service_type
                )
            )
            self.v1.create_namespaced_service(namespace=namespace, body=body)
            logging.info(f"Created service {name} in namespace {namespace}")
            return True
        except Exception as e:
            logging.error(f"Error creating service: {e}")
            return False

    def delete_service(self, name: str, namespace: str = 'default') -> bool:
        try:
            self.v1.delete_namespaced_service(
                name=name,
                namespace=namespace
            )
            logging.info(f"Deleted service {name} in namespace {namespace}")
            return True
        except Exception as e:
            logging.error(f"Error deleting service: {e}")
            return False

    def update_service(
        self,
        name: str,
        namespace: str,
        port: Optional[int] = None,
        target_port: Optional[int] = None,
        selector: Optional[Dict[str, str]] = None
    ) -> bool:
        try:
            service = self.v1.read_namespaced_service(name, namespace)

            if port is not None:
                service.spec.ports[0].port = port
            if target_port is not None:
                service.spec.ports[0].target_port = target_port
            if selector is not None:
                service.spec.selector = selector

            self.v1.patch_namespaced_service(
                name=name,
                namespace=namespace,
                body=service
            )
            logging.info(f"Updated service {name} in namespace {namespace}")
            return True
        except Exception as e:
            logging.error(f"Error updating service: {e}")
            return False
