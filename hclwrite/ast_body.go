package hclwrite

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

type Body struct {
	inTree

	items nodeSet
}

func newBody() *Body {
	return &Body{
		inTree: newInTree(),
		items:  newNodeSet(),
	}
}

func (b *Body) appendItem(c nodeContent) *node {
	nn := b.children.Append(c)
	b.items.Add(nn)
	return nn
}

func (b *Body) appendItemNode(nn *node) *node {
	nn.assertUnattached()
	b.children.AppendNode(nn)
	b.items.Add(nn)
	return nn
}

func (b *Body) AppendUnstructuredTokens(ts Tokens) {
	b.inTree.children.Append(ts)
}

// GetAttribute returns the attribute from the body that has the given name,
// or returns nil if there is currently no matching attribute.
func (b *Body) GetAttribute(name string) *Attribute {
	for n := range b.items {
		if attr, isAttr := n.content.(*Attribute); isAttr {
			nameObj := attr.name.content.(*identifier)
			if nameObj.hasName(name) {
				// We've found it!
				return attr
			}
		}
	}

	return nil
}

// SetAttributeValue either replaces the expression of an existing attribute
// of the given name or adds a new attribute definition to the end of the block.
//
// The value is given as a cty.Value, and must therefore be a literal. To set
// a variable reference or other traversal, use SetAttributeTraversal.
//
// The return value is the attribute that was either modified in-place or
// created.
func (b *Body) SetAttributeValue(name string, val cty.Value) *Attribute {
	attr := b.GetAttribute(name)
	expr := NewExpressionLiteral(val)
	if attr != nil {
		attr.expr = attr.expr.ReplaceWith(expr)
	} else {
		attr := newAttribute()
		attr.init(name, expr)
		b.appendItem(attr)
	}
	return attr
}

// SetAttributeTraversal either replaces the expression of an existing attribute
// of the given name or adds a new attribute definition to the end of the body.
//
// The new expression is given as a hcl.Traversal, which must be an absolute
// traversal. To set a literal value, use SetAttributeValue.
//
// The return value is the attribute that was either modified in-place or
// created.
func (b *Body) SetAttributeTraversal(name string, traversal hcl.Traversal) *Attribute {
	panic("Body.SetAttributeTraversal not yet implemented")
}

// AppendBlock appends a new nested block to the end of the receiving body.
//
// If blankLine is set, an additional empty line is added before the block
// for separation. Usual HCL style suggests that we group together blocks of
// the same type without intervening blank lines and then put blank lines
// between blocks of different types. In some languages, some different block
// types may be conceptually related and so may still be grouped together.
// It is the caller's responsibility to respect the usual conventions of the
// language being generated.
func (b *Body) AppendBlock(typeName string, labels []string, blankLine bool) *Block {
	block := newBlock()
	block.init(typeName, labels)
	if blankLine {
		b.AppendUnstructuredTokens(Tokens{
			{
				Type:  hclsyntax.TokenNewline,
				Bytes: []byte{'\n'},
			},
		})
	}
	b.appendItem(block)
	return block
}
