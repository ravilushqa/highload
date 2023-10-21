package e2e

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

var composeStack tc.ComposeStack

func TestMain(m *testing.M) {
	var status int
	// recover from panic if one occurred. Set exit code to 1 if there was a panic.
	defer func() {
		if r := recover(); r != nil {
			status = 1
			teardown(&status)
		}
	}()
	setup()
	defer teardown(&status)

	if st := m.Run(); st > status {
		status = st
	}
}

func setup() {
	var err error
	composeStack, err = tc.NewDockerCompose("../docker-compose-test.yml")
	if err != nil {
		fmt.Printf("\033[1;31m%s: %v\033[0m", "> Compose stack setup failed\n", err)
		os.Exit(1)
	}

	err = composeStack.
		WaitForService("app", wait.ForListeningPort("8080/tcp").WithStartupTimeout(10*time.Second)).
		WaitForService("users-api", wait.ForLog("starting grpc server").WithStartupTimeout(10*time.Second)).
		Up(context.Background(), tc.Wait(true))
	if err != nil {
		fmt.Printf("\033[1;31m%s: %v\033[0m", "> Compose stack start failed\n", err)
		code := 1
		teardown(&code)
	}
	time.Sleep(1 * time.Second) // ensure that all services are up
}

func teardown(i *int) {
	err := composeStack.Down(context.Background(), tc.RemoveImagesLocal, tc.RemoveOrphans(true), tc.RemoveVolumes(true))
	if err != nil {
		fmt.Printf("\033[1;31m%s: %v\033[0m", "> Teardown failed\n", err)
	}

	os.Exit(*i)
}
