package themes

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/lipgloss"
)

// GenerateChromaStyle creates a Chroma style from the current theme
func GenerateChromaStyle() (*chroma.Style, error) {
	t := GetCurrentTheme()
	
	// Convert lipgloss colors to Chroma format
	toChroma := func(c lipgloss.Color) string {
		s := string(c)
		// Remove # prefix if present
		if strings.HasPrefix(s, "#") {
			s = s[1:]
		}
		return s
	}
	
	// Generate Chroma style XML
	styleXML := fmt.Sprintf(`
<style name="pretty-diff">
    <!-- Base -->
    <entry type="Background" style="bg:%s"/>
    <entry type="Text" style="%s"/>
    <entry type="Error" style="%s bold"/>
    
    <!-- Keywords -->
    <entry type="Keyword" style="%s"/>
    <entry type="KeywordConstant" style="%s"/>
    <entry type="KeywordDeclaration" style="%s"/>
    <entry type="KeywordNamespace" style="%s"/>
    <entry type="KeywordPseudo" style="%s"/>
    <entry type="KeywordReserved" style="%s"/>
    <entry type="KeywordType" style="%s"/>
    
    <!-- Names -->
    <entry type="NameAttribute" style="%s"/>
    <entry type="NameBuiltin" style="%s"/>
    <entry type="NameBuiltinPseudo" style="%s"/>
    <entry type="NameClass" style="%s"/>
    <entry type="NameConstant" style="%s"/>
    <entry type="NameDecorator" style="%s"/>
    <entry type="NameEntity" style="%s"/>
    <entry type="NameException" style="%s"/>
    <entry type="NameFunction" style="%s"/>
    <entry type="NameFunctionMagic" style="%s"/>
    <entry type="NameLabel" style="%s"/>
    <entry type="NameNamespace" style="%s"/>
    <entry type="NameOther" style="%s"/>
    <entry type="NameProperty" style="%s"/>
    <entry type="NameTag" style="%s"/>
    <entry type="NameVariable" style="%s"/>
    <entry type="NameVariableClass" style="%s"/>
    <entry type="NameVariableGlobal" style="%s"/>
    <entry type="NameVariableInstance" style="%s"/>
    <entry type="NameVariableMagic" style="%s"/>
    
    <!-- Literals -->
    <entry type="LiteralString" style="%s"/>
    <entry type="LiteralStringAffix" style="%s"/>
    <entry type="LiteralStringBacktick" style="%s"/>
    <entry type="LiteralStringChar" style="%s"/>
    <entry type="LiteralStringDelimiter" style="%s"/>
    <entry type="LiteralStringDoc" style="%s"/>
    <entry type="LiteralStringDouble" style="%s"/>
    <entry type="LiteralStringEscape" style="%s"/>
    <entry type="LiteralStringHeredoc" style="%s"/>
    <entry type="LiteralStringInterpol" style="%s"/>
    <entry type="LiteralStringOther" style="%s"/>
    <entry type="LiteralStringRegex" style="%s"/>
    <entry type="LiteralStringSingle" style="%s"/>
    <entry type="LiteralStringSymbol" style="%s"/>
    
    <entry type="LiteralNumber" style="%s"/>
    <entry type="LiteralNumberBin" style="%s"/>
    <entry type="LiteralNumberFloat" style="%s"/>
    <entry type="LiteralNumberHex" style="%s"/>
    <entry type="LiteralNumberInteger" style="%s"/>
    <entry type="LiteralNumberIntegerLong" style="%s"/>
    <entry type="LiteralNumberOct" style="%s"/>
    
    <!-- Comments -->
    <entry type="Comment" style="%s"/>
    <entry type="CommentHashbang" style="%s"/>
    <entry type="CommentMultiline" style="%s"/>
    <entry type="CommentPreproc" style="%s"/>
    <entry type="CommentPreprocFile" style="%s"/>
    <entry type="CommentSingle" style="%s"/>
    <entry type="CommentSpecial" style="%s"/>
    
    <!-- Operators & Punctuation -->
    <entry type="Operator" style="%s"/>
    <entry type="OperatorWord" style="%s"/>
    <entry type="Punctuation" style="%s"/>
    
    <!-- Generic (for diffs, etc) -->
    <entry type="GenericDeleted" style="%s bg:%s"/>
    <entry type="GenericInserted" style="%s bg:%s"/>
    <entry type="GenericHeading" style="%s bold"/>
    <entry type="GenericSubheading" style="%s bold"/>
    <entry type="GenericStrong" style="bold"/>
    <entry type="GenericEmph" style="italic"/>
</style>`,
		toChroma(t.Background),
		toChroma(t.Text),
		toChroma(t.Error),
		// Keywords
		toChroma(t.SyntaxKeyword),
		toChroma(t.SyntaxKeyword),
		toChroma(t.SyntaxKeyword),
		toChroma(t.SyntaxKeyword),
		toChroma(t.SyntaxKeyword),
		toChroma(t.SyntaxKeyword),
		toChroma(t.SyntaxType),
		// Names
		toChroma(t.SyntaxVariable),
		toChroma(t.SyntaxFunction),
		toChroma(t.SyntaxFunction),
		toChroma(t.SyntaxType),
		toChroma(t.SyntaxVariable),
		toChroma(t.SyntaxFunction),
		toChroma(t.SyntaxVariable),
		toChroma(t.SyntaxType),
		toChroma(t.SyntaxFunction),
		toChroma(t.SyntaxFunction),
		toChroma(t.SyntaxVariable),
		toChroma(t.SyntaxType),
		toChroma(t.SyntaxVariable),
		toChroma(t.SyntaxVariable),
		toChroma(t.SyntaxKeyword),
		toChroma(t.SyntaxVariable),
		toChroma(t.SyntaxVariable),
		toChroma(t.SyntaxVariable),
		toChroma(t.SyntaxVariable),
		toChroma(t.SyntaxVariable),
		// Strings
		toChroma(t.SyntaxString),
		toChroma(t.SyntaxString),
		toChroma(t.SyntaxString),
		toChroma(t.SyntaxString),
		toChroma(t.SyntaxString),
		toChroma(t.SyntaxString),
		toChroma(t.SyntaxString),
		toChroma(t.SyntaxString),
		toChroma(t.SyntaxString),
		toChroma(t.SyntaxString),
		toChroma(t.SyntaxString),
		toChroma(t.SyntaxString),
		toChroma(t.SyntaxString),
		toChroma(t.SyntaxString),
		// Numbers
		toChroma(t.SyntaxNumber),
		toChroma(t.SyntaxNumber),
		toChroma(t.SyntaxNumber),
		toChroma(t.SyntaxNumber),
		toChroma(t.SyntaxNumber),
		toChroma(t.SyntaxNumber),
		toChroma(t.SyntaxNumber),
		// Comments
		toChroma(t.SyntaxComment),
		toChroma(t.SyntaxComment),
		toChroma(t.SyntaxComment),
		toChroma(t.SyntaxComment),
		toChroma(t.SyntaxComment),
		toChroma(t.SyntaxComment),
		toChroma(t.SyntaxComment),
		// Operators
		toChroma(t.SyntaxOperator),
		toChroma(t.SyntaxOperator),
		toChroma(t.SyntaxPunctuation),
		// Generic (diff)
		toChroma(t.DiffRemoved), toChroma(t.DiffRemovedBg),
		toChroma(t.DiffAdded), toChroma(t.DiffAddedBg),
		toChroma(t.Text),
		toChroma(t.TextMuted),
	)
	
	// Create style from XML
	style, err := chroma.NewXMLStyle(strings.NewReader(styleXML))
	if err != nil {
		return nil, fmt.Errorf("failed to create Chroma style: %w", err)
	}
	
	return style, nil
}

