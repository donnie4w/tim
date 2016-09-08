/*
 * This auto-generated skeleton file illustrates how to build a server. If you
 * intend to customize it, you should edit a copy with another file name to 
 * avoid overwriting it when running the generator again.
 */
module ITim_server;

import std.stdio;
import thrift.codegen.processor;
import thrift.protocol.binary;
import thrift.server.simple;
import thrift.server.transport.socket;
import thrift.transport.buffered;
import thrift.util.hashset;

import ITim;
import tim_types;


class ITimHandler : ITim {
  this() {
    // Your initialization goes here.
  }

  void timStream(ref const(TimParam) param) {
    // Your implementation goes here.
    writeln("timStream called");
  }

  void timStarttls() {
    // Your implementation goes here.
    writeln("timStarttls called");
  }

  void timLogin(ref const(Tid) tid, string pwd) {
    // Your implementation goes here.
    writeln("timLogin called");
  }

  void timAck(ref const(TimAckBean) ab) {
    // Your implementation goes here.
    writeln("timAck called");
  }

  void timPresence(ref const(TimPBean) pbean) {
    // Your implementation goes here.
    writeln("timPresence called");
  }

  void timMessage(ref const(TimMBean) mbean) {
    // Your implementation goes here.
    writeln("timMessage called");
  }

  void timPing(string threadId) {
    // Your implementation goes here.
    writeln("timPing called");
  }

  void timError(ref const(TimError) e) {
    // Your implementation goes here.
    writeln("timError called");
  }

  void timLogout() {
    // Your implementation goes here.
    writeln("timLogout called");
  }

  void timRegist(ref const(Tid) tid, string auth) {
    // Your implementation goes here.
    writeln("timRegist called");
  }

  void timRoser(ref const(TimRoster) roster) {
    // Your implementation goes here.
    writeln("timRoser called");
  }

  void timMessageList(ref const(TimMBeanList) mbeanList) {
    // Your implementation goes here.
    writeln("timMessageList called");
  }

  void timPresenceList(ref const(TimPBeanList) pbeanList) {
    // Your implementation goes here.
    writeln("timPresenceList called");
  }

  void timMessageIq(ref const(TimMessageIq) timMsgIq, string iqType) {
    // Your implementation goes here.
    writeln("timMessageIq called");
  }

  void timMessageResult(ref const(TimMBean) mbean) {
    // Your implementation goes here.
    writeln("timMessageResult called");
  }

  void timProperty(ref const(TimPropertyBean) tpb) {
    // Your implementation goes here.
    writeln("timProperty called");
  }

  TimRemoteUserBean timRemoteUserAuth(ref const(Tid) tid, string pwd, ref const(TimAuth) auth) {
    // Your implementation goes here.
    writeln("timRemoteUserAuth called");
    return typeof(return).init;
  }

  TimRemoteUserBean timRemoteUserGet(ref const(Tid) tid, ref const(TimAuth) auth) {
    // Your implementation goes here.
    writeln("timRemoteUserGet called");
    return typeof(return).init;
  }

  TimRemoteUserBean timRemoteUserEdit(ref const(Tid) tid, ref const(TimUserBean) ub, ref const(TimAuth) auth) {
    // Your implementation goes here.
    writeln("timRemoteUserEdit called");
    return typeof(return).init;
  }

  TimResponseBean timResponsePresence(ref const(TimPBean) pbean, ref const(TimAuth) auth) {
    // Your implementation goes here.
    writeln("timResponsePresence called");
    return typeof(return).init;
  }

  TimResponseBean timResponseMessage(ref const(TimMBean) mbean, ref const(TimAuth) auth) {
    // Your implementation goes here.
    writeln("timResponseMessage called");
    return typeof(return).init;
  }

  TimMBeanList timResponseMessageIq(ref const(TimMessageIq) timMsgIq, string iqType, ref const(TimAuth) auth) {
    // Your implementation goes here.
    writeln("timResponseMessageIq called");
    return typeof(return).init;
  }

  TimResponseBean timResponsePresenceList(ref const(TimPBeanList) pbeanList, ref const(TimAuth) auth) {
    // Your implementation goes here.
    writeln("timResponsePresenceList called");
    return typeof(return).init;
  }

  TimResponseBean timResponseMessageList(ref const(TimMBeanList) mbeanList, ref const(TimAuth) auth) {
    // Your implementation goes here.
    writeln("timResponseMessageList called");
    return typeof(return).init;
  }

}

void main() {
  auto protocolFactory = new TBinaryProtocolFactory!();
  auto processor = new TServiceProcessor!ITim(new ITimHandler);
  auto serverTransport = new TServerSocket(9090);
  auto transportFactory = new TBufferedTransportFactory;
  auto server = new TSimpleServer(
    processor, serverTransport, transportFactory, protocolFactory);
  server.serve();
}
