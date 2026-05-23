export type SummaryMetric = {
  label: string;
  value: string;
  detail: string;
};

export type SettingsLink = {
  label: string;
  description: string;
  path: string;
};

export type SettingsGroup = {
  title: string;
  description: string;
  links: SettingsLink[];
};
