<script lang="ts">
	import { uniqueId } from '$lib/utils';
	import * as hoverCard from '@zag-js/hover-card';
	import { portal, useMachine, normalizeProps } from '@zag-js/svelte';
	import type { Snippet } from 'svelte';

	type Props = {
		placement?: hoverCard.Placement;
		trigger: Snippet;
		content: Snippet;
	};

	let { trigger, content, placement }: Props = $props();

	const [snapshot, send] = useMachine(
		hoverCard.machine({
			id: uniqueId(),
			openDelay: 1300,
			closeDelay: 200,
			positioning: {
				offset: {
					mainAxis: 2
				},
				placement: placement
			}
		})
	);

	const api = $derived(hoverCard.connect(snapshot, send, normalizeProps));
</script>

<div {...api.getTriggerProps()}>
	{@render trigger()}
</div>

{#if api.open}
	<div use:portal {...api.getPositionerProps()}>
		<div {...api.getContentProps()}>
			{@render content()}
		</div>
	</div>
{/if}

<style lang="postcss">
	[data-scope='hover-card'] {
		&[data-part='trigger'] {
			@apply inline-block;
		}

		&[data-part='content'] {
			@apply rounded-btn bg-base-100 px-3 py-2 border-2 border-secondary max-w-80;
		}
	}
</style>
