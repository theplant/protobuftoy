node_modules/.bin/pbts:
	yarn

node_modules/.bin/pbjs:
	yarn

testproto.d.ts: test.proto testproto.js node_modules/.bin/pbts
	yarn run pbts testproto.js -- -o testproto.d.ts

testproto.js: test.proto node_modules/.bin/pbjs
	yarn run pbjs -- -t static-module -w es6 -p . ./test.proto -o testproto.js

test.pb.go: test.proto
	protoc --go_out=plugins=grpc:. test.proto

runjs: harness.ts testproto.d.ts testproto.js
	yarn run ts-node harness.ts

rungo: test.pb.go
	go run examplemain/main.go
