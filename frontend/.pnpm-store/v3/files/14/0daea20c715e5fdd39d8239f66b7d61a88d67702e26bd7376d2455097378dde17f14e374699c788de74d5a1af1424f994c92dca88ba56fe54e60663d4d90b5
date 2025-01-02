import { type ESTreeAST } from '../../../parser/ast';
import type { CompiledASTNodes } from '../../../parser/extract/compiled/nodes';
import type { SvelteASTNodes } from '../../../parser/extract/svelte/nodes';
interface Params {
    nodes: {
        compiled: CompiledASTNodes;
        svelte: SvelteASTNodes;
    };
    filename?: string;
}
/**
 * Replaces `defineMeta({ ... })` with `defineMeta(meta)`,
 * and also it returns {@link ESTreeASTAST.ObjectExpression} which was replaced with {@link ESTreeAST.Identifier}
 */
export declare function replaceDefineMetaArgument(params: Params): ESTreeAST.ObjectExpression;
export {};
