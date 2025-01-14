<script lang="ts">
	import { uniqueId } from '$lib/utils';
	import * as hoverCard from '@zag-js/hover-card';
	import { portal, useMachine, normalizeProps } from '@zag-js/svelte';
	import type { Snippet } from 'svelte';

	type Props = {
		trigger: Snippet;
		content: Snippet;
	};

	let { trigger, content }: Props = $props();

	const [snapshot, send] = useMachine(
		hoverCard.machine({
			id: uniqueId()
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
			<div {...api.getArrowProps()}>
				<div {...api.getArrowTipProps()}></div>
			</div>

			{@render content()}
		</div>
	</div>
{/if}
