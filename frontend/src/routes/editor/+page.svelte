<script lang="ts">
	import { SignalEnumPanel, SignalTypePanel } from '$lib/panel';
	import { BusPanel } from '$lib/panel/bus';
	import MessagePanel from '$lib/panel/message-panel.svelte';
	import SignalUnitPanel from '$lib/panel/signal-unit-panel.svelte';
	import layout from '$lib/state/layout-state.svelte';
</script>

{#if layout.openPanelType === 'bus'}
	<BusPanel entityId={layout.openPanelId} />
{:else if layout.openPanelType === 'message'}
	<MessagePanel entityId={layout.openPanelId} />
{:else if layout.openPanelType === 'signal_type'}
	<SignalTypePanel entityId={layout.openPanelId} />
{:else if layout.openPanelType === 'signal_unit'}
	<SignalUnitPanel entityId={layout.openPanelId} />
{:else if layout.openPanelType === 'signal_enum'}
	<SignalEnumPanel entityId={layout.openPanelId} />
{/if}
