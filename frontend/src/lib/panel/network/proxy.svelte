<script lang="ts">
	import Panel from '../panel.svelte';
	import Buses from './buses.svelte';
	import Heading from './heading.svelte';
	import { loadNetwork } from './state.svelte';
</script>

<Panel>
	{#await loadNetwork() then}
		<Heading />

		<Buses />
	{/await}
</Panel>
