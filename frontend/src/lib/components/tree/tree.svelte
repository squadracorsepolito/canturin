<script lang="ts">
	import TreeNode from '$lib/components/tree/tree-node.svelte';
	import type { TreeNode as TreeNodeType } from './types';

	type Props = {
		rootNode: TreeNodeType;
		defaultOpen?: boolean;
	};

	let { rootNode, defaultOpen }: Props = $props();
</script>

<ul class="menu menu-sm">
	<li>
		<TreeNode node={rootNode} defaultOpen={defaultOpen ? defaultOpen : false} />
	</li>
</ul>
