#!/usr/bin/env python3

from datetime import datetime
from typing import Dict, List, Optional

from kubernetes import client
from rich.console import Console
from rich.table import Table
from rich import print as rprint


class BaseK8sManager:
    def __init__(
        self,
        v1_client: client.CoreV1Api,
        apps_v1_client: client.AppsV1Api,
        batch_v1_client: client.BatchV1Api,
        custom_objects_client: client.CustomObjectsApi,
    ) -> None:
        self.v1 = v1_client
        self.apps_v1 = apps_v1_client
        self.batch_v1 = batch_v1_client
        self.custom_objects = custom_objects_client
        self.console = Console()

    def _get_age(self, timestamp: Optional[datetime]) -> str:
        if not timestamp:
            return 'Unknown'
        age = datetime.now(timestamp.tzinfo) - timestamp
        if age.days > 0:
            return f"{age.days}d"
        hours = age.seconds // 3600
        if hours > 0:
            return f"{hours}h"
        minutes = (age.seconds % 3600) // 60
        return f"{minutes}m"

    def display_resources(self, resources: List[Dict], resource_type: str) -> None:
        table = Table(title=f"{resource_type} Resources", show_lines=True)
        if resources:
            for key in resources[0].keys():
                if key in ['status', 'ready']:
                    table.add_column(key.capitalize(), style='bold')
                elif key in ['age']:
                    table.add_column(key.capitalize(), style='dim')
                else:
                    table.add_column(key.capitalize())

            for resource in resources:
                row = [str(v) for v in resource.values()]
                if 'status' in resource:
                    status_index = list(resource.keys()).index('status')
                    if row[status_index] == 'Running':
                        row[status_index] = f"[green]{row[status_index]}[/green]"
                    elif row[status_index] in ['Pending', 'ContainerCreating']:
                        row[status_index] = f"[yellow]{row[status_index]}[/yellow]"
                    elif row[status_index] in ['Failed', 'Error']:
                        row[status_index] = f"[red]{row[status_index]}[/red]"

                table.add_row(*row)

            self.console.print(table)
        else:
            rprint(f"[yellow]No {resource_type.lower()} found[/yellow]")
