package wasm

import (
	"context"
	_ "embed"
	"fmt"
	sloth "github.com/slok/sloth/pkg/prometheus/api/v1"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/assemblyscript"
	"github.com/tfadeyi/sloth-simple-comments/internal/logging"
	"log"
	"os"
)

type parser struct {
	Spec              *sloth.Spec
	GeneralInfoSource string
	IncludedDirs      []string
	Logger            *logging.Logger
}

const (
	defaultSourceFile = "main.go"
)

// asWasm compiled using `npm install && npm run build`
//
//go:embed modules/index.wasm
var asWasm []byte

// newParser client parser performs all checks at initialization time
func newParser(logger *logging.Logger, dirs ...string) *parser {
	return &parser{
		Spec: &sloth.Spec{
			Version: sloth.Version,
			Service: "",
			Labels:  nil,
			SLOs:    nil,
		},
		GeneralInfoSource: defaultSourceFile,
		IncludedDirs:      dirs,
		Logger:            logger,
	}
}

func (p parser) Parse(ctx context.Context) (*sloth.Spec, error) {
	// collect all aloe error comments from packages and add them to the spec struct
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)
	// Instantiate a module implementing functions used by AssemblyScript.
	// Thrown errors will be logged to os.Stderr
	_, err := assemblyscript.Instantiate(ctx, r)
	if err != nil {
		log.Panicln(err)
	}

	// Instantiate a WebAssembly module that imports the "abort" and "trace"
	// functions defined by assemblyscript.Instantiate and exports functions
	// we'll use in this example.
	mod, err := r.InstantiateWithConfig(ctx, asWasm,
		// Override the default module config that discards stdout and stderr.
		wazero.NewModuleConfig().WithStdout(os.Stdout).WithStderr(os.Stderr))
	if err != nil {
		log.Panicln(err)
	}

	// Get references to WebAssembly functions we'll use in this example.
	helloWorld := mod.ExportedFunction("hello_world")
	ptrSize, err := helloWorld.Call(ctx)
	if err != nil {
		log.Panicln(err)
	}

	greetingPtr := uint32(ptrSize[0])
	greetingSize := uint32(4294)
	//
	// The pointer is a linear memory offset, which is where we write the name.
	if bytes, ok := mod.Memory().Read(greetingPtr, greetingSize); !ok {
		log.Panicf("Memory.Read(%d, %d) out of range of memory size %d",
			greetingPtr, greetingSize, mod.Memory().Size())
	} else {
		fmt.Println("go >>", string(bytes))
	}

	//buf := []byte{}

	//for i := 0; i < 5; i++ {
	//	if byt, ok := mod.Memory().ReadByte(greetingPtr, uint32(ptrSize[0])); !ok {
	//		log.Panicf("Memory.Read(%d, %d) out of range of memory size %d",
	//			greetingPtr, greetingSize, mod.Memory().Size())
	//	} else {
	//		fmt.Println("go >>", string(byt))
	//		//buf = append(buf, byt...)
	//		//fmt.Println("go >>", string(bytes))
	//		//n := bytes.IndexByte(byt, 0)
	//		//if n < 0 {
	//		//	// Not found!
	//		//	break
	//		//}
	//	}
	//
	//	greetingPtr = greetingPtr + 32
	//}

	fmt.Printf("hello_world returned: %v", greetingSize)

	return p.Spec, nil
}
