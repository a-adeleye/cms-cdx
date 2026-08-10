import type { AiDraft } from "./contracts";

export type DraftProviderInput = {
  prompt: string;
  instructions?: string;
  context?: Record<string, unknown>;
};

export interface AiDraftProvider {
  generate(input: DraftProviderInput): Promise<AiDraft>;
}
