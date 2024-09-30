<script lang="ts">
	import MessagePanel from '$lib/panel/message-panel.svelte';
	import SignalTypePanel from '$lib/panel/signal-type/signal-type-panel.svelte';
	import SignalUnitPanel from '$lib/panel/signal-unit-panel.svelte';
	import { getLayoutState } from '$lib/state/layout-state.svelte';

	const layout = getLayoutState();
</script>

{#if layout.openPanelType === 'message'}
	<MessagePanel entityId={layout.openPanelId} />
{:else if layout.openPanelType === 'signal_type'}
	<SignalTypePanel entityId={layout.openPanelId} />
{:else if layout.openPanelType === 'signal_unit'}
	<SignalUnitPanel entityId={layout.openPanelId} />
{/if}
