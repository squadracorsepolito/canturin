<script lang="ts">
	import type { Snippet } from 'svelte';
	import TableField from './table-field.svelte';
	import { IconButton } from '../button';
	import { AltArrowIcon, CloseIcon } from '../icon';

	type Props = {
		children?: Snippet;
	};

	let { children }: Props = $props();

	let expanded = $state(false);
</script>

<TableField>
	<div class="flex gap-2">
		<div class:line-clamp-2={!expanded}>
			{@render children?.()}
		</div>

		<div>
			<IconButton onclick={() => (expanded = !expanded)}>
				{#if expanded}
					<CloseIcon height="20" width="20" />
				{:else}
					<AltArrowIcon height="20" width="20" />
				{/if}
			</IconButton>
		</div>
	</div>
</TableField>
