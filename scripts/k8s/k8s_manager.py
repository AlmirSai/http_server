#!/usr/bin/env python3

import argparse
import logging
from typing import Optional

from kubernetes import client, config
from rich.console import Console

from .pod_manager import PodManager
from .deployment_manager import DeploymentManager
from .service_manager import ServiceManager
from .config_manager import ConfigManager


class K8sManager:
    def __init__(self):
        try:
            config.load_kube_config()
            self.v1 = client.CoreV1Api()
            self.apps_v1 = client.AppsV1Api()
            self.batch_v1 = client.BatchV1Api()
            self.custom_objects = client.CustomObjectsApi()
            self.console = Console()

            self.pod_manager = PodManager(
                self.v1, self.apps_v1, self.batch_v1, self.custom_objects)
            self.deployment_manager = DeploymentManager(
                self.v1, self.apps_v1, self.batch_v1, self.custom_objects)
            self.service_manager = ServiceManager(
                self.v1, self.apps_v1, self.batch_v1, self.custom_objects)
            self.config_manager = ConfigManager(
                self.v1, self.apps_v1, self.batch_v1, self.custom_objects)
        except Exception as e:
            logging.error(f"Error initializing K8s client: {e}")
            raise

    def setup_logging(self, verbose: bool = False):
        level = logging.DEBUG if verbose else logging.INFO
        logging.basicConfig(
            level=level,
            format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
        )


