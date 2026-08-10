import type { User } from "@/types/user";

type UserLike = Pick<User, "githubUsername" | "isLocalDevelopment"> | null | undefined;

const hasGithubIdentity = (user: UserLike): boolean => Boolean(user?.githubUsername);

const canManageConfiguration = (user: UserLike): boolean =>
  hasGithubIdentity(user) || Boolean(user?.isLocalDevelopment);

const assertGithubIdentity = (
  user: UserLike,
  message = "Only GitHub users can perform this action.",
) => {
  if (!hasGithubIdentity(user)) {
    throw new Error(message);
  }
};

export { hasGithubIdentity, assertGithubIdentity };
export { canManageConfiguration };
