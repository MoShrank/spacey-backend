# spacey-backend
This repository contains all spacey backend services including the API Gateway which were written in Go.

## Services
### api
The api gateway is responsible for routing all requests from any client side facing app
to the corresponding microservice. It is also responsible for common tasks such as authentication, rate limiting and cors. The api is the only public facing service.

### user-service
The user service is responsible for handling all tasks related to issuing authentication tokens and managing user accounts.

### config-service
The config service provides simple access to configuration values for back and front end applications.

### deck-management-service
The deck management service is responsible for handling all tasks related to managing decks and their corresponding cards. It acts as a simple crud interface.

### learning-service
The learning service is responsible for handling all tasks related to learning and simple statistics such as a score of how well the user remembers the cards in a deck.
Each review of a card is stored as an event which therefore allows to track progress of a user over time.

### card-generation-service
The card generation service is not contained in this repository since it is written in python. It can be found [here](https://github.com/MoShrank/card-generation-service).


## Architecture
![user-flow](./.github/images/user-flow.png)

## Files and Folders

[`/config`](./config) <br>
package for handling configuration values. Each value has a default value that can be overwritten by a .env file.<br>
[`/pkg`](./pkg) <br>
package folder that contains all packages used across the project. <br>
[`/services`](./services) <br>
folder which contains microservices.<br>
[`docker-compose.yml`](./docker-compose.yml)<br>
simple docker compose file to start a local docker environment which also sets up
a docker network and the database.<br>
[`mongo-init.js`](./mongo-init.js)<br>
init script for database setup to insert test user into db

## Run dev environment locally

### Prerequisites
- GO 1.17 needed
- docker and make needed

If make is not installed, the commands can also be executed manually and found inside the [Makefile](./Makefile).

### Environment Variables
The following environment variables should be declared to run the backend locally. There are a few additional config values that can be set via environment variables, which can also be found in the config package. Those are not important for running the backend locally.

```
MONGO_DB_CONNECTION=<mongo_db_uri>
DB_NAME=<name_of_database>
PORT=<port_for_server>
```

### Serving the backend
- `make serve`

### Shutting the backend down
- `make cleanup`

### Running Tests
- `make test`

## API Routes

## Database
MongoDB was chosen because of its maturity, scallability and flexibility to use. Although it is currently deployed on a single VPS instance together with all services, which makes it difficult to scale and less resistant to failures, it can theoretically be deployed on a managed AWS instance and therefore be scaled up using either shards or replicas.
Since we do not have any real users and therefore not a lot of data, we use mongodump to backup the data once a day at night and push the exported bson dumps to a s3 bucket.
A full description of collections and its corresponding indices can be found [here](./docs/Collections.md).

## Security
![threat-model](./.github/images/threat-model.png)
A full list of security measures that we use can be found [here](./docs/Security.md).


## Deployment
The backend is deployed on a virtual private server uisng a [docker compose file](https://github.com/MoShrank/spacey-docker-services). A github workflow is used to run tests, build and push the docker image to a AWS docker registry on each push/merge to master. The VPS runs [watchtower](https://github.com/containrrr/watchtower) which continuously polls the latest docker image and updates it if a newer version is available.
![CI/CD Pipeline](./.github/images/pipeline.png)