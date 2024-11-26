<script lang="ts">
	import type { PanelSectionProps } from '../types';
	import Interface from './interface.svelte';
	import { getNodeState } from './state.svelte';

	let { entityId }: PanelSectionProps = $props();

	const ns = getNodeState(entityId);
</script>

<section>
	<h3 class="pb-5">Interfaces</h3>

	{#if ns.entity.interfaces && ns.entity.interfaces.length > 0}
		{#await ns.getBuses() then buses}
			{#each ns.entity.interfaces as int}
				<Interface {buses} {int} {entityId} />
			{/each}
		{/await}
	{/if}
</section>
