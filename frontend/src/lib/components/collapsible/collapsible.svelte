<script lang="ts">
	import type { Snippet } from 'svelte';
	import { AltArrowIcon, CloseIcon } from '../icon';

	type Props = {
		initialCollapsed?: boolean;
		trigger: Snippet;
		content: Snippet;
	};

	let { initialCollapsed, trigger, content }: Props = $props();

	let collapsed = $state(initialCollapsed ?? false);

	function toggle() {
		collapsed = !collapsed;
	}
</script>

<button
	onclick={toggle}
	class="w-full flex justify-between rounded-btn items-center gap-2 p-3 hover:bg-base-content/20 transition-colors"
>
	{@render trigger()}

	{#if collapsed}
		<AltArrowIcon />
	{:else}
		<CloseIcon />
	{/if}
</button>

{#if !collapsed}
	<div class="pt-5">
		{@render content()}
	</div>
{/if}
