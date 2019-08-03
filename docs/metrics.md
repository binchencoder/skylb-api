Here summarizes skylb api metrics, for both goland and java.

Items marked with \* means slight difference in naming.

 

Client side metrics:

| Help                                                                                         | Name (golang)                           | Name (java)                                     |
|----------------------------------------------------------------------------------------------|-----------------------------------------|-------------------------------------------------|
| Total number of RPCs completed by the client, regardless of success or failure.              | skylb\_client\_handled\_total           | skylb\_client\_completed (\*)                   |
| Histogram of response latency (seconds) of the gRPC until it is finished by the application. | skylb\_client\_handling\_seconds        | skylb\_client\_completed\_latency\_seconds (\*) |
| Total number of RPC stream messages received by the client.                                  | skylb\_client\_msg\_received\_total     | skylb\_client\_msg\_received\_total             |
| Total number of gRPC stream messages sent by the client.                                     | skylb\_client\_msg\_sent\_total         | skylb\_client\_msg\_sent\_total                 |
| Number of service endpoints.                                                                 | skylb\_client\_service\_endpoint\_gauge | skylb\_client\_service\_endpoint\_gauge         |
| Total number of RPCs started on the client.                                                  | skylb\_client\_started\_total           | skylb\_client\_started\_total                   |

 

Server side metrics:

| Help                                                                                                   | Name (golang)                       | Name (java)                                   |
|--------------------------------------------------------------------------------------------------------|-------------------------------------|-----------------------------------------------|
| Total number of RPCs completed on the server, regardless of success or failure.                        | skylb\_server\_handled\_total       | skylb\_server\_handled\_total                 |
| Histogram of response latency (seconds) of gRPC that had been application-level handled by the server. | skylb\_server\_handling\_seconds    | skylb\_server\_handled\_latency\_seconds (\*) |
| Total number of RPC stream messages received on the server.                                            | skylb\_server\_msg\_received\_total | skylb\_server\_msg\_received\_total           |
| Total number of gRPC stream messages sent by the server.                                               | skylb\_server\_msg\_sent\_total     | skylb\_server\_msg\_sent\_total               |
| Total number of RPCs started on the server.                                                            | skylb\_server\_started\_total       | skylb\_server\_started\_total                 |