// SyntaxHighlight applies syntax highlighting to source code
func SyntaxHighlight(source, filename string) (string, error) {
	// Determine lexer
	var lexer chroma.Lexer
	if filename != "" {
		lexer = lexers.Match(filename)
	}
	if lexer == nil {
		lexer = lexers.Analyse(source)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	
	// Coalesce lexer
	lexer = chroma.Coalesce(lexer)
	
	// Generate Chroma style
	style, err := GenerateChromaStyle()
	if err != nil {
		// Fall back to default style
		style = styles.Get("monokai")
	}
	
	// Create formatter
	formatter := formatters.Get("terminal16m")
	if formatter == nil {
		formatter = formatters.Fallback
	}
	
	// Tokenize
	tokens, err := lexer.Tokenise(nil, source)
	if err != nil {
		return source, err
	}
	
	// Format
	var buf bytes.Buffer
	err = formatter.Format(&buf, style, tokens)
	if err != nil {
		return source, err
	}
	
	return buf.String(), nil
}

// SyntaxHighlightLine highlights a single line with proper ANSI handling
func SyntaxHighlightLine(line, filename string) string {
	// Don't highlight empty lines
	if strings.TrimSpace(line) == "" {
		return line
	}
	
	highlighted, err := SyntaxHighlight(line, filename)
	if err != nil {
		return line
	}
	
	// Remove trailing newline that Chroma adds
	return strings.TrimSuffix(highlighted, "\n")
}

// ApplySyntaxHighlighting applies highlighting to a writer with background color
func ApplySyntaxHighlighting(w io.Writer, source, filename string) error {
	// Determine lexer
	var lexer chroma.Lexer
	if filename != "" {
		lexer = lexers.Match(filename)
	}
	if lexer == nil {
		lexer = lexers.Analyse(source)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	
	// Generate style
	style, err := GenerateChromaStyle()
	if err != nil {
		style = styles.Get("monokai")
	}
	
	// Get formatter
	formatter := formatters.Get("terminal16m")
	if formatter == nil {
		formatter = formatters.Fallback
	}
	
	// Tokenize and format
	tokens, err := lexer.Tokenise(nil, source)
	if err != nil {
		return err
	}
	
	return formatter.Format(w, style, tokens)
}