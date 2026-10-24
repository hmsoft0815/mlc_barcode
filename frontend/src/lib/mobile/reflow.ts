// reflow turns a plain-text file hard-wrapped at ~72–80 columns (LICENSE,
// THIRD_PARTY_NOTICES.txt) into lines that the page can wrap itself: on a
// phone the fixed breaks plus the page's own wrapping gave ragged text.
//
// A line is joined to the paragraph before when it sits at that
// paragraph's text indent (for a list item: where the text after "1." or
// "(a)" starts), the line before was long — so it was only wrapped — and
// ends mid-sentence, and the line does not start a new list item. Lines
// ending in a link or path, rules (=====) and blank lines stay as they are.

const LIST_ITEM = /^\s*(\d+\.|\([a-z0-9]+\)|[-*•])\s+/i;
const RULE = /^\s*[=\-_*]{5,}\s*$/;
const WRAPPED_MIN = 50;

function indentOf(line: string): number {
  return line.length - line.trimStart().length;
}

// textIndent is the column where a paragraph's text starts.
function textIndent(line: string): number {
  return LIST_ITEM.exec(line)?.[0].length ?? indentOf(line);
}

function continues(prev: string, line: string, indent: number): boolean {
  const p = prev.trimEnd();
  return (
    line.trim() !== '' &&
    p.trim() !== '' &&
    indentOf(line) === indent &&
    p.length >= WRAPPED_MIN &&
    !/[.:]$/.test(p) &&
    !/\S*\/\S*$/.test(p) && // ends in a link or path
    !RULE.test(p) &&
    !RULE.test(line) &&
    !LIST_ITEM.test(line)
  );
}

export function reflow(text: string): string {
  const out: string[] = [];
  let prev = ''; // the last source line, as written
  let indent = 0; // text indent of the paragraph in out[out.length - 1]
  for (const line of text.replace(/\r\n/g, '\n').split('\n')) {
    if (out.length > 0 && continues(prev, line, indent)) {
      out[out.length - 1] += ' ' + line.trim();
    } else {
      out.push(line);
      indent = textIndent(line);
    }
    prev = line;
  }
  return out.join('\n');
}

export interface Paragraph {
  kind: 'text' | 'gap' | 'rule'; // gap: blank line; rule: a ===== line
  text: string; // without indent and list marker
  marker: string; // "1.", "(a)" or "" — set in a column of its own
  level: number; // indent, in steps of about four columns, at most 3
  under: boolean; // continues the list item before: aligned with its text
}

// paragraphs is reflow split for rendering. The fixed-width layout of the
// files does not carry over to a proportional font, so indents become a
// few levels and a list marker gets a column of its own; centred title
// lines simply start at the left.
export function paragraphs(text: string): Paragraph[] {
  const blank = { text: '', marker: '', level: 0, under: false };
  let item: { hang: number; level: number } | null = null; // the open list item
  return reflow(text)
    .split('\n')
    .map((line): Paragraph => {
      if (line.trim() === '') return { kind: 'gap', ...blank };
      if (RULE.test(line)) return { kind: 'rule', ...blank };
      const m = LIST_ITEM.exec(line);
      const level = Math.min(3, Math.round(indentOf(line) / 4));
      if (m) {
        item = { hang: m[0].length, level };
        return { kind: 'text', text: line.slice(m[0].length).trim(), marker: m[1], level, under: false };
      }
      if (item && indentOf(line) === item.hang) {
        return { kind: 'text', text: line.trim(), marker: '', level: item.level, under: true };
      }
      item = null;
      return { kind: 'text', text: line.trim(), marker: '', level, under: false };
    });
}
