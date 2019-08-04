package com.binchencoder.skylb.demo;

import com.binchencoder.skylb.MonitorService;
import com.binchencoder.skylb.SkyLBServiceReporter;
import com.binchencoder.skylb.demo.proto.DemoGrpc.DemoImplBase;
import com.binchencoder.skylb.demo.proto.GreetingProtos;
import com.binchencoder.skylb.demo.proto.GreetingProtos.GreetingResponse;
import com.binchencoder.skylb.grpc.ServerTemplate;
import com.binchencoder.skylb.grpchealth.JinHealthServiceImpl;
import com.binchencoder.skylb.grpchealth.JinHealthServiceInterceptor;
import com.binchencoder.skylb.metrics.MetricsServerInterceptor;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Random;
import java.util.concurrent.TimeUnit;

import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.ServerInterceptors;
import io.grpc.stub.StreamObserver;

public class GreetingServer {
  private static final Logger logger = LoggerFactory.getLogger(GreetingServer.class);

  private String serviceName;
  private String portName;
  private int port;

  private SkyLBServiceReporter reporter;
  private Server server;

  void start(String skylbUri, String serviceName, String portName, int port) throws Exception {
    this.serviceName = serviceName;
    this.portName = portName;
    this.port = port;

    // Recommend using new way of initializing grpc service.
    boolean newWay = true;
    if (newWay) {
      this.server = ServerTemplate.create(this.port, new GreeterImpl(), serviceName)
          .build()
          .start();
    } else {
      this.server = ServerBuilder.forPort(this.port)
          .addService(ServerInterceptors.intercept(
              // Enable service provider.
              new GreeterImpl(),
              // Enable metrics.
              MetricsServerInterceptor.create(
                  com.binchencoder.skylb.metrics.Configuration.allMetrics(), serviceName)))
          // Enable grpc health checking.
          .addService(ServerInterceptors.intercept(new JinHealthServiceImpl(),
              new JinHealthServiceInterceptor()))
          .build()
          .start();
    }
    logger.info("Target server {}@:{} started.", serviceName, this.port);
    this.reporter = ServerTemplate.reportLoad(skylbUri, serviceName, portName, this.port);

    logger.info("Target server {}@:{}-{} registered into skylb service {}.",
        serviceName, port, portName, skylbUri);

    Runtime.getRuntime().addShutdownHook(new Thread() {
      @Override
      public void run() {
        GreetingServer.this.stop();
        System.err.println("Target server " + serviceName + "@" + port + " shut down");
      }
    });
  }

  void stop() {
    this.reporter.shutdown();
    if (server != null && !server.isShutdown()) {
      this.server.shutdownNow();
      logger.info("GreetingServer {}@:{}-{} stopped!", this.serviceName, this.portName, this.port);
    }
  }

  private class GreeterImpl extends DemoImplBase {
    private Random rand = new Random(System.currentTimeMillis());

    @Override
    public void greeting(GreetingProtos.GreetingRequest request,
                         StreamObserver<GreetingProtos.GreetingResponse> responseObserver) {
      logger.info("Got req:{}", request);

      // 随机耗时350~550毫秒.
      int elapse = 350 + rand.nextInt(200);
      try {
        TimeUnit.MILLISECONDS.sleep(elapse);
      } catch (InterruptedException e) {
        logger.warn("sleep interrupted");
      }
      GreetingResponse reply = GreetingResponse.newBuilder().setGreeting(
          "Hello " + request.getName() + ", from :" + port + ", elapse " + elapse).build();
      responseObserver.onNext(reply);
      responseObserver.onCompleted();
    }
  }

  public static void main(String[] args) throws Exception {
    // (Optional) Start prometheus.
    MonitorService.getInstance().startPrometheus(
        Configuration.METRICS_IP, Configuration.METRICS_PORT, Configuration.METRICS_PATH);

    // Start service.
    String skylbUri = Configuration.getSkylbUri(args);
    GreetingServer greetingServer50001 = new GreetingServer();
    greetingServer50001.start(skylbUri, Configuration.SERVER_SERVICE_NAME,
        Configuration.PORT_NAME, Configuration.SERVICE_PORT);
    Thread.sleep(TimeUnit.MINUTES.toMillis(10));
  }
}
