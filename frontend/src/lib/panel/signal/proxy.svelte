<script lang="ts">
	import Panel from '../panel.svelte';
	import type { PanelSectionProps } from '../types';
	import Attributes from './attributes.svelte';
	import Heading from './heading.svelte';
	import { loadSignal } from './state.svelte';

	let { entityId }: PanelSectionProps = $props();

	let promise = $derived(loadSignal(entityId));
</script>

<Panel>
	{#if entityId === 'draft'}
		draft
	{:else}
		{#await promise then}
			<Heading {entityId} />

			<Attributes {entityId} />
		{/await}
	{/if}
</Panel>
