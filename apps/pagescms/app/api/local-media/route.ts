import { getLocalRepositoryStore } from "@/lib/repository/local-store-server";
import { isLocalDevelopmentMode } from "@/lib/local-dev/mode";
import { requireApiUserSession } from "@/lib/session-server";
import { getLocalRepositoryConfig } from "@/lib/local-dev/config-server";
import { getConfig } from "@/lib/config-store";
import { isPathWithin } from "@/lib/utils/file";

const MIME_TYPES: Record<string, string> = {
  png: "image/png",
  jpg: "image/jpeg",
  jpeg: "image/jpeg",
  gif: "image/gif",
  webp: "image/webp",
};

export async function GET(request: Request) {
  if (!isLocalDevelopmentMode()) return new Response(null, { status: 404 });
  const session = await requireApiUserSession();
  if ("response" in session) return session.response;
  const repositoryPath = new URL(request.url).searchParams.get("path");
  if (!repositoryPath) return Response.json({ message: "Missing path." }, { status: 400 });

  try {
    const local = await getLocalRepositoryConfig();
    const config = await getConfig(local.owner, local.repo, local.branch);
    const mediaRoots = Array.isArray(config?.object?.media)
      ? config.object.media.map((item: any) => item.input).filter((item: unknown): item is string => typeof item === "string")
      : [];
    if (!mediaRoots.some((root) => isPathWithin(repositoryPath, root))) return new Response(null, { status: 403 });
    const file = await (await getLocalRepositoryStore()).read(repositoryPath);
    const extension = file.name.split(".").pop()?.toLowerCase() || "";
    if (!MIME_TYPES[extension]) return new Response(null, { status: 415 });
    return new Response(file.bytes, {
      headers: {
        "Content-Type": MIME_TYPES[extension],
        "Cache-Control": "no-store",
        "X-Content-Type-Options": "nosniff",
      },
    });
  } catch (error: any) {
    return new Response(null, { status: error?.status === 400 ? 400 : 404 });
  }
}
