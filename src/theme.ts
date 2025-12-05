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
  // Syntax highlighting
  syntax: {
    keyword: string;
    string: string;
    comment: string;
    number: string;
    function: string;
    type: string;
    variable: string;
    property: string;
    operator: string;
    punctuation: string;
  };
}

export const themes: Record<string, Theme> = {
  "github-dark": {
    name: "GitHub Dark",
    bg: "#0d1117",
    panelBg: "#161b22",
    text: "#e6edf3",
    textMuted: "#8b949e",
    addedBg: "#1a4d1a",
    addedFg: "#22c55e",
    addedHighlight: "#196c2e",
    removedBg: "#4d1a1a",
    removedFg: "#ef4444",
    removedHighlight: "#8b1c1c",
    contextBg: "#0d1117",
    lineNumberBg: "#161b22",
    lineNumberFg: "#6b7280",
    addedLineNumberBg: "#0d3a0d",
    removedLineNumberBg: "#3a0d0d",
    border: "#30363d",
    accent: "#58a6ff",
    hunkHeader: "#388bfd",
    syntax: {
      keyword: "#ff7b72",
      string: "#a5d6ff",
      comment: "#8b949e",
      number: "#79c0ff",
      function: "#d2a8ff",
      type: "#ffa657",
      variable: "#e6edf3",
      property: "#79c0ff",
      operator: "#ff7b72",
      punctuation: "#f0f6fc",
    },
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
    syntax: {
      keyword: "#ff79c6",
      string: "#f1fa8c",
      comment: "#6272a4",
      number: "#bd93f9",
      function: "#50fa7b",
      type: "#8be9fd",
      variable: "#f8f8f2",
      property: "#66d9ef",
      operator: "#ff79c6",
      punctuation: "#f8f8f2",
    },
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
    syntax: {
      keyword: "#f92672",
      string: "#e6db74",
      comment: "#75715e",
      number: "#ae81ff",
      function: "#a6e22e",
      type: "#66d9ef",
      variable: "#f8f8f2",
      property: "#66d9ef",
      operator: "#f92672",
      punctuation: "#f8f8f2",
    },
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
    syntax: {
      keyword: "#c678dd",
      string: "#98c379",
      comment: "#5c6370",
      number: "#d19a66",
      function: "#61afef",
      type: "#e5c07b",
      variable: "#e06c75",
      property: "#e5c07b",
      operator: "#56b6c2",
      punctuation: "#abb2bf",
    },
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
    syntax: {
      keyword: "#81a1c1",
      string: "#a3be8c",
      comment: "#616e88",
      number: "#b48ead",
      function: "#88c0d0",
      type: "#8fbcbb",
      variable: "#d8dee9",
      property: "#8fbcbb",
      operator: "#81a1c1",
      punctuation: "#eceff4",
    },
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
    syntax: {
      keyword: "#bb9af7",
      string: "#9ece6a",
      comment: "#565f89",
      number: "#ff9e64",
      function: "#7aa2f7",
      type: "#2ac3de",
      variable: "#c0caf5",
      property: "#73daca",
      operator: "#89ddff",
      punctuation: "#a9b1d6",
    },
  },
};

export const defaultTheme = "github-dark";
