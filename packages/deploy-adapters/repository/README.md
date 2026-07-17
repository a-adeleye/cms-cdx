# Repository deployment adapter

Use `git_repository` when the CMS blog must live inside an existing landing-site repository.

Configuration fields:

```json
{
  "repositoryUrl": "https://github.com/example/landing-site.git",
  "branch": "main",
  "contentPath": "public/blog",
  "tokenSecretRef": "GITHUB_PUSH_TOKEN",
  "publicUrl": "https://www.example.com/blog/"
}
```

`repositoryUrl` must be an HTTPS URL on the API server's Git host allowlist (GitHub and GitLab by default) and cannot contain credentials. `branch` is cloned and pushed without force. `contentPath` is the only directory replaced; it must be relative, remain inside the checkout, and cannot contain `.git`. `tokenSecretRef` names an API-server environment variable and may be omitted when the runtime already has a non-interactive Git credential. `publicUrl` is optional and is recorded in deployment history.

For a Cloudflare Pages Git project, configure Pages to build the same branch. A production push normally targets `main`; previews can target a separate branch such as `preview`. The Git host is responsible for starting and reporting the downstream Pages build.

Recovery is a normal Git revert of the recorded `deployRevision`, followed by a push to the configured branch. The CMS never force-pushes.
