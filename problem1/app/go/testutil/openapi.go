package testutil

import (
	"fmt"
	"sync"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

type OpenAPITester interface {
	CheckValidationBySchema(t *testing.T, schemaName string, target any) error
	ValidateBySchema(t *testing.T, schemaName string, target any)
}

type openAPITester struct {
	doc *openapi3.T
}

var testerMap = sync.Map{}

func NewOpenAPITester(t *testing.T, specPath string) OpenAPITester {
	t.Helper()

	if v, ok := testerMap.Load(specPath); ok {
		return v.(OpenAPITester)
	}

	doc, err := openapi3.NewLoader().LoadFromFile(specPath)
	if err != nil {
		t.Fatal(err)
	}

	ot := &openAPITester{doc: doc}
	testerMap.Store(specPath, ot)

	return ot
}

func (o *openAPITester) CheckValidationBySchema(t *testing.T, schemaName string, target any) error {
	t.Helper()
	var vv any
	I2V(t, target, &vv)

	s := o.doc.Components.Schemas[schemaName]
	if s == nil {
		return fmt.Errorf("schema %s does not exist", schemaName)
	}
	if s.Value == nil {
		return fmt.Errorf("schema %s does not have Value", schemaName)
	}

	return s.Value.VisitJSON(vv)
}

func (o *openAPITester) ValidateBySchema(t *testing.T, schemaName string, target any) {
	t.Helper()
	if err := o.CheckValidationBySchema(t, schemaName, target); err != nil {
		t.Fatal(err)
	}
}
