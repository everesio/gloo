
---
title: "plugins.proto"
weight: 5
---

<!-- Code generated by solo-kit. DO NOT EDIT. -->


### Package: `gloo.solo.io` 
#### Types:


- [ListenerPlugins](#listenerplugins)
- [VirtualHostPlugins](#virtualhostplugins)
- [RoutePlugins](#routeplugins)
- [DestinationSpec](#destinationspec)
- [UpstreamSpec](#upstreamspec)
  



##### Source File: [github.com/solo-io/gloo/projects/gloo/api/v1/plugins.proto](https://github.com/solo-io/gloo/blob/master/projects/gloo/api/v1/plugins.proto)





---
### ListenerPlugins

 
Plugin-specific configuration that lives on listeners
Each ListenerPlugin object contains configuration for a specific plugin
Note to developers: new Listener Plugins must be added to this struct
to be usable by Gloo.

```yaml
"grpcWeb": .grpc_web.plugins.gloo.solo.io.GrpcWeb

```

| Field | Type | Description | Default |
| ----- | ---- | ----------- |----------- | 
| `grpcWeb` | [.grpc_web.plugins.gloo.solo.io.GrpcWeb](../plugins/grpc_web/grpc_web.proto.sk#grpcweb) |  |  |




---
### VirtualHostPlugins

 
Plugin-specific configuration that lives on virtual hosts
Each VirtualHostPlugin object contains configuration for a specific plugin
Note to developers: new Virtual Host Plugins must be added to this struct
to be usable by Gloo.

```yaml
"extensions": .gloo.solo.io.Extensions

```

| Field | Type | Description | Default |
| ----- | ---- | ----------- |----------- | 
| `extensions` | [.gloo.solo.io.Extensions](../extensions.proto.sk#extensions) |  |  |




---
### RoutePlugins

 
Plugin-specific configuration that lives on routes
Each RoutePlugin object contains configuration for a specific plugin
Note to developers: new Route Plugins must be added to this struct
to be usable by Gloo.

```yaml
"transformations": .transformation.plugins.gloo.solo.io.RouteTransformations
"faults": .fault.plugins.gloo.solo.io.RouteFaults
"prefixRewrite": .transformation.plugins.gloo.solo.io.PrefixRewrite
"timeout": .google.protobuf.Duration
"retries": .retries.plugins.gloo.solo.io.RetryPolicy
"extensions": .gloo.solo.io.Extensions

```

| Field | Type | Description | Default |
| ----- | ---- | ----------- |----------- | 
| `transformations` | [.transformation.plugins.gloo.solo.io.RouteTransformations](../plugins/transformation/transformation.proto.sk#routetransformations) |  |  |
| `faults` | [.fault.plugins.gloo.solo.io.RouteFaults](../plugins/faultinjection/fault.proto.sk#routefaults) |  |  |
| `prefixRewrite` | [.transformation.plugins.gloo.solo.io.PrefixRewrite](../plugins/transformation/prefix_rewrite.proto.sk#prefixrewrite) |  |  |
| `timeout` | [.google.protobuf.Duration](https://developers.google.com/protocol-buffers/docs/reference/csharp/class/google/protobuf/well-known-types/duration) |  |  |
| `retries` | [.retries.plugins.gloo.solo.io.RetryPolicy](../plugins/retries/retries.proto.sk#retrypolicy) |  |  |
| `extensions` | [.gloo.solo.io.Extensions](../extensions.proto.sk#extensions) |  |  |




---
### DestinationSpec

 
Configuration for Destinations that are tied to the UpstreamSpec or ServiceSpec on that destination

```yaml
"aws": .aws.plugins.gloo.solo.io.DestinationSpec
"azure": .azure.plugins.gloo.solo.io.DestinationSpec
"rest": .rest.plugins.gloo.solo.io.DestinationSpec
"grpc": .grpc.plugins.gloo.solo.io.DestinationSpec

```

| Field | Type | Description | Default |
| ----- | ---- | ----------- |----------- | 
| `aws` | [.aws.plugins.gloo.solo.io.DestinationSpec](../plugins/aws/aws.proto.sk#destinationspec) |  |  |
| `azure` | [.azure.plugins.gloo.solo.io.DestinationSpec](../plugins/azure/azure.proto.sk#destinationspec) |  |  |
| `rest` | [.rest.plugins.gloo.solo.io.DestinationSpec](../plugins/rest/rest.proto.sk#destinationspec) |  |  |
| `grpc` | [.grpc.plugins.gloo.solo.io.DestinationSpec](../plugins/grpc/grpc.proto.sk#destinationspec) |  |  |




---
### UpstreamSpec

 
Each upstream in Gloo has a type. Supported types include `static`, `kubernetes`, `aws`, `consul`, and more.
Each upstream type is handled by a corresponding Gloo plugin.

```yaml
"sslConfig": .gloo.solo.io.UpstreamSslConfig
"kube": .kubernetes.plugins.gloo.solo.io.UpstreamSpec
"static": .static.plugins.gloo.solo.io.UpstreamSpec
"aws": .aws.plugins.gloo.solo.io.UpstreamSpec
"azure": .azure.plugins.gloo.solo.io.UpstreamSpec
"consul": .consul.plugins.gloo.solo.io.UpstreamSpec

```

| Field | Type | Description | Default |
| ----- | ---- | ----------- |----------- | 
| `sslConfig` | [.gloo.solo.io.UpstreamSslConfig](../ssl.proto.sk#upstreamsslconfig) |  |  |
| `kube` | [.kubernetes.plugins.gloo.solo.io.UpstreamSpec](../plugins/kubernetes/kubernetes.proto.sk#upstreamspec) |  |  |
| `static` | [.static.plugins.gloo.solo.io.UpstreamSpec](../plugins/static/static.proto.sk#upstreamspec) |  |  |
| `aws` | [.aws.plugins.gloo.solo.io.UpstreamSpec](../plugins/aws/aws.proto.sk#upstreamspec) |  |  |
| `azure` | [.azure.plugins.gloo.solo.io.UpstreamSpec](../plugins/azure/azure.proto.sk#upstreamspec) |  |  |
| `consul` | [.consul.plugins.gloo.solo.io.UpstreamSpec](../plugins/consul/consul.proto.sk#upstreamspec) |  |  |





<!-- Start of HubSpot Embed Code -->
<script type="text/javascript" id="hs-script-loader" async defer src="//js.hs-scripts.com/5130874.js"></script>
<!-- End of HubSpot Embed Code -->
