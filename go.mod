module github.com/uptrace/go-realworld-example-app

go 1.25.0

require (
	github.com/benbjohnson/clock v1.0.3
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/go-pg/migrations/v8 v8.0.1
	github.com/go-pg/pg/v10 v10.6.2
	github.com/go-pg/pgext v0.2.0
	github.com/go-pg/urlstruct v1.0.0
	github.com/go-redis/cache/v8 v8.2.1
	github.com/go-redis/redis/v8 v8.3.3
	github.com/go-redis/redis_rate/v9 v9.0.2
	github.com/go-redis/redisext v0.3.1
	github.com/gosimple/slug v1.9.0
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.43.0
	github.com/sirupsen/logrus v1.7.0
	github.com/uptrace/uptrace-go v0.4.2
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.13.0
	go.opentelemetry.io/otel v0.13.0
	go.opentelemetry.io/otel/sdk v0.13.0
	golang.org/x/crypto v0.53.0
	golang.org/x/exp v0.0.0-20201008143054-e3b2a7f2fdc7
	gopkg.in/yaml.v2 v2.3.0
)

require (
	cloud.google.com/go v0.26.0 // indirect
	dmitri.shuralyov.com/gpu/mtl v0.0.0-20190408044501-666a987793e9 // indirect
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/BurntSushi/xgb v0.0.0-20160522181843-27f122750802 // indirect
	github.com/DataDog/sketches-go v0.0.1 // indirect
	github.com/census-instrumentation/opencensus-proto v0.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/client9/misspell v0.3.4 // indirect
	github.com/cncf/udpa/go v0.0.0-20191209042840-269d4d468f6f // indirect
	github.com/codemodus/kace v0.5.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/envoyproxy/go-control-plane v0.9.4 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.1.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20200222043503-6f7a984d4dc4 // indirect
	github.com/go-pg/pg/v9 v9.1.6 // indirect
	github.com/go-pg/zerochecker v0.2.0 // indirect
	github.com/go-playground/assert/v2 v2.0.1 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b // indirect
	github.com/golang/mock v1.1.1 // indirect
	github.com/golang/protobuf v1.5.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/hpcloud/tail v1.0.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/klauspost/compress v1.11.2 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/kr/pty v1.1.1 // indirect
	github.com/kr/text v0.1.0 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/nxadm/tail v1.4.4 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4 // indirect
	github.com/rainycape/unidecode v0.0.0-20150907023854-cb7f23ec59be // indirect
	github.com/segmentio/encoding v0.2.2 // indirect
	github.com/stretchr/objx v0.1.0 // indirect
	github.com/stretchr/testify v1.6.1 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/ugorji/go v1.1.13 // indirect
	github.com/ugorji/go/codec v1.1.13 // indirect
	github.com/vmihailenco/bufpool v0.1.11 // indirect
	github.com/vmihailenco/go-tinylfu v0.1.0 // indirect
	github.com/vmihailenco/msgpack/v4 v4.3.11 // indirect
	github.com/vmihailenco/msgpack/v5 v5.0.0-rc.2 // indirect
	github.com/vmihailenco/tagparser v0.1.2 // indirect
	github.com/yuin/goldmark v1.4.13 // indirect
	go.opentelemetry.io/contrib v0.13.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.13.0 // indirect
	go.opentelemetry.io/contrib/propagators v0.13.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout v0.13.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/image v0.0.0-20190802002840-cff245a6509b // indirect
	golang.org/x/lint v0.0.0-20190313153728-d0100b6bd8b3 // indirect
	golang.org/x/mobile v0.0.0-20190719004257-d2bd2a29d028 // indirect
	golang.org/x/mod v0.36.0 // indirect
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/telemetry v0.0.0-20260508192327-42602be52be6 // indirect
	golang.org/x/term v0.44.0 // indirect
	golang.org/x/text v0.38.0 // indirect
	golang.org/x/tools v0.45.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/grpc v1.31.0 // indirect
	google.golang.org/protobuf v1.36.7 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
	honnef.co/go/tools v0.0.0-20190523083050-ea95bdfd59fc // indirect
	mellium.im/sasl v0.2.1 // indirect
)
