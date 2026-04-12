//
//  Generated code. Do not modify.
//  source: proto/ProxyCoreService.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

class StartCoreRequest extends $pb.GeneratedMessage {
  factory StartCoreRequest({
    $core.String? coreName,
    $core.String? dir,
    $core.String? config,
    $core.int? memory,
    $core.bool? isString,
    $core.bool? isVpnMode,
    $core.int? tunFD,
    $core.int? proxyPort,
  }) {
    final $result = create();
    if (coreName != null) {
      $result.coreName = coreName;
    }
    if (dir != null) {
      $result.dir = dir;
    }
    if (config != null) {
      $result.config = config;
    }
    if (memory != null) {
      $result.memory = memory;
    }
    if (isString != null) {
      $result.isString = isString;
    }
    if (isVpnMode != null) {
      $result.isVpnMode = isVpnMode;
    }
    if (tunFD != null) {
      $result.tunFD = tunFD;
    }
    if (proxyPort != null) {
      $result.proxyPort = proxyPort;
    }
    return $result;
  }
  StartCoreRequest._() : super();
  factory StartCoreRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory StartCoreRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'StartCoreRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'ProxyCore'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'coreName', protoName: 'coreName')
    ..aOS(2, _omitFieldNames ? '' : 'dir')
    ..aOS(3, _omitFieldNames ? '' : 'config')
    ..a<$core.int>(4, _omitFieldNames ? '' : 'memory', $pb.PbFieldType.O3)
    ..aOB(5, _omitFieldNames ? '' : 'isString', protoName: 'isString')
    ..aOB(6, _omitFieldNames ? '' : 'isVpnMode', protoName: 'isVpnMode')
    ..a<$core.int>(7, _omitFieldNames ? '' : 'tunFD', $pb.PbFieldType.OU3, protoName: 'tunFD')
    ..a<$core.int>(8, _omitFieldNames ? '' : 'proxyPort', $pb.PbFieldType.O3, protoName: 'proxyPort')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  StartCoreRequest clone() => StartCoreRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  StartCoreRequest copyWith(void Function(StartCoreRequest) updates) => super.copyWith((message) => updates(message as StartCoreRequest)) as StartCoreRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static StartCoreRequest create() => StartCoreRequest._();
  StartCoreRequest createEmptyInstance() => create();
  static $pb.PbList<StartCoreRequest> createRepeated() => $pb.PbList<StartCoreRequest>();
  @$core.pragma('dart2js:noInline')
  static StartCoreRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<StartCoreRequest>(create);
  static StartCoreRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get coreName => $_getSZ(0);
  @$pb.TagNumber(1)
  set coreName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasCoreName() => $_has(0);
  @$pb.TagNumber(1)
  void clearCoreName() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get dir => $_getSZ(1);
  @$pb.TagNumber(2)
  set dir($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasDir() => $_has(1);
  @$pb.TagNumber(2)
  void clearDir() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get config => $_getSZ(2);
  @$pb.TagNumber(3)
  set config($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasConfig() => $_has(2);
  @$pb.TagNumber(3)
  void clearConfig() => clearField(3);

  @$pb.TagNumber(4)
  $core.int get memory => $_getIZ(3);
  @$pb.TagNumber(4)
  set memory($core.int v) { $_setSignedInt32(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasMemory() => $_has(3);
  @$pb.TagNumber(4)
  void clearMemory() => clearField(4);

  @$pb.TagNumber(5)
  $core.bool get isString => $_getBF(4);
  @$pb.TagNumber(5)
  set isString($core.bool v) { $_setBool(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasIsString() => $_has(4);
  @$pb.TagNumber(5)
  void clearIsString() => clearField(5);

  @$pb.TagNumber(6)
  $core.bool get isVpnMode => $_getBF(5);
  @$pb.TagNumber(6)
  set isVpnMode($core.bool v) { $_setBool(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasIsVpnMode() => $_has(5);
  @$pb.TagNumber(6)
  void clearIsVpnMode() => clearField(6);

  @$pb.TagNumber(7)
  $core.int get tunFD => $_getIZ(6);
  @$pb.TagNumber(7)
  set tunFD($core.int v) { $_setUnsignedInt32(6, v); }
  @$pb.TagNumber(7)
  $core.bool hasTunFD() => $_has(6);
  @$pb.TagNumber(7)
  void clearTunFD() => clearField(7);

  @$pb.TagNumber(8)
  $core.int get proxyPort => $_getIZ(7);
  @$pb.TagNumber(8)
  set proxyPort($core.int v) { $_setSignedInt32(7, v); }
  @$pb.TagNumber(8)
  $core.bool hasProxyPort() => $_has(7);
  @$pb.TagNumber(8)
  void clearProxyPort() => clearField(8);
}

class MeasurePingRequest extends $pb.GeneratedMessage {
  factory MeasurePingRequest({
    $core.Iterable<$core.String>? url,
  }) {
    final $result = create();
    if (url != null) {
      $result.url.addAll(url);
    }
    return $result;
  }
  MeasurePingRequest._() : super();
  factory MeasurePingRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory MeasurePingRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'MeasurePingRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'ProxyCore'), createEmptyInstance: create)
    ..pPS(1, _omitFieldNames ? '' : 'url')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  MeasurePingRequest clone() => MeasurePingRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  MeasurePingRequest copyWith(void Function(MeasurePingRequest) updates) => super.copyWith((message) => updates(message as MeasurePingRequest)) as MeasurePingRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MeasurePingRequest create() => MeasurePingRequest._();
  MeasurePingRequest createEmptyInstance() => create();
  static $pb.PbList<MeasurePingRequest> createRepeated() => $pb.PbList<MeasurePingRequest>();
  @$core.pragma('dart2js:noInline')
  static MeasurePingRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<MeasurePingRequest>(create);
  static MeasurePingRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<$core.String> get url => $_getList(0);
}

class BooleanResponse extends $pb.GeneratedMessage {
  factory BooleanResponse({
    $core.bool? message,
  }) {
    final $result = create();
    if (message != null) {
      $result.message = message;
    }
    return $result;
  }
  BooleanResponse._() : super();
  factory BooleanResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory BooleanResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'BooleanResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'ProxyCore'), createEmptyInstance: create)
    ..aOB(1, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  BooleanResponse clone() => BooleanResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  BooleanResponse copyWith(void Function(BooleanResponse) updates) => super.copyWith((message) => updates(message as BooleanResponse)) as BooleanResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BooleanResponse create() => BooleanResponse._();
  BooleanResponse createEmptyInstance() => create();
  static $pb.PbList<BooleanResponse> createRepeated() => $pb.PbList<BooleanResponse>();
  @$core.pragma('dart2js:noInline')
  static BooleanResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<BooleanResponse>(create);
  static BooleanResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.bool get message => $_getBF(0);
  @$pb.TagNumber(1)
  set message($core.bool v) { $_setBool(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(1)
  void clearMessage() => clearField(1);
}

class VersionResponse extends $pb.GeneratedMessage {
  factory VersionResponse({
    $core.String? message,
  }) {
    final $result = create();
    if (message != null) {
      $result.message = message;
    }
    return $result;
  }
  VersionResponse._() : super();
  factory VersionResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory VersionResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'VersionResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'ProxyCore'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  VersionResponse clone() => VersionResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  VersionResponse copyWith(void Function(VersionResponse) updates) => super.copyWith((message) => updates(message as VersionResponse)) as VersionResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static VersionResponse create() => VersionResponse._();
  VersionResponse createEmptyInstance() => create();
  static $pb.PbList<VersionResponse> createRepeated() => $pb.PbList<VersionResponse>();
  @$core.pragma('dart2js:noInline')
  static VersionResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<VersionResponse>(create);
  static VersionResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(1)
  set message($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(1)
  void clearMessage() => clearField(1);
}

class LogResponse extends $pb.GeneratedMessage {
  factory LogResponse({
    $core.String? logs,
  }) {
    final $result = create();
    if (logs != null) {
      $result.logs = logs;
    }
    return $result;
  }
  LogResponse._() : super();
  factory LogResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory LogResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'LogResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'ProxyCore'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'logs')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  LogResponse clone() => LogResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  LogResponse copyWith(void Function(LogResponse) updates) => super.copyWith((message) => updates(message as LogResponse)) as LogResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LogResponse create() => LogResponse._();
  LogResponse createEmptyInstance() => create();
  static $pb.PbList<LogResponse> createRepeated() => $pb.PbList<LogResponse>();
  @$core.pragma('dart2js:noInline')
  static LogResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<LogResponse>(create);
  static LogResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get logs => $_getSZ(0);
  @$pb.TagNumber(1)
  set logs($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasLogs() => $_has(0);
  @$pb.TagNumber(1)
  void clearLogs() => clearField(1);
}

class MeasurePingResponse extends $pb.GeneratedMessage {
  factory MeasurePingResponse({
    $core.Iterable<PingResult>? results,
  }) {
    final $result = create();
    if (results != null) {
      $result.results.addAll(results);
    }
    return $result;
  }
  MeasurePingResponse._() : super();
  factory MeasurePingResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory MeasurePingResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'MeasurePingResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'ProxyCore'), createEmptyInstance: create)
    ..pc<PingResult>(1, _omitFieldNames ? '' : 'results', $pb.PbFieldType.PM, subBuilder: PingResult.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  MeasurePingResponse clone() => MeasurePingResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  MeasurePingResponse copyWith(void Function(MeasurePingResponse) updates) => super.copyWith((message) => updates(message as MeasurePingResponse)) as MeasurePingResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MeasurePingResponse create() => MeasurePingResponse._();
  MeasurePingResponse createEmptyInstance() => create();
  static $pb.PbList<MeasurePingResponse> createRepeated() => $pb.PbList<MeasurePingResponse>();
  @$core.pragma('dart2js:noInline')
  static MeasurePingResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<MeasurePingResponse>(create);
  static MeasurePingResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<PingResult> get results => $_getList(0);
}

class PingResult extends $pb.GeneratedMessage {
  factory PingResult({
    $core.String? url,
    $fixnum.Int64? delay,
  }) {
    final $result = create();
    if (url != null) {
      $result.url = url;
    }
    if (delay != null) {
      $result.delay = delay;
    }
    return $result;
  }
  PingResult._() : super();
  factory PingResult.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory PingResult.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'PingResult', package: const $pb.PackageName(_omitMessageNames ? '' : 'ProxyCore'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'url')
    ..aInt64(2, _omitFieldNames ? '' : 'delay')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  PingResult clone() => PingResult()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  PingResult copyWith(void Function(PingResult) updates) => super.copyWith((message) => updates(message as PingResult)) as PingResult;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PingResult create() => PingResult._();
  PingResult createEmptyInstance() => create();
  static $pb.PbList<PingResult> createRepeated() => $pb.PbList<PingResult>();
  @$core.pragma('dart2js:noInline')
  static PingResult getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<PingResult>(create);
  static PingResult? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get url => $_getSZ(0);
  @$pb.TagNumber(1)
  set url($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasUrl() => $_has(0);
  @$pb.TagNumber(1)
  void clearUrl() => clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get delay => $_getI64(1);
  @$pb.TagNumber(2)
  set delay($fixnum.Int64 v) { $_setInt64(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasDelay() => $_has(1);
  @$pb.TagNumber(2)
  void clearDelay() => clearField(2);
}

class TrafficStatsResponse extends $pb.GeneratedMessage {
  factory TrafficStatsResponse({
    $fixnum.Int64? uplinkTotal,
    $fixnum.Int64? downlinkTotal,
  }) {
    final $result = create();
    if (uplinkTotal != null) {
      $result.uplinkTotal = uplinkTotal;
    }
    if (downlinkTotal != null) {
      $result.downlinkTotal = downlinkTotal;
    }
    return $result;
  }
  TrafficStatsResponse._() : super();
  factory TrafficStatsResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory TrafficStatsResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'TrafficStatsResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'ProxyCore'), createEmptyInstance: create)
    ..aInt64(1, _omitFieldNames ? '' : 'uplinkTotal', protoName: 'uplinkTotal')
    ..aInt64(2, _omitFieldNames ? '' : 'downlinkTotal', protoName: 'downlinkTotal')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  TrafficStatsResponse clone() => TrafficStatsResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  TrafficStatsResponse copyWith(void Function(TrafficStatsResponse) updates) => super.copyWith((message) => updates(message as TrafficStatsResponse)) as TrafficStatsResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TrafficStatsResponse create() => TrafficStatsResponse._();
  TrafficStatsResponse createEmptyInstance() => create();
  static $pb.PbList<TrafficStatsResponse> createRepeated() => $pb.PbList<TrafficStatsResponse>();
  @$core.pragma('dart2js:noInline')
  static TrafficStatsResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<TrafficStatsResponse>(create);
  static TrafficStatsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $fixnum.Int64 get uplinkTotal => $_getI64(0);
  @$pb.TagNumber(1)
  set uplinkTotal($fixnum.Int64 v) { $_setInt64(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasUplinkTotal() => $_has(0);
  @$pb.TagNumber(1)
  void clearUplinkTotal() => clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get downlinkTotal => $_getI64(1);
  @$pb.TagNumber(2)
  set downlinkTotal($fixnum.Int64 v) { $_setInt64(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasDownlinkTotal() => $_has(1);
  @$pb.TagNumber(2)
  void clearDownlinkTotal() => clearField(2);
}

class Empty extends $pb.GeneratedMessage {
  factory Empty() => create();
  Empty._() : super();
  factory Empty.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory Empty.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'Empty', package: const $pb.PackageName(_omitMessageNames ? '' : 'ProxyCore'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  Empty clone() => Empty()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  Empty copyWith(void Function(Empty) updates) => super.copyWith((message) => updates(message as Empty)) as Empty;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Empty create() => Empty._();
  Empty createEmptyInstance() => create();
  static $pb.PbList<Empty> createRepeated() => $pb.PbList<Empty>();
  @$core.pragma('dart2js:noInline')
  static Empty getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Empty>(create);
  static Empty? _defaultInstance;
}


const _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
