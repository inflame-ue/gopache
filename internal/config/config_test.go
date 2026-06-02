package config

import (
	"reflect"
	"testing"
)

func TestNewProxyConfigFromFlags(t *testing.T) {
	type testCase struct {
		port       int
		origin     string
		clearCache bool
		cachePath  string
		want       *ProxyConfig
		wantErr    bool
	}

	tests := map[string]testCase{
		"valid proxy config, no clear cache": {
			port:       3000,
			origin:     "https://httpbin.com",
			clearCache: false,
			cachePath:  "test.json",
			want: &ProxyConfig{
				Port:       3000,
				Origin:     "https://httpbin.com",
				FlushCache: false,
				CachePath:  "test.json",
			},
			wantErr: false,
		},
		"empty origin, clear cache": {
			port:       3000,
			origin:     "",
			clearCache: true,
			cachePath:  "test.json",
			want: &ProxyConfig{
				Port:       3000,
				Origin:     "",
				FlushCache: true,
				CachePath:  "test.json",
			},
			wantErr: false,
		},
		"empty origin, no clear cache": {
			port:       3000,
			origin:     "",
			clearCache: false,
			cachePath:  "",
			want:       nil,
			wantErr:    true,
		},
		"origin is not a valid url": {
			port:       3000,
			origin:     "not a url",
			clearCache: false,
			cachePath:  "",
			want:       nil,
			wantErr:    true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewProxyConfigFromFlags(tc.port, tc.origin, tc.clearCache, tc.cachePath)

			if tc.wantErr && err == nil {
				t.Errorf("expected err, got nil instead")
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
