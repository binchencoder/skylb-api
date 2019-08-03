package com.jingoal.skylb.util;

import com.jingoal.vexillary.grpc.DataProtos.ServiceId;

public class ServiceNameUtil {

  private static final String UNDERLINE = "_";

  private static final String STRIKETHROUGH = "-";

  public static String toString(ServiceId serviceId) {
    return serviceId.name().toLowerCase().replace(UNDERLINE, STRIKETHROUGH);
  }

  public static ServiceId toServiceId(String serviceName) {
    return ServiceId.valueOf(serviceName.toUpperCase().replace(STRIKETHROUGH, UNDERLINE));
  }
}
