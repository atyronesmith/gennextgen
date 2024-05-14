
# Table of Contents

- [Section 1](#section-1)
- [Section 2](#section-2)
- [Section 3](#section-3)

## Section 1

In the nncp CR file, configure the interfaces for each isolated network on each worker node in the RHOCP cluster.
There is an IP address per interface for each isolated network on the worker node hosting OSP services.  There is
a nncp configuration for each worker node that is going to run OSP control place services.  There must be >= 2
worker nodes to run the OSP services.  The workers run services that would typically run on the primary OSP controllers
as well as Networker, Storage, OVN, nodes, etc...

In the nad CR file, configure a nad resource for each isolated network to attach a service deployment pod to the network.

In the IPAddressPool CR file, configure an IPAddressPool resource on the isolated network to specify the IP address ranges over which MetalLB has authority:

In the L2Advertisement CR file, configure a L2Advertisement resource to define which node advertises a service to the local network. Create one L2Advertisement resource for each network.

In the openstacknetconfig.yaml file, define the topology for each data plane network. 

## Section 2

This is the content of section 2.

## Section 3

This is the content of section 3.