# Cloudflare Pages Deploy Adapter

The API builds the selected site's published CMS content into static HTML and uploads the resulting directory to a Cloudflare Pages Direct Upload project through Wrangler.

## Site Configuration

Set the site's deployment provider to `cloudflare_pages` and use JSON containing the existing Pages project name. No token or account identifier belongs in this JSON.

```json
{
  "projectName": "example-site",
  "productionBranch": "main"
}
```

`productionBranch` defaults to `main`. Preview builds upload to the `cms-preview` branch. The Pages project must already be a Direct Upload project; create it once with `wrangler pages project create <project-name> --production-branch main`.

Cloudflare Pages builds currently publish the CMS's `/articles` route. Use `/articles` as the site's blog path when configuring this provider.

## Runtime Configuration

Configure these values in the API deployment platform, not in site settings or source control:

- `CLOUDFLARE_API_TOKEN`: Cloudflare API token with **Account / Cloudflare Pages / Edit** permission for the target account.
- `CLOUDFLARE_ACCOUNT_ID`: Cloudflare account ID that owns the Pages project.
- `WRANGLER_COMMAND`: optional Wrangler executable path. The API container defaults to its pinned bundled Wrangler binary.

The token is passed only to Wrangler's child process. Deployment logs redact that value before they are stored with the build record.

## Verify and Roll Back

After a published build, the API checks the Pages deployment URL for the home page, `/articles/`, and one generated article route before recording success. Verify a custom domain separately after DNS propagation. Cloudflare Pages retains prior deployments; roll back by promoting the prior production deployment in the Pages dashboard.

