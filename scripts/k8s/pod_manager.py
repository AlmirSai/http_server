#!/usr/bin/env python3

from typing import Dict, List, Optional
import logging

from kubernetes import client, stream
from .base_manager import BaseK8sManager


class PodManager(BaseK8sManager):
    def list_pods(
        self,
        namespace: str = 'default',
        show_logs: bool = False) -> List[Dict[str, str]]:
        try:
            pods = self.v1.list_namespaced_pod(namespace)
            pod_list = []
            for pod in pods.items:
                pod_info = {
                    'name': pod.metadata.name,
                    'namespace': pod.metadata.namespace,
                    'status': pod.status.phase,
                    'ip': pod.status.pod_ip,
                    'node': pod.spec.node_name,
                    'age': self._get_age(pod.metadata.creation_timestamp),
                    'ready': self._get_pod_ready_status(pod)
                }
                if show_logs:
                    try:
                        logs = self.v1.read_namespaced_pod_log(
                            pod.metadata.name, namespace, tail_lines=50)
                        pod_info['recent_logs'] = logs
                    except Exception as e:
                        pod_info['recent_logs'] = f"Error fetching logs: {e}"
                pod_list.append(pod_info)
            return pod_list
        except Exception as e:
            logging.error(f"Error listing pods: {e}")
            return []

    def _get_pod_ready_status(self, pod: client.V1Pod) -> str:
        if not pod.status.container_statuses:
            return 'Unknown'
        ready_containers = sum(
            1 for container in pod.status.container_statuses if container.ready)
        total_containers = len(pod.status.container_statuses)
        return f"{ready_containers}/{total_containers}"

    def get_pod_logs(
        self,
        name: str,
        namespace: str = 'default',
        tail_lines: int = 100) -> str:
        try:
            return self.v1.read_namespaced_pod_log(
                name=name,
                namespace=namespace,
                tail_lines=tail_lines
            )
        except Exception as e:
            logging.error(f"Error getting pod logs: {e}")
            return f"Error: {str(e)}"

    def exec_command(
        self,
        name: str,
        namespace: str,
        command: List[str],
        container: Optional[str] = None) -> str:
        try:
            exec_command = [
                '/bin/sh',
                '-c',
                ' '.join(command)
            ]
            resp = stream.stream(
                self.v1.connect_get_namespaced_pod_exec,
                name,
                namespace,
                command=exec_command,
                container=container,
                stderr=True,
                stdin=False,
                stdout=True,
                tty=False
            )
            return resp
        except Exception as e:
            logging.error(f"Error executing command in pod: {e}")
            return f"Error: {str(e)}"

    def delete_pod(self, name: str, namespace: str = 'default') -> bool:
        try:
            self.v1.delete_namespaced_pod(
                name=name,
                namespace=namespace,
                body=client.V1DeleteOptions()
            )
            logging.info(f"Deleted pod {name} in namespace {namespace}")
            return True
        except Exception as e:
            logging.error(f"Error deleting pod: {e}")
            return False
