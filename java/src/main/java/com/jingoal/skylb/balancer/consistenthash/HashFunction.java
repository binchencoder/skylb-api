package com.jingoal.skylb.balancer.consistenthash;

public interface HashFunction {
  public int hash(String key);
}
