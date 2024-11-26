<script lang="ts">
	import Panel from '../panel.svelte';
	import type { PanelSectionProps } from '../types';
	import Attributes from './attributes.svelte';
	import Heading from './heading.svelte';
	import Interfaces from './interfaces.svelte';
	import { loadNode } from './state.svelte';

	let { entityId }: PanelSectionProps = $props();

	let promise = $derived(loadNode(entityId));
</script>

<Panel>
	{#await promise then}
		<Heading {entityId} />

		<Attributes {entityId} />

		<Interfaces {entityId} />
	{/await}
</Panel>
