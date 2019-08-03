package com.jingoal.skytest;

import com.jingoal.skylb.MonitorService;
import com.jingoal.skylb.grpc.Channels;
import com.jingoal.skylb.grpc.ClientTemplate;
import com.jingoal.skylb.skytest.proto.SkytestGrpc;
import com.jingoal.skylb.skytest.proto.SkytestProtos.GreetingRequest;
import com.jingoal.skylb.skytest.proto.SkytestProtos.GreetingResponse;

import java.util.Calendar;
import java.util.concurrent.TimeUnit;

import io.grpc.Status;
import io.grpc.StatusRuntimeException;

public class SkytestClient {
  private static final String callerServiceName = "shared-test-client-service";
  private static final String calleeServiceName = "shared-test-server-service";

  public static void main(String[] args) {
    String port = System.getProperty("prometheus.port", "15950");
    MonitorService.getInstance().startPrometheus(
        "0.0.0.0", Integer.parseInt(port), "/_/metrics");

    new Thread(new RpcThread("shared-test-server-service")).start();
    new Thread(new RpcThread("idm-service")).start();
    new Thread(new RpcThread("pruinae-avatar-service")).start();
    new Thread(new RpcThread("merged-server")).start();
    new Thread(new RpcThread("reminder-rest-api-service")).start();
    new Thread(new RpcThread("notification-push-service")).start();
    new Thread(new RpcThread("workbench-service")).start();
    new Thread(new RpcThread("module-service")).start();
    new Thread(new RpcThread("web-common-service")).start();
    new Thread(new RpcThread("netdisk-service")).start();

    // channels.getOriginChannel().shutdown();
  }

  static class RpcThread extends Thread {
    private String svcName;

    public RpcThread(String svcName) {
      this.svcName = svcName;
    }

    public void run() {
      String target = "skylb://192.168.221.104:11900,192.168.221.105:11900,192.168.221.106:11900";
      Channels channels = ClientTemplate.createChannel(target, svcName, "grpc", null, callerServiceName);
      SkytestGrpc.SkytestBlockingStub blockingStub = SkytestGrpc.newBlockingStub(channels.getChannel());

      while (true) {
        try {
          GreetingRequest request = GreetingRequest.newBuilder()
              .setName("GC " + Calendar.getInstance().get(Calendar.SECOND)).build();
          GreetingResponse response = blockingStub
              .withDeadlineAfter(100, TimeUnit.MILLISECONDS)
              .greeting(request);
          // System.out.println("=======>      " + svcName + " | " + response);
          // System.out.print(".");
        } catch (StatusRuntimeException e) {
          if (e.getStatus() == Status.DEADLINE_EXCEEDED) {
            System.out.println("Hello exceeded deadline. " + svcName);
            e.printStackTrace(System.out);
          }
        } catch (RuntimeException e) {
          System.out.println("Hello Error" + e);
        } catch (Throwable e) {
          System.out.println("Err, " + e);
        } finally {
          try {
            Thread.sleep(200);
            // To test idle timeout, sleep 3s or longer, combined with
            // ManagedChannelBuilder...idleTimeout(2, TimeUnit.SECONDS).
          } catch (InterruptedException e) {
            System.out.println("sleep was interrupted.");
          }
        }
      }
    }
  }
}
