package cache

import (
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	cacheMap := NewCacheMap()
	cachedResp := &CachedResponse{
		Status: 200,
		Body:   []byte("success"),
	}
	cacheMap.Set("https://test.com", cachedResp)

	type testCase struct {
		key    string
		want   *CachedResponse
		wantOk bool
	}

	tests := map[string]testCase{
		"valid existing key": {
			key:    "https://test.com",
			want:   cachedResp,
			wantOk: true,
		},
		"valid existing upper-case key": {
			key:    "HTTPS://TEST.COM",
			want:   cachedResp,
			wantOk: true,
		},
		"non-existing key": {
			key:    "bad key",
			want:   nil,
			wantOk: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, ok := cacheMap.Get(tc.key)

			if tc.wantOk != ok {
				t.Errorf("expected ok = %v, got ok = %v", tc.wantOk, ok)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestSet(t *testing.T) {
	cacheMap := NewCacheMap()
	cachedResp := &CachedResponse{
		Status: 200,
		Body:   []byte("success"),
	}
	cachedResp2 := &CachedResponse{
		Status: 200,
		Body:   []byte("new success"),
	}
	
	type testCase struct {
		key       string
		value     *CachedResponse
		want 	  *CachedResponse
		wantOk    bool
	}

	tests := map[string]testCase{
		"set new key, value": {
			key: "https://httpbin.com",
			value: cachedResp,
			want: cachedResp,
			wantOk: true,
		},
		"set overwrite uppercase key, value": {
			key: "HTTPS://HTTPBIN.COM",
			value: cachedResp2,
			want: cachedResp2,
			wantOk: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cacheMap.Set(tc.key, tc.value)
			val, ok := cacheMap.Get(tc.key)

			if tc.wantOk != ok {
				t.Errorf("want ok = %v, got ok = %v", tc.wantOk, ok)
			}

			if !reflect.DeepEqual(tc.want, val) {
				t.Errorf("expected: %v, got: %v", tc.want, val)
			}
		})
	}
}
