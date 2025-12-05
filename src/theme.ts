// Theme definitions for diff viewer
export interface Theme {
  name: string;
  // Backgrounds
  bg: string;
  panelBg: string;
  // Text
  text: string;
  textMuted: string;
  // Diff colors
  addedBg: string;
  addedFg: string;
  addedHighlight: string;
  removedBg: string;
  removedFg: string;
  removedHighlight: string;
  contextBg: string;
  // Line numbers
  lineNumberBg: string;
  lineNumberFg: string;
  addedLineNumberBg: string;
  removedLineNumberBg: string;
  // UI
  border: string;
  accent: string;
  hunkHeader: string;
  // Shiki theme name
  shikiTheme: string;
}

export const themes: Record<string, Theme> = {
  "github-dark": {
    name: "GitHub Dark",
    bg: "#0d1117",
    panelBg: "#161b22",
    text: "#c9d1d9",
    textMuted: "#8b949e",
    addedBg: "#12261e",
    addedFg: "#3fb950",
    addedHighlight: "#196c2e",
    removedBg: "#2d1b1e",
    removedFg: "#f85149",
    removedHighlight: "#8b1c1c",
    contextBg: "#0d1117",
    lineNumberBg: "#161b22",
    lineNumberFg: "#6e7681",
    addedLineNumberBg: "#1b4721",
    removedLineNumberBg: "#67060c",
    border: "#30363d",
    accent: "#58a6ff",
    hunkHeader: "#388bfd",
    shikiTheme: "github-dark-default",
  },
  dracula: {
    name: "Dracula",
    bg: "#282a36",
    panelBg: "#21222c",
    text: "#f8f8f2",
    textMuted: "#6272a4",
    addedBg: "#1e3a29",
    addedFg: "#50fa7b",
    addedHighlight: "#2d5a3d",
    removedBg: "#3a1e24",
    removedFg: "#ff5555",
    removedHighlight: "#5a2d34",
    contextBg: "#282a36",
    lineNumberBg: "#21222c",
    lineNumberFg: "#6272a4",
    addedLineNumberBg: "#264d35",
    removedLineNumberBg: "#4d2630",
    border: "#44475a",
    accent: "#bd93f9",
    hunkHeader: "#8be9fd",
    shikiTheme: "dracula",
  },
  monokai: {
    name: "Monokai",
    bg: "#272822",
    panelBg: "#1e1f1c",
    text: "#f8f8f2",
    textMuted: "#75715e",
    addedBg: "#1e3a1e",
    addedFg: "#a6e22e",
    addedHighlight: "#2d5a2d",
    removedBg: "#3a1e1e",
    removedFg: "#f92672",
    removedHighlight: "#5a2d2d",
    contextBg: "#272822",
    lineNumberBg: "#1e1f1c",
    lineNumberFg: "#75715e",
    addedLineNumberBg: "#2a4f2a",
    removedLineNumberBg: "#4f2a35",
    border: "#3e3d32",
    accent: "#66d9ef",
    hunkHeader: "#ae81ff",
    shikiTheme: "monokai",
  },
  "one-dark": {
    name: "One Dark",
    bg: "#282c34",
    panelBg: "#21252b",
    text: "#abb2bf",
    textMuted: "#5c6370",
    addedBg: "#1e3327",
    addedFg: "#98c379",
    addedHighlight: "#2d4d39",
    removedBg: "#352129",
    removedFg: "#e06c75",
    removedHighlight: "#4d3039",
    contextBg: "#282c34",
    lineNumberBg: "#21252b",
    lineNumberFg: "#5c6370",
    addedLineNumberBg: "#2a4030",
    removedLineNumberBg: "#402a30",
    border: "#3e4451",
    accent: "#61afef",
    hunkHeader: "#c678dd",
    shikiTheme: "one-dark-pro",
  },
  nord: {
    name: "Nord",
    bg: "#2e3440",
    panelBg: "#3b4252",
    text: "#d8dee9",
    textMuted: "#4c566a",
    addedBg: "#243330",
    addedFg: "#a3be8c",
    addedHighlight: "#354d45",
    removedBg: "#3d2f38",
    removedFg: "#bf616a",
    removedHighlight: "#5a424a",
    contextBg: "#2e3440",
    lineNumberBg: "#3b4252",
    lineNumberFg: "#4c566a",
    addedLineNumberBg: "#2f4038",
    removedLineNumberBg: "#4a353d",
    border: "#4c566a",
    accent: "#88c0d0",
    hunkHeader: "#81a1c1",
    shikiTheme: "nord",
  },
  tokyonight: {
    name: "Tokyo Night",
    bg: "#1a1b26",
    panelBg: "#16161e",
    text: "#a9b1d6",
    textMuted: "#565f89",
    addedBg: "#1a2f29",
    addedFg: "#9ece6a",
    addedHighlight: "#2a4a3d",
    removedBg: "#2d1f2f",
    removedFg: "#f7768e",
    removedHighlight: "#4a2f42",
    contextBg: "#1a1b26",
    lineNumberBg: "#16161e",
    lineNumberFg: "#3b4261",
    addedLineNumberBg: "#243d33",
    removedLineNumberBg: "#3d2438",
    border: "#414868",
    accent: "#7aa2f7",
    hunkHeader: "#bb9af7",
    shikiTheme: "tokyo-night",
  },
};

export const defaultTheme = "github-dark";
