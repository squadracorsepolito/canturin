<script lang="ts">
	import Divider from '$lib/components/divider/divider.svelte';
	import Select from '$lib/components/select/select.svelte';
	import SortableList from '$lib/components/sortable/sortable-list.svelte';
	import { reorder } from '@atlaskit/pragmatic-drag-and-drop/reorder';
	import {
		SignalEnumPanel,
		SignalTypePanel,
		BusPanel,
		NodePanel,
		SignalUnitPanel,
		MessagePanel
	} from '$lib/panel';
	// import MessagePanel from '$lib/panel/message-panel.svelte';
	import layout from '$lib/state/layout-state.svelte';
	import Combobox from '$lib/components/combobox/combobox.svelte';

	let items = $state([
		{
			id: 'item-1',
			label: 'Item 1'
		},
		{
			id: 'item-2',
			label: 'Item 2'
		},
		{
			id: 'item-3',
			label: 'Item 3'
		},
		{
			id: 'item-4',
			label: 'Item 4'
		},
		{
			id: 'item-5',
			label: 'Item 5'
		}
	]);

	function reorderItems(_id: string, startIndex: number, finishIndex: number) {
		console.log({ startIndex, finishIndex });

		const updatedItems = reorder({
			list: items,
			startIndex,
			finishIndex
		});
		items = updatedItems;
	}

	let selectData = [
		{ label: 'Nigeria', value: 'NG' },
		{ label: 'Japan', value: 'JP' },
		{ label: 'Korea', value: 'KO' },
		{ label: 'Kenya', value: 'KE' },
		{ label: 'United Kingdom', value: 'UK' },
		{ label: 'Ghana', value: 'GH' },
		{ label: 'Uganda', value: 'UG' }
	];

	let selected = $state('');

	let comboboxData = [
		{ label: 'Nigeria', value: 'NG' },
		{ label: 'Japan', value: 'JP' },
		{ label: 'Korea', value: 'KO' },
		{ label: 'Kenya', value: 'KE' },
		{ label: 'United Kingdom', value: 'UK' },
		{ label: 'Ghana', value: 'GH' },
		{ label: 'Uganda', value: 'UG' }
	];

	let comboboxSelected = $state('GH');
</script>

{#if layout.openPanelType === 'bus'}
	<BusPanel entityId={layout.openPanelId} />
{:else if layout.openPanelType === 'node'}
	<NodePanel entityId={layout.openPanelId} />
{:else if layout.openPanelType === 'message'}
	<MessagePanel entityId={layout.openPanelId} />
{:else if layout.openPanelType === 'signal_type'}
	<SignalTypePanel entityId={layout.openPanelId} />
{:else if layout.openPanelType === 'signal_unit'}
	<SignalUnitPanel entityId={layout.openPanelId} />
{:else if layout.openPanelType === 'signal_enum'}
	<SignalEnumPanel entityId={layout.openPanelId} />
{:else}
	<div>
		<Select name="select-test" items={selectData} bind:selected labelKey="label" valueKey="value" />

		<Divider></Divider>

		<SortableList {items} instanceId="items" reorder={reorderItems}>
			{#snippet itemBody({ item: { label } })}
				<div class="p-3">{label}</div>
			{/snippet}
		</SortableList>

		<Divider></Divider>

		<Combobox
			items={comboboxData}
			bind:selected={comboboxSelected}
			name="combobox-test"
			labelKey="label"
			valueKey="value"
		/>
	</div>
{/if}
