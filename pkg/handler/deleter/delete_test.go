package deleter_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tomMoulard/htransformation/pkg/handler/deleter"
	"github.com/tomMoulard/htransformation/pkg/types"
)

func TestDeleteHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		rule           types.Rule
		requestHeaders map[string]string
		want           map[string]string
	}{
		{
			name: "Remove not existing header",
			rule: types.Rule{
				Header: "X-Test",
			},
			requestHeaders: map[string]string{
				"Foo": "Bar",
			},
			want: map[string]string{
				"Foo": "Bar",
			},
		},
		{
			name: "Remove one header",
			rule: types.Rule{
				Header: "X-Test",
			},
			requestHeaders: map[string]string{
				"Foo":    "Bar",
				"X-Test": "Bar",
			},
			want: map[string]string{
				"Foo": "Bar",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
			require.NoError(t, err)

			for hName, hVal := range test.requestHeaders {
				req.Header.Add(hName, hVal)
			}

			deleter.Handle(nil, req, test.rule)

			for hName, hVal := range test.want {
				assert.Equal(t, hVal, req.Header.Get(hName))
			}
		})
	}
}

func TestValidation(t *testing.T) {
	testCases := []struct {
		name   string
		rule   types.Rule
		expect assert.ErrorAssertionFunc
	}{
		{
			name:   "no rules",
			expect: assert.NoError,
		},
		{
			name: "valid rule",
			rule: types.Rule{
				Name:   "Delete Rule",
				Header: "not-empty",
				Type:   types.Delete,
			},
			expect: assert.NoError,
		},
	}

	for _, test := range testCases {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := deleter.Validate(test.rule)
			t.Log(err)
			test.expect(t, err)
		})
	}
}
