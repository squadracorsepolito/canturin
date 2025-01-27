<script lang="ts">
	import { colorByName } from '$lib/actions/color-name.svelte';
	import { ReferenceKind, type Reference } from '$lib/api/canturin';
	import layoutStateSvelte, { type PanelType } from '$lib/state/layout-state.svelte';

	type Props = {
		nodes: Reference[];
		depth: number;
	};

	let { nodes, depth }: Props = $props();

	function handleClick(node: Reference) {
		let panelType: PanelType = 'bus';
		let entityId = node.entityId;

		switch (node.kind) {
			case ReferenceKind.ReferenceKindBus:
				panelType = 'bus';
				break;

			case ReferenceKind.ReferenceKindNode:
				panelType = 'node';
				break;

			case ReferenceKind.ReferenceKindMessage:
				panelType = 'message';
				break;

			case ReferenceKind.ReferenceKindSignal:
				panelType = 'signal';
				console.log(entityId);
				break;
		}

		layoutStateSvelte.openPanel(panelType, entityId);
	}
</script>

{#snippet treeNodeItem(node: Reference)}
	<div>
		<button
			use:colorByName={{ name: node.name }}
			onclick={() => handleClick(node)}
			class="text-sm py-3 truncate overflow-x-hidden w-full h-full rounded-box hover:ring-2 ring-primary"
		>
			{node.name}
		</button>
	</div>
{/snippet}

{#snippet treeNode(node: Reference, depth: number)}
	{#if depth > 1}
		<div style:grid-column="span {depth} / span {depth}" class="grid gap-3 grid-cols-subgrid py-2">
			{@render treeNodeItem(node)}

			<div
				style:grid-column="span {depth - 1} / span {depth - 1}"
				style:grid-template-columns="repeat({depth - 1}, minmax(0, 1fr))"
				class="grid"
			>
				{#if node.children}
					{#each node.children as child}
						{@render treeNode(child, depth - 1)}
					{/each}
				{/if}
			</div>
		</div>
	{:else}
		<div class="py-2">
			{@render treeNodeItem(node)}
		</div>
	{/if}
{/snippet}

<div style:grid-template-columns="repeat({depth}, minmax(0, 1fr))" class="grid gap-y-5">
	{#each nodes as node}
		{@render treeNode(node, depth)}
	{/each}
</div>
