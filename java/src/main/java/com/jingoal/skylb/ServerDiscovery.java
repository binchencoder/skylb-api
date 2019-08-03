package com.jingoal.skylb;

import com.jingoal.skylb.metrics.Configuration;
import com.jingoal.skylb.metrics.MetricsClientInterceptor;
import com.jingoal.skylb.proto.ClientProtos.ResolveRequest;
import com.jingoal.skylb.proto.ClientProtos.ResolveResponse;
import com.jingoal.skylb.proto.ClientProtos.ServiceSpec;
import com.jingoal.skylb.proto.SkylbGrpc;

import java.net.InetSocketAddress;
import java.util.Iterator;
import java.util.List;
import java.util.concurrent.CopyOnWriteArrayList;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.internal.DnsNameResolverProvider;

/**
 * ServerDiscovery resolves a set of ServiceSpec.
 */
public class ServerDiscovery {
  private List<ServiceSpec> serviceSpecs = new CopyOnWriteArrayList<ServiceSpec>();

  private final InetSocketAddress addr;
  private final ManagedChannel channel;

  public ServerDiscovery(final InetSocketAddress addr) {
    this.addr = addr;
    this.channel = ManagedChannelBuilder.forAddress(this.addr.getHostName(), this.addr.getPort())
        .nameResolverFactory(DnsNameResolverProvider.asFactory())
        .usePlaintext(true)
        .intercept(MetricsClientInterceptor.create(Configuration.allMetrics(),
            SkyLBConst.SKYLB_SERVER, SkyLBConst.SKYLB_DISCOVERY))
        // Add .idleTimeout(2, TimeUnit.SECONDS) to verify idleTimeout doesn't harm.
        .build();
  }

  /**
   * shutdown releases grpc resource.
   */
  public void shutdown() {
    if (this.channel != null && !this.channel.isShutdown()) {
      this.channel.shutdownNow();
    }
  }

  /**
   * A loop to wait service list change from skylb server and refresh the given
   * listener.
   */
  public void refreshLoop(String callerServiceName, ServiceSpec serviceSpec, ServerListener listener) {
    if (!serviceSpecs.contains(serviceSpec)) {
      serviceSpecs.add(serviceSpec);
    }

    Iterator<ResolveResponse> iter = SkylbGrpc.newBlockingStub(this.channel)
        .resolve(ResolveRequest.newBuilder()
            // .setCallerServiceId(csId) // (fuyc): will not pass id now.
            .setCallerServiceName(callerServiceName)
            .addAllServices(serviceSpecs)
            .setResolveFullEndpoints(true)
            .build());

    while (iter.hasNext()) {
      listener.onChange(iter.next().getSvcEndpoints());
    }
  }

  public InetSocketAddress getAddress() {
    return this.addr;
  }
}
