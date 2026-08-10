import "server-only";
import { z } from "zod";
import type { AiDraftProvider } from "./draft-provider";
import { MockAiDraftProvider } from "./providers/mock";
import { OpenAiResponsesProvider } from "./providers/openai-responses";

const ConfigSchema = z.object({
  provider: z.enum(["mock", "openai"]),
  model: z.string().min(1),
  timeoutMs: z.number().int().min(1_000).max(120_000),
});

export const getDraftProvider = (): AiDraftProvider => {
  const config = ConfigSchema.parse({
    provider: process.env.AI_PROVIDER || "mock",
    model: process.env.OPENAI_MODEL || "gpt-5.6-luna",
    timeoutMs: Number(process.env.AI_REQUEST_TIMEOUT_MS || "30000"),
  });
  if (config.provider === "mock") return new MockAiDraftProvider();
  if (!process.env.OPENAI_API_KEY) throw new Error("OPENAI_API_KEY is required when AI_PROVIDER=openai.");
  return new OpenAiResponsesProvider(process.env.OPENAI_API_KEY, config.model, config.timeoutMs);
};
