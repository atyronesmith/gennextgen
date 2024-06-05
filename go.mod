module github.com/atyronesmith/gennextgen

go 1.22.0

toolchain go1.22.2

require (
	github.com/go-git/go-git/v5 v5.12.0
	github.com/k8snetworkplumbingwg/network-attachment-definition-client v1.7.0
	github.com/k8snetworkplumbingwg/whereabouts v0.7.0
	github.com/openstack-k8s-operators/dataplane-operator/api v0.3.1-0.20240523171901-e556179f4d44
	github.com/wasilibs/go-re2 v1.5.2
	gopkg.in/ini.v1 v1.67.0
	k8s.io/api v0.30.1
	k8s.io/apimachinery v0.30.1
	sigs.k8s.io/yaml v1.4.0
)

require (
	dario.cat/mergo v1.0.0 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/ProtonMail/go-crypto v1.0.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cert-manager/cert-manager v1.13.6 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudflare/circl v1.3.7 // indirect
	github.com/cyphar/filepath-securejoin v0.2.4 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/emicklei/go-restful/v3 v3.12.0 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/evanphx/json-patch/v5 v5.9.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.5.0 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.19.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/gnostic-models v0.6.9-0.20230804172637-c7be7c783f49 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gophercloud/gophercloud v1.11.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/magefile/mage v1.14.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/metal3-io/baremetal-operator/apis v0.6.0 // indirect
	github.com/metal3-io/baremetal-operator/pkg/hardwareutils v0.6.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/openshift/api v3.9.0+incompatible // indirect
	github.com/openstack-k8s-operators/barbican-operator/api v0.0.0-20240523081006-98dc23995369 // indirect
	github.com/openstack-k8s-operators/cinder-operator/api v0.3.1-0.20240514185824-359c4b1e9655 // indirect
	github.com/openstack-k8s-operators/designate-operator/api v0.0.0-20240523092716-301e74907910 // indirect
	github.com/openstack-k8s-operators/glance-operator/api v0.3.1-0.20240522202242-daeab50a9031 // indirect
	github.com/openstack-k8s-operators/heat-operator/api v0.3.1-0.20240520005136-22ce42e2173f // indirect
	github.com/openstack-k8s-operators/horizon-operator/api v0.3.1-0.20240515093015-9d2609bb0015 // indirect
	github.com/openstack-k8s-operators/infra-operator/apis v0.3.1-0.20240429104248-25176c735750
	github.com/openstack-k8s-operators/ironic-operator/api v0.3.1-0.20240523093522-7a25580a799d // indirect
	github.com/openstack-k8s-operators/keystone-operator/api v0.3.1-0.20240523092718-53aebcc693a5 // indirect
	github.com/openstack-k8s-operators/lib-common/modules/common v0.3.1-0.20240523061957-78c308214ae3 // indirect
	github.com/openstack-k8s-operators/lib-common/modules/openstack v0.3.1-0.20240523061957-78c308214ae3 // indirect
	github.com/openstack-k8s-operators/lib-common/modules/storage v0.3.1-0.20240523061957-78c308214ae3 // indirect
	github.com/openstack-k8s-operators/manila-operator/api v0.3.1-0.20240523092949-4ce774522565 // indirect
	github.com/openstack-k8s-operators/mariadb-operator/api v0.3.1-0.20240514152407-b2bea62f05db // indirect
	github.com/openstack-k8s-operators/neutron-operator/api v0.3.1-0.20240523093521-639e83314a88 // indirect
	github.com/openstack-k8s-operators/nova-operator/api v0.3.1-0.20240522182146-9eeda511bbe1 // indirect
	github.com/openstack-k8s-operators/octavia-operator/api v0.3.1-0.20240523092718-22ea0a3c5100 // indirect
	github.com/openstack-k8s-operators/openstack-baremetal-operator/api v0.3.1-0.20240523005930-c332745390c8
	github.com/openstack-k8s-operators/openstack-operator/apis v0.0.0-20240523103020-86d8ae842b80 // indirect
	github.com/openstack-k8s-operators/ovn-operator/api v0.3.1-0.20240522132732-6d6cd2162ea3 // indirect
	github.com/openstack-k8s-operators/placement-operator/api v0.3.1-0.20240522174844-950b72813aa5 // indirect
	github.com/openstack-k8s-operators/swift-operator/api v0.3.1-0.20240514095108-f59b40442e91 // indirect
	github.com/openstack-k8s-operators/telemetry-operator/api v0.3.1-0.20240523083214-397d3fd0da29 // indirect
	github.com/pjbgf/sha1cd v0.3.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.19.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.53.0 // indirect
	github.com/prometheus/procfs v0.15.0 // indirect
	github.com/rabbitmq/cluster-operator/v2 v2.6.0 // indirect
	github.com/rhobs/obo-prometheus-operator/pkg/apis/monitoring v0.69.0-rhobs1 // indirect
	github.com/rhobs/observability-operator v0.0.28 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/sergi/go-diff v1.3.2-0.20230802210424-5b0b94c5c0d3 // indirect
	github.com/skeema/knownhosts v1.2.2 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tetratelabs/wazero v1.7.0 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/oauth2 v0.20.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/term v0.20.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	golang.org/x/tools v0.21.0 // indirect
	gomodules.xyz/jsonpatch/v2 v2.4.0 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	k8s.io/apiextensions-apiserver v0.30.1 // indirect
	k8s.io/client-go v0.30.1 // indirect
	k8s.io/kube-openapi v0.0.0-20240521193020-835d969ad83a // indirect
	sigs.k8s.io/controller-runtime v0.18.2 // indirect
	sigs.k8s.io/gateway-api v0.8.0 // indirect
)

require (
	github.com/containernetworking/cni v1.2.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/klog/v2 v2.120.1 // indirect
	k8s.io/utils v0.0.0-20240502163921-fe8a2dddb1d0 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
)

// mschuppert: map to latest commit from release-4.13 tag
// must consistent within modules and service operators
replace github.com/openshift/api => github.com/openshift/api v0.0.0-20230414143018-3367bc7e6ac7 //allow-merging

// custom RabbitmqClusterSpecCore for OpenStackControlplane (v2.6.0_patches_tag)
replace github.com/rabbitmq/cluster-operator/v2 => github.com/openstack-k8s-operators/rabbitmq-cluster-operator/v2 v2.6.1-0.20240313124519-961a0ee8bf7f //allow-merging
