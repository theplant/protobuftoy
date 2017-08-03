package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gogo/protobuf/proto"

	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

type search struct{}

func (s search) Search(ctx context.Context, r *SearchRequest) (*SearchResponse, error) {
	fmt.Println("SEARCH", r)
	return &SearchResponse{
		Result: []*Result{&Result{Url: r.Query}}}, nil

}

func (s search) SearchAlt(ctx context.Context, r *SearchRequest) (*SearchResponse, error) {
	fmt.Println("SEARCH ALT", r)
	return &SearchResponse{
		Result: []*Result{&Result{Url: "Search alt"}}}, nil

}

func wrapper(m grpc.MethodDesc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//func _SearchService_Search_Handler(
		dec := func(i interface{}) error {
			buff, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return err
			}
			return proto.Unmarshal(buff, i.(proto.Message))
		}

		var interceptor grpc.UnaryServerInterceptor
		//		interceptor := func(ctx context.Context, req interface{}, info *UnaryServerInfo, handler UnaryHandler) (resp interface{}, err error) {
		//			fmt.Println("req", req)
		//		}

		//) (interface{}, error) {
		resp, err := m.Handler(
			//srv interface{},
			search{},
			//ctx context.Context,
			r.Context(),
			//dec func(interface{}) error,
			dec,
			//interceptor grpc.UnaryServerInterceptor
			interceptor)

		fmt.Println("handler error", err)

		b, err := proto.Marshal(resp.(proto.Message))

		fmt.Println("marshal error", err)

		w.Write(b)
	})
}

func main() {
	mux := http.NewServeMux()
	for _, desc := range _SearchService_serviceDesc.Methods {
		mux.Handle("/"+desc.MethodName, wrapper(desc))
	}
	fmt.Println("OK, GO")
	http.ListenAndServe(":8080", mux)
}
