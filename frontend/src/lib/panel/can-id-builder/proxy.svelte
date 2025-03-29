<script lang="ts">
	import PanelHeading from '../panel-heading.svelte';
	import PanelRefs from '../panel-refs.svelte';
	import Panel from '../panel.svelte';
	import type { PanelSectionProps } from '../types';
	import Operations from './operations.svelte';
	import { loadCanIdBuilder } from './state.svete';

	let { entityId }: PanelSectionProps = $props();
</script>

<Panel>
	{#await loadCanIdBuilder(entityId) then s}
		<PanelHeading entityState={s} />

		<Operations {s} />

		<PanelRefs reference={s.entity.reference} />
	{/await}
</Panel>
