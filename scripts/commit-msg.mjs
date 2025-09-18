#!/usr/bin/env node
import { readFileSync } from "node:fs";

const PROJECT_KEY = (process.env.PROJECT_KEY || "SIK6").toUpperCase();

const msgFile = process.argv[2];
const raw = readFileSync(msgFile, "utf8");

// First non-empty, non-comment line
const firstLine = raw
  .split(/\r?\n/)
  .map((line) => line.trim())
  .find((line) => line.length > 0 && !line.startsWith("#"));

// [SIK6-123] type(optional-scope): short message #123
const regex = new RegExp(
  `^\\[${PROJECT_KEY}-\\d+\\]\\s+` +     // [SIK6-123]
  `[a-z]+(?:\\([^)]+\\))?:\\s+` +        // type or type(scope):
  `[^#\\r\\n]+\\s+#\\d+\\s*$`            // message then #ref (digits)
);

if (!firstLine || !regex.test(firstLine)) {
  console.error(
    `\n❌ Invalid commit message.\n` +
    `Expected: [${PROJECT_KEY}-123] type(optional-scope): short message #123\n`
  );
  process.exit(1);
}
