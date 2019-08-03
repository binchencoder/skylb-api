package com.jingoal.skylb.healthcheck;

public interface SizerUser {
  // Called back by lb factory.newLoadBalancer.
  void setSizer(Sizer sizer);
}
