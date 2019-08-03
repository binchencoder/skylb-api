package com.jingoal.skylb.demo;

import java.net.URI;
import java.util.Calendar;
import java.util.concurrent.TimeUnit;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.jingoal.skylb.MonitorService;
import com.jingoal.skylb.SkyLBConst;
import com.jingoal.skylb.SkyLBNameResolverFactory;
import com.jingoal.skylb.balancer.consistenthash.ConsistentHashLoadBalancerFactory;
import com.jingoal.skylb.demo.proto.DemoGrpc;
import com.jingoal.skylb.demo.proto.GreetingProtos.GreetingRequest;
import com.jingoal.skylb.demo.proto.GreetingProtos.GreetingResponse;
import com.jingoal.skylb.metrics.MetricsClientInterceptor;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.Status;
import io.grpc.StatusRuntimeException;

public class ConsistentHashClientMultiSvc {
  private static final Logger logger = LoggerFactory.getLogger(ConsistentHashClientMultiSvc.class);

  private ManagedChannel channel;
  private URI uri;

  // skylb://skylb-server1:port1,skylb-server2:port2,.../serviceName?portName=myPort
  public ConsistentHashClientMultiSvc(String uri, String callerServiceName) {
    this.uri = URI.create(uri);
    String calleeServiceName = this.uri.getPath();
    calleeServiceName = calleeServiceName.substring(1);  // remove the leading '/'.

    this.channel = ManagedChannelBuilder
        .forTarget(uri)
        .nameResolverFactory(SkyLBNameResolverFactory.getInstance(callerServiceName))
        // Use consistent hash load balancer.
        .defaultLoadBalancingPolicy(SkyLBConst.JG_CONSISTEN_HASH)
        .intercept(MetricsClientInterceptor.create(
            com.jingoal.skylb.metrics.Configuration.allMetrics(),
            calleeServiceName, callerServiceName))
        .usePlaintext()
        .build();
  }

  public void testService(String valueForHash) {
    DemoGrpc.DemoBlockingStub blockingStub = DemoGrpc.newBlockingStub(this.channel);

    for (int i = 0; i < Configuration.TEST_COUNT; i++) {
      try {
        GreetingRequest request = GreetingRequest.newBuilder()
            .setName("GC " + Calendar.getInstance().get(Calendar.SECOND)).build();
        logger.info("Hello request: {}", request.getName());
        GreetingResponse response = blockingStub
            // Provide value for hashing.
            .withOption(ConsistentHashLoadBalancerFactory.HASHKEY, valueForHash)
            .withDeadlineAfter(Configuration.DEADLINE_FOR_TEST, TimeUnit.MILLISECONDS)
            .greeting(request);
        logger.info("Hello response: {}", response.getGreeting());
      } catch (StatusRuntimeException e) {
        if (e.getStatus() == Status.DEADLINE_EXCEEDED) {
          logger.warn("Hello exceeded deadline.");
        }
      } catch (RuntimeException e) {
        String errMsg = e.getMessage() != null ?
            e.getMessage() : (e.getCause() != null) ?
            e.getCause().getMessage() : "";
        logger.error("Hello Error: {}", errMsg);

        i--;
      } catch (Throwable e) {
        logger.error("Err", e);
      } finally {
        try {
          Thread.sleep(1000);
        } catch (InterruptedException e) {
          logger.warn("sleep was interrupted.");
        }
      }
    }

    this.channel.shutdown();
    logger.info("Greeting client shutdown.");
  }

  public static void main(String[] args) {
    // Start prometheus.
    MonitorService.getInstance().startPrometheus(
        Configuration.METRICS_IP,
        // Make metrics port different from GreetingServer so as to void port conflict.
        Configuration.METRICS_PORT + 1,
        Configuration.METRICS_PATH);

    String skylbUri = Configuration.getSkylbUri(args);
    ConsistentHashClientMultiSvc client = new ConsistentHashClientMultiSvc(
        skylbUri + "/" + Configuration.SERVER_SERVICE_NAME + "?portName="
            + Configuration.PORT_NAME,
        Configuration.CLIENT_SERVICE_NAME);

    ConsistentHashClientMultiSvc client2 = new ConsistentHashClientMultiSvc(
        skylbUri + "/" + "vexillary-test-service" + "?portName="
            + Configuration.PORT_NAME,
        Configuration.CLIENT_SERVICE_NAME);

    // Wait a while for resolving to complete.
    try {
      Thread.sleep(1000);
    } catch (Exception e) {
      e.printStackTrace();
    }

    new Thread(new Runnable() {
      @Override
      public void run() {
        try {
          Thread.sleep(100);
        } catch (InterruptedException e) {
          e.printStackTrace();
        }
        client2.testService("anotherValue");
      }
    }).start();

    client.testService("valueOfYourKey");
  }
}
