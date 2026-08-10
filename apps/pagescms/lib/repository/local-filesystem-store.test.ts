import { mkdtemp, readFile, rm, symlink, writeFile } from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import { afterEach, describe, expect, it } from "vitest";
import { LocalFilesystemRepositoryStore } from "./local-filesystem-store";

const roots: string[] = [];
const makeStore = async () => {
  const root = await mkdtemp(path.join(os.tmpdir(), "pagescms-store-"));
  roots.push(root);
  return { root, store: new LocalFilesystemRepositoryStore(root) };
};

afterEach(async () => {
  await Promise.all(roots.splice(0).map((root) => rm(root, { recursive: true, force: true })));
});

describe("LocalFilesystemRepositoryStore", () => {
  it("writes, reads, and updates with the expected content digest", async () => {
    const { root, store } = await makeStore();
    const created = await store.write("content/article.json", Buffer.from('{"title":"First"}'));
    expect((await store.read("content/article.json")).content).toBe('{"title":"First"}');
    await store.write("content/article.json", Buffer.from('{"title":"Second"}'), created.file.sha);
    expect(await readFile(path.join(root, "content/article.json"), "utf8")).toBe('{"title":"Second"}');
  });

  it("rejects stale writes", async () => {
    const { store } = await makeStore();
    await store.write("content/article.json", Buffer.from("one"));
    await expect(store.write("content/article.json", Buffer.from("two"), "stale")).rejects.toMatchObject({ status: 409 });
  });

  it("rejects files larger than the local upload limit", async () => {
    const { store } = await makeStore();
    await expect(store.write("media/large.png", Buffer.alloc(10 * 1024 * 1024 + 1))).rejects.toMatchObject({ status: 413 });
  });

  it.each(["../outside", "/absolute", "C:/windows", "folder\\escape", "folder/../escape"])(
    "rejects unsafe path %s",
    async (unsafePath) => {
      const { store } = await makeStore();
      await expect(store.write(unsafePath, Buffer.from("x"))).rejects.toMatchObject({ status: 400 });
    },
  );

  it("rejects symlinks that leave the repository root", async () => {
    const { root, store } = await makeStore();
    const outside = await mkdtemp(path.join(os.tmpdir(), "pagescms-outside-"));
    roots.push(outside);
    await writeFile(path.join(outside, "secret.txt"), "secret");
    await symlink(outside, path.join(root, "linked"), "junction");
    await expect(store.read("linked/secret.txt")).rejects.toMatchObject({ status: 400 });
  });
});
