package rest_proxy

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/asoorm/todo-grpc/pkg/log"
	"github.com/fullstorydev/grpcurl"
	"github.com/go-chi/chi"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

func Run(listenPort, grpcServicePort int) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := grpc.Dial(
		fmt.Sprintf("localhost:%d", grpcServicePort),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithBlock(),
	)
	log.FatalOnError(err)
	defer conn.Close()

	reflectionClient := grpcreflect.NewClient(
		context.Background(),
		grpc_reflection_v1alpha.NewServerReflectionClient(conn),
	)

	services, err := reflectionClient.ListServices()
	log.Info("reflected services: %#v", services)

	mainRouter := chi.NewMux()
	for i, svc := range services {
		serviceDescriptor, err := reflectionClient.ResolveService(svc)
		if err != nil {
			log.Error("service resolution error %s", err.Error())
			continue
		}

		log.Info("service %d: %s", i, serviceDescriptor.GetName())

		serviceRouter := chi.NewMux()

		servicePath := fmt.Sprintf("/%s", serviceDescriptor.GetName())
		log.Info("%s", servicePath)

		for _, methodDescriptor := range serviceDescriptor.GetMethods() {
			listenPath := fmt.Sprintf("/%s", methodDescriptor.GetName())
			log.Info("\t%s%s", servicePath, listenPath)

			serviceRouter.Handle(listenPath, requestHandler(reflectionClient, conn, methodDescriptor.GetFullyQualifiedName()))
		}

		mainRouter.Mount(servicePath, serviceRouter)

		log.Info("serviceDescriptor %s", serviceDescriptor.String())
	}

	listenAddress := fmt.Sprintf(":%d", listenPort)
	log.Info("starting rest server on %s", listenAddress)
	log.FatalOnError(http.ListenAndServe(listenAddress, mainRouter))
}

func requestHandler(refClient *grpcreflect.Client, cc *grpc.ClientConn, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		bodyBytes, _ := ioutil.ReadAll(r.Body)

		defer r.Body.Close()

		format := "json"
		includeSeparators := true
		emitDefaults := false
		verbose := true
		descriptorSource := grpcurl.DescriptorSourceFromServer(context.Background(), refClient)
		rf, formatter, err := grpcurl.RequestParserAndFormatterFor(
			grpcurl.Format(format),
			descriptorSource,
			emitDefaults,
			includeSeparators,
			bytes.NewBuffer(bodyBytes), // Something buggy here
		)
		if err != nil {
			log.Error("failed to construct request formatter for %s: %s", format, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			return
		}

		// Shouldn't write directly to http.ResponseWriter
		// Need to format output for REST-based calls, but this is ok for PoC
		handler := grpcurl.NewDefaultEventHandler(w, descriptorSource, formatter, verbose)

		_ = grpcurl.InvokeRPC(context.Background(), descriptorSource, cc, method, []string{}, handler, rf.Next)
	}
}
