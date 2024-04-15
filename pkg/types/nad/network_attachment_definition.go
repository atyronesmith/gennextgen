package nad

import (
	networkv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewNetworkAttachmentDefinition() *networkv1.NetworkAttachmentDefinition {
	return &networkv1.NetworkAttachmentDefinition{
		TypeMeta: v1.TypeMeta{
			APIVersion: "k8s.cni.cncf.io/v1",
			Kind:       "NetworkAttachmentDefinition",
		},
	}
}
