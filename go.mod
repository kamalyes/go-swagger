module github.com/kamalyes/go-swagger

go 1.25.0

require (
	github.com/fsnotify/fsnotify v1.9.0
	github.com/kamalyes/go-config v0.21.14
	github.com/kamalyes/go-logger v0.6.0
	github.com/kamalyes/go-toolbox v0.16.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/kamalyes/go-argus v0.3.1 // indirect
	go.opentelemetry.io/otel v1.44.0 // indirect
	go.opentelemetry.io/otel/trace v1.44.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	google.golang.org/grpc v1.82.1 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

// 本地开发替换
// replace github.com/kamalyes/go-toolbox => ../go-toolbox

// replace github.com/kamalyes/go-logger => ../go-logger

// replace github.com/kamalyes/go-config => ../go-config
