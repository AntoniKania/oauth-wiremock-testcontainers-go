package testutil

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/wiremock/go-wiremock"
	"github.com/wiremock/wiremock-testcontainers-go"
	"os"
	"strings"
	"testing"
)

type WireMock struct {
	Client  *wiremock.Client
	Address string
}

func SetupWireMock(t *testing.T) *WireMock {
	wmContainer, ctx := setupWireMockTestContainer(t)

	ip, _ := wmContainer.Container.Host(ctx)
	port, _ := wmContainer.Container.MappedPort(ctx, "8080")
	wmAdd := fmt.Sprintf("http://%s:%s", ip, port.Port())

	wmClient := wmContainer.Client
	stubWellKnown(wmClient, port)

	return &WireMock{
		Client:  wmClient,
		Address: wmAdd,
	}
}

func setupWireMockTestContainer(t *testing.T) (*testcontainers_wiremock.WireMockContainer, context.Context) {
	ctx := context.Background()
	container, err := testcontainers_wiremock.RunDefaultContainerAndStopOnCleanup(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	return container, ctx
}

func ReadFile(fileName string) string {
	b, _ := os.ReadFile("../testdata/" + fileName)
	return string(b)
}

func stubWellKnown(wm *wiremock.Client, port nat.Port) {
	portPlaceholder := "${PORT}"
	wellKnownResponse := strings.ReplaceAll(ReadFile("well-known-response.json"), portPlaceholder, port.Port())

	wm.StubFor(wiremock.Get(wiremock.URLPathEqualTo("/.well-known/openid-configuration")).
		WillReturnResponse(wiremock.NewResponse().WithStatus(200).WithBody(wellKnownResponse)))
}


