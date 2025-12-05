import React, { useState, useMemo } from "react";
import {
  useTerminalDimensions,
  useKeyboard,
  useRenderer,
} from "@opentui/react";
import { SyntaxStyle, pathToFiletype, parseColor } from "@opentui/core";
import { parsePatch, type ParsedDiff } from "diff";
import { themes, defaultTheme, type Theme } from "./theme";
import { extractFilename } from "./utils";

interface AppProps {
  diffContent: string;
  initialTheme?: string;
  initialViewMode?: "unified" | "split";
  initialShowLineNumbers?: boolean;
}

// Convert theme to SyntaxStyle for OpenTUI (matching diff-demo.ts structure exactly)
function createSyntaxStyle(theme: Theme): SyntaxStyle {
  return SyntaxStyle.fromStyles({
    keyword: { fg: parseColor(theme.syntax.keyword), bold: true },
    "keyword.import": { fg: parseColor(theme.syntax.keyword), bold: true },
    string: { fg: parseColor(theme.syntax.string) },
    comment: { fg: parseColor(theme.syntax.comment), italic: true },
    number: { fg: parseColor(theme.syntax.number) },
    boolean: { fg: parseColor(theme.syntax.number) },
    constant: { fg: parseColor(theme.syntax.number) },
    function: { fg: parseColor(theme.syntax.function) },
    "function.call": { fg: parseColor(theme.syntax.function) },
    constructor: { fg: parseColor(theme.syntax.type) },
    type: { fg: parseColor(theme.syntax.type) },
    operator: { fg: parseColor(theme.syntax.operator) },
    variable: { fg: parseColor(theme.syntax.variable) },
    property: { fg: parseColor(theme.syntax.property) },
    bracket: { fg: parseColor(theme.syntax.punctuation) },
    punctuation: { fg: parseColor(theme.syntax.punctuation) },
    default: { fg: parseColor(theme.text) },
  });
}

