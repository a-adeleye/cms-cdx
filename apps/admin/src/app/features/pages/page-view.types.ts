export type SummaryMetric = {
  label: string;
  value: string;
  detail: string;
};

export type SettingsLink = {
  label: string;
  path: string;
};

export type SettingsGroup = {
  title: string;
  links: SettingsLink[];
};
