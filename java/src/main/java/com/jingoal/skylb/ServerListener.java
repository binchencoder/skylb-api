package com.jingoal.skylb;

import com.jingoal.skylb.proto.ClientProtos.ServiceEndpoints;

public interface ServerListener {
  void onChange(ServiceEndpoints endpoints);

  String serverInfoToString();
}
