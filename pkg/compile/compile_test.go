package compile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/go-util/assertfile"
	"github.com/kalo-build/plugin-openapi-morphe/internal/testutils"
	"github.com/kalo-build/plugin-openapi-morphe/pkg/compile"
	"github.com/kalo-build/plugin-openapi-morphe/pkg/compile/cfg"
)

type CompileTestSuite struct {
	assertfile.FileSuite

	TestDirPath            string
	TestGroundTruthDirPath string
}

func TestCompileTestSuite(t *testing.T) {
	suite.Run(t, new(CompileTestSuite))
}

func (suite *CompileTestSuite) SetupTest() {
	suite.TestDirPath = testutils.GetTestDirPath()
	suite.TestGroundTruthDirPath = filepath.Join(suite.TestDirPath, "ground-truth", "compile-minimal")
}

func (suite *CompileTestSuite) TearDownTest() {
	suite.TestDirPath = ""
}

func (suite *CompileTestSuite) TestOpenAPIToMorphe() {
	workingDirPath := filepath.Join(suite.TestDirPath, "working")
	suite.Nil(os.Mkdir(workingDirPath, 0755))
	defer os.RemoveAll(workingDirPath)

	config := cfg.OpenAPIToMorpheConfig{
		InputDir:   filepath.Join(suite.TestDirPath, "input"),
		OutputPath: workingDirPath,
	}

	compileErr := compile.OpenAPIToMorphe(config)
	suite.NoError(compileErr)

	structuresDirPath := filepath.Join(workingDirPath, "structures")
	gtStructuresDirPath := filepath.Join(suite.TestGroundTruthDirPath, "structures")
	suite.DirExists(structuresDirPath)

	suite.FileExists(filepath.Join(structuresDirPath, "customer.str"))
	suite.FileEquals(
		filepath.Join(structuresDirPath, "customer.str"),
		filepath.Join(gtStructuresDirPath, "customer.str"),
	)

	suite.FileExists(filepath.Join(structuresDirPath, "customer_create.str"))
	suite.FileEquals(
		filepath.Join(structuresDirPath, "customer_create.str"),
		filepath.Join(gtStructuresDirPath, "customer_create.str"),
	)

	suite.FileExists(filepath.Join(structuresDirPath, "invoice.str"))
	suite.FileEquals(
		filepath.Join(structuresDirPath, "invoice.str"),
		filepath.Join(gtStructuresDirPath, "invoice.str"),
	)

	suite.FileExists(filepath.Join(structuresDirPath, "payment.str"))
	suite.FileEquals(
		filepath.Join(structuresDirPath, "payment.str"),
		filepath.Join(gtStructuresDirPath, "payment.str"),
	)

	suite.FileExists(filepath.Join(structuresDirPath, "payment_create.str"))
	suite.FileEquals(
		filepath.Join(structuresDirPath, "payment_create.str"),
		filepath.Join(gtStructuresDirPath, "payment_create.str"),
	)

	enumsDirPath := filepath.Join(workingDirPath, "enums")
	gtEnumsDirPath := filepath.Join(suite.TestGroundTruthDirPath, "enums")
	suite.DirExists(enumsDirPath)

	suite.FileExists(filepath.Join(enumsDirPath, "invoice_status.enum"))
	suite.FileEquals(
		filepath.Join(enumsDirPath, "invoice_status.enum"),
		filepath.Join(gtEnumsDirPath, "invoice_status.enum"),
	)

	suite.FileExists(filepath.Join(enumsDirPath, "payment_method.enum"))
	suite.FileEquals(
		filepath.Join(enumsDirPath, "payment_method.enum"),
		filepath.Join(gtEnumsDirPath, "payment_method.enum"),
	)

	suite.FileExists(filepath.Join(enumsDirPath, "payment_create_method.enum"))
	suite.FileEquals(
		filepath.Join(enumsDirPath, "payment_create_method.enum"),
		filepath.Join(gtEnumsDirPath, "payment_create_method.enum"),
	)

	suite.FileExists(filepath.Join(enumsDirPath, "payment_status.enum"))
	suite.FileEquals(
		filepath.Join(enumsDirPath, "payment_status.enum"),
		filepath.Join(gtEnumsDirPath, "payment_status.enum"),
	)
}
