module github.com/evcc-io/evcc

go 1.18

require (
	github.com/AlecAivazis/survey/v2 v2.3.5
	github.com/BurntSushi/toml v1.2.0
	github.com/Masterminds/sprig/v3 v3.2.2
	github.com/PuerkitoBio/goquery v1.8.0
	github.com/asaskevich/EventBus v0.0.0-20200907212545-49d423059eef
	github.com/avast/retry-go/v3 v3.1.1
	github.com/aws/aws-sdk-go v1.44.84
	github.com/basgys/goxml2json v1.1.0
	github.com/basvdlei/gotsmart v0.0.3
	github.com/benbjohnson/clock v1.3.0
	github.com/bogosj/tesla v1.0.2
	github.com/cjrd/allocate v0.0.0-20220510215731-986f24f0fb18
	github.com/cloudfoundry/jibber_jabber v0.0.0-20151120183258-bcc4c8345a21
	github.com/containrrr/shoutrrr v0.6.1
	github.com/coreos/go-oidc/v3 v3.2.0
	github.com/denisbrodbeck/machineid v1.0.1
	github.com/dmarkham/enumer v1.5.6
	github.com/dustin/go-humanize v1.0.0
	github.com/dylanmei/iso8601 v0.1.0
	github.com/eclipse/paho.mqtt.golang v1.4.1
	github.com/emirpasic/gods v1.18.1
	github.com/evcc-io/eebus v0.0.0-20220815121332-bbb3cd4260e0
	github.com/fatih/structs v1.1.0
	github.com/foogod/go-powerwall v0.2.0
	github.com/go-ping/ping v1.1.0
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/godbus/dbus/v5 v5.1.0
	github.com/gokrazy/updater v0.0.0-20211121155532-30ae8cd650ea
	github.com/golang/mock v1.6.0
	github.com/google/go-github/v32 v32.1.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.5.0
	github.com/gregdel/pushover v1.1.0
	github.com/grid-x/modbus v0.0.0-20220419073012-0daecbb3900f
	github.com/hashicorp/go-version v1.6.0
	github.com/imdario/mergo v0.3.13
	github.com/influxdata/influxdb-client-go/v2 v2.9.2
	github.com/itchyny/gojq v0.12.8
	github.com/jeremywohl/flatten v1.0.1
	github.com/jinzhu/copier v0.3.5
	github.com/joeshaw/carwings v0.0.0-20220711204617-d90801cf2848
	github.com/jpfielding/go-http-digest v0.0.0-20211006141426-fbc93758452e
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/koron/go-ssdp v0.0.3
	github.com/korylprince/ipnetgen v1.0.1
	github.com/libp2p/zeroconf/v2 v2.2.0
	github.com/lorenzodonini/ocpp-go v0.16.0
	github.com/lunixbochs/struc v0.0.0-20200707160740-784aaebc1d40
	github.com/manifoldco/promptui v0.9.0
	github.com/mergermarket/go-pkcs7 v0.0.0-20170926155232-153b18ea13c9
	github.com/mitchellh/mapstructure v1.5.0
	github.com/mlnoga/rct v0.1.2-0.20220320164346-9f2daa4d6734
	github.com/muka/go-bluetooth v0.0.0-20220219050759-674a63b8741a
	github.com/mxschmitt/golang-combinations v1.1.0
	github.com/nicksnyder/go-i18n/v2 v2.2.0
	github.com/nirasan/go-oauth-pkce-code-verifier v0.0.0-20220510032225-4f9f17eaec4c
	github.com/olekukonko/tablewriter v0.0.5
	github.com/philippseith/signalr v0.5.3-0.20211205201131-d57b5a34379a
	github.com/prometheus/client_golang v1.13.0
	github.com/prometheus/common v0.37.0
	github.com/robertkrimen/otto v0.0.0-20211024170158-b87d35c0b86f
	github.com/samber/lo v1.27.0
	github.com/shurcooL/graphql v0.0.0-20220606043923-3cf50f8a0a29
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	github.com/spf13/cobra v1.5.0
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.12.0
	github.com/stretchr/testify v1.8.0
	github.com/tv42/httpunix v0.0.0-20191220191345-2ba4b9c3382c
	github.com/volkszaehler/mbmd v0.0.0-20220713075154-0069b8b73e4f
	github.com/writeas/go-strip-markdown v2.0.1+incompatible
	gitlab.com/bboehmke/sunny v0.15.1-0.20211022160056-2fba1c86ade6
	golang.org/x/exp v0.0.0-20220823124025-807a23277127
	golang.org/x/net v0.0.0-20220822230855-b0a4917ee28c
	golang.org/x/oauth2 v0.0.0-20220822191816-0ebed06d0094
	golang.org/x/text v0.3.7
	google.golang.org/api v0.94.0
	google.golang.org/grpc v1.49.0
	google.golang.org/protobuf v1.28.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	cloud.google.com/go/compute v1.9.0 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.1.1 // indirect
	github.com/andig/gosunspec v0.0.0-20211108155140-af2e73b86e71 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/chzyer/readline v1.5.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/deepmap/oapi-codegen v1.11.0 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.1.0 // indirect
	github.com/googleapis/gax-go/v2 v2.5.1 // indirect
	github.com/grid-x/serial v0.0.0-20211107191517-583c7356b3aa // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/influxdata/line-protocol v0.0.0-20210922203350-b1ad95c89adf // indirect
	github.com/itchyny/timefmt-go v0.1.3 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/miekg/dns v1.1.50 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/pascaldekloe/name v1.0.1 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/rickb777/date v1.20.0 // indirect
	github.com/rickb777/plural v1.4.1 // indirect
	github.com/rivo/uniseg v0.3.4 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/smartystreets/goconvey v1.7.2 // indirect
	github.com/spf13/afero v1.9.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/subosito/gotenv v1.4.1 // indirect
	github.com/teivah/onecontext v1.3.0 // indirect
	github.com/thoas/go-funk v0.9.2 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.5 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/crypto v0.0.0-20220824171710-5757bc0c5503 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/sync v0.0.0-20220819030929-7fc1605a5dde // indirect
	golang.org/x/sys v0.0.0-20220823224334-20c2bfdbfe24 // indirect
	golang.org/x/term v0.0.0-20220722155259-a9ba230a4035 // indirect
	golang.org/x/tools v0.1.12 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220822174746-9e6da59bd2fc // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	nhooyr.io/websocket v1.8.7 // indirect
)

replace github.com/foogod/go-powerwall => github.com/andig/go-powerwall v0.2.1-0.20220205120646-e5220ad9a9a0
