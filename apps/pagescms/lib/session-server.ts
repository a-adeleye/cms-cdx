import { cache } from "react";
import { headers } from "next/headers";
import { isLocalDevelopmentMode } from "@/lib/local-dev/mode";

const localSession = {
  session: {
    id: "local-development-session",
    userId: "local-development-user",
    expiresAt: new Date("2099-01-01T00:00:00.000Z"),
  },
  user: {
    id: "local-development-user",
    email: "local@pagescms.test",
    name: "Local editor",
    image: null,
    emailVerified: true,
    githubUsername: null,
    isLocalDevelopment: true,
  },
};

const getAuthSession = async () => {
  const { auth } = await import("@/lib/auth");
  return auth.api.getSession({ headers: await headers() });
};

const getServerSession = cache(async () => {
  if (isLocalDevelopmentMode()) return localSession;
  return getAuthSession();
});

const requireApiUserSession = async () => {
  const session = isLocalDevelopmentMode() ? localSession : await getAuthSession();
  if (!session?.user) {
    return { response: new Response(null, { status: 401 }) };
  }

  return { user: session.user };
};

export { getServerSession, requireApiUserSession };
