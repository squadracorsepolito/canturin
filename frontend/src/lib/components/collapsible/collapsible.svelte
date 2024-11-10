<script lang="ts">
	import type { Snippet } from 'svelte';
	import { IconButton } from '../button';
	import { AltArrowIcon, CloseIcon } from '../icon';

	type Props = {
		children: Snippet;
	};

	let { children }: Props = $props();

	let collapsed = $state(false);
</script>

<div class="flex gap-2">
	<div class:line-clamp-2={collapsed}>
		{@render children()}
	</div>

	<div>
		<IconButton onclick={() => (collapsed = !collapsed)}>
			{#if collapsed}
				<AltArrowIcon height="20" width="20" />
			{:else}
				<CloseIcon height="20" width="20" />
			{/if}
		</IconButton>
	</div>
</div>
