<script lang="ts">
	import type { PanelSectionProps } from '../types';
	import { loadSignalType } from './state.svelte';
	import Panel from '../panel.svelte';
	import Heading from './heading.svelte';
	import Refs from './refs.svelte';
	import SignalTypeDraft from './draft.svelte';
	import Attributes from './attributes.svelte';

	let { entityId }: PanelSectionProps = $props();

	let promise = $derived(loadSignalType(entityId));
</script>

<Panel>
	{#if entityId === 'draft'}
		<SignalTypeDraft />
	{:else}
		{#await promise then}
			<Heading {entityId} />

			<Attributes {entityId} />

			<Refs {entityId} />
		{/await}
	{/if}
</Panel>
