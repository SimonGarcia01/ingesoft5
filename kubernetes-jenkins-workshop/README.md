# CICD-DEMO

This project aims to be the basic skeleton to apply continuous integration and continuous delivery.

## Topology

CICD Demo uses some kubernetes primitives to deploy:

- Deployment
- Services
- Ingress ( with TLS )

```bash
     internet
        |
   [ Ingress ]
   --|-----|--
   [ Services ]
   --|-----|--
   [   Pods   ]

```

This project includes:

- Spring Boot java app
- Jenkinsfile integration to run pipelines
- Dockerfile containing the base image to run java apps
- Makefile and docker-compose to make the pipeline steps much simpler
- Kubernetes deployment file demonstrating how to deploy this app in a simple Kubernetes cluster

## Pipeline Setup

Pipelines exist at Travis.

Some pipelines are configured by **GitHub/Projects**. If you have created a repository in one of these, your project will be **automatically** built if it has a Jenkinsfile/Travis/Gitlab/CircleCI.

Other pipelines are configured manually under folders. You can create a project manually with the following steps:

How to run the app:

```make
make
```

## Jenkins Pipeline (Quality + Security)

The Jenkins pipeline runs these stages:

- Build & Test (Maven)
- Static Analysis (SonarQube) with Quality Gate
- Docker Build
- Container Security Scan (Trivy) failing on CRITICAL
- Deploy (only on main)

### SonarQube (Local)

Start SonarQube in the cluster:

```bash
kubectl apply -f k8s/sonarqube.yml
kubectl rollout status deploy/sonarqube
kubectl port-forward svc/sonarqube 9000:9000
```

Open http://localhost:9000 and create a token:

- User menu -> My Account -> Security -> Generate Token

Add the token to Jenkins as a Secret text credential with ID `sonar-token`.

Quality Gate:

- Configure a Quality Gate in SonarQube that fails when Security Hotspots are detected
  (for example, Security Hotspots Reviewed < 100% or Security Hotspots > 0).
- The pipeline uses `sonar.qualitygate.wait=true` to fail the build when the gate fails.

### Trivy

The pipeline downloads Trivy into the workspace and runs:

```bash
trivy image --exit-code 1 --severity CRITICAL --no-progress mi-app:latest
```

### Deploy

Deploy runs only on the `main` branch:

```bash
docker run -d --name mi-app -p 80:80 -e SERVER_PORT=80 mi-app:latest
```

If port 80 is already in use, update the port mapping in the Jenkinsfile.

## Testing

Unit tests and integrations tests are separated using [JUnit Categories][].

[JUnit Categories]: https://maven.apache.org/surefire/maven-surefire-plugin/examples/junit.html

### Unit Tests

```java
mvn test -Dgroups=UnitTest
```

Or using Docker:

```bash
make build
```

### Integration Tests

```java
mvn integration-test -Dgroups=IntegrationTests
```

Or using Docker:

```bash
make integrationTest
```

### System Tests

System tests run with Selenium using docker-compose to run a [Selenium standalone container][] with Chrome.

[Selenium standalone container]: https://github.com/SeleniumHQ/docker-selenium

Using Docker:

- If you are running locally, make sure the `$APP_URL` is populated and points to a valid instance of your application. This variable is populated automatically in Jenkins.

```bash
APP_URL=http://dev-cicd-demo-master.anzcd.internal/ make systemTest
```
