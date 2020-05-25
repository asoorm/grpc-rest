package rest_proxy

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/asoorm/oas3"
	"github.com/asoorm/todo-grpc/pkg/log"
	"github.com/fullstorydev/grpcurl"
	"github.com/go-chi/chi"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoprint"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

type Property struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

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

	docs := &oas3.Oas3{}

	mainRouter := chi.NewMux()
	for i, svc := range services {
		// assumption that 1st service is the reflection service - Investigate further
		if i == 0 {
			continue
		}

		serviceDescriptor, err := reflectionClient.ResolveService(svc)
		if err != nil {
			log.Error("service resolution error %s", err.Error())
			continue
		}

		//printServiceDescriptor(serviceDescriptor)

		serviceRouter := chi.NewMux()

		servicePath := fmt.Sprintf("/%s", serviceDescriptor.GetName())

		for _, methodDescriptor := range serviceDescriptor.GetMethods() {

			//printMethodDescriptor(methodDescriptor)

			listenPath := fmt.Sprintf("/%s", methodDescriptor.GetName())

			docs.Add(servicePath+listenPath, http.MethodPost, &oas3.Operation{
				OperationId: listenPath,
				Description: "This gets lost in translation",
				Parameters:  nil,
				Summary:     "This gets lost in translation",
				Responses:   nil,
			})

			log.Info("****************** listening on: %s: %s", listenPath, methodDescriptor.GetFullyQualifiedName())
			serviceRouter.Handle(listenPath, requestHandler(reflectionClient, conn, methodDescriptor.GetFullyQualifiedName()))
		}

		//fileDescriptor := serviceDescriptor.GetFile()
		//printFileDescriptor(fileDescriptor)
		//
		//messageDescriptors := fileDescriptor.GetMessageTypes()
		//printMessageDescriptors(messageDescriptors)

		mainRouter.Mount(servicePath, serviceRouter)
	}

	mainRouter.Handle("/swagger.json", documentationHandler(docs))

	listenAddress := fmt.Sprintf(":%d", listenPort)
	log.Info("starting rest server on %s", listenAddress)
	log.FatalOnError(http.ListenAndServe(listenAddress, mainRouter))
}

func documentationHandler(docs *oas3.Oas3) http.Handler {

	docs.Openapi = "3.0.0"
	docs.Info = oas3.Info{
		// Don't know if tyk cares about this stuff???
		//Title:          "",
		//Version:        "",
		//Description:    "",
		//TermsOfService: "",
		//Contact: struct {
		//	Name  string `json:"name"`
		//	Email string `json:"email"`
		//	URL   string `json:"url"`
		//}{},
		//License: struct {
		//	Name string `json:"name"`
		//	URL  string `json:"url"`
		//}{},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		w.Header().Set("access-control-allow-origin", "*")
		data, _ := json.Marshal(docs)
		w.Write(data)
	})
}

func requestHandler(refClient *grpcreflect.Client, cc *grpc.ClientConn, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := ioutil.ReadAll(r.Body)

		defer r.Body.Close()

		log.Info("body: %s", string(bodyBytes))

		format := "json"
		includeSeparators := true
		emitDefaults := false
		verbose := true
		descriptorSource := grpcurl.DescriptorSourceFromServer(context.Background(), refClient)
		requestFormatter, formatter, err := grpcurl.RequestParserAndFormatterFor(
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
		handler := grpcurl.NewDefaultEventHandler(io.MultiWriter(os.Stdout, w), descriptorSource, formatter, verbose)

		err = grpcurl.InvokeRPC(context.Background(), descriptorSource, cc, method, []string{}, handler, requestFormatter.Next)
		if err != nil {
			log.Error("unable to invoke RPC: %#v", err)
		}
	}
}

func printFileDescriptor(fileDescriptor *desc.FileDescriptor) {
	log.Info("********** fileDescriptor: ***********")
	printer := protoprint.Printer{}
	fd, _ := printer.PrintProtoToString(fileDescriptor)
	log.Info("proto: %s", fd)
}

func printMessageDescriptors(messageDescriptors []*desc.MessageDescriptor) {
	log.Info("********** messageDescriptors: ***********")
	printer := protoprint.Printer{}
	for i, md := range messageDescriptors {
		log.Info("  [%d]: %s", i, md.GetName())
		log.Info("  [%d]: String: %s", i, md.String())
		log.Info("  [%d]: GetFullyQualifiedName: %s", i, md.GetFullyQualifiedName())
		fd, _ := printer.PrintProtoToString(md)
		log.Info("  [%d]: proto: %s", i, fd)
	}
}

func printServiceDescriptor(serviceDescriptor *desc.ServiceDescriptor) {
	log.Info("********** serviceDescriptor: ***********")
	printer := protoprint.Printer{}
	proto, err := printer.PrintProtoToString(serviceDescriptor)
	log.FatalOnError(err)
	log.Info("proto: %s", proto)
}

func printMethodDescriptor(methodDescriptor *desc.MethodDescriptor) {
	log.Info("********** methodDescriptor: ***********")
	printer := protoprint.Printer{}
	proto, err := printer.PrintProtoToString(methodDescriptor)
	log.FatalOnError(err)
	log.Info("proto: %s", proto)
}
