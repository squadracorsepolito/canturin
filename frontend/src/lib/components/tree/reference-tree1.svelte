<script lang="ts">
	import { colorByName } from '$lib/actions/color-name.svelte';
	import type { Ref } from '$lib/api/canturin';
	import { openPanelFromEntity } from '$lib/panel/panel-stack-state.svelte';

	type Props = {
		root: Ref;
	};

	let { root }: Props = $props();

	let depth = $derived.by(() => {
		const stack = [root];

		let longest = 0;
		let tmpDepth = 0;

		while (stack.length > 0) {
			const node = stack.pop();
			if (!node) continue;

			tmpDepth++;

			if (!node.children || node.children.length === 0) {
				if (tmpDepth > longest) {
					longest = tmpDepth;
				}

				tmpDepth = 0;

				continue;
			}

			for (const child of node.children) {
				stack.push(child);
			}
		}

		return longest;
	});

	function handleNodeClick(node: Ref) {
		openPanelFromEntity(node);
	}
</script>

{#snippet treeNodeItem(node: Ref)}
	<div>
		<button
			use:colorByName={{ name: node.name }}
			onclick={() => handleNodeClick(node)}
			class="text-sm py-3 truncate overflow-x-hidden w-full h-full rounded-box hover:ring-2 ring-primary"
		>
			{node.name}
		</button>
	</div>
{/snippet}

{#snippet treeNode(node: Ref, currDepth: number)}
	{#if currDepth > 1}
		<div
			style:grid-column="span {currDepth} / span {currDepth}"
			class="grid gap-3 grid-cols-subgrid py-2"
		>
			{@render treeNodeItem(node)}

			<div
				style:grid-column="span {currDepth - 1} / span {currDepth - 1}"
				style:grid-template-columns="repeat({currDepth - 1}, minmax(0, 1fr))"
				class="grid"
			>
				{#each node.children || [] as child}
					{@render treeNode(child, currDepth - 1)}
				{/each}
			</div>
		</div>
	{:else}
		<div class="py-2">
			{@render treeNodeItem(node)}
		</div>
	{/if}
{/snippet}

<div style:grid-template-columns="repeat({depth}, minmax(0, 1fr))" class="grid gap-y-5">
	{@render treeNode(root, depth)}
</div>
