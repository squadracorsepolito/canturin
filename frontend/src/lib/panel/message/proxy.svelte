<script lang="ts">
	import Panel from '../panel.svelte';
	import type { PanelSectionProps } from '../types';
	import Attributes from './attributes.svelte';
	import Draft from './draft.svelte';
	import Heading from './heading.svelte';
	import Signals from './signals.svelte';
	import { loadMessage } from './state.svelte';

	let { entityId }: PanelSectionProps = $props();

	let promise = $derived(loadMessage(entityId));
</script>

<Panel>
	{#if entityId === 'draft'}
		<Draft />
	{:else}
		{#await promise then}
			<Heading {entityId} />

			<Attributes {entityId} />

			<Signals {entityId} />
		{/await}
	{/if}
</Panel>
