export type InlineNode =
  | { kind: 'text'; value: string }
  | { kind: 'strong' | 'emphasis'; children: InlineNode[] }
  | { kind: 'code'; value: string }
  | { kind: 'link'; href: string; isExternal: boolean; children: InlineNode[] };

export type MarkdownBlock =
  | { kind: 'heading'; depth: 2 | 3; id: string; children: InlineNode[]; text: string }
  | { kind: 'paragraph'; children: InlineNode[] }
  | { kind: 'quote'; children: InlineNode[] }
  | { kind: 'list'; ordered: boolean; entries: InlineNode[][] };

export type TableOfContentsEntry = {
  depth: 2 | 3;
  label: string;
  id: string;
};

export type ParsedMarkdown = {
  blocks: MarkdownBlock[];
  tableOfContents: TableOfContentsEntry[];
};

type UnparsedBlock =
  | { kind: 'heading'; depth: 1 | 2 | 3; text: string }
  | { kind: 'paragraph'; text: string }
  | { kind: 'quote'; text: string }
  | { kind: 'list'; ordered: boolean; entries: string[] };

const MAX_MARKDOWN_LENGTH = 250_000;

export function parseArticleMarkdown(value: unknown, articleTitle = ''): ParsedMarkdown {
  const markdown = normalizeMarkdown(value);
  const headingIds = new Map<string, number>();
  const unparsedBlocks = omitDuplicateTitle(parseBlocks(markdown), articleTitle);
  const blocks = unparsedBlocks.map((block) => finalizeBlock(block, headingIds));
  const tableOfContents = blocks
    .filter((block): block is Extract<MarkdownBlock, { kind: 'heading' }> => block.kind === 'heading')
    .map((heading) => ({ depth: heading.depth, label: heading.text, id: heading.id }));
  return { blocks, tableOfContents };
}

function omitDuplicateTitle(blocks: UnparsedBlock[], articleTitle: string): UnparsedBlock[] {
  const first = blocks[0];
  if (first?.kind !== 'heading' || first.depth !== 1) return blocks;
  const normalize = (value: string) => value.trim().replace(/\s+/g, ' ').toLowerCase();
  return normalize(first.text) === normalize(articleTitle) ? blocks.slice(1) : blocks;
}

export function slugifyHeading(value: string): string {
  return value
    .toLowerCase()
    .normalize('NFKD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-|-$/g, '') || 'section';
}

function normalizeMarkdown(value: unknown): string {
  if (typeof value !== 'string') return '';
  if (value.length > MAX_MARKDOWN_LENGTH) {
    throw new Error('Article content exceeds the supported length.');
  }
  if (value.includes('\0')) {
    throw new Error('Article content contains an unsupported control character.');
  }
  return value.replace(/\r\n?/g, '\n');
}

function parseBlocks(markdown: string): UnparsedBlock[] {
  const lines = markdown.split('\n');
  const blocks: UnparsedBlock[] = [];
  let cursor = 0;
  while (cursor < lines.length) {
    if (!lines[cursor].trim()) {
      cursor += 1;
      continue;
    }
    const next = readBlock(lines, cursor);
    blocks.push(next.block);
    cursor = next.cursor;
  }
  return blocks;
}

function readBlock(lines: string[], cursor: number): { block: UnparsedBlock; cursor: number } {
  const line = lines[cursor].trim();
  const heading = /^(#{1,3})\s+(.+)$/.exec(line);
  if (heading) {
    return {
      block: { kind: 'heading', depth: heading[1].length as 1 | 2 | 3, text: heading[2].trim() },
      cursor: cursor + 1,
    };
  }
  if (line.startsWith('>')) return readQuote(lines, cursor);
  if (/^(?:[-*]|\d+\.)\s+/.test(line)) return readList(lines, cursor);
  return readParagraph(lines, cursor);
}

function readQuote(lines: string[], cursor: number): { block: UnparsedBlock; cursor: number } {
  const quote: string[] = [];
  while (cursor < lines.length && lines[cursor].trim().startsWith('>')) {
    quote.push(lines[cursor].trim().replace(/^>\s?/, ''));
    cursor += 1;
  }
  return { block: { kind: 'quote', text: quote.join(' ') }, cursor };
}

function readList(lines: string[], cursor: number): { block: UnparsedBlock; cursor: number } {
  const ordered = /^\d+\.\s+/.test(lines[cursor].trim());
  const pattern = ordered ? /^\d+\.\s+(.+)$/ : /^[-*]\s+(.+)$/;
  const entries: string[] = [];
  while (cursor < lines.length) {
    const match = pattern.exec(lines[cursor].trim());
    if (!match) break;
    entries.push(match[1]);
    cursor += 1;
  }
  return { block: { kind: 'list', ordered, entries }, cursor };
}

function readParagraph(lines: string[], cursor: number): { block: UnparsedBlock; cursor: number } {
  const paragraph: string[] = [];
  while (cursor < lines.length) {
    const line = lines[cursor].trim();
    if (!line || /^(?:#{1,3}\s+|>|[-*]\s+|\d+\.\s+)/.test(line)) break;
    paragraph.push(line);
    cursor += 1;
  }
  return { block: { kind: 'paragraph', text: paragraph.join(' ') }, cursor };
}

function finalizeBlock(block: UnparsedBlock, headingIds: Map<string, number>): MarkdownBlock {
  if (block.kind === 'heading') {
    const baseId = slugifyHeading(block.text);
    const count = (headingIds.get(baseId) ?? 0) + 1;
    headingIds.set(baseId, count);
    return {
      kind: 'heading',
      depth: block.depth === 3 ? 3 : 2,
      id: count === 1 ? baseId : `${baseId}-${count}`,
      text: block.text,
      children: parseInline(block.text),
    };
  }
  if (block.kind === 'list') {
    return { kind: 'list', ordered: block.ordered, entries: block.entries.map(parseInline) };
  }
  return { kind: block.kind, children: parseInline(block.text) };
}

function parseInline(value: string): InlineNode[] {
  const nodes: InlineNode[] = [];
  const pattern = /(\[([^\]]+)\]\(([^)]+)\)|\*\*([^*]+)\*\*|\*([^*]+)\*|`([^`]+)`)/g;
  let cursor = 0;
  for (const match of value.matchAll(pattern)) {
    const index = match.index ?? 0;
    if (index > cursor) nodes.push({ kind: 'text', value: value.slice(cursor, index) });
    nodes.push(readInlineMatch(match));
    cursor = index + match[0].length;
  }
  if (cursor < value.length) nodes.push({ kind: 'text', value: value.slice(cursor) });
  return nodes;
}

function readInlineMatch(match: RegExpMatchArray): InlineNode {
  if (match[2] !== undefined) {
    const href = normalizeLink(match[3]);
    return href
      ? { kind: 'link', href, isExternal: /^https?:\/\//i.test(href), children: [{ kind: 'text', value: match[2] }] }
      : { kind: 'text', value: match[2] };
  }
  if (match[4] !== undefined) return { kind: 'strong', children: [{ kind: 'text', value: match[4] }] };
  if (match[5] !== undefined) return { kind: 'emphasis', children: [{ kind: 'text', value: match[5] }] };
  return { kind: 'code', value: match[6] ?? '' };
}

function normalizeLink(value: string): string | undefined {
  const href = value.trim();
  if (/^(?:https?:\/\/|mailto:|\/|#)[^\s\\]*$/i.test(href) && !href.startsWith('//')) {
    return href;
  }
  return undefined;
}
