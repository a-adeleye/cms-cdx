import { describe, expect, it } from "vitest";
import { AiDraftSchema } from "../contracts";
import { MockAiDraftProvider } from "./mock";

describe("MockAiDraftProvider", () => {
  it("returns the same valid draft for the same prompt", async () => {
    const provider = new MockAiDraftProvider();
    const first = await provider.generate({ prompt: "local first publishing" });
    const second = await provider.generate({ prompt: "local first publishing" });
    expect(second).toEqual(first);
    expect(AiDraftSchema.parse(first)).toEqual(first);
  });
});
