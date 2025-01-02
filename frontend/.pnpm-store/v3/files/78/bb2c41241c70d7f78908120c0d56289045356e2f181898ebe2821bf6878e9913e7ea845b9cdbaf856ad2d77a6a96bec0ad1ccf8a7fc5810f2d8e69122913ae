import { extractModuleNodes } from './module-nodes.js';
import { extractFragmentNodes } from './fragment-nodes.js';
import { extractInstanceNodes } from './instance-nodes.js';
/**
 * Pick only required Svelte AST nodes for further usage in this addon.
 */
export async function extractSvelteASTNodes(params) {
    const { ast, filename } = params;
    const { module, fragment, instance } = ast;
    const moduleNodes = await extractModuleNodes({ module, filename });
    const instanceNodes = await extractInstanceNodes({
        instance,
        filename,
        moduleNodes,
    });
    const fragmentNodes = await extractFragmentNodes({
        fragment,
        filename,
        moduleNodes,
        instanceNodes,
    });
    return {
        ...moduleNodes,
        ...instanceNodes,
        ...fragmentNodes,
    };
}
