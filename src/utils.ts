import { execSync } from "child_process";

// Language detection from file extension
export const languageMap: Record<string, string> = {
  ts: "typescript",
  tsx: "tsx",
  js: "javascript",
  jsx: "jsx",
  mjs: "javascript",
  cjs: "javascript",
  py: "python",
  rb: "ruby",
  go: "go",
  rs: "rust",
  c: "c",
  cpp: "cpp",
  cc: "cpp",
  h: "c",
  hpp: "cpp",
  java: "java",
  kt: "kotlin",
  scala: "scala",
  json: "json",
  yaml: "yaml",
  yml: "yaml",
  toml: "toml",
  md: "markdown",
  mdx: "mdx",
  sh: "bash",
  bash: "bash",
  zsh: "bash",
  fish: "fish",
  css: "css",
  scss: "scss",
  less: "less",
  html: "html",
  xml: "xml",
  svg: "xml",
  vue: "vue",
  svelte: "svelte",
  sql: "sql",
  graphql: "graphql",
  gql: "graphql",
  dockerfile: "dockerfile",
  makefile: "makefile",
  zig: "zig",
  lua: "lua",
  vim: "viml",
  swift: "swift",
  php: "php",
  r: "r",
  cs: "csharp",
  fs: "fsharp",
  ex: "elixir",
  exs: "elixir",
  erl: "erlang",
  clj: "clojure",
  hs: "haskell",
  ml: "ocaml",
  elm: "elm",
  dart: "dart",
  tf: "hcl",
  hcl: "hcl",
  nix: "nix",
  prisma: "prisma",
};

export function detectLanguage(filename: string): string | undefined {
  const ext = filename.split(".").pop()?.toLowerCase();
  if (!ext) return undefined;
  return languageMap[ext];
}

// Extract filename from diff header or parsed filename
export function extractFilename(diffHeader: string): string {
  // Handle git diff format: "diff --git a/file b/file"
  const gitMatch = diffHeader.match(/^diff --git a\/(.+) b\/(.+)$/m);
  if (gitMatch) return gitMatch[2];

  // Handle standard diff format: "+++ b/file" or "+++ file"
  const plusMatch = diffHeader.match(/^\+\+\+ (?:b\/)?(.+?)(?:\t|$)/m);
  if (plusMatch) return plusMatch[1];

  // Handle "--- a/file"
  const minusMatch = diffHeader.match(/^--- (?:a\/)?(.+?)(?:\t|$)/m);
  if (minusMatch) return minusMatch[1];

  // If no patterns match, the input is likely already a parsed filename
  // Just extract the basename and return it
  const parts = diffHeader.split("/");
  return parts[parts.length - 1] || diffHeader || "unknown";
}

// Run a command and return output
export function runCommand(cmd: string): string {
  try {
    return execSync(cmd, { encoding: "utf-8", maxBuffer: 50 * 1024 * 1024 });
  } catch (e: any) {
    // diff command exits with 1 when files differ
    if (e.stdout) return e.stdout;
    if (e.output) return e.output.filter(Boolean).join("");
    return "";
  }
}

// Calculate Levenshtein distance ratio for word-level diff
export function similarityRatio(a: string, b: string): number {
  if (a === b) return 1;
  if (a.length === 0 || b.length === 0) return 0;

  const matrix: number[][] = [];

  for (let i = 0; i <= a.length; i++) {
    matrix[i] = [i];
  }
  for (let j = 0; j <= b.length; j++) {
    matrix[0][j] = j;
  }

  for (let i = 1; i <= a.length; i++) {
    for (let j = 1; j <= b.length; j++) {
      const cost = a[i - 1] === b[j - 1] ? 0 : 1;
      matrix[i][j] = Math.min(
        matrix[i - 1][j] + 1,
        matrix[i][j - 1] + 1,
        matrix[i - 1][j - 1] + cost
      );
    }
  }

  const distance = matrix[a.length][b.length];
  const maxLen = Math.max(a.length, b.length);
  return 1 - distance / maxLen;
}

// Pad number for line display
export function padLineNumber(num: number | undefined, width: number): string {
  if (num === undefined) return " ".repeat(width);
  return String(num).padStart(width, " ");
}

// ANSI color helpers
export function hexToRgb(hex: string): { r: number; g: number; b: number } {
  const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
  if (!result) return { r: 0, g: 0, b: 0 };
  return {
    r: parseInt(result[1], 16),
    g: parseInt(result[2], 16),
    b: parseInt(result[3], 16),
  };
}

export function rgbToAnsi(r: number, g: number, b: number): string {
  return `\x1b[38;2;${r};${g};${b}m`;
}

export function bgRgbToAnsi(r: number, g: number, b: number): string {
  return `\x1b[48;2;${r};${g};${b}m`;
}

export const RESET = "\x1b[0m";

export function colorize(text: string, fg: string, bg?: string): string {
  const fgRgb = hexToRgb(fg);
  let result = rgbToAnsi(fgRgb.r, fgRgb.g, fgRgb.b);
  if (bg) {
    const bgRgb = hexToRgb(bg);
    result += bgRgbToAnsi(bgRgb.r, bgRgb.g, bgRgb.b);
  }
  return result + text + RESET;
}
