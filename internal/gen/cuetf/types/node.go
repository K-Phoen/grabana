package types

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"github.com/K-Phoen/grabana/internal/gen/cuetf/internal/utils"
)

type Node struct {
	Name             string
	Kind             cue.Kind
	SubKind          cue.Kind // For list only, kind of its elements
	IsMap            bool
	IsDisjunction    bool
	Optional         bool
	Default          string
	Doc              string
	Children         []Node
	Parent           *Node
	DisjunctionKinds []cue.Kind
}

func (n *Node) TerraformModelField(structName string) string {
	kind := kindMappings[n.Kind]
	subKind := kindMappings[n.SubKind]
	typeStr := ""
	switch true {
	case n.Kind == cue.ListKind && n.SubKind == cue.StructKind:
		typeStr = "[]" + structName + "_" + utils.Title(n.Name)
	case n.Kind == cue.ListKind && subKind != nil:
		typeStr = "types.List"
	case n.Kind == cue.StructKind:
		if n.IsMap {
			typeStr = "map[string]"
		}
		typeStr += structName + "_" + utils.Title(n.Name)
		if n.Optional {
			typeStr = "*" + typeStr
		}
	default:
		typeStr = "types." + kind.terraformType
	}

	return fmt.Sprintf("%s %s `tfsdk:\"%s\"`", utils.ToCamelCase(n.Name), typeStr, utils.ToSnakeCase(n.Name))
}

func (n *Node) JSONModelField() string {
	kind := kindMappings[n.Kind]
	subKind := kindMappings[n.SubKind]
	golangType := ""
	switch true {
	case n.Kind == cue.ListKind && n.SubKind == cue.StructKind:
		golangType = "[]interface{}"
	case n.Kind == cue.ListKind && subKind != nil:
		golangType = "[]" + subKind.golangType
	case n.Kind == cue.StructKind || len(n.DisjunctionKinds) > 0:
		golangType = "interface{}"
	default:
		golangType = kind.golangType
	}

	omitStr := ""
	if n.Optional {
		if !strings.HasPrefix(golangType, "[]") && golangType != "interface{}" {
			golangType = "*" + golangType
		}
		omitStr = ",omitempty"
	}

	return fmt.Sprintf("%s %s `json:\"%s%s\"`", utils.ToCamelCase(n.Name), golangType, n.Name, omitStr)
}

func (n *Node) TerraformType() string {
	if kindMappings[n.Kind] == nil {
		return ""
	}

	return kindMappings[n.Kind].terraformType
}

func (n *Node) SubTerraformType() string {
	if kindMappings[n.SubKind] == nil {
		return ""
	}

	return kindMappings[n.SubKind].terraformType
}

func (n *Node) subGolangType() string {
	if kindMappings[n.SubKind] == nil {
		return ""
	}

	return kindMappings[n.SubKind].golangType
}

func (n *Node) terraformFunc() string {
	if kindMappings[n.Kind] == nil {
		return ""
	}

	terraformFunc := kindMappings[n.Kind].terraformFunc
	if n.Optional {
		return terraformFunc + "Pointer()"
	}

	return terraformFunc + "()"
}

func (n *Node) subTerraformFunc() string {
	if kindMappings[n.SubKind] == nil {
		return ""
	}

	return kindMappings[n.SubKind].terraformFunc + "()"
}

func (n *Node) IsGenerated() bool {
	return kindMappings[n.Kind] != nil ||
		(n.Kind == cue.ListKind && kindMappings[n.SubKind] != nil) ||
		(n.Kind == cue.ListKind && n.SubKind == cue.StructKind) ||
		(n.Kind == cue.StructKind)
}
