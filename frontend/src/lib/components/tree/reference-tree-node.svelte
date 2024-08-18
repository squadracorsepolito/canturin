<script lang="ts">
	import { getColorByName } from '$lib/utils';
	import type { ReferenceTreeNode } from './types';

	type Props = {
		node: ReferenceTreeNode;
		depth: number;
	};

	let { node, depth }: Props = $props();
</script>

{#snippet item(name: string)}
	{@const color = getColorByName(name)}

	<div
		style:background-color={color.bgColor}
		style:color={color.textColor}
		class="flex items-center justify-center rounded-box"
	>
		<span class="text-sm m-3 truncate text-ellipsis">{name}</span>
	</div>
{/snippet}

{#if depth > 1}
	<div style:grid-column="span {depth} / span {depth}" class="grid gap-3 grid-cols-subgrid py-2">
		{@render item(node.name)}

		<div
			style:grid-column="span {depth - 1} / span {depth - 1}"
			style:grid-template-columns="repeat({depth - 1}, minmax(0, 1fr))"
			class="grid"
		>
			{#each node.childNodes as child}
				<svelte:self node={child} depth={depth - 1} />
			{/each}
		</div>
	</div>
{:else}
	<div class="py-2">
		{@render item(node.name)}
	</div>
{/if}
