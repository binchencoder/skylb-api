package com.jingoal.skylb.demo;

import java.io.IOException;
import java.util.Random;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.jingoal.skylb.demo.proto.DemoGrpc.DemoImplBase;
import com.jingoal.skylb.demo.proto.GreetingProtos;
import com.jingoal.skylb.demo.proto.GreetingProtos.GreetingResponse;

import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.stub.StreamObserver;

/**
 * Plain grpc server.
 */
public class SingleGreetingServer {
  private static final Logger logger = LoggerFactory.getLogger(SingleGreetingServer.class);
  private Server server;

  public static void main(String[] args) throws IOException {
    SingleGreetingServer singleServer = new SingleGreetingServer();
    logger.info("Starts singleServer");
    singleServer.start();

    Runtime.getRuntime().addShutdownHook(new Thread() {
      @Override
      public void run() {
        singleServer.stop();
      }
    });

    CountDownLatch latch = new CountDownLatch(1);
    try {
      latch.await();
    } catch (InterruptedException e) {
      System.exit(0);
    }
  }

  private void start() throws IOException {
    server = ServerBuilder.forPort(8080)
        .addService(new GreeterImpl())
        .build()
        .start();
  }

  private void stop() {
    server.shutdown();
  }

  private class GreeterImpl extends DemoImplBase {
    private Random rand = new Random(System.currentTimeMillis());

    @Override
    public void greeting(GreetingProtos.GreetingRequest request,
        StreamObserver<GreetingProtos.GreetingResponse> responseObserver) {

      // 随机耗时350~550毫秒.
      int elapse = 350 + rand.nextInt(200);
      try {
        TimeUnit.MILLISECONDS.sleep(elapse);
      } catch (InterruptedException e) {
        logger.warn("sleep interrupted.");
      }
      logger.info("A request reached.");
      GreetingResponse reply = GreetingResponse.newBuilder().setGreeting(
        "Hello " + request.getName() + ", from " + 10000 + ", elapse " + elapse + ", in java.").build();
      responseObserver.onNext(reply);
      responseObserver.onCompleted();
    }
  }
}
