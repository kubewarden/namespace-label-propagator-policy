//go:build !wasi
// +build !wasi

package main

import (
	capabilities "github.com/kubewarden/policy-sdk-go/pkg/capabilities"
)

var wapcClient *capabilities.MockWapcClient

func getWapcHost() capabilities.Host {
	return capabilities.Host{
		Client: wapcClient,
	}
}
