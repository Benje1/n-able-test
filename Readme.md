This is a basic service health monitoring programme

It is built so that it can have all that it needs to set itself up. When using docker the user should mount the services.yaml file (if this was pruduction I would not have had it there as it does not make sense to have a default service)

I has a basic untested github CI pipline that should run all the tests and if they succed, and the branch is main, build an image and in a docker container

I did not know how to add AWS services in so I just gave a basic curl call lambda.

Assumption I made:
Each server does not need to change too often, but as it does not have a database or any need for non live data restarting it with a new yaml file should not be a problem.


```
To get the go mods:
go mod downlaod

To run as a non dockerd server:
go run main.go

To run the tests:
go test .\ServiceMonitor\
```