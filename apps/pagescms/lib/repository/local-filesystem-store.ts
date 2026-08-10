import { createHash, randomUUID } from "node:crypto";
import path from "node:path";
import {
  mkdir,
  readFile,
  readdir,
  realpath,
  rename,
  rm,
  stat,
  writeFile,
} from "node:fs/promises";
import type {
  RepositoryEntry,
  RepositoryFile,
  RepositoryStore,
  RepositoryWriteResult,
} from "./contracts";

type StatusError = Error & { status?: number; response?: { data: { message: string } } };

const statusError = (message: string, status: number): StatusError => {
  const error = new Error(message) as StatusError;
  error.status = status;
  error.response = { data: { message } };
  return error;
};

const digest = (content: Buffer): string => createHash("sha256").update(content).digest("hex");
const MAX_LOCAL_FILE_BYTES = 10 * 1024 * 1024;

export class LocalFilesystemRepositoryStore implements RepositoryStore {
  constructor(private readonly root: string) {}

  private async resolvePath(repositoryPath: string, allowMissing = false): Promise<string> {
    if (
      !repositoryPath ||
      repositoryPath.includes("\0") ||
      repositoryPath.includes("\\") ||
      path.posix.isAbsolute(repositoryPath) ||
      /^[a-zA-Z]:/.test(repositoryPath) ||
      repositoryPath.split("/").some((segment) => segment === "..")
    ) {
      throw statusError("Invalid repository path.", 400);
    }

    const normalized = path.posix.normalize(repositoryPath).replace(/^\.\//, "");
    const candidate = path.resolve(this.root, ...normalized.split("/"));
    const relative = path.relative(this.root, candidate);
    if (!relative || relative.startsWith("..") || path.isAbsolute(relative)) {
      if (relative === "") return candidate;
      throw statusError("Repository path escapes the configured root.", 400);
    }

    let ancestor = candidate;
    if (allowMissing) {
      while (ancestor !== this.root) {
        try {
          await stat(ancestor);
          break;
        } catch {
          ancestor = path.dirname(ancestor);
        }
      }
    }

    try {
      const canonical = await realpath(ancestor);
      const canonicalRelative = path.relative(this.root, canonical);
      if (canonicalRelative.startsWith("..") || path.isAbsolute(canonicalRelative)) {
        throw statusError("Repository path resolves outside the configured root.", 400);
      }
    } catch (error) {
      if ((error as StatusError).status) throw error;
      if (!allowMissing) throw statusError("Repository entry not found.", 404);
    }

    return candidate;
  }

  private async toFile(repositoryPath: string, content?: Buffer): Promise<RepositoryFile> {
    const absolutePath = await this.resolvePath(repositoryPath);
    const value = content ?? await readFile(absolutePath);
    return {
      type: "file",
      name: path.posix.basename(repositoryPath),
      path: repositoryPath,
      sha: digest(value),
      size: value.byteLength,
      content: value.toString("utf8"),
      bytes: value,
      downloadUrl: `/api/local-media?path=${encodeURIComponent(repositoryPath)}`,
    };
  }

  async read(repositoryPath: string): Promise<RepositoryFile> {
    return this.toFile(repositoryPath);
  }

  async list(repositoryPath: string): Promise<RepositoryEntry[]> {
    const cleanPath = repositoryPath === "." ? "" : repositoryPath.replace(/\/$/, "");
    const absolutePath = cleanPath ? await this.resolvePath(cleanPath) : this.root;
    let entries;
    try {
      entries = await readdir(absolutePath, { withFileTypes: true });
    } catch {
      throw statusError("Repository directory not found.", 404);
    }

    return Promise.all(entries.map(async (entry): Promise<RepositoryEntry> => {
      const entryPath = cleanPath ? `${cleanPath}/${entry.name}` : entry.name;
      if (entry.isDirectory()) {
        return { name: entry.name, path: entryPath, parentPath: cleanPath, type: "dir" };
      }
      if (!entry.isFile()) {
        return { name: entry.name, path: entryPath, parentPath: cleanPath, type: "dir" };
      }
      const file = await this.toFile(entryPath);
      return { ...file, parentPath: cleanPath };
    }));
  }

  async write(repositoryPath: string, content: Buffer, expectedSha?: string): Promise<RepositoryWriteResult> {
    if (content.byteLength > MAX_LOCAL_FILE_BYTES) throw statusError("File exceeds the 10 MB local limit.", 413);
    const absolutePath = await this.resolvePath(repositoryPath, true);
    let current: Buffer | undefined;
    try {
      current = await readFile(absolutePath);
    } catch {
      current = undefined;
    }

    if (current && !expectedSha) throw statusError("File already exists.", 422);
    if (!current && expectedSha) throw statusError("File no longer exists.", 409);
    if (current && expectedSha !== digest(current)) throw statusError("File has changed since it was loaded.", 409);

    await mkdir(path.dirname(absolutePath), { recursive: true });
    const temporaryPath = `${absolutePath}.pagescms-${randomUUID()}.tmp`;
    await writeFile(temporaryPath, content, { flag: "wx" });
    await rename(temporaryPath, absolutePath);
    const file = await this.toFile(repositoryPath, content);
    const committedAt = new Date().toISOString();
    return { file, commitSha: file.sha, committedAt };
  }

  async delete(repositoryPath: string, expectedSha: string): Promise<RepositoryWriteResult> {
    const file = await this.read(repositoryPath);
    if (file.sha !== expectedSha) throw statusError("File has changed since it was loaded.", 409);
    const absolutePath = await this.resolvePath(repositoryPath);
    await rm(absolutePath);
    return { file, commitSha: file.sha, committedAt: new Date().toISOString() };
  }

  async rename(repositoryPath: string, newRepositoryPath: string, expectedSha: string): Promise<RepositoryWriteResult> {
    const file = await this.read(repositoryPath);
    if (file.sha !== expectedSha) throw statusError("File has changed since it was loaded.", 409);
    const source = await this.resolvePath(repositoryPath);
    const destination = await this.resolvePath(newRepositoryPath, true);
    try {
      await stat(destination);
      throw statusError("Destination already exists.", 409);
    } catch (error) {
      if ((error as StatusError).status) throw error;
    }
    await mkdir(path.dirname(destination), { recursive: true });
    await rename(source, destination);
    const renamed = await this.read(newRepositoryPath);
    return { file: renamed, commitSha: renamed.sha, committedAt: new Date().toISOString() };
  }
}
