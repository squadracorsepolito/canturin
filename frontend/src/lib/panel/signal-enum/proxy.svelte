<script lang="ts">
	import { loadSignalEnum } from '$lib/panel/signal-enum/state.svelte';
	import Panel from '../panel.svelte';
	import type { PanelSectionProps } from '../types';
	import Attributes from './attributes.svelte';
	import Heading from './heading.svelte';
	import Refs from './refs.svelte';
	import Values from './values.svelte';

	let { entityId }: PanelSectionProps = $props();

	let promise = $derived(loadSignalEnum(entityId));
</script>

<Panel>
	{#await promise then}
		<Heading {entityId} />

		<Attributes {entityId} />

		<Values {entityId} />

		<Refs {entityId} />
	{/await}
</Panel>