def main():
    parser = argparse.ArgumentParser(description='Kubernetes Resource Manager')
    parser.add_argument('-v', '--verbose', action='store_true',
                        help='Enable verbose logging')
    parser.add_argument('-n', '--namespace', default='default',
                        help='Kubernetes namespace')

    subparsers = parser.add_subparsers(dest='resource', help='Resource type')

    # Pod commands
    pod_parser = subparsers.add_parser('pod', help='Pod operations')
    pod_subparsers = pod_parser.add_subparsers(dest='action')

    list_pods = pod_subparsers.add_parser('list', help='List pods')
    list_pods.add_argument('--show-logs', action='store_true', help='Show pod logs')

    get_logs = pod_subparsers.add_parser('logs', help='Get pod logs')
    get_logs.add_argument('name', help='Pod name')
    get_logs.add_argument('--tail', type=int, default=100,
                          help='Number of lines to show')

    exec_cmd = pod_subparsers.add_parser('exec', help='Execute command in pod')
    exec_cmd.add_argument('name', help='Pod name')
    exec_cmd.add_argument('command', nargs='+', help='Command to execute')
    exec_cmd.add_argument('--container', help='Container name')

    delete_pod = pod_subparsers.add_parser('delete', help='Delete pod')
    delete_pod.add_argument('name', help='Pod name')

    # Deployment commands
    deploy_parser = subparsers.add_parser('deployment', help='Deployment operations')
    deploy_subparsers = deploy_parser.add_subparsers(dest='action')

    deploy_subparsers.add_parser('list', help='List deployments')

    scale = deploy_subparsers.add_parser('scale', help='Scale deployment')
    scale.add_argument('name', help='Deployment name')
    scale.add_argument('replicas', type=int, help='Number of replicas')

    rollback = deploy_subparsers.add_parser('rollback', help='Rollback deployment')
    rollback.add_argument('name', help='Deployment name')
    rollback.add_argument('--revision', type=int, help='Revision to rollback to')

    update_image = deploy_subparsers.add_parser(
        'update-image', help='Update deployment image')
    update_image.add_argument('name', help='Deployment name')
    update_image.add_argument('container', help='Container name')
    update_image.add_argument('image', help='New image')

    # Service commands
    svc_parser = subparsers.add_parser('service', help='Service operations')
    svc_subparsers = svc_parser.add_subparsers(dest='action')

    svc_subparsers.add_parser('list', help='List services')

    create_svc = svc_subparsers.add_parser('create', help='Create service')
    create_svc.add_argument('name', help='Service name')
    create_svc.add_argument('--port', type=int, required=True, help='Service port')
    create_svc.add_argument('--target-port', type=int,
                            required=True, help='Target port')
    create_svc.add_argument('--selector', required=True,
                            help='Pod selector (key=value)')
    create_svc.add_argument('--type', default='ClusterIP', help='Service type')

    delete_svc = svc_subparsers.add_parser('delete', help='Delete service')
    delete_svc.add_argument('name', help='Service name')

    # Config commands
    config_parser = subparsers.add_parser(
        'config', help='ConfigMap and Secret operations')
    config_subparsers = config_parser.add_subparsers(dest='action')

    config_subparsers.add_parser('list-configmaps', help='List ConfigMaps')
    config_subparsers.add_parser('list-secrets', help='List Secrets')

    create_cm = config_subparsers.add_parser(
        'create-configmap', help='Create ConfigMap')
    create_cm.add_argument('name', help='ConfigMap name')
    create_cm.add_argument('--from-literal', nargs=2, action='append', metavar=('key', 'value'),
                           help='Add literal data to the ConfigMap')

    create_secret = config_subparsers.add_parser('create-secret', help='Create Secret')
    create_secret.add_argument('name', help='Secret name')
    create_secret.add_argument('--from-literal', nargs=2, action='append', metavar=('key', 'value'),
                               help='Add literal data to the Secret')
    create_secret.add_argument('--type', default='Opaque', help='Secret type')

    args = parser.parse_args()

    try:
        k8s = K8sManager()
        k8s.setup_logging(args.verbose)

        if args.resource == 'pod':
            if args.action == 'list':
                pods = k8s.pod_manager.list_pods(args.namespace, args.show_logs)
                k8s.pod_manager.display_resources(pods, 'Pod')
            elif args.action == 'logs':
                logs = k8s.pod_manager.get_pod_logs(
                    args.name, args.namespace, args.tail)
                print(logs)
            elif args.action == 'exec':
                result = k8s.pod_manager.exec_command(
                    args.name, args.namespace, args.command, args.container)
                print(result)
            elif args.action == 'delete':
                success = k8s.pod_manager.delete_pod(args.name, args.namespace)
                if success:
                    print(f"Successfully deleted pod {args.name}")

        elif args.resource == 'deployment':
            if args.action == 'list':
                deployments = k8s.deployment_manager.list_deployments(args.namespace)
                k8s.deployment_manager.display_resources(deployments, 'Deployment')
            elif args.action == 'scale':
                success = k8s.deployment_manager.scale_deployment(
                    args.name, args.namespace, args.replicas)
                if success:
                    print(f"Successfully scaled deployment {
                          args.name} to {args.replicas} replicas")
            elif args.action == 'rollback':
                success = k8s.deployment_manager.rollback_deployment(
                    args.name, args.namespace, args.revision)
                if success:
                    print(f"Successfully rolled back deployment {args.name}")
            elif args.action == 'update-image':
                success = k8s.deployment_manager.update_deployment_image(
                    args.name, args.namespace, args.container, args.image)
                if success:
                    print(f"Successfully updated image for deployment {args.name}")

        elif args.resource == 'service':
            if args.action == 'list':
                services = k8s.service_manager.list_services(args.namespace)
                k8s.service_manager.display_resources(services, 'Service')
            elif args.action == 'create':
                selector = dict(s.split('=') for s in args.selector.split(','))
                success = k8s.service_manager.create_service(
                    args.name, args.namespace, args.port, args.target_port, selector, args.type)
                if success:
                    print(f"Successfully created service {args.name}")
            elif args.action == 'delete':
                success = k8s.service_manager.delete_service(args.name, args.namespace)
                if success:
                    print(f"Successfully deleted service {args.name}")

        elif args.resource == 'config':
            if args.action == 'list-configmaps':
                configmaps = k8s.config_manager.list_configmaps(args.namespace)
                k8s.config_manager.display_resources(configmaps, 'ConfigMap')
            elif args.action == 'list-secrets':
                secrets = k8s.config_manager.list_secrets(args.namespace)
                k8s.config_manager.display_resources(secrets, 'Secret')
            elif args.action == 'create-configmap':
                data = {k: v for k, v in args.from_literal} if args.from_literal else {}
                success = k8s.config_manager.create_configmap(
                    args.name, args.namespace, data)
                if success:
                    print(f"Successfully created ConfigMap {args.name}")
            elif args.action == 'create-secret':
                data = {k: v for k, v in args.from_literal} if args.from_literal else {}
                success = k8s.config_manager.create_secret(
                    args.name, args.namespace, data, args.type)
                if success:
                    print(f"Successfully created Secret {args.name}")

    except Exception as e:
        logging.error(f"Error: {e}")
        return 1

    return 0


if __name__ == '__main__':
    exit(main())
