package diff

import (
	"bytes"
	"fmt"
	"image/color"
	"io"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	chromastyles "github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/compat"

	"github.com/avgvstvs96/differential/pkg/styles"
	"github.com/avgvstvs96/differential/pkg/theme"
)

// SyntaxHighlight applies syntax highlighting to text based on file extension
func SyntaxHighlight(w io.Writer, source, fileName, formatter string, bg color.Color) error {
	t := theme.CurrentTheme()

	// Determine the language lexer to use
	l := lexers.Match(fileName)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	// Get the formatter
	f := formatters.Get(formatter)
	if f == nil {
		f = formatters.Fallback
	}

	// Dynamic theme based on current theme values
	syntaxThemeXml := fmt.Sprintf(`
	<style name="opencode-theme">
	<!-- Base colors -->
	<entry type="Background" style="bg:%s"/>
	<entry type="Text" style="%s"/>
	<entry type="Other" style="%s"/>
	<entry type="Error" style="%s"/>
	<!-- Keywords -->
	<entry type="Keyword" style="%s"/>
	<entry type="KeywordConstant" style="%s"/>
	<entry type="KeywordDeclaration" style="%s"/>
	<entry type="KeywordNamespace" style="%s"/>
	<entry type="KeywordPseudo" style="%s"/>
	<entry type="KeywordReserved" style="%s"/>
	<entry type="KeywordType" style="%s"/>
	<!-- Names -->
	<entry type="Name" style="%s"/>
	<entry type="NameAttribute" style="%s"/>
	<entry type="NameBuiltin" style="%s"/>
	<entry type="NameBuiltinPseudo" style="%s"/>
	<entry type="NameClass" style="%s"/>
	<entry type="NameConstant" style="%s"/>
	<entry type="NameDecorator" style="%s"/>
	<entry type="NameEntity" style="%s"/>
	<entry type="NameException" style="%s"/>
	<entry type="NameFunction" style="%s"/>
	<entry type="NameLabel" style="%s"/>
	<entry type="NameNamespace" style="%s"/>
	<entry type="NameOther" style="%s"/>
	<entry type="NameTag" style="%s"/>
	<entry type="NameVariable" style="%s"/>
	<entry type="NameVariableClass" style="%s"/>
	<entry type="NameVariableGlobal" style="%s"/>
	<entry type="NameVariableInstance" style="%s"/>
	<!-- Literals -->
	<entry type="Literal" style="%s"/>
	<entry type="LiteralDate" style="%s"/>
	<entry type="LiteralString" style="%s"/>
	<entry type="LiteralStringBacktick" style="%s"/>
	<entry type="LiteralStringChar" style="%s"/>
	<entry type="LiteralStringDoc" style="%s"/>
	<entry type="LiteralStringDouble" style="%s"/>
	<entry type="LiteralStringEscape" style="%s"/>
	<entry type="LiteralStringHeredoc" style="%s"/>
	<entry type="LiteralStringInterpol" style="%s"/>
	<entry type="LiteralStringOther" style="%s"/>
	<entry type="LiteralStringRegex" style="%s"/>
	<entry type="LiteralStringSingle" style="%s"/>
	<entry type="LiteralStringSymbol" style="%s"/>
	<!-- Numbers -->
	<entry type="LiteralNumber" style="%s"/>
	<entry type="LiteralNumberBin" style="%s"/>
	<entry type="LiteralNumberFloat" style="%s"/>
	<entry type="LiteralNumberHex" style="%s"/>
	<entry type="LiteralNumberInteger" style="%s"/>
	<entry type="LiteralNumberIntegerLong" style="%s"/>
	<entry type="LiteralNumberOct" style="%s"/>
	<!-- Operators -->
	<entry type="Operator" style="%s"/>
	<entry type="OperatorWord" style="%s"/>
	<entry type="Punctuation" style="%s"/>
	<!-- Comments -->
	<entry type="Comment" style="%s"/>
	<entry type="CommentHashbang" style="%s"/>
	<entry type="CommentMultiline" style="%s"/>
	<entry type="CommentSingle" style="%s"/>
	<entry type="CommentSpecial" style="%s"/>
	<entry type="CommentPreproc" style="%s"/>
	<!-- Generic styles -->
	<entry type="Generic" style="%s"/>
	<entry type="GenericDeleted" style="%s"/>
	<entry type="GenericEmph" style="italic %s"/>
	<entry type="GenericError" style="%s"/>
	<entry type="GenericHeading" style="bold %s"/>
	<entry type="GenericInserted" style="%s"/>
	<entry type="GenericOutput" style="%s"/>
	<entry type="GenericPrompt" style="%s"/>
	<entry type="GenericStrong" style="bold %s"/>
	<entry type="GenericSubheading" style="bold %s"/>
	<entry type="GenericTraceback" style="%s"/>
	<entry type="GenericUnderline" style="underline"/>
	<entry type="TextWhitespace" style="%s"/>
</style>
`,
		getChromaColor(t.BackgroundPanel()), // Background
		getChromaColor(t.Text()),            // Text
		getChromaColor(t.Text()),            // Other
		getChromaColor(t.Error()),           // Error

		getChromaColor(t.SyntaxKeyword()), // Keyword
		getChromaColor(t.SyntaxKeyword()), // KeywordConstant
		getChromaColor(t.SyntaxKeyword()), // KeywordDeclaration
		getChromaColor(t.SyntaxKeyword()), // KeywordNamespace
		getChromaColor(t.SyntaxKeyword()), // KeywordPseudo
		getChromaColor(t.SyntaxKeyword()), // KeywordReserved
		getChromaColor(t.SyntaxType()),    // KeywordType

		getChromaColor(t.Text()),           // Name
		getChromaColor(t.SyntaxVariable()), // NameAttribute
		getChromaColor(t.SyntaxType()),     // NameBuiltin
		getChromaColor(t.SyntaxVariable()), // NameBuiltinPseudo
		getChromaColor(t.SyntaxType()),     // NameClass
		getChromaColor(t.SyntaxVariable()), // NameConstant
		getChromaColor(t.SyntaxFunction()), // NameDecorator
		getChromaColor(t.SyntaxVariable()), // NameEntity
		getChromaColor(t.SyntaxType()),     // NameException
		getChromaColor(t.SyntaxFunction()), // NameFunction
		getChromaColor(t.Text()),           // NameLabel
		getChromaColor(t.SyntaxType()),     // NameNamespace
		getChromaColor(t.SyntaxVariable()), // NameOther
		getChromaColor(t.SyntaxKeyword()),  // NameTag
		getChromaColor(t.SyntaxVariable()), // NameVariable
		getChromaColor(t.SyntaxVariable()), // NameVariableClass
		getChromaColor(t.SyntaxVariable()), // NameVariableGlobal
		getChromaColor(t.SyntaxVariable()), // NameVariableInstance

		getChromaColor(t.SyntaxString()), // Literal
		getChromaColor(t.SyntaxString()), // LiteralDate
		getChromaColor(t.SyntaxString()), // LiteralString
		getChromaColor(t.SyntaxString()), // LiteralStringBacktick
		getChromaColor(t.SyntaxString()), // LiteralStringChar
		getChromaColor(t.SyntaxString()), // LiteralStringDoc
		getChromaColor(t.SyntaxString()), // LiteralStringDouble
		getChromaColor(t.SyntaxString()), // LiteralStringEscape
		getChromaColor(t.SyntaxString()), // LiteralStringHeredoc
		getChromaColor(t.SyntaxString()), // LiteralStringInterpol
		getChromaColor(t.SyntaxString()), // LiteralStringOther
		getChromaColor(t.SyntaxString()), // LiteralStringRegex
		getChromaColor(t.SyntaxString()), // LiteralStringSingle
		getChromaColor(t.SyntaxString()), // LiteralStringSymbol

		getChromaColor(t.SyntaxNumber()), // LiteralNumber
		getChromaColor(t.SyntaxNumber()), // LiteralNumberBin
		getChromaColor(t.SyntaxNumber()), // LiteralNumberFloat
		getChromaColor(t.SyntaxNumber()), // LiteralNumberHex
		getChromaColor(t.SyntaxNumber()), // LiteralNumberInteger
		getChromaColor(t.SyntaxNumber()), // LiteralNumberIntegerLong
		getChromaColor(t.SyntaxNumber()), // LiteralNumberOct

		getChromaColor(t.SyntaxOperator()),    // Operator
		getChromaColor(t.SyntaxKeyword()),     // OperatorWord
		getChromaColor(t.SyntaxPunctuation()), // Punctuation

		getChromaColor(t.SyntaxComment()), // Comment
		getChromaColor(t.SyntaxComment()), // CommentHashbang
		getChromaColor(t.SyntaxComment()), // CommentMultiline
		getChromaColor(t.SyntaxComment()), // CommentSingle
		getChromaColor(t.SyntaxComment()), // CommentSpecial
		getChromaColor(t.SyntaxKeyword()), // CommentPreproc

		getChromaColor(t.Text()),      // Generic
		getChromaColor(t.Error()),     // GenericDeleted
		getChromaColor(t.Text()),      // GenericEmph
		getChromaColor(t.Error()),     // GenericError
		getChromaColor(t.Text()),      // GenericHeading
		getChromaColor(t.Success()),   // GenericInserted
		getChromaColor(t.TextMuted()), // GenericOutput
		getChromaColor(t.Text()),      // GenericPrompt
		getChromaColor(t.Text()),      // GenericStrong
		getChromaColor(t.Text()),      // GenericSubheading
		getChromaColor(t.Error()),     // GenericTraceback
		getChromaColor(t.Text()),      // TextWhitespace
	)

	r := strings.NewReader(syntaxThemeXml)
	style := chroma.MustNewXMLStyle(r)

	// Modify the style to use the provided background
	s, err := style.Builder().Transform(
		func(t chroma.StyleEntry) chroma.StyleEntry {
			if _, ok := bg.(lipgloss.NoColor); ok {
				// For NoColor, explicitly clear the background
				t.Background = chroma.Colour(0)
				return t
			}
			r, g, b, _ := bg.RGBA()
			t.Background = chroma.NewColour(uint8(r>>8), uint8(g>>8), uint8(b>>8))
			return t
		},
	).Build()
	if err != nil {
		s = chromastyles.Fallback
	}

	// Tokenize and format
	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}

	return f.Format(w, s, it)
}

// getColor returns the appropriate hex color string based on terminal background
func getColor(adaptiveColor compat.AdaptiveColor) *string {
	return styles.AdaptiveColorToString(adaptiveColor)
}

func getChromaColor(adaptiveColor compat.AdaptiveColor) string {
	color := styles.AdaptiveColorToString(adaptiveColor)
	if color == nil {
		return ""
	}
	return *color
}

// highlightLine applies syntax highlighting to a single line
func highlightLine(fileName string, line string, bg color.Color) string {
	var buf bytes.Buffer
	err := SyntaxHighlight(&buf, line, fileName, "terminal16m", bg)
	if err != nil {
		return line
	}
	return buf.String()
}
