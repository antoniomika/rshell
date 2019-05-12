rshell
=========
A simple reverse shell service.

Preface
-------
This project was based on the project [reverse-shell](https://github.com/lukechilds/reverse-shell) by Luke Childs. It was mainly written in Go to test out [Google Cloud Run](https://cloud.google.com/run/) and attempting to make fully GitOps deployable projects as part of the free offerings on Google Cloud.

So What?
-------
Any push to this repository is automatically deployed to the Cloud Run service. All of the business logic for building and deploying is handled in the `cloudbuild.yaml` file. This makes use of [Google Cloud Build](https://cloud.google.com/cloud-build/) to build the container and pushes it to gcr.io

Now What?
---------
Any use of a service like this should only be done for educational purposes.