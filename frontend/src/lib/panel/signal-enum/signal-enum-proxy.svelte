<script lang="ts">
	import { loadSignalEnum } from '$lib/state/signal-enum-state.svelte';
	import Panel from '../panel.svelte';
	import type { PanelSectionProps } from '../types';
	import SignalEnumAttributes from './signal-enum-attributes.svelte';
	import SignalEnumDraft from './signal-enum-draft.svelte';
	import SignalEnumHeading from './signal-enum-heading.svelte';
	import SignalEnumRefs from './signal-enum-refs.svelte';
	import SignalEnumValues from './signal-enum-values.svelte';

	let { entityId }: PanelSectionProps = $props();

	let promise = $derived(loadSignalEnum(entityId));
</script>

<Panel>
	{#if entityId === 'draft'}
		<SignalEnumDraft />
	{:else}
		{#await promise then}
			<SignalEnumHeading {entityId} />

			<SignalEnumAttributes {entityId} />

			<SignalEnumValues {entityId} />

			<SignalEnumRefs {entityId} />
		{/await}
	{/if}
</Panel>
