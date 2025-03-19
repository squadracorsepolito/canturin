<script lang="ts">
	import type { BaseSignal, Signal, SignalKind } from '$lib/api/canturin';
	import { openPanel } from '$lib/panel/panel-stack-state.svelte';
	import { SignalKindBadge } from '../badge';
	import { LinkButton } from '../button';
	import HoverPreview from './hover-preview.svelte';

	type Props = {
		signal: BaseSignal;
	};

	let { signal }: Props = $props();
</script>

<HoverPreview>
	{#snippet trigger()}
		<LinkButton
			label={signal.name}
			onclick={() => openPanel('signal', signal.entityId, signal.name)}
		/>
	{/snippet}

	{#snippet content()}
		<div>
			<span class="font-medium text-sm pr-1">{signal.name}</span>

			<span>
				<SignalKindBadge signalKind={signal.kind} />
			</span>
		</div>

		{#if signal.desc}
			<div class="text-xs text-dimmed pt-1">{signal.desc}</div>
		{/if}
	{/snippet}
</HoverPreview>
