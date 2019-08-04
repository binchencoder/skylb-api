SkyLB API FAQ and best practice
===============================

(When reading this FAQ, it is assumed that you have read SkyLB API user guide
and had basic understanding of the concepts and API usage)

Why the clientServiceId parameter in NewServiceCli?
---------------------------------------------------

Basically, this parameter is used for the gRPC client telling SkyLB server “who
I am”, so SkyLB can track service dependencies.

Based on this info, we will be able to see the whole picture of service
dependency graph, and effectively establish further service governance like
performance monitoring, tuning and throttling.

What is the type of clientServiceId?
------------------------------------

For golang API, clientServiceId is strictly confined to a valid service id, as
defined in vexillary-client/proto/data/data.proto.

For Java API, service name (a string literal) is accepted, yet you should make
sure it’s consistent with its service id counterpart. (e.g. use
“vexillary-client”, to represent VEXILLARY\_CLIENT in data.proto)

What value should I use as clientServiceId?
-------------------------------------------

It’s a question of naming discretion.

Consider a service, whose service id is “SOMEWORK\_SERVICE”. Its “client” can
have two possible types: a pure client, or a client which in turn acts as a
service to serve other various clients.

For a pure client, which devotedly calls “SOMEWORK\_SERVICE” only, it is OK to
name the client as “SOMEWORK\_CLIENT”. (or choose a name which describes this
client’s nature best)

For a client which is also a service, you can name it as
“SOMEOTHERWORK\_SERVICE”, based on what this service does.

Example:

~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
[pc program] --> [message relay server] --> [message storage server]
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

For the three parties in the call hierarchy above, we can use “PC\_CLIENT”,
“MESSAGE\_RELAY\_SERVICE” and “MESSAGE\_STORAGE\_SERVICE” as their
clientServiceId respectively.

Must I explicitly define service id for a pure client?
------------------------------------------------------

Yes. Even if it’s a “MY\_PRIVATEWORK\_CLIENT” used exclusively by yourself, you
should still define it in vexillary-client/proto/data/data.proto and use it.

Do I need to notify anyone if I added a new service id?
-------------------------------------------------------

Yes, if any of the following scenarios applies:

-   You are a service, and some other clients needs to call you on this new
    service id, those clients must be updated too. (”update” means recompile and
    restart)

-   No matter you are a service or a client, if this new service id relies on
    vexillary server to fetch config or feature flag, vexillary scheduler must
    be updated.

Note that SkyLB server doesn’t require any updates.

Can I call NewServiceCli() multiple times?
------------------------------------------

NewServiceCli simply initializes a struct. There’s no point in calling it
multiple times and assigning to the same variable.

But in case you do need multiple ServiceCli instances, you can call
NewServiceCli multiple times: for each call, assign it to a different ServiceCli
variable, and call its Resolve() and Start(), respectively.

Should I create multiple ServiceCli instances?
----------------------------------------------

In reality, unless you have to announce yourself with different “client
identities” to different bundle of services - which is rare - you will only need
to create one ServiceCli instance.

Can I call Resolve() multiple times?
------------------------------------

Yes. Each Resolve() call means a different service to be associated with this
client.

Can I call Start() multiple times?
----------------------------------

For a given ServiceCli instance, Start() can only be called once.

How do I register multiple services in one gRPC server?
-------------------------------------------------------

While Register() and Start() register one single service only, RegisterMulti()
and StartMulti() register multiple services and start them all in once.

Demo code: skylb-api/cmd/demo/servermulti.go

How do I write a program which is both a client and a server?
-------------------------------------------------------------

Demo code: src/github.com/binchencoder/skylb-api/cmd/demo/serverclient.go

It largely combines demo/client.go and demo/server.go together.

What should I know about prometheus metrics?
--------------------------------------------

The SkyLB server API automatically starts a HTTP server which serves both
Prometheus and Golang pprof for debugging purpose. The server reuses the
same port as the gRPC service. For example, if a gRPC service is running
on port 8000, the prometheus metrics can be accessed with:

http://localhost:8000/_/metrics

and the pprof debugging page can be accessed with:

http://localhost:8000/_/debug/pprof/

Why there is -within-k8s CLI parameter?
---------------------------------------

SkyLB server and skylb-api-enabled grpc services have different behaviors
between running inside Kubernetes and outside, so when your are running them
outside k8s (e.g. locally), you must specify **-within-k8s=false** to skylb
server and skylb-api-enabled grpc services.

This parameter is not needed for skylb-api-enabled grpc clients.
