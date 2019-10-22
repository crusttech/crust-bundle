module github.com/crusttech/crust-bundle

go 1.12

require (
	github.com/cortezaproject/corteza-server v0.0.0-20191021084042-941ae38cb626
	github.com/go-chi/chi v3.3.4+incompatible
	github.com/joho/godotenv v1.3.0
	github.com/lib/pq v1.1.1 // indirect
	github.com/prometheus/client_golang v0.9.3 // indirect
	github.com/prometheus/common v0.4.1 // indirect
	github.com/prometheus/procfs v0.0.0-20190523193104-a7aeb8df3389 // indirect
	github.com/spf13/cobra v0.0.4 // indirect
)

replace gopkg.in/Masterminds/squirrel.v1 => github.com/Masterminds/squirrel v1.1.0
