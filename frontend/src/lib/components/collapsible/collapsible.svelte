<script lang="ts">
	import type { Snippet } from 'svelte';
	import { AltArrowIcon, CloseIcon } from '../icon';

	type Props = {
		initialCollapsed?: boolean;
		raw?: boolean;
		trigger: Snippet<[{ collapsed: boolean }]>;
		content: Snippet;
	};

	let { initialCollapsed, raw, trigger, content }: Props = $props();

	let collapsed = $state(initialCollapsed ?? false);

	function toggle() {
		collapsed = !collapsed;
	}
</script>

<button
	onclick={toggle}
	class={[
		'w-full rounded-btn hover:bg-base-content/20 transition-colors',
		!raw && 'p-3 flex justify-between items-center gap-2'
	]}
>
	{@render trigger({ collapsed })}

	{#if !raw}
		{#if collapsed}
			<AltArrowIcon />
		{:else}
			<CloseIcon />
		{/if}
	{/if}
</button>

{#if !collapsed}
	<div class={[!raw && 'pt-5']}>
		{@render content()}
	</div>
{/if}
