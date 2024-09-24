<script lang="ts">
	import { AltArrowIcon } from '../icon';
	import type { TreeNode } from './types';
	import Self from './tree-node.svelte';

	type Props = {
		node: TreeNode;
		defaultOpen: boolean;
	};

	let { node, defaultOpen }: Props = $props();

	let open = $state(defaultOpen);
</script>

{#snippet item(name: TreeNode['name'], icon: TreeNode['icon'])}
	{@const Icon = icon}

	<div class="flex items-center gap-2">
		<span>
			<Icon height="16" width="16" />
		</span>
		<span class="text-xs">{name}</span>
	</div>
{/snippet}

{#if node.childNodes.length === 0}
	<button onclick={node.onclick}>
		{@render item(node.name, node.icon)}
	</button>
{:else}
	<button onclick={() => (open = !open)} ondblclick={node.onclick}>
		<span class:-rotate-90={!open}>
			<AltArrowIcon height="16" width="16" />
		</span>

		{@render item(node.name, node.icon)}
	</button>

	{#if open}
		<ul class="menu menu-sm">
			{#each node.childNodes as childNode}
				<li>
					<Self node={childNode} {defaultOpen} />
				</li>
			{/each}
		</ul>
	{/if}
{/if}
