import path from "node:path";
import { realpath } from "node:fs/promises";
import { isLocalDevelopmentMode } from "./mode";

export type LocalRepositoryConfig = {
  root: string;
  owner: string;
  repo: string;
  branch: string;
};

let cachedConfig: LocalRepositoryConfig | undefined;

const validateSegment = (label: string, value: string): string => {
  if (!/^[a-zA-Z0-9._-]+$/.test(value) || value === "." || value === "..") {
    throw new Error(`Invalid ${label} for local Pages CMS mode.`);
  }
  return value;
};

export const getLocalRepositoryConfig = async (): Promise<LocalRepositoryConfig> => {
  if (!isLocalDevelopmentMode()) throw new Error("Local Pages CMS mode is disabled.");
  if (cachedConfig) return cachedConfig;

  const configuredRoot = process.env.PAGESCMS_LOCAL_ROOT;
  if (!configuredRoot || !path.isAbsolute(configuredRoot)) {
    throw new Error("PAGESCMS_LOCAL_ROOT must be an absolute path in local Pages CMS mode.");
  }

  const root = await realpath(configuredRoot);
  cachedConfig = {
    root,
    owner: validateSegment("owner", process.env.PAGESCMS_LOCAL_OWNER || "local"),
    repo: validateSegment("repository", process.env.PAGESCMS_LOCAL_REPO || path.basename(root)),
    branch: validateSegment("branch", process.env.PAGESCMS_LOCAL_BRANCH || "working-tree"),
  };
  return cachedConfig;
};

export const assertLocalRepositoryRef = async (
  owner: string,
  repo: string,
  branch?: string,
): Promise<LocalRepositoryConfig> => {
  const config = await getLocalRepositoryConfig();
  if (owner !== config.owner || repo !== config.repo || (branch && branch !== config.branch)) {
    const error = new Error("Local repository not found.") as Error & { status?: number };
    error.status = 404;
    throw error;
  }
  return config;
};

export const resetLocalRepositoryConfigForTests = (): void => {
  cachedConfig = undefined;
};
