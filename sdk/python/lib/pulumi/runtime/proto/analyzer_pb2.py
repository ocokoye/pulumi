# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: pulumi/analyzer.proto
"""Generated protocol buffer code."""
from google.protobuf.internal import builder as _builder
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from . import plugin_pb2 as pulumi_dot_plugin__pb2
from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2
from google.protobuf import struct_pb2 as google_dot_protobuf_dot_struct__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x15pulumi/analyzer.proto\x12\tpulumirpc\x1a\x13pulumi/plugin.proto\x1a\x1bgoogle/protobuf/empty.proto\x1a\x1cgoogle/protobuf/struct.proto\"\xd2\x01\n\x0e\x41nalyzeRequest\x12\x0c\n\x04type\x18\x01 \x01(\t\x12+\n\nproperties\x18\x02 \x01(\x0b\x32\x17.google.protobuf.Struct\x12\x0b\n\x03urn\x18\x03 \x01(\t\x12\x0c\n\x04name\x18\x04 \x01(\t\x12\x33\n\x07options\x18\x05 \x01(\x0b\x32\".pulumirpc.AnalyzerResourceOptions\x12\x35\n\x08provider\x18\x06 \x01(\x0b\x32#.pulumirpc.AnalyzerProviderResource\"\xb5\x03\n\x10\x41nalyzerResource\x12\x0c\n\x04type\x18\x01 \x01(\t\x12+\n\nproperties\x18\x02 \x01(\x0b\x32\x17.google.protobuf.Struct\x12\x0b\n\x03urn\x18\x03 \x01(\t\x12\x0c\n\x04name\x18\x04 \x01(\t\x12\x33\n\x07options\x18\x05 \x01(\x0b\x32\".pulumirpc.AnalyzerResourceOptions\x12\x35\n\x08provider\x18\x06 \x01(\x0b\x32#.pulumirpc.AnalyzerProviderResource\x12\x0e\n\x06parent\x18\x07 \x01(\t\x12\x14\n\x0c\x64\x65pendencies\x18\x08 \x03(\t\x12S\n\x14propertyDependencies\x18\t \x03(\x0b\x32\x35.pulumirpc.AnalyzerResource.PropertyDependenciesEntry\x1a\x64\n\x19PropertyDependenciesEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\x36\n\x05value\x18\x02 \x01(\x0b\x32\'.pulumirpc.AnalyzerPropertyDependencies:\x02\x38\x01\"\xc1\x02\n\x17\x41nalyzerResourceOptions\x12\x0f\n\x07protect\x18\x01 \x01(\x08\x12\x15\n\rignoreChanges\x18\x02 \x03(\t\x12\x1b\n\x13\x64\x65leteBeforeReplace\x18\x03 \x01(\x08\x12\"\n\x1a\x64\x65leteBeforeReplaceDefined\x18\x04 \x01(\x08\x12\x1f\n\x17\x61\x64\x64itionalSecretOutputs\x18\x05 \x03(\t\x12\x0f\n\x07\x61liases\x18\x06 \x03(\t\x12I\n\x0e\x63ustomTimeouts\x18\x07 \x01(\x0b\x32\x31.pulumirpc.AnalyzerResourceOptions.CustomTimeouts\x1a@\n\x0e\x43ustomTimeouts\x12\x0e\n\x06\x63reate\x18\x01 \x01(\x01\x12\x0e\n\x06update\x18\x02 \x01(\x01\x12\x0e\n\x06\x64\x65lete\x18\x03 \x01(\x01\"p\n\x18\x41nalyzerProviderResource\x12\x0c\n\x04type\x18\x01 \x01(\t\x12+\n\nproperties\x18\x02 \x01(\x0b\x32\x17.google.protobuf.Struct\x12\x0b\n\x03urn\x18\x03 \x01(\t\x12\x0c\n\x04name\x18\x04 \x01(\t\",\n\x1c\x41nalyzerPropertyDependencies\x12\x0c\n\x04urns\x18\x01 \x03(\t\"E\n\x13\x41nalyzeStackRequest\x12.\n\tresources\x18\x01 \x03(\x0b\x32\x1b.pulumirpc.AnalyzerResource\"D\n\x0f\x41nalyzeResponse\x12\x31\n\x0b\x64iagnostics\x18\x02 \x03(\x0b\x32\x1c.pulumirpc.AnalyzeDiagnostic\"\xd2\x01\n\x11\x41nalyzeDiagnostic\x12\x12\n\npolicyName\x18\x01 \x01(\t\x12\x16\n\x0epolicyPackName\x18\x02 \x01(\t\x12\x19\n\x11policyPackVersion\x18\x03 \x01(\t\x12\x13\n\x0b\x64\x65scription\x18\x04 \x01(\t\x12\x0f\n\x07message\x18\x05 \x01(\t\x12\x0c\n\x04tags\x18\x06 \x03(\t\x12\x35\n\x10\x65nforcementLevel\x18\x07 \x01(\x0e\x32\x1b.pulumirpc.EnforcementLevel\x12\x0b\n\x03urn\x18\x08 \x01(\t\"\x95\x02\n\x0c\x41nalyzerInfo\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x13\n\x0b\x64isplayName\x18\x02 \x01(\t\x12\'\n\x08policies\x18\x03 \x03(\x0b\x32\x15.pulumirpc.PolicyInfo\x12\x0f\n\x07version\x18\x04 \x01(\t\x12\x16\n\x0esupportsConfig\x18\x05 \x01(\x08\x12\x41\n\rinitialConfig\x18\x06 \x03(\x0b\x32*.pulumirpc.AnalyzerInfo.InitialConfigEntry\x1aM\n\x12InitialConfigEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12&\n\x05value\x18\x02 \x01(\x0b\x32\x17.pulumirpc.PolicyConfig:\x02\x38\x01\"\xc1\x01\n\nPolicyInfo\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x13\n\x0b\x64isplayName\x18\x02 \x01(\t\x12\x13\n\x0b\x64\x65scription\x18\x03 \x01(\t\x12\x0f\n\x07message\x18\x04 \x01(\t\x12\x35\n\x10\x65nforcementLevel\x18\x05 \x01(\x0e\x32\x1b.pulumirpc.EnforcementLevel\x12\x33\n\x0c\x63onfigSchema\x18\x06 \x01(\x0b\x32\x1d.pulumirpc.PolicyConfigSchema\"S\n\x12PolicyConfigSchema\x12+\n\nproperties\x18\x01 \x01(\x0b\x32\x17.google.protobuf.Struct\x12\x10\n\x08required\x18\x02 \x03(\t\"r\n\x0cPolicyConfig\x12\x35\n\x10\x65nforcementLevel\x18\x01 \x01(\x0e\x32\x1b.pulumirpc.EnforcementLevel\x12+\n\nproperties\x18\x02 \x01(\x0b\x32\x17.google.protobuf.Struct\"\xb5\x01\n\x18\x43onfigureAnalyzerRequest\x12K\n\x0cpolicyConfig\x18\x01 \x03(\x0b\x32\x35.pulumirpc.ConfigureAnalyzerRequest.PolicyConfigEntry\x1aL\n\x11PolicyConfigEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12&\n\x05value\x18\x02 \x01(\x0b\x32\x17.pulumirpc.PolicyConfig:\x02\x38\x01*=\n\x10\x45nforcementLevel\x12\x0c\n\x08\x41\x44VISORY\x10\x00\x12\r\n\tMANDATORY\x10\x01\x12\x0c\n\x08\x44ISABLED\x10\x02\x32\xf0\x02\n\x08\x41nalyzer\x12\x42\n\x07\x41nalyze\x12\x19.pulumirpc.AnalyzeRequest\x1a\x1a.pulumirpc.AnalyzeResponse\"\x00\x12L\n\x0c\x41nalyzeStack\x12\x1e.pulumirpc.AnalyzeStackRequest\x1a\x1a.pulumirpc.AnalyzeResponse\"\x00\x12\x44\n\x0fGetAnalyzerInfo\x12\x16.google.protobuf.Empty\x1a\x17.pulumirpc.AnalyzerInfo\"\x00\x12@\n\rGetPluginInfo\x12\x16.google.protobuf.Empty\x1a\x15.pulumirpc.PluginInfo\"\x00\x12J\n\tConfigure\x12#.pulumirpc.ConfigureAnalyzerRequest\x1a\x16.google.protobuf.Empty\"\x00\x42\x34Z2github.com/pulumi/pulumi/sdk/v3/proto/go;pulumirpcb\x06proto3')

