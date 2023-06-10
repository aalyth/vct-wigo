# WiGo - Wikipedia web crawler
## Requirements
In order to set up the project you need:
* Linux machine / VM / WSL
* [Docker](https://www.docker.com/) set up with [Kubernetes](https://docs.docker.com/desktop/kubernetes/)
* [Kubectl](https://kubernetes.io/docs/tasks/tools/)
## Setup
In order to run the project run the script `setup.sh` and it will automatically deploy everything for you.

## How to use
Once the project has been deployed you can go to localhost:8080 and acess the web app. Enter the name of the wikipedia article you want to crawl and the depth (must be between 1-3). The site will then display to you the found articles with short a short to each one.
<p align="center">
  <img src="https://github.com/aalyth/vct-wigo/assets/61279622/543212be-29e4-4c96-9f79-e5a5f00fb72d" width="850" height="800"/>
</p>
