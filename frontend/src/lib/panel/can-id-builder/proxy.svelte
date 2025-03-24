<script lang="ts">
	import PanelHeading from '../panel-heading.svelte';
	import Panel from '../panel.svelte';
	import type { PanelSectionProps } from '../types';
	import Operations from './operations.svelte';
	import { loadCanIdBuilder } from './state.svete';

	let { entityId }: PanelSectionProps = $props();

	let promise = $derived(loadCanIdBuilder(entityId));
</script>

<Panel>
	{#await promise then s}
		<PanelHeading entityState={s} />

		<Operations {s} />
	{/await}
</Panel>
