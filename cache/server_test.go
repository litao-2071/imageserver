package cache_test

import (
	"crypto/sha256"
	"testing"

	"github.com/pierrre/imageserver"
	. "github.com/pierrre/imageserver/cache"
	cachetest "github.com/pierrre/imageserver/cache/_test"
	"github.com/pierrre/imageserver/testdata"
)

func TestCacheImageServer(t *testing.T) {
	s := CacheImageServer{
		ImageServer: imageserver.ImageServerFunc(func(parameters imageserver.Parameters) (*imageserver.Image, error) {
			return testdata.Medium, nil
		}),
		Cache: cachetest.NewMapCache(),
		KeyGenerator: KeyGeneratorFunc(func(parameters imageserver.Parameters) string {
			return "test"
		}),
	}
	image1, err := s.Get(imageserver.Parameters{})
	if err != nil {
		t.Fatal(err)
	}
	image2, err := s.Get(imageserver.Parameters{})
	if err != nil {
		t.Fatal(err)
	}
	if !imageserver.ImageEqual(image1, image2) {
		t.Fatal("not equal")
	}
}

func TestCacheImageServerErrorImageServer(t *testing.T) {
	s := CacheImageServer{
		ImageServer: imageserver.ImageServerFunc(func(parameters imageserver.Parameters) (*imageserver.Image, error) {
			return nil, imageserver.NewError("error")
		}),
		Cache: cachetest.NewMapCache(),
		KeyGenerator: KeyGeneratorFunc(func(parameters imageserver.Parameters) string {
			return "test"
		}),
	}
	_, err := s.Get(imageserver.Parameters{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestCacheImageServerErrorCacheSet(t *testing.T) {
	s := CacheImageServer{
		ImageServer: imageserver.ImageServerFunc(func(parameters imageserver.Parameters) (*imageserver.Image, error) {
			return testdata.Medium, nil
		}),
		Cache: &cachetest.FuncCache{
			GetFunc: func(key string, parameters imageserver.Parameters) (*imageserver.Image, error) {
				return nil, imageserver.NewError("error")
			},
			SetFunc: func(key string, image *imageserver.Image, parameters imageserver.Parameters) error {
				return imageserver.NewError("error")
			},
		},
		KeyGenerator: KeyGeneratorFunc(func(parameters imageserver.Parameters) string {
			return "test"
		}),
	}
	_, err := s.Get(imageserver.Parameters{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestNewParametersHashKeyGeneratorFunc(t *testing.T) {
	g := NewParametersHashKeyGeneratorFunc(sha256.New)
	parameters := imageserver.Parameters{
		"foo": "bar",
	}
	g.GetKey(parameters)
}
