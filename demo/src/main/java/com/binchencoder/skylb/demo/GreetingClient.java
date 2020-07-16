package com.binchencoder.skylb.demo;

import java.net.URI;
import java.util.Calendar;
import java.util.concurrent.TimeUnit;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.binchencoder.skylb.MonitorService;
import com.binchencoder.skylb.SkyLBConst;
import com.binchencoder.skylb.SkyLBNameResolverFactory;
import com.binchencoder.skylb.balancer.consistenthash.ConsistentHashLoadBalancerFactory;
import com.binchencoder.skylb.balancer.roundrobin.RoundRobinLoadBalancerFactory;
import com.binchencoder.skylb.demo.proto.DemoGrpc;
import com.binchencoder.skylb.demo.proto.GreetingProtos.GreetingRequest;
import com.binchencoder.skylb.demo.proto.GreetingProtos.GreetingResponse;
import com.binchencoder.skylb.grpc.Channels;
import com.binchencoder.skylb.grpc.ClientTemplate;
import com.binchencoder.skylb.metrics.MetricsClientInterceptor;

import io.grpc.LoadBalancer;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.Status;
import io.grpc.StatusRuntimeException;

/**
 * Demonstrates various methods to init grpc channel.
 *
 * Recommend using ClientTemplate.createChannel(skylbUri, ...
 */
enum DemoGrpcInitMethod {
  // Use raw grpc.
  RAW,

  //Use ClientTemplate.createChannel(target,...
  WRAP1,

  // Use ClientTemplate.createChannel(skylbUri, ...
  WRAP2
}

/**
 * Demonstrates usage of load balancers.
 */
enum DemoLBType {
  RoundRobin,
  ConsistentHash
}

/**
 * Demonstrates usage of skylb uri types.
 */
enum DemoAddrType {
  SkyLB,
  Direct,
  SkyLBAndDirect
}

public class GreetingClient {
  private static final Logger logger = LoggerFactory.getLogger(GreetingClient.class);

  private ManagedChannel channel;
  private Channels channels;
  private URI uri;

  DemoGrpcInitMethod demoInitMethod = DemoGrpcInitMethod.WRAP2;
  DemoLBType demoLbType = null;

  public GreetingClient(String[] args) {
    demoLbType = DemoLBType.valueOf(System.getProperty("demo-lb-type", DemoLBType.RoundRobin.toString()));
    logger.info("InitMethod:{} lbType:{}", demoInitMethod, demoLbType);

    // Choose load balancer (optional)
    LoadBalancer.Factory lb = null;
    String loadBalancerDesc = SkyLBConst.JG_ROUND_ROBIN;
    switch (demoLbType) {
      case RoundRobin:
        loadBalancerDesc = SkyLBConst.JG_ROUND_ROBIN;
        lb = RoundRobinLoadBalancerFactory.getInstance();
        break;
      case ConsistentHash:
        loadBalancerDesc = SkyLBConst.JG_CONSISTEN_HASH;
        lb = ConsistentHashLoadBalancerFactory.getInstance();
        break;
      // (Here and following switch-cases omit "default" branch for simplicity)
    }

    String callerServiceName = Configuration.CLIENT_SERVICE_NAME;
    String calleeServiceName = Configuration.SERVER_SERVICE_NAME;

    // Choose how to init grpc.
    switch (demoInitMethod) {

      case RAW: {
        String skylbUri = Configuration.getSkylbUri(args);
        String target = skylbUri + "/" + Configuration.SERVER_SERVICE_NAME + "?portName="
            + Configuration.PORT_NAME;
        // skylb://skylb-server1:port1,skylb-server2:port2,.../serviceName?portName=myPort

        this.channel = ManagedChannelBuilder
            // 1. forTarget
            .forTarget(target)
            // 2. use binchencoder resolver
            .nameResolverFactory(SkyLBNameResolverFactory.getInstance(callerServiceName))
            .defaultLoadBalancingPolicy(loadBalancerDesc)
            .intercept(MetricsClientInterceptor.create(
                com.binchencoder.skylb.metrics.Configuration.allMetrics(),
                calleeServiceName, callerServiceName))
            // Set idleTimeout to several seconds to reproduce channel idle behavior.
            //.idleTimeout(31, TimeUnit.DAYS)
            .usePlaintext()
            .build();
      }
      break;

      case WRAP1: {
        DemoAddrType demoAddrType = DemoAddrType.SkyLB;
        String target = null;
        switch (demoAddrType) {
          case SkyLB:
            target = Configuration.getSkylbUri(args) + "/" + Configuration.SERVER_SERVICE_NAME
                + "?portName=" + Configuration.PORT_NAME;
            // skylb://skylb-server1:port1,skylb-server2:port2,.../serviceName?portName=myPort
            break;
          case Direct:
            target = "direct://127.0.0.1:" + Configuration.SERVICE_PORT;
            // direct://127.0.0.1:50001
            break;
          case SkyLBAndDirect:
            logger.error("Not supported");
            System.exit(2);
            break;
        }

        this.channels = ClientTemplate.createChannel(target, callerServiceName, calleeServiceName, lb);
        // You can also omit lb, to use default RoundRobinLoadBalancerFactory, i.e:
        if (false) {
          this.channels = ClientTemplate.createChannel(target, callerServiceName, calleeServiceName);
          // (This line of code is to ensure the method compiles)
        }
      }
      break;

      case WRAP2: {
        DemoAddrType demoAddrType = DemoAddrType.SkyLB;
        String skylbUri = null;
        switch (demoAddrType) {
          case SkyLB:
            skylbUri = Configuration.getSkylbUri(args);
            // skylb://192.168.38.6:1900
            break;
          case Direct:
            skylbUri = "direct://" + calleeServiceName
                + ":127.0.0.1:" + Configuration.SERVICE_PORT;
            // direct://shared-test-server-service:127.0.0.1:50001
            break;
          case SkyLBAndDirect:
            skylbUri = Configuration.getSkylbUri(args) + ";"
                + "direct://" + calleeServiceName
                + ":127.0.0.1:" + Configuration.SERVICE_PORT;
            // skylb://192.168.38.6:1900;direct://shared-test-server-service:127.0.0.1:50001
            break;
        }
        this.channels = ClientTemplate.createChannel(skylbUri,
            calleeServiceName, null, null,
            callerServiceName, lb);
        // Parameter lb can be omitted, whose default value is RoundRobinLoadBalancerFactory.

        // (This code block is to guarantee the alternative createChannel method doesn't get broken
        // by accident)
        if (false) {
          this.channels = ClientTemplate.createChannel(skylbUri,
              calleeServiceName, null, null,
              callerServiceName);
        }
      }
      break;
    }
  }

