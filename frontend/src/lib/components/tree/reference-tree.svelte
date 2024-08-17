<script lang="ts">
	import ReferenceTreeNode from './reference-tree-node.svelte';
	import type { ReferenceTreeNode as ReferenceTreeNodeType } from './types';

	type Props = {
		siblingNodes: ReferenceTreeNodeType[];
		depth: number;
	};

	let { siblingNodes, depth }: Props = $props();
</script>

<div style:grid-template-columns="repeat({depth}, minmax(0, 1fr))" class="grid gap-y-5">
	{#each siblingNodes as node}
		<ReferenceTreeNode {node} {depth} />
	{/each}
</div>
