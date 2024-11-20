<script lang="ts">
	import type { PanelSectionProps } from '../types';
	import { loadSignalUnit } from '$lib/state/signal-unit-state.svelte';
	import Panel from '../panel.svelte';
	import SignalUnitPanel from './signal-unit-panel.svelte'

	let { entityId }: PanelSectionProps = $props();

	let promise = $derived(loadSignalUnit(entityId));
</script>

<Panel>
	{#await promise then}
		<SignalUnitPanel {entityId} />
	{/await}
</Panel>
