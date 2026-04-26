package interactive

import (
	"context"
	"testing"

	"fileforge/internal/convert"
)

type fakeConverter struct {
	lastReq convert.ImageConvertRequest
	called  bool
	err     error
}

func (f *fakeConverter) Convert(_ context.Context, req convert.ImageConvertRequest) error {
	f.called = true
	f.lastReq = req
	return f.err
}

func TestConvertWizardExecuteUsesService(t *testing.T) {
	svc := &fakeConverter{}
	wizard := NewConvertWizard(svc, nil)

	err := wizard.Execute(context.Background(), ConvertInput{
		InputPath:  `"./input.png"`,
		ToFormat:   "jpg",
		OutputPath: `'./output.jpg'`,
		Force:      true,
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !svc.called {
		t.Fatal("expected converter to be called")
	}
	if svc.lastReq.InputPath != NormalizePath(`"./input.png"`) {
		t.Fatalf("unexpected input path: %q", svc.lastReq.InputPath)
	}
	if svc.lastReq.OutputPath != NormalizePath(`'./output.jpg'`) {
		t.Fatalf("unexpected output path: %q", svc.lastReq.OutputPath)
	}
	if !svc.lastReq.Force {
		t.Fatal("expected force to be passed through")
	}
}
