# CubeFlow ChatOps Platform

CubeFlow is a ChatOps platform that integrates Slack with a Kubernetes environment, enabling various operations through Slack chat commands. This platform provides capabilities for performing operations such as ArgoCD synchronization, pod restarts, and ArgoCD rollout promotions. Additionally, it offers webhook operations for integrating with Grafana and Alertmanager to enable autorestarts based on metrics.

## Features

CubeFlow ChatOps platform provides the following features:

1. **Service Restart:** Restart pods of a specific service in a namespace.
2. **Rollback:** Rollback a service to a previous version.
3. **Backup Database:** Create backups of a specified database.
4. **Sync ArgoCD App:** Synchronize ArgoCD applications.
5. **Scaling Service:** Scale the number of replicas for a service.
6. **Node Scaling:** Scale the number of nodes in a node group.
7. **Get Pods:** Retrieve information about pods in a namespace.
8. **Ping:** Check if the application is running.

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

## Helm Chart Deployment

The deployment of the CubeFlow application is managed using a Helm chart. Below are the snippet values you can customize when deploying CubeFlow:

```yaml

nameOverride: "cubeflow"

config:
  argocdUrl: argocd-server.argocd.svc.cluster.local
  env: staging
  promptHelper:
    rolloutName: []
    argoAppName: []

ingressRoute:
  enabled: false
  match: [] # Host(`cubeflow.stag-innov8.danaventures.id`)
  middlewares:
    enabled: false

istioIngress:
  enabled: false
  hostname: [] # cubeflow.stag-innov8.danaventures.id
  gateway:
    httpPort: 80
    httpsPort: 443
  virtualservice:
    destination: [] # cubeflow.cubeflow.svc.cluster.local

ingress:
  enabled: false
  className: ""
  annotations:
    kubernetes.io/ingress.class: nginx
  hosts:
    - host: [] # cubeflow.stag-innov8.danaventures.id
      paths:
        - path: /v2
          pathType: ImplementationSpecific
  tls:
    - hosts:
      - [] # cubeflow.stag-innov8.danaventures.id
      secretName: cubeflow

acmeCertificate:
  enabled: true
  host: # cubeflow.stag-innov8.danaventures.id
  issuerRef:
    kind: ClusterIssuer
    name: # letsencrypt-staging-nginx
```

You can modify the values in the Helm chart according to your deployment requirements. Make sure to set the appropriate values for your environment and customize the configuration as needed.


## Vault Config Values

### How to Setup

1. Create new role on kubernetes auth method with
   - Bound service account names: **cubeflow-sa**
   - Bound service account namespaces: **{cubeflow namespace}**
   - Generated token type: **default**
   - Generated Token's Initial TTL: **86400**
   - Token policy: use one of your vault policy
2. Create new secret engine on vault web UI on secret path: {env}/data/cubeflow
   - Add vault key named gcp-credentials, and filled the value using GCP Service account secret 
   - Add vault key named config, and paste filled config below:

    ```bash
    # config.yaml

    env:                         # development, staging, production 

    argocd:
      server_url:                # https://argocd-server.argocd.svc.local
      username: 
      password: 
      appname:
        - appname-1
        - appname-2

    slack:
      signing_secret: 
      token: 
      channel:
        - name: channel-project1
          rollouts_name:
            - app1                
            - app2
          argo_app_name:
            - argo app1
            - argo app2
        - name: channel-project2
          rollouts_name:
            - app1
            - app2
          argo_app_name:
            - argo app1
            - argo app2

    db:
      name: 
      host: 
      port: 
      username: 
      password: 

    gcp:
      project_id: 
      credentials_path: 
      bucket_name: 
      gke:
        cluster_name: 
        zone: 
        namespace_to_scale:
          - namespace-1
          - namespace-2
    ```
 

## Getting Started
To get started with CubeFlow, follow these steps:

1. Deploy CubeFlow using the Helm chart with the desired configuration values.
2. Configure `Slack Command` in your Slack Bot App, to enable ChatOps capabilities.
3. Apply kubernetes secret config with the desired configuration values
4. Start using CubeFlow by issuing chat commands in your Slack workspace.

For detailed installation and configuration instructions, please refer to the CubeFlow documentation.
