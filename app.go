package goroot

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// App is the main application.
type App struct {
	name     string
	version  uint
	handlers []Handler
}

// GetName returns the name of the app.
func (app *App) GetName() string {
	return app.name
}

// GetVersion returns the version of the app.
func (app *App) GetVersion() uint {
	return app.version
}

// RegisterHandlers registers a handler to the app.
func (app *App) RegisterHandlers(handler ...Handler) error {
	for _, h := range handler {
		if !inArray(h, app.handlers) {
			app.handlers = append(app.handlers, h)
		}
	}

	return nil
}

func (app *App) Run(grpcAddress string) error {
	// generate protobuf
	goPackage := getGoPackagePathFromRuntime()
	protoContent, err := GenerateProtoContent(
		goPackage,
		getProtoPackagePathFromGoPackagePath(goPackage),
		app.GetName(),
		app.GetVersion(),
		app.handlers,
		ExtractAllStructsFromHandlers(app.handlers),
	)
	if err != nil {
		return err
	}

	// generate grpc
	protoFileName := strings.ReplaceAll(strings.ToLower(app.GetName()), " ", "_") + ".proto"
	err = writeToFile(protoFileName, protoContent)
	if err != nil {
		return err
	}

	// compile proto to go
	err = compileProtoToGo(protoFileName)
	if err != nil {
		return err
	}

	// generate protoTypes transformers (from proto to handler request/response types)
	// register handlers to grpc server
	// start grpc server
	// TODO
	return nil
}

// compileProtoToGo compiles a .proto file to go files
func compileProtoToGo(name string) error {
	// Get the absolute path of the .proto file.
	absPath, err := filepath.Abs(name)
	if err != nil {
		return err
	}

	// Check if the .proto file exists.
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", absPath)
	}

	// Prepare the protoc command.
	cmd := exec.Command("protoc", fmt.Sprintf("--go-grpc_out=. --go_out=. ./%s.proto", absPath))

	// Run the protoc command and capture the output and error streams.
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error compiling .proto file: %s", err.Error())
	}

	fmt.Printf("Compilation output:\n%s\n", output)
	return nil
}

func inArray(h Handler, handlers []Handler) bool {
	for _, handler := range handlers {
		if getTypeName(h) == getTypeName(handler) {
			return true
		}
	}
	return false
}

// NewApp creates a new app.
func NewApp(name string, version uint) *App {
	return &App{
		name:    name,
		version: version,
	}
}

func getGoPackagePathFromRuntime() string {
	_, file, _, _ := runtime.Caller(0)
	return path.Base(filepath.Dir(file))
}

func getProtoPackagePathFromGoPackagePath(goPackagePath string) string {
	return strings.ReplaceAll(goPackagePath, "/", ".")
}