  public void testService() {
    DemoGrpc.DemoBlockingStub blockingStub = null;
    switch (demoInitMethod) {
      case RAW:
        blockingStub = DemoGrpc.newBlockingStub(this.channel);
        break;
      case WRAP1:
      case WRAP2:
        blockingStub = DemoGrpc.newBlockingStub(this.channels.getChannel());
        break;
    }

    for (int i = 0; i < Configuration.TEST_COUNT; i++) {
      try {
        GreetingRequest request = GreetingRequest.newBuilder()
            .setName("GC " + Calendar.getInstance().get(Calendar.SECOND)).build();
        logger.info("Hello request: {}", request.getName());
        GreetingResponse response = null;
        switch (demoLbType) {
          case RoundRobin:
            response = blockingStub
                .withDeadlineAfter(Configuration.DEADLINE_FOR_TEST, TimeUnit.MILLISECONDS)
                .greeting(request);
            break;
          case ConsistentHash:
            response = blockingStub
                .withOption(ConsistentHashLoadBalancerFactory.HASHKEY, "valueOfYourKey")
                .withDeadlineAfter(Configuration.DEADLINE_FOR_TEST, TimeUnit.MILLISECONDS)
                .greeting(request);
            break;
        }
        logger.info("Hello response: {}", response.getGreeting());
      } catch (StatusRuntimeException e) {
        if (e.getStatus() == Status.DEADLINE_EXCEEDED) {
          logger.warn("Hello exceeded deadline.");
        }
      } catch (RuntimeException e) {
        logger.error("Hello Error", e);

        i--;
      } catch (Throwable e) {
        logger.error("Err", e);
      } finally {
        try {
          Thread.sleep(3000);
          // To test idle timeout, sleep 3s or longer, combined with
          // ManagedChannelBuilder...idleTimeout(2, TimeUnit.SECONDS).
        } catch (InterruptedException e) {
          logger.warn("sleep was interrupted.");
        }
      }
    }

    switch (demoInitMethod) {
      case RAW:
        this.channel.shutdown();
        break;
      case WRAP1:
      case WRAP2:
        this.channels.getOriginChannel().shutdown();
        break;
    }
    logger.info("Greeting client shutdown.");
  }

  public static void main(String[] args) {
    // Start prometheus (optional).
    MonitorService.getInstance().startPrometheus(
        Configuration.METRICS_IP,
        // Make metrics port different from GreetingServer so as to void port conflict.
        Configuration.METRICS_PORT + 1,
        Configuration.METRICS_PATH);

    GreetingClient client = new GreetingClient(args);

    // Wait a while for resolving to complete.
    try {
      Thread.sleep(1000);
    } catch (Exception e) {
      e.printStackTrace();
    }

    client.testService();
  }
}
