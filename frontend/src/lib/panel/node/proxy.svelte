<script lang="ts">
	import Panel from '../panel.svelte';
	import type { PanelSectionProps } from '../types';
	import Heading from './heading.svelte';
	import { loadNode } from './state.svelte';

	let { entityId }: PanelSectionProps = $props();

	let promise = $derived(loadNode(entityId));
</script>

<Panel>
	{#await promise then}
		<Heading {entityId} />
	{/await}
</Panel>
