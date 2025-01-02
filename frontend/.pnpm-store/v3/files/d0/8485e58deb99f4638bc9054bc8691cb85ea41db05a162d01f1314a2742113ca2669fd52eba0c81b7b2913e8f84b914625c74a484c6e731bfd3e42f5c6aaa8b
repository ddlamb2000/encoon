import { extractModuleNodes } from './module-nodes.js';
import { extractFragmentNodes } from './fragment-nodes.js';
import { extractInstanceNodes } from './instance-nodes.js';
import type { SvelteAST } from '../../ast.js';
/**
 * Selected nodes extracted from the Svelte AST via `svelte.compile`,
 * needed for further code analysis/transformation.
 */
export type SvelteASTNodes = Awaited<ReturnType<typeof extractModuleNodes>> & Awaited<ReturnType<typeof extractFragmentNodes>> & Awaited<ReturnType<typeof extractInstanceNodes>>;
interface Params {
    ast: SvelteAST.Root;
    filename?: string;
}
/**
 * Pick only required Svelte AST nodes for further usage in this addon.
 */
export declare function extractSvelteASTNodes(params: Params): Promise<SvelteASTNodes>;
export {};
