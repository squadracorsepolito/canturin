<script lang="ts">
	import Panel from '../panel.svelte';
	import Buses from './buses.svelte';
	import Heading from './heading.svelte';
	import { getNetwork } from './state.svelte';
</script>

<Panel>
	{#await getNetwork() then}
		<Heading />

		<Buses />
	{/await}
</Panel>
