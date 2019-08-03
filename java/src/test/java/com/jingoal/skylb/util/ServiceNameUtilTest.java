package com.jingoal.skylb.util;

import com.jingoal.vexillary.grpc.DataProtos;

import org.junit.Assert;
import org.junit.Test;

public class ServiceNameUtilTest {
  @Test
  public void testToString() throws Exception {
    Assert.assertEquals("avatar-service",
        ServiceNameUtil.toString(DataProtos.ServiceId.AVATAR_SERVICE));
  }

  @Test
  public void toServiceId() throws Exception {
    Assert.assertEquals(DataProtos.ServiceId.AVATAR_SERVICE,
        ServiceNameUtil.toServiceId("avatar-service"));
  }
}