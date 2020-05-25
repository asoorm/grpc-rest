package v1

import (
	"context"
	"fmt"
	"strings"

	"github.com/asoorm/todo-grpc/pkg/log"

	"github.com/asoorm/todo-grpc/pkg/model/v1/address_formatter"
)

func NewAddressFormatterService() *AddressFormatterService {
	return &AddressFormatterService{}
}

type AddressFormatterService struct{}

func (s AddressFormatterService) Format(ctx context.Context, in *address_formatter.AddressRequest) (*address_formatter.AddressResponse, error) {
	log.Info("format: %s", in.String())

	if err := checkAPIVersion(in.ApiVersion); err != nil {
		return nil, err
	}

	return &address_formatter.AddressResponse{
		ApiVersion:      apiVersion,
		BillingAddress:  strings.TrimRight(fmt.Sprintf("%s, %s, %s", in.BillingAddress.StreetAddress, in.BillingAddress.City, in.BillingAddress.State), ", "),
		ShippingAddress: strings.TrimRight(fmt.Sprintf("%s, %s, %s", in.ShippingAddress.StreetAddress, in.ShippingAddress.City, in.ShippingAddress.State), ", "),
	}, nil
}
