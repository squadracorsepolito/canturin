<script lang="ts">
	import { AltArrowIcon } from '../icon';
	import type { TreeNode } from './types';

	type Props = {
		node: TreeNode;
		defaultOpen: boolean;
	};

	let { node, defaultOpen }: Props = $props();

	let open = $state(defaultOpen);
</script>

{#snippet item(name: TreeNode['name'], icon: TreeNode['icon'])}
	<div class="flex items-center gap-2">
		<span>
			<svelte:component this={icon} height="16" width="16" />
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
					<svelte:self node={childNode} {defaultOpen} />
				</li>
			{/each}
		</ul>
	{/if}
{/if}
