import { getLocalRepositoryConfig } from "@/lib/local-dev/config-server";
import { LocalFilesystemRepositoryStore } from "./local-filesystem-store";

let store: LocalFilesystemRepositoryStore | undefined;

export const getLocalRepositoryStore = async (): Promise<LocalFilesystemRepositoryStore> => {
  if (store) return store;
  const config = await getLocalRepositoryConfig();
  store = new LocalFilesystemRepositoryStore(config.root);
  return store;
};
