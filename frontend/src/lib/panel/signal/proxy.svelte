<script lang="ts">
	import { SignalKind } from '$lib/api/canturin';
	import Panel from '../panel.svelte';
	import type { PanelSectionProps } from '../types';
	import Attributes from './attributes.svelte';
	import Heading from './heading.svelte';
	import MultiplexerSignals from './multiplexer-signals.svelte';
	import { loadSignal } from './state.svelte';

	let { entityId }: PanelSectionProps = $props();

	let promise = $derived(loadSignal(entityId));
</script>

<Panel>
	{#await promise then signal}
		<Heading {entityId} />

		<Attributes {entityId} />

		{#if signal.kind === SignalKind.SignalKindMultiplexer}
			<MultiplexerSignals {entityId} signal={signal.multiplexer} />
		{/if}
	{/await}
</Panel>
