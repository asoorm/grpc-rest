package schema

import (
	"fmt"
	"testing"

	"github.com/xeipuuv/gojsonschema"
)

const schema = `
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": " v1.PingMessage",

  "definitions": {
    "address": {
      "type": "object",
      "properties": {
        "street_address": {"type": "string"},
        "city": {"type": "string"},
        "state": {"type": "string"}
      },
      "required": ["street_address", "city"]
    }
  },

  "type": "object",

  "properties": {
    "billing_address": { "$ref": "#/definitions/address" },
    "shipping_address": { "$ref": "#/definitions/address" }
  },

  "required": ["billing_address", "shipping_address"]
}
`

func TestStaticSchema(t *testing.T) {

	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "empty array input",
			input: `[]`,
			valid: false,
		},
		{
			name:  "empty object input",
			input: `{}`,
			valid: false,
		},
		{
			name:  "missing required root level fields",
			input: `{"a": "foo", "b": "bar"}`,
			valid: false,
		},
		{
			name:  "expected object, given string (both)",
			input: `{"billing_address": "foo", "shipping_address": "bar"}`,
			valid: false,
		},
		{
			name:  "correct fields & types, missing required subfields",
			input: `{"billing_address": {}, "shipping_address": {}}`,
			valid: false,
		},
		{
			name:  "this should pass",
			input: `{"billing_address": {"street_address":"foo", "city": "bar"}, "shipping_address": {"street_address":"bar", "city": "baz"}}`,
			valid: true,
		},
	}

	schemaLoader := gojsonschema.NewStringLoader(schema)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputLoader := gojsonschema.NewStringLoader(tt.input)

			res, err := gojsonschema.Validate(schemaLoader, inputLoader)
			if err != nil {
				t.Fatal(err.Error())
			}
			if res.Valid() != tt.valid {
				t.Fatal(fmt.Errorf("input: %s\n\texp: %t got: %t", tt.input, tt.valid, res.Valid()).Error())
			} else {
				t.Logf("errors: %q", res.Errors())
			}
		})
	}
}
