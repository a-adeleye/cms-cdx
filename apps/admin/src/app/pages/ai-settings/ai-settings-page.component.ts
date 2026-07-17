import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, effect, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { createPageActionFeedback } from '../../features/pages/page-feedback';
import { AIProvider } from '../../features/pages/pages.model';
import { WorkspaceStateService } from '../../features/pages/workspace-state.service';

const DEFAULT_MODELS: Record<Exclude<AIProvider, 'none' | 'openai_compatible'>, string> = {
  openai: 'gpt-4.1-mini',
  anthropic: 'claude-sonnet-4-20250514',
  google: 'gemini-2.5-flash',
};

@Component({
  selector: 'app-ai-settings-page',
  templateUrl: './ai-settings-page.component.html',
  styleUrls: ['../../features/pages/page-view.component.css', './ai-settings-page.component.css'],
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AiSettingsPageComponent {
  private readonly fb = inject(FormBuilder);
  readonly state = inject(WorkspaceStateService);
  readonly feedback = createPageActionFeedback();
  readonly form = this.fb.nonNullable.group({
    provider: ['none' as AIProvider, [Validators.required]],
    model: [''],
    apiKeySecretRef: [''],
    baseUrl: [''],
  });

  constructor() {
    effect(() => {
      const config = parseAIConfig(this.state.selectedSite().aiConfig);
      this.form.reset({
        provider: config.provider,
        model: config.model,
        apiKeySecretRef: config.apiKeySecretRef,
        baseUrl: config.baseUrl,
      }, { emitEvent: false });
    });
  }

  isConfigured(): boolean {
    return this.form.controls.provider.value !== 'none';
  }

  isCompatibleProvider(): boolean {
    return this.form.controls.provider.value === 'openai_compatible';
  }

  setProvider(provider: AIProvider): void {
    const current = this.form.getRawValue();
    this.form.patchValue({
      provider,
      model: provider === 'none' ? '' : current.model || DEFAULT_MODELS[provider as keyof typeof DEFAULT_MODELS] || '',
      baseUrl: provider === 'openai_compatible' ? current.baseUrl : '',
    });
  }

  async save(): Promise<void> {
    const value = this.form.getRawValue();
    if (value.provider !== 'none' && (!value.model.trim() || !value.apiKeySecretRef.trim() || (value.provider === 'openai_compatible' && !value.baseUrl.trim()))) {
      this.form.markAllAsTouched();
      this.feedback.error('Choose a model and secret environment variable. Compatible providers also require a base URL.');
      return;
    }

    try {
      this.feedback.loading('Saving AI configuration...');
      const aiConfig = value.provider === 'none'
        ? '{}'
        : JSON.stringify({
            provider: value.provider,
            model: value.model.trim(),
            apiKeySecretRef: value.apiKeySecretRef.trim(),
            ...(value.provider === 'openai_compatible' ? { baseUrl: value.baseUrl.trim() } : {}),
          });
      await this.state.updateSelectedSite({ aiConfig });
      this.feedback.success('AI configuration saved.');
    } catch (error) {
      this.feedback.error(error instanceof Error ? error.message : 'Unable to save AI configuration.');
    }
  }
}

function parseAIConfig(raw: string): { provider: AIProvider; model: string; apiKeySecretRef: string; baseUrl: string } {
  try {
    const config = JSON.parse(raw || '{}') as Record<string, unknown>;
    const provider = config['provider'];
    return {
      provider: provider === 'openai' || provider === 'anthropic' || provider === 'google' || provider === 'openai_compatible' ? provider : 'none',
      model: typeof config['model'] === 'string' ? config['model'] : '',
      apiKeySecretRef: typeof config['apiKeySecretRef'] === 'string' ? config['apiKeySecretRef'] : '',
      baseUrl: typeof config['baseUrl'] === 'string' ? config['baseUrl'] : '',
    };
  } catch {
    return { provider: 'none', model: '', apiKeySecretRef: '', baseUrl: '' };
  }
}