_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, globals())
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'pulumi.analyzer_pb2', globals())
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z2github.com/pulumi/pulumi/sdk/v3/proto/go;pulumirpc'
  _ANALYZERRESOURCE_PROPERTYDEPENDENCIESENTRY._options = None
  _ANALYZERRESOURCE_PROPERTYDEPENDENCIESENTRY._serialized_options = b'8\001'
  _ANALYZERINFO_INITIALCONFIGENTRY._options = None
  _ANALYZERINFO_INITIALCONFIGENTRY._serialized_options = b'8\001'
  _CONFIGUREANALYZERREQUEST_POLICYCONFIGENTRY._options = None
  _CONFIGUREANALYZERREQUEST_POLICYCONFIGENTRY._serialized_options = b'8\001'
  _ENFORCEMENTLEVEL._serialized_start=2468
  _ENFORCEMENTLEVEL._serialized_end=2529
  _ANALYZEREQUEST._serialized_start=117
  _ANALYZEREQUEST._serialized_end=327
  _ANALYZERRESOURCE._serialized_start=330
  _ANALYZERRESOURCE._serialized_end=767
  _ANALYZERRESOURCE_PROPERTYDEPENDENCIESENTRY._serialized_start=667
  _ANALYZERRESOURCE_PROPERTYDEPENDENCIESENTRY._serialized_end=767
  _ANALYZERRESOURCEOPTIONS._serialized_start=770
  _ANALYZERRESOURCEOPTIONS._serialized_end=1091
  _ANALYZERRESOURCEOPTIONS_CUSTOMTIMEOUTS._serialized_start=1027
  _ANALYZERRESOURCEOPTIONS_CUSTOMTIMEOUTS._serialized_end=1091
  _ANALYZERPROVIDERRESOURCE._serialized_start=1093
  _ANALYZERPROVIDERRESOURCE._serialized_end=1205
  _ANALYZERPROPERTYDEPENDENCIES._serialized_start=1207
  _ANALYZERPROPERTYDEPENDENCIES._serialized_end=1251
  _ANALYZESTACKREQUEST._serialized_start=1253
  _ANALYZESTACKREQUEST._serialized_end=1322
  _ANALYZERESPONSE._serialized_start=1324
  _ANALYZERESPONSE._serialized_end=1392
  _ANALYZEDIAGNOSTIC._serialized_start=1395
  _ANALYZEDIAGNOSTIC._serialized_end=1605
  _ANALYZERINFO._serialized_start=1608
  _ANALYZERINFO._serialized_end=1885
  _ANALYZERINFO_INITIALCONFIGENTRY._serialized_start=1808
  _ANALYZERINFO_INITIALCONFIGENTRY._serialized_end=1885
  _POLICYINFO._serialized_start=1888
  _POLICYINFO._serialized_end=2081
  _POLICYCONFIGSCHEMA._serialized_start=2083
  _POLICYCONFIGSCHEMA._serialized_end=2166
  _POLICYCONFIG._serialized_start=2168
  _POLICYCONFIG._serialized_end=2282
  _CONFIGUREANALYZERREQUEST._serialized_start=2285
  _CONFIGUREANALYZERREQUEST._serialized_end=2466
  _CONFIGUREANALYZERREQUEST_POLICYCONFIGENTRY._serialized_start=2390
  _CONFIGUREANALYZERREQUEST_POLICYCONFIGENTRY._serialized_end=2466
  _ANALYZER._serialized_start=2532
  _ANALYZER._serialized_end=2900
# @@protoc_insertion_point(module_scope)
