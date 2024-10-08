<script lang="ts">
	import type { PanelSectionProps } from '../types';
	import { loadSignalType } from '../../state/signal-type-state.svelte';
	import Panel from '../panel.svelte';
	import SignalTypeHeading from './signal-type-heading.svelte';
	import SignalTypeRefs from './signal-type-refs.svelte';
	import SignalTypeDraft from './signal-type-draft.svelte';
	import SignalTypeAttributes from './signal-type-attributes.svelte';

	let { entityId }: PanelSectionProps = $props();

	let promise = $derived(loadSignalType(entityId));
</script>

<Panel>
	{#if entityId === 'draft'}
		<SignalTypeDraft />
	{:else}
		{#await promise then}
			<SignalTypeHeading {entityId} />

			<SignalTypeAttributes {entityId} />

			<SignalTypeRefs {entityId} />
		{/await}
	{/if}
</Panel>
