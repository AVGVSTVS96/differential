#!/usr/bin/env bun
import React from "react";
import { createCliRenderer } from "@opentui/core";
import { createRoot } from "@opentui/react";
import cac from "cac";
import { App } from "./app";
import { themes, defaultTheme } from "./theme";
import { runCommand, colorize } from "./utils";

const cli = cac("differential");

// Helper for pipe mode output
function outputPipeMode(diff: string, themeName: string) {
  const theme = themes[themeName] || themes[defaultTheme];
  const lines = diff.split("\n");

  for (const line of lines) {
    if (line.startsWith("+") && !line.startsWith("+++")) {
      console.log(colorize(line, theme.addedFg));
    } else if (line.startsWith("-") && !line.startsWith("---")) {
      console.log(colorize(line, theme.removedFg));
    } else if (line.startsWith("@@")) {
      console.log(colorize(line, theme.hunkHeader));
    } else if (
      line.startsWith("diff ") ||
      line.startsWith("index ") ||
      line.startsWith("---") ||
      line.startsWith("+++")
    ) {
      console.log(colorize(line, theme.accent));
    } else {
      console.log(colorize(line, theme.textMuted));
    }
  }
}

// Main command
cli
  .command("[...files]", "Compare files or view diff")
  .option("-v, --view <mode>", "View mode: unified or split", {
    default: "unified",
  })
  .option("-t, --theme <name>", `Theme: ${Object.keys(themes).join(", ")}`, {
    default: defaultTheme,
  })
  .option("-n, --no-line-numbers", "Hide line numbers")
  .option("-p, --pipe", "Output to stdout (non-interactive)")
  .option("--staged", "Show staged changes (git diff --cached)")
  .option("-c, --commit <ref>", "Show changes for a specific commit")
  .action(async (files: string[], options) => {
    let diffContent: string;

    // Determine input source
    if (files.length === 0) {
      // Check for piped input
      if (!process.stdin.isTTY) {
        const chunks: Uint8Array[] = [];
        for await (const chunk of Bun.stdin.stream()) {
          chunks.push(chunk);
        }
        diffContent = Buffer.concat(chunks).toString("utf-8");
      } else if (options.staged) {
        // Git staged changes
        diffContent = runCommand("git diff --cached");
      } else if (options.commit) {
        // Specific commit
        diffContent = runCommand(`git show ${options.commit}`);
      } else {
        // Default: unstaged changes
        diffContent = runCommand("git diff");
      }
    } else if (files.length === 1) {
      // Single file - git diff on it
      const file = files[0];
      diffContent = runCommand(`git diff -- "${file}"`);
      if (!diffContent.trim()) {
        diffContent = runCommand(`git diff HEAD -- "${file}"`);
      }
    } else if (files.length === 2) {
      // Two files - compare them
      const [file1, file2] = files;
      const f1 = Bun.file(file1);
      const f2 = Bun.file(file2);

      if (!(await f1.exists())) {
        console.error(`File not found: ${file1}`);
        process.exit(1);
      }
      if (!(await f2.exists())) {
        console.error(`File not found: ${file2}`);
        process.exit(1);
      }

      diffContent = runCommand(`diff -u "${file1}" "${file2}"`);
    } else {
      console.error("Too many arguments. Expected 0-2 files.");
      process.exit(1);
    }

    if (!diffContent.trim()) {
      console.log("No changes found.");
      process.exit(0);
    }

    // Pipe mode - just output colored diff
    if (options.pipe || !process.stdout.isTTY) {
      outputPipeMode(diffContent, options.theme);
      process.exit(0);
    }

    // Interactive TUI mode
    try {
      const renderer = await createCliRenderer({
        useAlternateScreen: true,
        useMouse: true,
        exitOnCtrlC: true,
      });

      createRoot(renderer).render(
        <App
          diffContent={diffContent}
          initialTheme={options.theme}
          initialViewMode={options.view as "unified" | "split"}
          initialShowLineNumbers={options.lineNumbers !== false}
        />
      );

      await renderer.setupTerminal();
      renderer.start();
    } catch (err: any) {
      // Fallback to pipe mode if TUI fails
      console.error("Failed to start TUI, falling back to pipe mode:", err.message);
      outputPipeMode(diffContent, options.theme);
    }
  });

// Help text
cli.help();
cli.version("0.1.0");

// Parse and run
cli.parse();
