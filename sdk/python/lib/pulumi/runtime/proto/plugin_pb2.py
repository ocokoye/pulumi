# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: pulumi/plugin.proto
"""Generated protocol buffer code."""
from google.protobuf.internal import builder as _builder
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x13pulumi/plugin.proto\x12\tpulumirpc\"\x1d\n\nPluginInfo\x12\x0f\n\x07version\x18\x01 \x01(\t\"O\n\x10PluginDependency\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x0c\n\x04kind\x18\x02 \x01(\t\x12\x0f\n\x07version\x18\x03 \x01(\t\x12\x0e\n\x06server\x18\x04 \x01(\t\"\x1f\n\x0cPluginAttach\x12\x0f\n\x07\x61\x64\x64ress\x18\x01 \x01(\tB4Z2github.com/pulumi/pulumi/sdk/v3/proto/go;pulumirpcb\x06proto3')

_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, globals())
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'pulumi.plugin_pb2', globals())
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z2github.com/pulumi/pulumi/sdk/v3/proto/go;pulumirpc'
  _PLUGININFO._serialized_start=34
  _PLUGININFO._serialized_end=63
  _PLUGINDEPENDENCY._serialized_start=65
  _PLUGINDEPENDENCY._serialized_end=144
  _PLUGINATTACH._serialized_start=146
  _PLUGINATTACH._serialized_end=177
# @@protoc_insertion_point(module_scope)
