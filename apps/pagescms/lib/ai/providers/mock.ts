import type { AiDraftProvider, DraftProviderInput } from "../draft-provider";

const sentence = (value: string): string => {
  const normalized = value.replace(/\s+/g, " ").trim().replace(/[.!?]+$/, "");
  return normalized.charAt(0).toUpperCase() + normalized.slice(1);
};

export class MockAiDraftProvider implements AiDraftProvider {
  async generate(input: DraftProviderInput) {
    const topic = sentence(input.prompt).slice(0, 140);
    return {
      title: topic,
      excerpt: `A practical introduction to ${topic.toLowerCase()}.`,
      body: `## ${topic}\n\nThis local draft was generated deterministically so the complete CMS workflow can be tested without an API key.\n\n### Start here\n\nFocus on the reader's goal, explain the essential ideas clearly, and finish with a concrete next step.`,
      seoTitle: topic.slice(0, 60),
      seoDescription: `Learn the essentials of ${topic.toLowerCase()} with a clear, practical guide.`.slice(0, 160),
    };
  }
}