export function App({
  diffContent,
  initialTheme = defaultTheme,
  initialViewMode = "unified",
  initialShowLineNumbers = true,
}: AppProps) {
  const { width, height } = useTerminalDimensions();
  const renderer = useRenderer();

  // Parse diff into files
  const files = useMemo(() => {
    try {
      return parsePatch(diffContent);
    } catch {
      return [];
    }
  }, [diffContent]);

  // State
  const [currentFileIndex, setCurrentFileIndex] = useState(0);
  const [viewMode, setViewMode] = useState<"unified" | "split">(initialViewMode);
  const [showLineNumbers, setShowLineNumbers] = useState(initialShowLineNumbers);
  const [themeName, setThemeName] = useState(initialTheme);
  const theme = themes[themeName] || themes[defaultTheme];
  const syntaxStyle = useMemo(() => createSyntaxStyle(theme), [theme]);

  const currentFile = files[currentFileIndex];

  // Calculate additions/deletions for current file
  const { additions, deletions, status } = useMemo(() => {
    if (!currentFile) return { additions: 0, deletions: 0, status: "modified" as const };

    let adds = 0;
    let dels = 0;

    for (const hunk of currentFile.hunks) {
      for (const line of hunk.lines) {
        if (line.startsWith("+") && !line.startsWith("+++")) adds++;
        else if (line.startsWith("-") && !line.startsWith("---")) dels++;
      }
    }

    const status: "created" | "deleted" | "modified" =
      currentFile.oldFileName === "/dev/null"
        ? "created"
        : currentFile.newFileName === "/dev/null"
          ? "deleted"
          : "modified";

    return { additions: adds, deletions: dels, status };
  }, [currentFile]);

  // Generate diff string for current file
  const currentDiff = useMemo(() => {
    if (!currentFile) return "";

    const lines: string[] = [];

    // Add header
    if (currentFile.oldFileName) {
      lines.push(`--- ${currentFile.oldFileName}`);
    }
    if (currentFile.newFileName) {
      lines.push(`+++ ${currentFile.newFileName}`);
    }

    // Add hunks
    for (const hunk of currentFile.hunks) {
      lines.push(`@@ -${hunk.oldStart},${hunk.oldLines} +${hunk.newStart},${hunk.newLines} @@`);
      lines.push(...hunk.lines);
    }

    return lines.join("\n");
  }, [currentFile]);

  // Keyboard handling
  useKeyboard((key) => {
    // Handle different key event formats
    const keyName = typeof key === "string" ? key : key.name || key.sequence || "";

    // Quit
    if (keyName === "q" || keyName === "escape") {
      renderer.destroy();
      process.exit(0);
    }

    // Navigation between files
    if (keyName === "n" || keyName === "right" || keyName === "l") {
      if (currentFileIndex < files.length - 1) {
        setCurrentFileIndex((i) => i + 1);
      }
    }
    if (keyName === "p" || keyName === "left" || keyName === "h") {
      if (currentFileIndex > 0) {
        setCurrentFileIndex((i) => i - 1);
      }
    }

    // Toggle view mode
    if (keyName === "v") {
      setViewMode((m) => (m === "unified" ? "split" : "unified"));
    }

    // Toggle line numbers
    if (keyName === "L") {
      setShowLineNumbers((s) => !s);
    }

    // Cycle themes
    if (keyName === "t") {
      const themeNames = Object.keys(themes);
      const currentIndex = themeNames.indexOf(themeName);
      const nextIndex = (currentIndex + 1) % themeNames.length;
      setThemeName(themeNames[nextIndex]);
    }
  });

  // Empty state
  if (files.length === 0 || !currentFile) {
    return (
      <box
        flexDirection="column"
        alignItems="center"
        justifyContent="center"
        height="100%"
        backgroundColor={theme.bg}
      >
        <text color={theme.textMuted}>No diff content to display</text>
        <text color={theme.textMuted}>Press q to quit</text>
      </box>
    );
  }

  const filename = extractFilename(
    currentFile.newFileName || currentFile.oldFileName || "unknown"
  );

  // Detect filetype from filename using OpenTUI's resolver
  const filetype = pathToFiletype(filename) || filename.split(".").pop() || "";

  return (
    <box flexDirection="column" height="100%" backgroundColor={theme.bg}>
      {/* Header */}
      <box
        flexDirection="row"
        justifyContent="space-between"
        backgroundColor={theme.panelBg}
        paddingLeft={1}
        paddingRight={1}
      >
        <box flexDirection="row">
          <text color={status === "created" ? theme.addedFg : status === "deleted" ? theme.removedFg : theme.accent} bold>
            {status === "created" ? "Created" : status === "deleted" ? "Deleted" : "Modified"}
          </text>
          <text color={theme.text}> </text>
          <text color={theme.text} bold>{filename}</text>
          <text color={theme.textMuted}> </text>
          {additions > 0 && <text color={theme.addedFg}>+{additions}</text>}
          <text color={theme.textMuted}> </text>
          {deletions > 0 && <text color={theme.removedFg}>-{deletions}</text>}
        </box>
        <box flexDirection="row">
          <text color={theme.textMuted}>
            {files.length > 1 ? `${currentFileIndex + 1}/${files.length}` : ""}
          </text>
          <text color={theme.textMuted}> │ </text>
          <text color={theme.textMuted}>{theme.name}</text>
          <text color={theme.textMuted}> │ </text>
          <text color={theme.textMuted}>{viewMode}</text>
        </box>
      </box>

      {/* Diff content using built-in diff component */}
      <diff
        diff={currentDiff}
        view={viewMode}
        filetype={filetype}
        syntaxStyle={syntaxStyle}
        showLineNumbers={showLineNumbers}
        wrapMode="word"
        addedBg={theme.addedBg}
        removedBg={theme.removedBg}
        contextBg={theme.contextBg}
        addedSignColor={theme.addedFg}
        removedSignColor={theme.removedFg}
        lineNumberFg={theme.lineNumberFg}
        lineNumberBg={theme.lineNumberBg}
        addedLineNumberBg={theme.addedLineNumberBg}
        removedLineNumberBg={theme.removedLineNumberBg}
        flexGrow={1}
      />

      {/* Footer with keybindings */}
      <box
        flexDirection="row"
        justifyContent="space-between"
        backgroundColor={theme.panelBg}
        paddingLeft={1}
        paddingRight={1}
      >
        <box flexDirection="row">
          <text color={theme.accent}>q</text>
          <text color={theme.textMuted}>:quit </text>
          {files.length > 1 && (
            <>
              <text color={theme.accent}>n/p</text>
              <text color={theme.textMuted}>:files </text>
            </>
          )}
          <text color={theme.accent}>v</text>
          <text color={theme.textMuted}>:view </text>
          <text color={theme.accent}>t</text>
          <text color={theme.textMuted}>:theme </text>
          <text color={theme.accent}>L</text>
          <text color={theme.textMuted}>:lines</text>
        </box>
        {files.length > 1 && (
          <text color={theme.textMuted}>←/→ navigate files</text>
        )}
      </box>
    </box>
  );
}
