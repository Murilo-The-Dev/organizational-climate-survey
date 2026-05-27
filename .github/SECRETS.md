# Required GitHub Secrets

Configure these in GitHub > Settings > Secrets and variables > Actions.

## CD Pipeline (production deploy)
| Secret | Description |
|---|---|
| REGISTRY_URL | Container registry URL (example: ghcr.io/org) |
| REGISTRY_USERNAME | Registry login username |
| REGISTRY_PASSWORD | Registry login password or token |
| DEPLOY_HOST | Production server IP or hostname |
| DEPLOY_USER | SSH user on production server |
| DEPLOY_SSH_KEY | Private SSH key for production server |
| DEPLOY_PATH | Absolute path to project on server (example: /opt/clima) |