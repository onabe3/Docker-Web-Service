package main

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type User struct {
	Name  string `json:"name" xml:"name"`
	Email string `json:"email" xml:"email"`
}

type Result struct {
	Error   string `json:"error" xml:"error"`
	Message string `json:"message" xml:"message"`
}

func tarBuildContext(dockerfilePath string, dockerfileName string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	// Dockerfileの内容を読み込む
	dockerfile, err := os.Open(dockerfilePath)
	if err != nil {
		return nil, err
	}
	defer dockerfile.Close()

	dockerfileContent, err := io.ReadAll(dockerfile)
	if err != nil {
		return nil, err
	}

	// tarにDockerfileを追加する
	tarHeader := &tar.Header{
		Name: dockerfileName,
		Size: int64(len(dockerfileContent)),
	}
	if err := tw.WriteHeader(tarHeader); err != nil {
		return nil, err
	}
	if _, err := tw.Write(dockerfileContent); err != nil {
		return nil, err
	}

	return buf, nil
}

func create_Container(containerName string, osName string) Result {
	// Docker clientの初期化
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return Result{Error: fmt.Sprintf("Initialize Docker client error: %v", err),
			Message: ""}
	}

	// Dockerfileの内容を定義
	dockerfileContent := `FROM ` + osName + `:latest`
	// Dockerfileに内容を書き込む
	dockerfile, err := os.Create("Dockerfile")
	if err != nil {
		return Result{Error: fmt.Sprintf("Failed to create Dockerfile: %v", err),
			Message: ""}
	}
	_, err = dockerfile.WriteString(dockerfileContent)
	if err != nil {
		return Result{Error: fmt.Sprintf("Failed to write Dockerfile: %v", err),
			Message: ""}
	}
	dockerfile.Close()

	// Dockerfileをtarアーカイブに変換
	buildContext, err := tarBuildContext("Dockerfile", "Dockerfile")
	if err != nil {
		return Result{Error: fmt.Sprintf("Failed to create tar build context: %v", err),
			Message: ""}
	}

	// イメージ名を定義
	imageName := "my-" + osName + "-image"

	// Dockerイメージをビルドする
	buildResp, err := cli.ImageBuild(
		context.Background(),
		io.NopCloser(buildContext),
		types.ImageBuildOptions{
			Tags:       []string{imageName},
			Dockerfile: "Dockerfile",
		},
	)
	if err != nil {
		return Result{Error: fmt.Sprintf("Image build error: %v", err),
			Message: ""}
	}
	defer buildResp.Body.Close()

	// Buildのログを表示する
	response, err := io.ReadAll(buildResp.Body)
	if err != nil {
		return Result{Error: fmt.Sprintf("Failed to read build response: %v", err),
			Message: ""}
	}
	fmt.Println(string(response))

	// イメージを元に新しいコンテナを作成する
	config := &container.Config{
		Image: imageName,
		Cmd:   []string{"/bin/bash"},
		Tty:   true,
	}

	resp, err := cli.ContainerCreate(context.Background(), config, nil, nil, nil, containerName)
	if err != nil {
		return Result{Error: fmt.Sprintf("Container create error: %v", err),
			Message: ""}
	}

	// コンテナを起動する
	if err := cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		return Result{Error: fmt.Sprintf("Container start error: %v", err),
			Message: ""}
	}

	return Result{Error: "", Message: fmt.Sprintf("Container %s started\n", resp.ID)}
}

func main() {

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// アクセスを許可するオリジンを指定
		AllowOrigins: []string{"http://localhost:3000"},
		// アクセスを許可するメソッドを指定
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		// アクセスを許可するヘッダーを指定
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization, "X-CSRF-Header"},
		AllowCredentials: true,
	}))
	u := &User{
		Name:  "John",
		Email: "jon@labstack.com",
	}
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, u)
	})
	e.POST("/create/:createOs/:containerName", func(c echo.Context) error {
		createOs := c.Param("createOs")
		createContainerName := c.Param("containerName")
		result := create_Container(createContainerName, createOs)
		return c.JSON(http.StatusOK, result)
	})
	e.Logger.Fatal(e.Start("127.0.0.1:1323"))
}
