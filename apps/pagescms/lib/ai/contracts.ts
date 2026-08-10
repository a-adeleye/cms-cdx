import { z } from "zod";

const safeFieldPath = z.string().min(1).max(200).regex(
  /^(?!.*(?:^|\.)(?:__proto__|prototype|constructor)(?:\.|$))[A-Za-z_][A-Za-z0-9_]*(?:\.[A-Za-z_][A-Za-z0-9_]*)*$/,
  "Invalid target field path.",
);

export const AiDraftSchema = z.object({
  title: z.string().min(1).max(200),
  excerpt: z.string().min(1).max(1000),
  body: z.string().min(1).max(50_000),
  seoTitle: z.string().min(1).max(200),
  seoDescription: z.string().min(1).max(1000),
}).strict();

export const AiDraftRequestSchema = z.object({
  schemaName: z.string().min(1).max(100),
  prompt: z.string().trim().min(3).max(4_000),
  context: z.record(z.string(), z.unknown()).optional(),
}).strict();

export const AiDraftOptionsSchema = z.object({
  instructions: z.string().max(2_000).optional(),
  targets: z.object({
    title: safeFieldPath,
    excerpt: safeFieldPath,
    body: safeFieldPath,
    seoTitle: safeFieldPath,
    seoDescription: safeFieldPath,
  }).strict(),
}).strict();

export type AiDraft = z.infer<typeof AiDraftSchema>;
export type AiDraftRequest = z.infer<typeof AiDraftRequestSchema>;
export type AiDraftOptions = z.infer<typeof AiDraftOptionsSchema>;
