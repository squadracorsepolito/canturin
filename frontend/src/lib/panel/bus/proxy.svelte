<script lang="ts">
	import Panel from '../panel.svelte';
	import type { PanelSectionProps } from '../types';
	import AttachedNodes from './attached-nodes.svelte';
	import Attributes from './attributes.svelte';
	import Draft from './draft.svelte';
	import Heading from './heading.svelte';
	import { loadBus } from './state.svelte';

	let { entityId }: PanelSectionProps = $props();

	let promise = $derived(loadBus(entityId));
</script>

<Panel>
	{#if entityId === 'draft'}
		<Draft />
	{:else}
		{#await promise then}
			<Heading {entityId} />

			<Attributes {entityId} />

			<AttachedNodes {entityId} />
		{/await}
	{/if}
</Panel>
