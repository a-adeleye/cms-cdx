const LOCAL_MODE_ENV = "PAGESCMS_LOCAL_MODE";

export const isLocalDevelopmentMode = (): boolean => {
  const requested = process.env[LOCAL_MODE_ENV] === "true";
  if (!requested) return false;

  if (process.env.NODE_ENV === "production") {
    throw new Error(`${LOCAL_MODE_ENV}=true is only allowed when NODE_ENV=development.`);
  }

  return process.env.NODE_ENV === "development" || process.env.NODE_ENV === "test";
};

export const LOCAL_DEVELOPMENT_TOKEN = "pagescms-local-development";
