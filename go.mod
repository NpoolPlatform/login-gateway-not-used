module github.com/NpoolPlatform/login-gateway

go 1.16

require (
	entgo.io/ent v0.10.1
	github.com/NpoolPlatform/api-manager v0.0.0-20220505084652-c0caff45e937
	github.com/NpoolPlatform/appuser-manager v0.0.0-20220526091819-514238cc3520
	github.com/NpoolPlatform/go-service-framework v0.0.0-20220404143809-82c40930388a
	github.com/NpoolPlatform/libent-cruder v0.0.0-20220503164840-0962cb617722
	github.com/NpoolPlatform/message v0.0.0-20220520115554-2c9bbb3dbe95
	github.com/NpoolPlatform/third-gateway v0.0.0-20220205094704-a585d53bd025
	github.com/NpoolPlatform/third-login-gateway v0.0.0-20220511110308-4d1de890428c
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-redis/redis/v8 v8.11.4
	github.com/go-resty/resty/v2 v2.7.0
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.0
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.7.1
	github.com/urfave/cli/v2 v2.5.1
	google.golang.org/grpc v1.46.2
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.2.0
	google.golang.org/protobuf v1.28.0
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.41.0

replace (
	github.com/NpoolPlatform/appuser-manager => github.com/NpoolPlatform/appuser-manager v0.0.0-20220525131713-a1779ffa6288
	github.com/NpoolPlatform/message => github.com/NpoolPlatform/message v0.0.0-20220526101845-3393acf1ee8f
	github.com/NpoolPlatform/third-login-gateway => github.com/NpoolPlatform/third-login-gateway v0.0.2-0.20220526125409-c75fab8bc628
)
