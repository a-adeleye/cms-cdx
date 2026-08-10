export type RepositoryEntry = {
  name: string;
  path: string;
  parentPath: string;
  type: "file" | "dir";
  sha?: string;
  size?: number;
  content?: string;
  downloadUrl?: string;
};

export type RepositoryFile = Required<Pick<RepositoryEntry, "name" | "path" | "sha" | "size">> & {
  type: "file";
  content: string;
  bytes: Buffer;
  downloadUrl?: string;
};

export type RepositoryWriteResult = {
  file: RepositoryFile;
  commitSha: string;
  committedAt: string;
};

export interface RepositoryStore {
  read(path: string): Promise<RepositoryFile>;
  list(path: string): Promise<RepositoryEntry[]>;
  write(path: string, content: Buffer, expectedSha?: string): Promise<RepositoryWriteResult>;
  delete(path: string, expectedSha: string): Promise<RepositoryWriteResult>;
  rename(path: string, newPath: string, expectedSha: string): Promise<RepositoryWriteResult>;
}
