# cicd-helper

Small little cicd helper that can forward harbor webhooks to restart deployments on [`kthcloud`](https://cloud.cbh.kth.se/) by wrapping the [`kthcloud` api](https://api.cloud.cbh.kth.se/deploy/v2/docs/index.html).

## Usage (of existing deployment)

### Automatically restart all your images that use the image that has been updated on a harbor project

Tutorial
1. Head over to your project on [the kthcloud harbor registry](https://registry.cloud.cbh.kth.se/)
2. Click `Webhooks`
3. Click `NEW WEBHOOK`
4. Choose a name for the webhook
5. Deselect all `Event Type`s except `Artifact pushed`
6. Add `https://cicd.app.cloud.cbh.kth.se/harbor/restart` as the `Endpoint URL`
7. Add your [kthcloud api token](https://docs.cloud.cbh.kth.se/usage/api/#accessing-the-api) as `Auth Header`
8. It should look like this:
   
   ![image](https://github.com/user-attachments/assets/ff9e3c05-748e-46bb-aa35-a5245561174e)

### Manually specifying image

The following command will restart the deployment with the `<deployment-id-here>`
```bash
curl -X POST  https://cicd.app.cloud.cbh.kth.se/forward?deploymentid=<deployment-id-here> \
  -H "Content-Type: application/json" -H "Authorization: <kthcloud-api-token-here>"

```

## Usage (setting up own deployment)

### Build the dockerfile

#### Prerequisites for building docker image

- Docker
- Docker buildx extension

#### Command

You can build the Docker image using the following command

```bash
docker buildx build --file Dockerfile --tag cicd-helper:latest .
```
***Note:** must be ran in the root directory of this repo*

### Build a binary

#### Prerequisites for building a binary

- go
- make

#### Command

You can build it to a binary using the following command

```bash
make
```
***Note:** must be ran in the root directory of this repo*

### Configuration

| Name    | Description                        | Default                      |
|---------|------------------------------------|------------------------------|
| PORT    | Set the port to use                | 8080                         |
| API_URL | Set the base url of the api to use | https://api.cloud.cbh.kth.se |
