# WiGo - Wikipedia web crawler
## Requirements
In order to set up the project you need:
* Linux machine / VM / WSL
* [Docker](https://www.docker.com/) set up with [Kubernetes](https://docs.docker.com/desktop/kubernetes/)
* [Kubectl](https://kubernetes.io/docs/tasks/tools/)
## Setup
In order to run the project run the script `setup.sh` and it will automatically deploy everything for you.

## How to use
Once the project has been deployed you can go to localhost:8080 and acess the web app. Enter the name of the wikipedia article you want to crawl and the depth (must be between 1-3). The site will then display to you the found articles with a short summary to each one.
<p align="center">
  <img src="https://github.com/aalyth/vct-wigo/assets/61279622/543212be-29e4-4c96-9f79-e5a5f00fb72d" width="850" height="800"/>
</p>

## Architecture
### Pods
Each pod is a WiGo server instance that is able to process request. The minimal number of WiGo pods running is 3, but when traffic gets bigger the HPA (Horizontal Pod Autoscale) can scale up the number of pods up to 15. Kubernetes periodically performs Health Checks (HC) in order to respawn failed pods.

### Kubernetes (K8s) Load Balancer (LB)
In order to split the traffic between all running pods, the WiGo service is deployed with the `--type=LoadBalancer`. This way Kubernetes automatically creates the LoadBalancer and [performs periodic HC](https://kubernetes.io/docs/tasks/access-application-cluster/create-external-load-balancer/#external-load-balancer-providers) in order to keep it running.

### Security
Because users have access to the `/api/wiki` endpoint, the server has a request limit - any client can only use the endpoint once every 2 seconds.

<p align="center">
  <img src="https://github.com/aalyth/vct-wigo/blob/main/others/architecture.png"/>
</p>

## Why did I choose Go?
The reason I chose Go is because I wanted this project to be as efficient as possible. In order to achieve maximum efficiency while crawling, I made use of the Go goroutines and the blazingly fast Gin web server framework.
