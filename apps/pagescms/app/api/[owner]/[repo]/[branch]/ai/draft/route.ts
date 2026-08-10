import { AiDraftOptionsSchema, AiDraftRequestSchema } from "@/lib/ai/contracts";
import { getDraftProvider } from "@/lib/ai/get-draft-provider";
import { getRepoReadContext } from "@/lib/api-repo-context";
import { createHttpError, toErrorResponse } from "@/lib/api-error";
import { getFieldByPath, getSchemaByName } from "@/lib/schema";

const ALLOWED_TARGET_TYPES = new Set(["string", "text", "rich-text", "code"]);

export async function POST(request: Request, context: { params: Promise<{ owner: string; repo: string; branch: string }> }) {
  try {
    const params = await context.params;
    const { config } = await getRepoReadContext(params);
    const body = AiDraftRequestSchema.safeParse(await request.json());
    if (!body.success) throw createHttpError("Invalid AI draft request.", 400);
    if (JSON.stringify(body.data.context || {}).length > 50_000) throw createHttpError("AI draft context is too large.", 413);
    const schema = getSchemaByName(config.object, body.data.schemaName);
    if (!schema) throw createHttpError("Content schema not found.", 404);
    const assistantField = schema.fields?.find((field: any) => field.type === "ai-draft");
    if (!assistantField) throw createHttpError("AI drafting is not configured for this content type.", 400);
    const options = AiDraftOptionsSchema.safeParse(assistantField.options);
    if (!options.success) throw createHttpError("AI drafting has invalid field mappings.", 400);
    for (const target of Object.values(options.data.targets)) {
      const field = getFieldByPath(schema.fields, target);
      if (!field || field.list || field.readonly || !ALLOWED_TARGET_TYPES.has(field.type)) {
        throw createHttpError("AI drafting has an unsupported target field.", 400);
      }
    }
    const draft = await getDraftProvider().generate({
      prompt: body.data.prompt,
      context: body.data.context,
      instructions: options.data.instructions,
    });
    return Response.json({ status: "success", data: draft });
  } catch (error: any) {
    if (error?.name === "AbortError") return Response.json({ status: "error", message: "AI generation timed out." }, { status: 504 });
    return toErrorResponse(error);
  }
}
