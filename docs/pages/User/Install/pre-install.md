---
title: Pre-Install
weight: 2
---

### Introduction

Liqo can be installed either in private or local clusters. Its configuration depends on the type of connectivity between the two clusters. Before installing Liqo, you have to consider how your clusters can connect to each other and can peer together.

### Peering Requirements

Liqo requires the following services to be reciprocally accessible to perform cluster peerings:

* **Authentication server** used to authenticate other clusters (i.e. `liqo-auth`).
* The **Kubernetes API server** you want to peer.
* **Network gateway** used to establish interconnection between clusters (i.e. `liqo-gateway`)

Those services have to be accessible from the other clusters to peer with them. This may change the way you would like to have them exposed.
Below it is possible to find some common scenarios that Liqo can handle. Once you identify yours, you can go have to the *table* of each section to find the right values you should specify when installing Liqo using the chart.

### Cloud to cloud

![](/images/scenarios/cloud-to-cloud.svg)

Two managed clusters peered together through the internet. It is possible to have a multi-cloud setup (AKS to AKS, GKE to GKE, and AKS to GKE).

|  | Cluster A (Cloud) | Cluster B (Cloud) |
| --------- | -------- |  ---------       |
| **Auth Server** |  LoadBalancer/ingress | LoadBalancer/ingress |
| **API server** | Provided | Provided |
| **Network gateway** | LoadBalancer | LoadBalancer |

### On-premise to cloud

![](/images/scenarios/on-prem-to-cloud.svg)

On-premise cluster (K3s or K8s) exposed through the Internet peered with a Managed cluster (AKS or GKE).

|  | Cluster A (On-prem) | Cluster B (Cloud) |
| --------- | -------- |  ---------       |
| **Auth Server** |  LoadBalancer/ingress | LoadBalancer/ingress |
| **API server** | Ingress/Public IP | Provided |
| **Network gateway** | LoadBalancer | LoadBalancer |

### On-premise to on-premise

![](/images/scenarios/on-prem-to-on-prem.svg)

On-premise cluster (K3s or K8s) peered with another on-premise cluster (K3s or K8s) in the same LAN.
From the discovery perspective, if the clusters you would like to connect are in the same L2 broadcast domain, the Liqo discovery mechanism based on mDNS will handle the discovery automatically. If you have your clusters in different L3 domains, you have to manually [create](/user/post-install/discovery#forging-the-foreigncluster) a *foreign_cluster* resource or rely on [DNS discovery](/user/post-install/discovery#manual-configuration).

|  | Cluster A (On-prem) | Cluster B (On-prem) |
| --------- | -------- |  ---------       |
| **Auth Server** |  NodePort | NodePort |
| **API server** | Exposed | Exposed |
| **Network gateway** | NodePort | NodePort |

### On-premise behind NAT to cloud

![](/images/scenarios/on-prem-nat-to-cloud.svg)

On-premise cluster (K3s or K8s) exposed through a NAT over the Internet peered with a managed cluster (AKS or GKE).

|  | Cluster A (On-prem behind NAT) | Cluster B (Cloud) |
| --------- | -------- |  ---------       |
| **Auth Server** |  NodePort with port-forwarding | LoadBalancer/ingress |
| **API server** | Port-forwarding | Provided |
| **Network gateway** | NodePort with port-forwarding | LoadBalancer |
