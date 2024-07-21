# cicd-helper

Small little cicd helper that can forward harbor webhooks to restart deployments on `kthcloud` by wrapping the kthcloud api.

## usage
The following command will restart the deployment with the `<deployment-id-here>`
```bash
curl -X POST  https://cicd.app.cloud.cbh.kth.se/forward?deploymentid=<deployment-id-here> \
  -H "Content-Type: application/json" -H "Authorization: <kthcloud-api-token-here>"

```
