"use client";

import { useRef, useState } from "react";
import { useFormContext } from "react-hook-form";
import { Sparkles } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { useConfig } from "@/contexts/config-context";
import { AiDraftOptionsSchema, AiDraftSchema } from "@/lib/ai/contracts";

export function EditComponent({ field, contentName, runBeforeSubmitHooks }: any) {
  const [prompt, setPrompt] = useState("");
  const [isGenerating, setIsGenerating] = useState(false);
  const requestSequence = useRef(0);
  const form = useFormContext();
  const { config } = useConfig();

  const generate = async () => {
    const options = AiDraftOptionsSchema.safeParse(field.options);
    if (!options.success || !config || !contentName) {
      toast.error("AI drafting is not configured correctly.");
      return;
    }
    if (prompt.trim().length < 3) {
      toast.error("Describe the article you want to draft.");
      return;
    }

    setIsGenerating(true);
    const sequence = ++requestSequence.current;
    const controller = new AbortController();
    const timeout = window.setTimeout(() => controller.abort(), 35_000);
    try {
      await runBeforeSubmitHooks?.();
      const targets = options.data.targets;
      const paths = Object.values(targets);
      const before = Object.fromEntries(paths.map((path) => [path, form.getValues(path)]));
      const response = await fetch(`/api/${config.owner}/${config.repo}/${encodeURIComponent(config.branch)}/ai/draft`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        signal: controller.signal,
        body: JSON.stringify({ schemaName: contentName, prompt: prompt.trim(), context: before }),
      });
      const payload = await response.json();
      if (!response.ok) throw new Error(payload?.message || "Could not generate a draft.");
      const draft = AiDraftSchema.parse(payload.data);
      if (sequence !== requestSequence.current) return;
      if (paths.some((path) => JSON.stringify(form.getValues(path)) !== JSON.stringify(before[path]))) {
        toast.error("A target field changed while the draft was being generated. Generate again to avoid overwriting it.");
        return;
      }
      const values: Record<string, string> = {
        [targets.title]: draft.title,
        [targets.excerpt]: draft.excerpt,
        [targets.body]: draft.body,
        [targets.seoTitle]: draft.seoTitle,
        [targets.seoDescription]: draft.seoDescription,
      };
      for (const [path, value] of Object.entries(values)) {
        form.setValue(path, value, { shouldDirty: true, shouldTouch: true });
      }
      await form.trigger(paths);
      toast.success("Draft added to the form. Review it before saving.");
    } catch (error: any) {
      toast.error(error?.name === "AbortError" ? "AI generation timed out." : error?.message || "Could not generate a draft.");
    } finally {
      window.clearTimeout(timeout);
      if (sequence === requestSequence.current) setIsGenerating(false);
    }
  };

  return (
    <div className="rounded-lg border bg-muted/30 p-4 space-y-3">
      <Textarea
        value={prompt}
        onChange={(event) => setPrompt(event.target.value)}
        placeholder="Describe the article, audience, and angle…"
        maxLength={4000}
        disabled={isGenerating}
      />
      <div className="flex items-center justify-between gap-3">
        <p className="text-xs text-muted-foreground">Generated content stays editable and is never saved automatically.</p>
        <Button type="button" onClick={generate} disabled={isGenerating || prompt.trim().length < 3}>
          <Sparkles /> {isGenerating ? "Generating…" : "Generate draft"}
        </Button>
      </div>
    </div>
  );
}
