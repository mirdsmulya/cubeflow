# Cubeflow ChatOps Application

Cubeflow is a ChatOps application written in Golang, designed to streamline self service infrastructure operations on Kubernetes. With Cubeflow, you can perform various tasks such as restarting services, rolling back deployments, backing up database, scheduled scaling node & services, and syncing ArgoCD applications, all through Slack commands/CronJob/Internal API calls. For now this app features are still limited to GCP & Slack API environment.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Release Note](#release-note)

## Features

1. **Service Restart:** Restart pods of a specific service in a namespace.
2. **Rollback:** Rollback a service to a previous version.
3. **Backup Database:** Create backups of a specified database.
4. **Sync ArgoCD App:** Synchronize ArgoCD applications.
5. **Scaling Service:** Scale the number of replicas for a service.
6. **Node Scaling:** Scale the number of nodes in a node group.
7. **Get Pods:** Retrieve information about pods in a namespace.
8. **Ping:** Check if the application is running.

## Prerequisites

1. Kubernetes Cluster
2. Slack Workspace
3. Golang installed
4. ArgoCD installed (optional, required for syncing applications)

## Installation

1. Clone the Cubeflow repository:

    ```bash
    git clone https://github.com/mirdsmulya/cubeflow.git
    cd cubeflow
    ```

2. Fill Cubeflow configuration on config/config.yaml:
   
3. Connect to local/remote kubernetes cluster until you able to perform kubectl command.

4. Run on local:
    ```bash
    go run cmd/main.go
    ```

5. Make your local env reachable from public IP using [ngrok](https://ngrok.com/)

6. Setup local Cubeflow API endpoint on Slack API Application using slash command feature.

7. All set!


### API Endpoints

#### Internal API Routes
Internal API is only accessible using local/direct connection

1. **Service Restart:** `/v1/service/restart/:namespace/:rolloutName`
2. **Service Scaling:** `/v1/service/scale/:replicaCount`
3. **Get Pods:** `/v1/service/getpods/:namespace`
4. **Node Scaling:** `/v1/cluster/scale/:nodeGroupName/:nodeCount`
5. **Database Backup:** `/v1/db/backup/:dbname`

#### External API Routes
External API is only accessible using Public connection via Slack

1. **Service Restart:** `/v2/service/restart/:namespace`
2. **Promote Rollout:** `/v2/service/promote/:namespace`
3. **Rollback Service:** `/v2/service/rollback/:namespace`
4. **Sync Argo App:** `/v2/service/sync`
5. **Database Backup:** `/v2/db/backup`
6. **Ping:** `/v2/ping`



# CubeFlow Release Notes

## Version 0.6.0
- **New Endpoint Path Update and Handler Refactoring**
  - Added a new endpoint path update.
  - Refactored handlers for improved performance.
  - Docker Hub Tag: `mirdsmulya/cubeflow:0.6.0`

## Version 0.5.6
- **Multiple Slack Channel Handlers with YAML-Based Env Config**
  - Introduced support for multiple Slack channel handlers.
  - Implemented YAML-based environment configuration.
  - Docker Hub Tag: `mirdsmulya/cubeflow:0.5.6`

## Version 0.5.5
- **Rollback Service (Image Tag)**
  - Added rollback service functionality using image tags.
  - Docker Hub Tag: `mirdsmulya/cubeflow:0.5.5`

## Version 0.5.4
- **Multiple Service Restart**
  - Enabled the ability to restart multiple services.
  - Docker Hub Tag: `mirdsmulya/cubeflow:0.5.4`

## Version 0.5.3
- **Restart with CLBO Handler**
  - Implemented service restart with CLBO handler.
  - Docker Hub Tag: `mirdsmulya/cubeflow:0.5.3`

## Version 0.5.2
- **DB Backup with Node Scaler**
  - Added database backup functionality with node scaler.
  - Docker Hub Tag: `mirdsmulya/cubeflow:0.5.2`

## Version 0.5.1
- **DB Backup with DBName Command**
  - Introduced database backup with a specific dbname command.
  - Docker Hub Tag: `mirdsmulya/cubeflow:0.5.1`

## Version 0.5.0
- **DB Backup (Internal Webhook and Slack), Refactor ISPodsReady, Add CLBO Handler for Promote**
  - Implemented database backup with internal webhook and Slack integration.
  - Refactored ISPodsReady for better functionality.
  - Added CLBO handler when performing promote operations.
  - Docker Hub Tag: `mirdsmulya/cubeflow:0.5.0`

## Version 0.4.10
- **Restart, Promote, Sync**
  - Included features for restarting, promoting, and syncing operations.
  - Docker Hub Tag: `mirdsmulya/cubeflow:0.4.10`
