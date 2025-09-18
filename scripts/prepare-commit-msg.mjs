#!/usr/bin/env node
import { writeFileSync, readFileSync } from "node:fs";
import { execSync } from "node:child_process";

const PROJECT_KEY = (process.env.PROJECT_KEY || "SIK6").toUpperCase();

const msgFile = process.argv[2];
const source = process.argv[3];

// Skip auto commits we don't want to prefill
if (["merge", "squash"].includes(source)) process.exit(0);

// Get current branch
let branch = "";
try {
  branch = execSync("git symbolic-ref --short HEAD", {
    stdio: ["ignore", "pipe", "ignore"],
  }).toString().trim();
} catch {
  // Detached HEAD → skip
  process.exit(0);
}

// Find "<PROJECT_KEY>-<number>" anywhere (case-insensitive, global)
const re = new RegExp(`${PROJECT_KEY}-(\\d+)`, "gi"); // ← change 2 (kept 'gi', tolerates "sik6")
let m, lastNum = null;
while ((m = re.exec(branch)) !== null) lastNum = m[1];
if (!lastNum) process.exit(0);

// Normalize ticket number (e.g., 007 → 7)
const ticketId = String(parseInt(lastNum, 10));

// Build default template (no behavior change)
const line = `[${PROJECT_KEY}-${ticketId}] type: description... #<ref>\n`;

const current = readFileSync(msgFile, "utf8");

// Only prefill if no user content (ignoring comments)
const nonComment = current
  .split(/\r?\n/)
  .filter((l) => !l.trim().startsWith("#"))
  .join("\n")
  .trim();

if (!nonComment) {
  writeFileSync(msgFile, line + current, "utf8");
}
