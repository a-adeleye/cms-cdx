import { afterEach, describe, expect, it, vi } from "vitest";
import { isLocalDevelopmentMode } from "./mode";

afterEach(() => {
  vi.unstubAllEnvs();
});

describe("isLocalDevelopmentMode", () => {
  it("is disabled unless explicitly requested", () => {
    vi.stubEnv("PAGESCMS_LOCAL_MODE", "");
    expect(isLocalDevelopmentMode()).toBe(false);
  });

  it("fails closed when requested in production", () => {
    vi.stubEnv("PAGESCMS_LOCAL_MODE", "true");
    vi.stubEnv("NODE_ENV", "production");
    expect(() => isLocalDevelopmentMode()).toThrow(/only allowed/);
  });
});
