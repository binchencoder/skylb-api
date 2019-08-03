package com.jingoal.skylb.demo;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class Configuration {

  static Logger logger = LoggerFactory.getLogger(Configuration.class);

  // Pre-requisite: Use docker-compose to start skylb server locally.
  public static final String SKYLB_URI = "skylb://localhost:1900";
  public static final String PORT_NAME = "grpc";
  public static final int SERVICE_PORT = 50001;

  public static final String METRICS_IP = "0.0.0.0";
  public static final int METRICS_PORT = 1949;
  public static final String METRICS_PATH = "/_/metrics";

  //(fuyc): As han suggested, use string literals instead of enums defined in
  // vexillary-client/proto/data, to void enforced dependency and
  // tight-coupling for the time being.
  public static final String CLIENT_SERVICE_NAME = "shared-test-client-service";
  public static final String SERVER_SERVICE_NAME = "shared-test-server-service";

  public static final int TEST_COUNT = 20000;
  public static final int DEADLINE_FOR_TEST = 700; // 毫秒

  public static String getSkylbUri(String[] args) {
    String skylbUri = SKYLB_URI;
    if (args.length > 0) {
      skylbUri = args[0];
    }
    logger.info("skylb uri: {}", skylbUri);
    return skylbUri;
  }
}