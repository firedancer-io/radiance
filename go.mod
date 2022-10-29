module go.firedancer.io/radiance

go 1.19

// Vote instruction support
replace github.com/gagliardetto/solana-go => github.com/ripatel-jump/solana-go v0.0.0-20221019012627-2c3d04688523

require (
	cloud.google.com/go/bigtable v1.16.0
	github.com/LiamHaworth/go-tproxy v0.0.0-20190726054950-ef7efd7f24ed
	github.com/VividCortex/ewma v1.2.0
	github.com/coreos/go-systemd/v22 v22.4.0
	github.com/filecoin-project/go-leb128 v0.0.0-20190212224330-8d79a5489543
	github.com/gagliardetto/binary v0.7.2
	github.com/gagliardetto/solana-go v1.6.1-0.20221018174950-475b9d64e462
	github.com/go-logr/logr v1.2.3
	github.com/golang/protobuf v1.5.2
	github.com/google/gopacket v1.1.19
	github.com/google/nftables v0.0.0-20220611213346-a346d51f53b3
	github.com/ipfs/go-cid v0.3.2
	github.com/ipfs/go-datastore v0.5.1
	github.com/ipfs/go-ipfs-blockstore v1.1.2
	github.com/ipfs/go-ipld-cbor v0.0.5
	github.com/ipld/go-car v0.5.0
	github.com/ipld/go-ipld-prime v0.16.0
	github.com/klauspost/compress v1.15.11
	github.com/linxGnu/grocksdb v1.7.7
	github.com/lucas-clemente/quic-go v0.29.1
	github.com/mattn/go-isatty v0.0.16
	github.com/minio/sha256-simd v1.0.0
	github.com/mr-tron/base58 v1.2.0
	github.com/multiformats/go-multicodec v0.6.0
	github.com/multiformats/go-multihash v0.1.0
	github.com/novifinancial/serde-reflection/serde-generate/runtime/golang v0.0.0-20220519162058-e5cd3c3b3f3a
	github.com/pkg/profile v1.2.1
	github.com/prometheus/client_golang v1.13.0
	github.com/segmentio/textio v1.2.0
	github.com/slok/go-http-metrics v0.10.0
	github.com/spaolacci/murmur3 v1.1.0
	github.com/spf13/cobra v1.5.0
	github.com/stretchr/testify v1.8.0
	github.com/twmb/franz-go v1.7.1
	github.com/vbauerster/mpb/v8 v8.0.2
	golang.org/x/crypto v0.1.0
	golang.org/x/sync v0.1.0
	golang.org/x/sys v0.1.0
	google.golang.org/protobuf v1.28.1
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/klog/v2 v2.80.1
)

require (
	cloud.google.com/go v0.102.1 // indirect
	cloud.google.com/go/compute v1.7.0 // indirect
	cloud.google.com/go/iam v0.3.0 // indirect
	contrib.go.opencensus.io/exporter/stackdriver v0.13.14 // indirect
	filippo.io/edwards25519 v1.0.0 // indirect
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d // indirect
	github.com/andres-erbsen/clock v0.0.0-20160526145045-9e14626cd129 // indirect
	github.com/aybabtme/rgbterm v0.0.0-20170906152045-cc83f3b3ce59 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blendle/zapdriver v1.3.1 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/census-instrumentation/opencensus-proto v0.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cncf/udpa/go v0.0.0-20220112060539-c52dc94e7fbe // indirect
	github.com/cncf/xds/go v0.0.0-20220520190051-1e77728a1eaa // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dfuse-io/logging v0.0.0-20210109005628-b97a57253f70 // indirect
	github.com/envoyproxy/go-control-plane v0.10.2-0.20220325020618-49ff273808a1 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.6.7 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/francoispqt/gojay v1.2.13 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/gagliardetto/treeout v0.1.4 // indirect
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.1.0 // indirect
	github.com/googleapis/gax-go/v2 v2.4.0 // indirect
	github.com/gorilla/rpc v1.2.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/ipfs/bbloom v0.0.4 // indirect
	github.com/ipfs/go-block-format v0.0.3 // indirect
	github.com/ipfs/go-blockservice v0.2.1 // indirect
	github.com/ipfs/go-ipfs-ds-help v1.1.0 // indirect
	github.com/ipfs/go-ipfs-exchange-interface v0.1.0 // indirect
	github.com/ipfs/go-ipfs-util v0.0.2 // indirect
	github.com/ipfs/go-ipld-format v0.2.0 // indirect
	github.com/ipfs/go-ipld-legacy v0.1.0 // indirect
	github.com/ipfs/go-log v1.0.5 // indirect
	github.com/ipfs/go-log/v2 v2.3.0 // indirect
	github.com/ipfs/go-merkledag v0.5.1 // indirect
	github.com/ipfs/go-metrics-interface v0.0.1 // indirect
	github.com/ipfs/go-verifcid v0.0.1 // indirect
	github.com/ipld/go-codec-dagpb v1.3.1 // indirect
	github.com/jbenet/goprocess v0.1.4 // indirect
	github.com/josharian/native v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/logrusorgru/aurora v2.0.3+incompatible // indirect
	github.com/marten-seemann/qtls-go1-18 v0.1.2 // indirect
	github.com/marten-seemann/qtls-go1-19 v0.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/mdlayher/netlink v1.6.0 // indirect
	github.com/mdlayher/socket v0.2.3 // indirect
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mostynb/zstdpool-freelist v0.0.0-20201229113212-927304c0c3b1 // indirect
	github.com/multiformats/go-base32 v0.0.3 // indirect
	github.com/multiformats/go-base36 v0.1.0 // indirect
	github.com/multiformats/go-multibase v0.0.3 // indirect
	github.com/multiformats/go-varint v0.0.6 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.15.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/polydawn/refmt v0.0.0-20201211092308-30ac6d18308e // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/streamingfast/logging v0.0.0-20220813175024-b4fbb0e893df // indirect
	github.com/teris-io/shortid v0.0.0-20220617161101-71ec9f2aa569 // indirect
	github.com/tidwall/gjson v1.14.3 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/twmb/franz-go/pkg/kmsg v1.2.0 // indirect
	github.com/vishvananda/netns v0.0.0-20210104183010-2eb08e3e575f // indirect
	github.com/whyrusleeping/cbor-gen v0.0.0-20200123233031-1cdf64d27158 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/goleak v1.1.12 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/ratelimit v0.2.0 // indirect
	go.uber.org/zap v1.23.0 // indirect
	golang.org/x/exp v0.0.0-20221028150844-83b7d23a625f // indirect
	golang.org/x/mod v0.6.0 // indirect
	golang.org/x/net v0.1.0 // indirect
	golang.org/x/oauth2 v0.1.0 // indirect
	golang.org/x/term v0.1.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	golang.org/x/time v0.1.0 // indirect
	golang.org/x/tools v0.2.0 // indirect
	golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f // indirect
	google.golang.org/api v0.85.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220725144611-272f38e5d71b // indirect
	google.golang.org/grpc v1.48.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	lukechampine.com/blake3 v1.1.6 // indirect
)
