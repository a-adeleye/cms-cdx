import "server-only";
import { AiDraftSchema, type AiDraft } from "../contracts";
import type { AiDraftProvider, DraftProviderInput } from "../draft-provider";

export class OpenAiResponsesProvider implements AiDraftProvider {
  constructor(
    private readonly apiKey: string,
    private readonly model: string,
    private readonly timeoutMs: number,
    private readonly fetchImpl: typeof fetch = fetch,
  ) {}

  async generate(input: DraftProviderInput): Promise<AiDraft> {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), this.timeoutMs);
    try {
      const response = await this.fetchImpl("https://api.openai.com/v1/responses", {
        method: "POST",
        headers: { "Content-Type": "application/json", Authorization: `Bearer ${this.apiKey}` },
        signal: controller.signal,
        body: JSON.stringify({
          model: this.model,
          store: false,
          instructions: input.instructions || "Draft clear, accurate editorial content. Return only the requested JSON fields.",
          input: JSON.stringify({ prompt: input.prompt, context: input.context || {} }),
          text: {
            format: {
              type: "json_schema",
              name: "article_draft",
              strict: true,
              schema: {
                type: "object",
                additionalProperties: false,
                required: ["title", "excerpt", "body", "seoTitle", "seoDescription"],
                properties: {
                  title: { type: "string" }, excerpt: { type: "string" }, body: { type: "string" },
                  seoTitle: { type: "string" }, seoDescription: { type: "string" },
                },
              },
            },
          },
        }),
      });
      if (!response.ok) throw new Error("The AI provider could not generate a draft.");
      const payload: any = await response.json();
      if (payload.status === "incomplete") throw new Error("The AI response was incomplete.");
      const content = payload.output?.flatMap((item: any) => item.content || []) || [];
      if (content.some((item: any) => item.type === "refusal")) throw new Error("The AI provider declined this request.");
      const outputText = content.find((item: any) => item.type === "output_text")?.text;
      if (typeof outputText !== "string") throw new Error("The AI provider returned an invalid response.");
      return AiDraftSchema.parse(JSON.parse(outputText));
    } finally {
      clearTimeout(timeout);
    }
  }
}
