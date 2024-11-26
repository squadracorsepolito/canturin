<script lang="ts">
	import { colorByName } from '$lib/actions/color-name.svelte';
	import type { BusBase, NodeInterface } from '$lib/api/canturin';
	import { Attribute } from '$lib/components/attribute';
	import { UnderlinedButton } from '$lib/components/button';
	import { Collapsible } from '$lib/components/collapsible';
	import Divider from '$lib/components/divider/divider.svelte';
	import { Readonly } from '$lib/components/readonly';
	import { Select } from '$lib/components/select';
	import { Table, TableField, TableTitle } from '$lib/components/table';
	import type { PanelSectionProps } from '../types';
	import { getNodeState } from './state.svelte';

	type Props = {
		buses: BusBase[];
		int: NodeInterface;
	} & PanelSectionProps;

	let { buses, int, entityId }: Props = $props();

	const ns = getNodeState(entityId);

	function handleAttachBus(bus: BusBase) {
		ns.attachBus(int.number, bus.entityId);
	}

	let invalidBusIds = $derived.by(() => {
		if (!ns.entity.interfaces) return [];

		const ids = [];
		for (const int of ns.entity.interfaces) {
			if (int.number !== int.number) {
				ids.push(int.attachedBus.entityId);
			}
		}
		return ids;
	});

	function getCount(arr: any[] | null) {
		if (!arr) return 0;
		return arr.length;
	}
</script>

<div class="flex gap-4">
	<div use:colorByName={{ name: int.attachedBus.name }} class="block w-2 rounded-btn"></div>

	<div class="flex-1">
		<div class="grid grid-cols-2 gap-5">
			<Attribute label="Number" desc="The number of the interface">
				<Readonly>
					<span class="font-medium">{int.number}</span>
				</Readonly>
			</Attribute>

			<Attribute label="Attached Bus" desc="The bus the interface is attached to">
				<Select
					items={buses}
					valueKey="entityId"
					labelKey="name"
					bind:selected={int.attachedBus}
					filter={(item) => {
						return invalidBusIds.includes(item.entityId);
					}}
					onselect={handleAttachBus}
				/>
			</Attribute>
		</div>

		<Divider />

		<Collapsible>
			{#snippet trigger()}
				<h4>Sent Messages ({getCount(int.sentMessages)})</h4>
			{/snippet}

			{#snippet content()}
				{#if int.sentMessages}
					<Table items={int.sentMessages} idKey="entityId">
						{#snippet bulkActions({ selectedCount, selectedItems })}
							bulkactions
						{/snippet}

						{#snippet header()}
							<TableTitle>Name</TableTitle>
						{/snippet}

						{#snippet row(msg)}
							<TableField>
								<UnderlinedButton
									label={msg.name}
									onclick={() => {
										// TODO! open message panel
										console.log('open message panel');
									}}
								/>
							</TableField>
						{/snippet}

						{#snippet rowActions(msg)}
							actions
						{/snippet}
					</Table>
				{/if}
			{/snippet}
		</Collapsible>

		<Divider />

		<Collapsible>
			{#snippet trigger()}
				<h4>Received Messages ({getCount(int.receivedMessages)})</h4>
			{/snippet}

			{#snippet content()}
				{#if int.receivedMessages}
					<Table items={int.receivedMessages} idKey="entityId">
						{#snippet bulkActions({ selectedCount, selectedItems })}
							bulkactions
						{/snippet}

						{#snippet header()}
							<TableTitle>Name</TableTitle>
						{/snippet}

						{#snippet row(msg)}
							<TableField>
								<UnderlinedButton
									label={msg.name}
									onclick={() => {
										// TODO! open message panel
										console.log('open message panel');
									}}
								/>
							</TableField>
						{/snippet}

						{#snippet rowActions(msg)}
							actions
						{/snippet}
					</Table>
				{/if}
			{/snippet}
		</Collapsible>
	</div>
</div>
