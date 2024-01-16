package methodhandler

import (
	"bytes"
	"encoding/json"
	"html/template"
	"os"

	"github.com/TheLeeeo/grpc-hole/service"
	"github.com/hashicorp/go-hclog"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/grpc"
)

// A static handler returns a response based on a template file.
type staticHandler struct {
	method *desc.MethodDescriptor
	lg     hclog.Logger
}

func NewStaticHandler(method *desc.MethodDescriptor, logger hclog.Logger) Handler {
	return &staticHandler{
		method: method,
		lg:     logger,
	}
}

func (h *staticHandler) Name() string {
	return h.method.GetName()
}

func (h *staticHandler) Handle(stream grpc.ServerStream) error {
	inputMsg := dynamic.NewMessage(h.method.GetInputType())
	if err := stream.RecvMsg(inputMsg); err != nil {
		return err
	}

	inputJSON, _ := inputMsg.MarshalJSON()
	h.lg.Info("Received request", "Method", h.method.GetName(), "Input", string(inputJSON))

	outType := h.method.GetOutputType()
	var out *dynamic.Message

	respTemplate, err := service.LoadResponse(h.method.GetService().GetName(), h.method.GetName())
	if err != nil {
		// If the error is something else than "file not found", return it.
		if !os.IsNotExist(err) {
			return err
		}
		out = CreatePopulatedMessage(outType, 0)
	} else {
		var inputMap map[string]any
		if err := json.Unmarshal(inputJSON, &inputMap); err != nil {
			return err
		}

		tmpl, err := template.New("inputParser").Funcs(funcMap).Parse(string(respTemplate))
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, inputMap)
		if err != nil {
			return err
		}

		out = dynamic.NewMessage(outType)
		if err := out.UnmarshalJSON(buf.Bytes()); err != nil {
			h.lg.Error("Failed to unmarshal json", "Method", h.method.GetName(), "Error", err)
			return err
		}
	}

	return stream.SendMsg(out)
}
