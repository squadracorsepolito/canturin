<script lang="ts">
	import type { PanelSectionProps } from '../types';
	import { loadSignalUnit } from '$lib/panel/signal-unit/signal-unit-state.svelte';
	import Panel from '../panel.svelte';
	import Heading from './heading.svelte';
	import Attributes from './attributes.svelte';
	import Refs from './refs.svelte';

	let { entityId }: PanelSectionProps = $props();

	let promise = $derived(loadSignalUnit(entityId));
</script>

<Panel>
	{#await promise then}
		<Heading {entityId} />

		<Attributes {entityId} />

		<Refs {entityId} />
	{/await}
</Panel>
