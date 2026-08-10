import { z } from "zod";
import { EditComponent } from "./edit-component";

const label = "AI draft assistant";
const schema = () => z.null().optional();
const defaultValue = null;
const write = () => undefined;

export { label, schema, defaultValue, write, EditComponent };
